package aster

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = false
	// DefaultWebsocketURL is the default websocket base URL
	DefaultWebsocketURL = "wss://fstream.asterdex.com"
)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint string
}

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConn represents a websocket connection
type WsConn struct {
	sync.Mutex
	conn             *websocket.Conn
	ctx              context.Context
	cancel           context.CancelFunc
	hub              *WsHub
	closed           bool
	wg               sync.WaitGroup
}

// WsHub manages all websocket connections
type WsHub struct {
	connections map[string]*WsConn
	mu          sync.RWMutex
}

// NewWsHub creates a new websocket hub
func NewWsHub() *WsHub {
	return &WsHub{
		connections: make(map[string]*WsConn),
	}
}

// GetConnection gets a connection by ID
func (h *WsHub) GetConnection(id string) (*WsConn, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conn, ok := h.connections[id]
	return conn, ok
}

// AddConnection adds a new connection
func (h *WsHub) AddConnection(id string, conn *WsConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[id] = conn
}

// RemoveConnection removes a connection
func (h *WsHub) RemoveConnection(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.connections, id)
}

// newWsConfig creates a new websocket configuration
func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

// newWsConn creates a new websocket connection
func newWsConn(ctx context.Context, hub *WsHub, endpoint string) (*WsConn, error) {
	dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		HandshakeTimeout:  45 * time.Second,
	}
	
	conn, _, err := dialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	
	return &WsConn{
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
		hub:    hub,
	}, nil
}

// sendMessage sends a message to the websocket connection
func (c *WsConn) sendMessage(msg interface{}) error {
	c.Lock()
	defer c.Unlock()
	
	if c.closed {
		return errors.New("connection is closed")
	}
	
	return c.conn.WriteJSON(msg)
}

// close closes the websocket connection
func (c *WsConn) close() error {
	c.Lock()
	defer c.Unlock()
	
	if c.closed {
		return nil
	}
	
	c.closed = true
	c.cancel()
	
	// Send close message
	deadline := time.Now().Add(time.Second)
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	c.conn.WriteControl(websocket.CloseMessage, msg, deadline)
	
	return c.conn.Close()
}

// keepAlive sends ping messages to keep the connection alive
func (c *WsConn) keepAlive() {
	ticker := time.NewTicker(WebsocketTimeout)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.Lock()
			if c.closed {
				c.Unlock()
				return
			}
			deadline := time.Now().Add(10 * time.Second)
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, deadline); err != nil {
				c.Unlock()
				return
			}
			c.Unlock()
		case <-c.ctx.Done():
			return
		}
	}
}

// readMessages reads messages from the websocket connection
func (c *WsConn) readMessages(handler WsHandler, errHandler ErrHandler) {
	defer c.wg.Done()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msgType, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					errHandler(fmt.Errorf("websocket error: %v", err))
				}
				return
			}
			
			if msgType == websocket.TextMessage {
				handler(message)
			}
		}
	}
}

// WsServe starts a websocket connection and handles messages
func WsServe(ctx context.Context, endpoint string, handler WsHandler, errHandler ErrHandler) (*WsConn, error) {
	hub := NewWsHub()
	conn, err := newWsConn(ctx, hub, endpoint)
	if err != nil {
		return nil, err
	}
	
	conn.wg.Add(1)
	go conn.readMessages(handler, errHandler)
	
	if WebsocketKeepalive {
		go conn.keepAlive()
	}
	
	return conn, nil
}

// Close closes the websocket connection
func (c *WsConn) Close() error {
	err := c.close()
	c.wg.Wait()
	return err
}