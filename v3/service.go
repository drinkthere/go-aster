package aster

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func newJSON(data []byte) (j *jsoniter.Any, err error) {
	j = new(jsoniter.Any)
	err = json.Unmarshal(data, j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// NewPingService init ping service
func (c *Client) NewPingService() *PingService {
	return &PingService{c: c}
}

// NewServerTimeService init server time service
func (c *Client) NewServerTimeService() *ServerTimeService {
	return &ServerTimeService{c: c}
}

// NewExchangeInfoService init exchange info service
func (c *Client) NewExchangeInfoService() *ExchangeInfoService {
	return &ExchangeInfoService{c: c}
}

// NewDepthService init depth service
func (c *Client) NewDepthService() *DepthService {
	return &DepthService{c: c}
}

// NewRecentTradesListService init recent trades list service
func (c *Client) NewRecentTradesListService() *RecentTradesListService {
	return &RecentTradesListService{c: c}
}

// NewAggTradesService init agg trades service
func (c *Client) NewAggTradesService() *AggTradesService {
	return &AggTradesService{c: c}
}

// NewKlinesService init klines service
func (c *Client) NewKlinesService() *KlinesService {
	return &KlinesService{c: c}
}

// NewTickerPriceService init ticker price service
func (c *Client) NewTickerPriceService() *TickerPriceService {
	return &TickerPriceService{c: c}
}

// NewBookTickerService init book ticker service
func (c *Client) NewBookTickerService() *BookTickerService {
	return &BookTickerService{c: c}
}

// NewCreateOrderService init create order service
func (c *Client) NewCreateOrderService() *CreateOrderService {
	return &CreateOrderService{c: c}
}

// NewGetOrderService init get order service
func (c *Client) NewGetOrderService() *GetOrderService {
	return &GetOrderService{c: c}
}

// NewCancelOrderService init cancel order service
func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

// NewCancelAllOpenOrdersService init cancel all open orders service
func (c *Client) NewCancelAllOpenOrdersService() *CancelAllOpenOrdersService {
	return &CancelAllOpenOrdersService{c: c}
}

// NewListOpenOrdersService init list open orders service
func (c *Client) NewListOpenOrdersService() *ListOpenOrdersService {
	return &ListOpenOrdersService{c: c}
}

// NewGetAccountService init get account service
func (c *Client) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

// NewGetBalanceService init get balance service
func (c *Client) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c}
}

// NewGetPositionsService init get positions service
func (c *Client) NewGetPositionsService() *GetPositionsService {
	return &GetPositionsService{c: c}
}

// NewStartUserStreamService init start user stream service
func (c *Client) NewStartUserStreamService() *StartUserStreamService {
	return &StartUserStreamService{c: c}
}

// NewKeepaliveUserStreamService init keepalive user stream service
func (c *Client) NewKeepaliveUserStreamService() *KeepaliveUserStreamService {
	return &KeepaliveUserStreamService{c: c}
}

// NewCloseUserStreamService init close user stream service
func (c *Client) NewCloseUserStreamService() *CloseUserStreamService {
	return &CloseUserStreamService{c: c}
}