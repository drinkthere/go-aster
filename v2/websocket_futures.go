package aster

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Futures WebSocket services

// WsFuturesDepthServe serves websocket depth stream for futures
func WsFuturesDepthServe(symbol string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@depth", getWsEndpoint(true, false), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		err := JSON.Unmarshal(message, event)
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

// WsFuturesPartialDepthServe serves websocket partial depth stream for futures
func WsFuturesPartialDepthServe(symbol string, levels int, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@depth%d@100ms", getWsEndpoint(true, false), strings.ToLower(symbol), levels)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		err := JSON.Unmarshal(message, event)
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

// WsFuturesKlineServe serves websocket kline stream for futures
func WsFuturesKlineServe(symbol string, interval string, handler WsFuturesKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@kline_%s", getWsEndpoint(true, false), strings.ToLower(symbol), interval)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsFuturesKlineEvent)
		err := JSON.Unmarshal(message, event)
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

// WsFuturesAggTradeServe serves websocket aggregate trade stream for futures
func WsFuturesAggTradeServe(symbol string, handler WsFuturesAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@aggTrade", getWsEndpoint(true, false), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsFuturesAggTradeEvent)
		err := JSON.Unmarshal(message, event)
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

// WsFuturesMarkPriceServe serves websocket mark price stream for futures
func WsFuturesMarkPriceServe(symbol string, handler WsFuturesMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@markPrice", getWsEndpoint(true, false), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsFuturesMarkPriceEvent)
		err := JSON.Unmarshal(message, event)
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

// WsFuturesAllMarkPriceServe serves websocket all mark price stream for futures
func WsFuturesAllMarkPriceServe(handler WsFuturesMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/!markPrice@arr", getWsEndpoint(true, false))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var events []*WsFuturesMarkPriceEvent
		err := JSON.Unmarshal(message, &events)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}
		for _, event := range events {
			handler(event)
		}
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsFuturesBookTickerServe serves websocket book ticker stream for futures
func WsFuturesBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@bookTicker", getWsEndpoint(true, false), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := JSON.Unmarshal(message, event)
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

// WsFuturesAllBookTickerServe serves websocket all book tickers stream for futures
func WsFuturesAllBookTickerServe(handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/!bookTicker", getWsEndpoint(true, false))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := JSON.Unmarshal(message, event)
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

// WsFuturesUserDataServe serves websocket user data stream for futures
func WsFuturesUserDataServe(listenKey string, handler WsFuturesUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s", getWsEndpoint(true, false), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		// First check the event type
		var eventType struct {
			Event string `json:"e"`
		}
		err := JSON.Unmarshal(message, &eventType)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}

		event := new(WsFuturesUserDataEvent)
		event.Event = eventType.Event

		switch eventType.Event {
		case "ACCOUNT_UPDATE":
			var accountUpdate WsFuturesAccountUpdate
			err = JSON.Unmarshal(message, &accountUpdate)
			if err == nil {
				event.Time = accountUpdate.Time
				event.TransactionTime = accountUpdate.TransactionTime
				event.AccountUpdate = &accountUpdate
			}
		case "ORDER_TRADE_UPDATE":
			var orderData struct {
				Event           string              `json:"e"`
				Time            int64               `json:"E"`
				TransactionTime int64               `json:"T"`
				Order           WsFuturesOrderUpdate `json:"o"`
			}
			err = JSON.Unmarshal(message, &orderData)
			if err == nil {
				event.Time = orderData.Time
				event.TransactionTime = orderData.TransactionTime
				event.OrderUpdate = &orderData.Order
			}
		case "ACCOUNT_CONFIG_UPDATE":
			var configData struct {
				Event                string                       `json:"e"`
				Time                 int64                        `json:"E"`
				TransactionTime      int64                        `json:"T"`
				AccountConfigUpdate  WsFuturesAccountConfigUpdate `json:"ac"`
			}
			err = JSON.Unmarshal(message, &configData)
			if err == nil {
				event.Time = configData.Time
				event.TransactionTime = configData.TransactionTime
				event.AccountConfigUpdate = &configData.AccountConfigUpdate
			}
		case "MARGIN_CALL":
			var marginCall WsFuturesMarginCall
			err = JSON.Unmarshal(message, &marginCall)
			if err == nil {
				event.MarginCall = &marginCall
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

// Combined streams for futures

// WsCombinedFuturesDepthServe serves websocket combined depth stream for futures
func WsCombinedFuturesDepthServe(symbols []string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var streams []string
	for _, s := range symbols {
		streams = append(streams, fmt.Sprintf("%s@depth", strings.ToLower(s)))
	}
	endpoint := fmt.Sprintf("%s?streams=%s", combinedFuturesBaseURL, strings.Join(streams, "/"))
	return wsCombinedFuturesDepthServe(endpoint, handler, errHandler)
}

// Internal function for combined futures depth
func wsCombinedFuturesDepthServe(endpoint string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var combinedEvent struct {
			Stream string          `json:"stream"`
			Data   json.RawMessage `json:"data"`
		}
		err := JSON.Unmarshal(message, &combinedEvent)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
			}
			return
		}

		event := new(WsDepthEvent)
		err = JSON.Unmarshal(combinedEvent.Data, event)
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