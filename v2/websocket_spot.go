package aster

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Spot WebSocket services

// WsSpotDepthServe serves websocket depth stream
func WsSpotDepthServe(symbol string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@depth", getWsEndpoint(false, false), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotPartialDepthServe serves websocket partial depth stream
func WsSpotPartialDepthServe(symbol string, levels int, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@depth%d", getWsEndpoint(false, false), strings.ToLower(symbol), levels)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotKlineServe serves websocket kline stream
func WsSpotKlineServe(symbol string, interval string, handler WsSpotKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@kline_%s", getWsEndpoint(false, false), strings.ToLower(symbol), interval)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsSpotKlineEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotAggTradeServe serves websocket aggregate trade stream
func WsSpotAggTradeServe(symbol string, handler WsSpotAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@aggTrade", getWsEndpoint(false, false), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsSpotAggTradeEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotBookTickerServe serves websocket book ticker stream
func WsSpotBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@bookTicker", getWsEndpoint(false, false), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotAllBookTickerServe serves websocket all book tickers stream
func WsSpotAllBookTickerServe(handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/!bookTicker", getWsEndpoint(false, false))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotAllMarketsStatServe serves websocket 24hr statistics stream for all markets
func WsSpotAllMarketsStatServe(handler WsSpotAllMarketsStatHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/!ticker@arr", getWsEndpoint(false, false))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsSpotAllMarketsStatEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotUserDataServe serves websocket user data stream
func WsSpotUserDataServe(listenKey string, handler WsSpotUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s", getWsEndpoint(false, false), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		// First check the event type
		var eventType struct {
			Event string `json:"e"`
		}
		err := json.Unmarshal(message, &eventType)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}

		event := new(WsSpotUserDataEvent)
		event.Event = eventType.Event
		
		switch eventType.Event {
		case "outboundAccountPosition":
			var accountUpdate struct {
				Event         string          `json:"e"`
				Time          int64           `json:"E"`
				AccountUpdate WsSpotAccountUpdate `json:"u"`
			}
			err = json.Unmarshal(message, &accountUpdate)
			if err == nil {
				event.Time = accountUpdate.Time
				event.AccountUpdate = &accountUpdate.AccountUpdate
			}
		case "executionReport":
			var orderUpdate WsSpotOrderUpdate
			err = json.Unmarshal(message, &orderUpdate)
			if err == nil {
				event.OrderUpdate = &orderUpdate
			}
		}

		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}

		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// Combined streams

// WsCombinedSpotDepthServe serves websocket combined depth stream
func WsCombinedSpotDepthServe(symbols []string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var streams []string
	for _, s := range symbols {
		streams = append(streams, fmt.Sprintf("%s@depth", strings.ToLower(s)))
	}
	endpoint := fmt.Sprintf("%s?streams=%s", combinedBaseURL, strings.Join(streams, "/"))
	return wsCombinedSpotDepthServe(endpoint, handler, errHandler)
}

// WsCombinedSpotBookTickerServe serves websocket combined book ticker stream
func WsCombinedSpotBookTickerServe(symbols []string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var streams []string
	for _, s := range symbols {
		streams = append(streams, fmt.Sprintf("%s@bookTicker", strings.ToLower(s)))
	}
	endpoint := fmt.Sprintf("%s?streams=%s", combinedBaseURL, strings.Join(streams, "/"))
	return wsCombinedSpotBookTickerServe(endpoint, handler, errHandler)
}

// Internal function for combined depth
func wsCombinedSpotDepthServe(endpoint string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var combinedEvent struct {
			Stream string          `json:"stream"`
			Data   json.RawMessage `json:"data"`
		}
		err := json.Unmarshal(message, &combinedEvent)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}

		event := new(WsDepthEvent)
		err = json.Unmarshal(combinedEvent.Data, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// Internal function for combined book ticker
func wsCombinedSpotBookTickerServe(endpoint string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var combinedEvent struct {
			Stream string          `json:"stream"`
			Data   json.RawMessage `json:"data"`
		}
		err := json.Unmarshal(message, &combinedEvent)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}

		event := new(WsBookTickerEvent)
		err = json.Unmarshal(combinedEvent.Data, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WebSocket functions with LocalAddress support

// WsSpotDepthServeWithLocalAddr serves websocket depth stream with local address binding
func WsSpotDepthServeWithLocalAddr(symbol string, handler WsDepthHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@depth", getWsEndpoint(false, false), strings.ToLower(symbol))
	var cfg *WsConfig
	if localAddr != "" {
		cfg = newWsConfigWithIP(endpoint, localAddr)
	} else {
		cfg = newWsConfig(endpoint)
	}
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotKlineServeWithLocalAddr serves websocket kline stream with local address binding
func WsSpotKlineServeWithLocalAddr(symbol string, interval string, handler WsSpotKlineHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@kline_%s", getWsEndpoint(false, false), strings.ToLower(symbol), interval)
	var cfg *WsConfig
	if localAddr != "" {
		cfg = newWsConfigWithIP(endpoint, localAddr)
	} else {
		cfg = newWsConfig(endpoint)
	}
	wsHandler := func(message []byte) {
		event := new(WsSpotKlineEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotAggTradeServeWithLocalAddr serves websocket aggregate trade stream with local address binding
func WsSpotAggTradeServeWithLocalAddr(symbol string, handler WsSpotAggTradeHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@aggTrade", getWsEndpoint(false, false), strings.ToLower(symbol))
	var cfg *WsConfig
	if localAddr != "" {
		cfg = newWsConfigWithIP(endpoint, localAddr)
	} else {
		cfg = newWsConfig(endpoint)
	}
	wsHandler := func(message []byte) {
		event := new(WsSpotAggTradeEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotBookTickerServeWithLocalAddr serves websocket book ticker stream with local address binding
func WsSpotBookTickerServeWithLocalAddr(symbol string, handler WsBookTickerHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@bookTicker", getWsEndpoint(false, false), strings.ToLower(symbol))
	var cfg *WsConfig
	if localAddr != "" {
		cfg = newWsConfigWithIP(endpoint, localAddr)
	} else {
		cfg = newWsConfig(endpoint)
	}
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotAllMarketsStatServeWithLocalAddr serves websocket all markets statistics stream with local address binding
func WsSpotAllMarketsStatServeWithLocalAddr(handler WsSpotAllMarketsStatHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/!ticker@arr", getWsEndpoint(false, false))
	var cfg *WsConfig
	if localAddr != "" {
		cfg = newWsConfigWithIP(endpoint, localAddr)
	} else {
		cfg = newWsConfig(endpoint)
	}
	wsHandler := func(message []byte) {
		var event WsSpotAllMarketsStatEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsSpotUserDataServeWithLocalAddr serves websocket user data stream with local address binding
func WsSpotUserDataServeWithLocalAddr(listenKey string, handler WsSpotUserDataHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s", getWsEndpoint(false, false), listenKey)
	var cfg *WsConfig
	if localAddr != "" {
		cfg = newWsConfigWithIP(endpoint, localAddr)
	} else {
		cfg = newWsConfig(endpoint)
	}
	wsHandler := func(message []byte) {
		event := new(WsSpotUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedSpotBookTickerServeWithLocalAddr serves websocket combined book ticker stream with local address binding
func WsCombinedSpotBookTickerServeWithLocalAddr(symbols []string, handler WsBookTickerHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	var streams []string
	for _, s := range symbols {
		streams = append(streams, fmt.Sprintf("%s@bookTicker", strings.ToLower(s)))
	}
	endpoint := fmt.Sprintf("%s?streams=%s", combinedBaseURL, strings.Join(streams, "/"))
	var cfg *WsConfig
	if localAddr != "" {
		cfg = newWsConfigWithIP(endpoint, localAddr)
	} else {
		cfg = newWsConfig(endpoint)
	}
	wsHandler := func(message []byte) {
		var combinedEvent struct {
			Stream string          `json:"stream"`
			Data   json.RawMessage `json:"data"`
		}
		err := json.Unmarshal(message, &combinedEvent)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}

		event := new(WsBookTickerEvent)
		err = json.Unmarshal(combinedEvent.Data, event)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}