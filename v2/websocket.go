package aster

import (
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint string
	IP       string
	Resolver *net.Resolver
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

func newWsConfigWithIP(endpoint string, localIP string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
		IP:       localIP,
	}
}

func (cfg *WsConfig) WithIP(ip string) {
	cfg.IP = ip
}

func (cfg *WsConfig) WithResolver(resolver *net.Resolver) {
	cfg.Resolver = resolver
}

var wsServe = func(cfg *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var Dialer websocket.Dialer
	if cfg.IP == "" {
		Dialer = websocket.Dialer{
			Proxy:             http.ProxyFromEnvironment,
			HandshakeTimeout:  45 * time.Second,
			EnableCompression: false,
		}
	} else {
		Dialer = websocket.Dialer{
			NetDial: func(network, addr string) (net.Conn, error) {
				localAddr, err := net.ResolveTCPAddr("tcp", cfg.IP+":0")
				if err != nil {
					return nil, err
				}
				d := net.Dialer{LocalAddr: localAddr}
				return d.Dial(network, addr)
			},
			Proxy:             http.ProxyFromEnvironment,
			HandshakeTimeout:  45 * time.Second,
			EnableCompression: false,
		}
	}

	c, _, err := Dialer.Dial(cfg.Endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		defer func() {
			c.Close()
			close(doneC)
		}()
		
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(10 * time.Second))
			return nil
		})
		
		for {
			select {
			case <-stopC:
				return
			default:
				messageType, message, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						if errHandler != nil {
							errHandler(err)
						}
					}
					return
				}
				
				if messageType == websocket.TextMessage {
					handler(message)
				}
			}
		}
	}()
	
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				err := c.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					return
				}
			case <-stopC:
				return
			}
		}
	}()
	
	return doneC, stopC, nil
}