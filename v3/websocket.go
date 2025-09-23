package aster

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handle error
type ErrHandler func(err error)

// WsConfig websocket configuration
type WsConfig struct {
	Endpoint string
}

// WsDepthHandler websocket depth handler
type WsDepthHandler func(event *WsDepthEvent)

// WsDepthEvent websocket depth event
type WsDepthEvent struct {
	Event         string          `json:"e"`
	Time          int64           `json:"E"`
	Symbol        string          `json:"s"`
	FirstUpdateID int64           `json:"U"`
	LastUpdateID  int64           `json:"u"`
	Bids          [][]string      `json:"b"`
	Asks          [][]string      `json:"a"`
}

// WsAggTradeHandler websocket aggregate trade handler
type WsAggTradeHandler func(event *WsAggTradeEvent)

// WsAggTradeEvent websocket aggregate trade event
type WsAggTradeEvent struct {
	Event            string `json:"e"`
	Time             int64  `json:"E"`
	Symbol           string `json:"s"`
	AggregateTradeID int64  `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeID     int64  `json:"f"`
	LastTradeID      int64  `json:"l"`
	TradeTime        int64  `json:"T"`
	IsBuyerMaker     bool   `json:"m"`
}

// WsKlineHandler websocket kline handler
type WsKlineHandler func(event *WsKlineEvent)

// WsKlineEvent websocket kline event
type WsKlineEvent struct {
	Event  string    `json:"e"`
	Time   int64     `json:"E"`
	Symbol string    `json:"s"`
	Kline  WsKline   `json:"k"`
}

// WsKline websocket kline
type WsKline struct {
	StartTime                int64  `json:"t"`
	EndTime                  int64  `json:"T"`
	Symbol                   string `json:"s"`
	Interval                 string `json:"i"`
	FirstTradeID             int64  `json:"f"`
	LastTradeID              int64  `json:"L"`
	Open                     string `json:"o"`
	Close                    string `json:"c"`
	High                     string `json:"h"`
	Low                      string `json:"l"`
	Volume                   string `json:"v"`
	TradeNum                 int64  `json:"n"`
	IsFinal                  bool   `json:"x"`
	QuoteVolume              string `json:"q"`
	ActiveBuyVolume          string `json:"V"`
	ActiveBuyQuoteVolume     string `json:"Q"`
}

// WsDepthServe serve websocket depth handler
func WsDepthServe(wsURL string, symbol string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@depth", wsURL, symbol)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAggTradeServe serve websocket aggregate trade handler
func WsAggTradeServe(wsURL string, symbol string, handler WsAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@aggTrade", wsURL, symbol)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsAggTradeEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsKlineServe serve websocket kline handler
func WsKlineServe(wsURL string, symbol string, interval string, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@kline_%s", wsURL, symbol, interval)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsKlineEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

func wsServe(cfg *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	c, _, err := websocket.DefaultDialer.Dial(cfg.Endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer func() {
			c.Close()
			close(doneC)
		}()
		if err := c.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
			errHandler(err)
			return
		}
		c.SetPongHandler(func(string) error {
			if err := c.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
				errHandler(err)
			}
			return nil
		})
		for {
			select {
			case <-stopC:
				return
			default:
				msgType, message, err := c.ReadMessage()
				if err != nil {
					errHandler(err)
					return
				}
				if msgType != websocket.TextMessage {
					continue
				}
				handler(message)
			}
		}
	}()
	return doneC, stopC, nil
}