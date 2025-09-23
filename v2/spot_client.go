package aster

import (
	"github.com/drinkthere/go-aster/v2/common"
)

// SpotClient defines spot client
type SpotClient struct {
	*BaseClient
}

// NewSpot creates a new spot client
func NewSpot(apiKey, secretKey string, opts ...ClientOption) *SpotClient {
	baseClient := NewSpotClient(apiKey, secretKey, opts...)
	return &SpotClient{BaseClient: baseClient}
}

// Account and trading endpoints
func (c *SpotClient) NewCreateOrderService() *CreateSpotOrderService {
	return &CreateSpotOrderService{c: c.BaseClient}
}

func (c *SpotClient) NewGetOrderService() *GetSpotOrderService {
	return &GetSpotOrderService{c: c.BaseClient}
}

func (c *SpotClient) NewCancelOrderService() *CancelSpotOrderService {
	return &CancelSpotOrderService{c: c.BaseClient}
}

func (c *SpotClient) NewCancelOpenOrdersService() *CancelOpenSpotOrdersService {
	return &CancelOpenSpotOrdersService{c: c.BaseClient}
}

func (c *SpotClient) NewListOpenOrdersService() *ListOpenSpotOrdersService {
	return &ListOpenSpotOrdersService{c: c.BaseClient}
}

func (c *SpotClient) NewListOrdersService() *ListSpotOrdersService {
	return &ListSpotOrdersService{c: c.BaseClient}
}

func (c *SpotClient) NewGetAccountService() *GetSpotAccountService {
	return &GetSpotAccountService{c: c.BaseClient}
}

func (c *SpotClient) NewListTradesService() *ListSpotTradesService {
	return &ListSpotTradesService{c: c.BaseClient}
}

// Market data endpoints
func (c *SpotClient) NewPingService() *SpotPingService {
	return &SpotPingService{c: c.BaseClient}
}

func (c *SpotClient) NewServerTimeService() *SpotServerTimeService {
	return &SpotServerTimeService{c: c.BaseClient}
}

func (c *SpotClient) NewExchangeInfoService() *SpotExchangeInfoService {
	return &SpotExchangeInfoService{c: c.BaseClient}
}

func (c *SpotClient) NewDepthService() *SpotDepthService {
	return &SpotDepthService{c: c.BaseClient}
}

func (c *SpotClient) NewAggTradesService() *SpotAggTradesService {
	return &SpotAggTradesService{c: c.BaseClient}
}

func (c *SpotClient) NewRecentTradesListService() *SpotRecentTradesListService {
	return &SpotRecentTradesListService{c: c.BaseClient}
}

func (c *SpotClient) NewKlinesService() *SpotKlinesService {
	return &SpotKlinesService{c: c.BaseClient}
}

func (c *SpotClient) NewListPriceChangeStatsService() *SpotListPriceChangeStatsService {
	return &SpotListPriceChangeStatsService{c: c.BaseClient}
}

func (c *SpotClient) NewListPricesService() *SpotListPricesService {
	return &SpotListPricesService{c: c.BaseClient}
}

func (c *SpotClient) NewListBookTickersService() *SpotListBookTickersService {
	return &SpotListBookTickersService{c: c.BaseClient}
}

// User stream endpoints
func (c *SpotClient) NewStartUserStreamService() *StartSpotUserStreamService {
	return &StartSpotUserStreamService{c: c.BaseClient}
}

func (c *SpotClient) NewKeepaliveUserStreamService() *KeepaliveSpotUserStreamService {
	return &KeepaliveSpotUserStreamService{c: c.BaseClient}
}

func (c *SpotClient) NewCloseUserStreamService() *CloseSpotUserStreamService {
	return &CloseSpotUserStreamService{c: c.BaseClient}
}

// WebSocket streams
func (c *SpotClient) WsDepthServe(symbol string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return WsSpotDepthServe(symbol, handler, errHandler)
}

func (c *SpotClient) WsKlineServe(symbol string, interval common.Interval, handler WsSpotKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return WsSpotKlineServe(symbol, string(interval), handler, errHandler)
}

func (c *SpotClient) WsAggTradeServe(symbol string, handler WsSpotAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return WsSpotAggTradeServe(symbol, handler, errHandler)
}

func (c *SpotClient) WsBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return WsSpotBookTickerServe(symbol, handler, errHandler)
}

func (c *SpotClient) WsAllMarketsStatServe(handler WsSpotAllMarketsStatHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return WsSpotAllMarketsStatServe(handler, errHandler)
}

func (c *SpotClient) WsUserDataServe(listenKey string, handler WsSpotUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return WsSpotUserDataServe(listenKey, handler, errHandler)
}