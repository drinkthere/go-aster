package aster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetWsEndpoint returns the websocket endpoint
func GetWsEndpoint() string {
	return DefaultWebsocketURL
}

// WsDepthServe serve websocket depth handler with depth
func WsDepthServe(symbol string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@depth", GetWsEndpoint(), strings.ToLower(symbol))
	return wsDepthServe(endpoint, handler, errHandler)
}

// WsPartialDepthServe serve websocket depth handler with partial depth
func WsPartialDepthServe(symbol string, levels int, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@depth%d@100ms", GetWsEndpoint(), strings.ToLower(symbol), levels)
	return wsDepthServe(endpoint, handler, errHandler)
}

// WsCombinedDepthServe serves combined depth handler
func WsCombinedDepthServe(symbols []string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var streams []string
	for _, s := range symbols {
		streams = append(streams, fmt.Sprintf("%s@depth", strings.ToLower(s)))
	}
	endpoint := fmt.Sprintf("%s/stream?streams=%s", GetWsEndpoint(), strings.Join(streams, "/"))
	return wsDepthServe(endpoint, handler, errHandler)
}

func wsDepthServe(endpoint string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	
	conn, err := WsServe(ctx, endpoint, wsHandler, errHandler)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	
	go func() {
		select {
		case <-stopC:
			cancel()
			conn.Close()
		case <-ctx.Done():
			conn.Close()
		}
		close(doneC)
	}()
	
	return doneC, stopC, nil
}

// WsAggTradeServe serve websocket aggregate trade handler
func WsAggTradeServe(symbol string, handler WsAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@aggTrade", GetWsEndpoint(), strings.ToLower(symbol))
	return wsAggTradeServe(endpoint, handler, errHandler)
}

// WsCombinedAggTradeServe serves combined aggregate trade handler
func WsCombinedAggTradeServe(symbols []string, handler WsAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var streams []string
	for _, s := range symbols {
		streams = append(streams, fmt.Sprintf("%s@aggTrade", strings.ToLower(s)))
	}
	endpoint := fmt.Sprintf("%s/stream?streams=%s", GetWsEndpoint(), strings.Join(streams, "/"))
	return wsAggTradeServe(endpoint, handler, errHandler)
}

func wsAggTradeServe(endpoint string, handler WsAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	wsHandler := func(message []byte) {
		event := new(WsAggTradeEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	
	conn, err := WsServe(ctx, endpoint, wsHandler, errHandler)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	
	go func() {
		select {
		case <-stopC:
			cancel()
			conn.Close()
		case <-ctx.Done():
			conn.Close()
		}
		close(doneC)
	}()
	
	return doneC, stopC, nil
}

// WsKlineServe serve websocket kline handler
func WsKlineServe(symbol string, interval string, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@kline_%s", GetWsEndpoint(), strings.ToLower(symbol), interval)
	return wsKlineServe(endpoint, handler, errHandler)
}

// WsCombinedKlineServe serves combined kline handler
func WsCombinedKlineServe(symbolIntervalPairs map[string]string, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var streams []string
	for symbol, interval := range symbolIntervalPairs {
		streams = append(streams, fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval))
	}
	endpoint := fmt.Sprintf("%s/stream?streams=%s", GetWsEndpoint(), strings.Join(streams, "/"))
	return wsKlineServe(endpoint, handler, errHandler)
}

func wsKlineServe(endpoint string, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	wsHandler := func(message []byte) {
		event := new(WsKlineEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	
	conn, err := WsServe(ctx, endpoint, wsHandler, errHandler)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	
	go func() {
		select {
		case <-stopC:
			cancel()
			conn.Close()
		case <-ctx.Done():
			conn.Close()
		}
		close(doneC)
	}()
	
	return doneC, stopC, nil
}

// WsBookTickerServe serve websocket book ticker handler
func WsBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s@bookTicker", GetWsEndpoint(), strings.ToLower(symbol))
	return wsBookTickerServe(endpoint, handler, errHandler)
}

// WsAllBookTickerServe serve websocket all book ticker handler
func WsAllBookTickerServe(handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/!bookTicker", GetWsEndpoint())
	return wsBookTickerServe(endpoint, handler, errHandler)
}

// WsCombinedBookTickerServe serves combined book ticker handler
func WsCombinedBookTickerServe(symbols []string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var streams []string
	for _, s := range symbols {
		streams = append(streams, fmt.Sprintf("%s@bookTicker", strings.ToLower(s)))
	}
	endpoint := fmt.Sprintf("%s/stream?streams=%s", GetWsEndpoint(), strings.Join(streams, "/"))
	return wsBookTickerServe(endpoint, handler, errHandler)
}

func wsBookTickerServe(endpoint string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	
	conn, err := WsServe(ctx, endpoint, wsHandler, errHandler)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	
	go func() {
		select {
		case <-stopC:
			cancel()
			conn.Close()
		case <-ctx.Done():
			conn.Close()
		}
		close(doneC)
	}()
	
	return doneC, stopC, nil
}

// WsUserDataServe serve websocket user data handler
func WsUserDataServe(listenKey string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/ws/%s", GetWsEndpoint(), listenKey)
	return wsUserDataServe(endpoint, handler, errHandler)
}

// WsCombinedUserDataServe serves combined user data handler
func WsCombinedUserDataServe(listenKeys []string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/stream?streams=%s", GetWsEndpoint(), strings.Join(listenKeys, "/"))
	return wsUserDataServe(endpoint, handler, errHandler)
}

func wsUserDataServe(endpoint string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	wsHandler := func(message []byte) {
		event := new(WsUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	
	conn, err := WsServe(ctx, endpoint, wsHandler, errHandler)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	
	go func() {
		select {
		case <-stopC:
			cancel()
			conn.Close()
		case <-ctx.Done():
			conn.Close()
		}
		close(doneC)
	}()
	
	return doneC, stopC, nil
}