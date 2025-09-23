package aster

import (
	"github.com/drinkthere/go-aster/v2/common"
)

// FuturesClient defines futures client
type FuturesClient struct {
	*BaseClient
}

// NewFutures creates a new futures client
func NewFutures(apiKey, secretKey string, opts ...ClientOption) *FuturesClient {
	baseClient := NewFuturesClient(apiKey, secretKey, opts...)
	return &FuturesClient{BaseClient: baseClient}
}

// WebSocket streams
func (c *FuturesClient) WsDepthServe(symbol string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return c.WsDepthServeWithLocalAddr(symbol, handler, errHandler, c.LocalAddress)
}

func (c *FuturesClient) WsKlineServe(symbol string, interval common.Interval, handler WsFuturesKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return c.WsKlineServeWithLocalAddr(symbol, string(interval), handler, errHandler, c.LocalAddress)
}

func (c *FuturesClient) WsAggTradeServe(symbol string, handler WsFuturesAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return c.WsAggTradeServeWithLocalAddr(symbol, handler, errHandler, c.LocalAddress)
}

func (c *FuturesClient) WsBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return c.WsBookTickerServeWithLocalAddr(symbol, handler, errHandler, c.LocalAddress)
}

func (c *FuturesClient) WsMarkPriceServe(symbol string, handler WsFuturesMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return c.WsMarkPriceServeWithLocalAddr(symbol, handler, errHandler, c.LocalAddress)
}

func (c *FuturesClient) WsUserDataServe(listenKey string, handler WsFuturesUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return c.WsUserDataServeWithLocalAddr(listenKey, handler, errHandler, c.LocalAddress)
}

func (c *FuturesClient) WsCombinedBookTickerServe(symbols []string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return WsCombinedFuturesBookTickerServeWithLocalAddr(symbols, handler, errHandler, c.LocalAddress)
}

// WebSocket streams with LocalAddress support
func (c *FuturesClient) WsDepthServeWithLocalAddr(symbol string, handler WsDepthHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	return WsFuturesDepthServeWithLocalAddr(symbol, handler, errHandler, localAddr)
}

func (c *FuturesClient) WsKlineServeWithLocalAddr(symbol string, interval string, handler WsFuturesKlineHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	return WsFuturesKlineServeWithLocalAddr(symbol, interval, handler, errHandler, localAddr)
}

func (c *FuturesClient) WsAggTradeServeWithLocalAddr(symbol string, handler WsFuturesAggTradeHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	return WsFuturesAggTradeServeWithLocalAddr(symbol, handler, errHandler, localAddr)
}

func (c *FuturesClient) WsBookTickerServeWithLocalAddr(symbol string, handler WsBookTickerHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	return WsFuturesBookTickerServeWithLocalAddr(symbol, handler, errHandler, localAddr)
}

func (c *FuturesClient) WsMarkPriceServeWithLocalAddr(symbol string, handler WsFuturesMarkPriceHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	return WsFuturesMarkPriceServeWithLocalAddr(symbol, handler, errHandler, localAddr)
}

func (c *FuturesClient) WsUserDataServeWithLocalAddr(listenKey string, handler WsFuturesUserDataHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	return WsFuturesUserDataServeWithLocalAddr(listenKey, handler, errHandler, localAddr)
}

func (c *FuturesClient) WsCombinedBookTickerServeWithLocalAddr(symbols []string, handler WsBookTickerHandler, errHandler ErrHandler, localAddr string) (doneC, stopC chan struct{}, err error) {
	return WsCombinedFuturesBookTickerServeWithLocalAddr(symbols, handler, errHandler, localAddr)
}