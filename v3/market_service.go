package aster

import (
	"context"
	"net/http"

	"github.com/your-org/go-aster/v3/futures"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/ping",
	}
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// ServerTimeService get server time
type ServerTimeService struct {
	c *Client
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...RequestOption) (serverTime int64, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/time",
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	j, err := newJSON(data)
	if err != nil {
		return 0, err
	}
	serverTime = j.Get("serverTime").MustInt64()
	return serverTime, nil
}

// ExchangeInfoService get exchange info
type ExchangeInfoService struct {
	c *Client
}

// Do send request
func (s *ExchangeInfoService) Do(ctx context.Context, opts ...RequestOption) (res *ExchangeInfo, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/exchangeInfo",
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(ExchangeInfo)
	err = json.Unmarshal(data, res)
	return res, err
}

// ExchangeInfo exchange info
type ExchangeInfo struct {
	Timezone   string                `json:"timezone"`
	ServerTime int64                 `json:"serverTime"`
	RateLimits []futures.RateLimit   `json:"rateLimits"`
	Symbols    []futures.Symbol      `json:"symbols"`
}

// DepthService get order book
type DepthService struct {
	c      *Client
	symbol string
	limit  *int
}

// Symbol set symbol
func (s *DepthService) Symbol(symbol string) *DepthService {
	s.symbol = symbol
	return s
}

// Limit set limit
func (s *DepthService) Limit(limit int) *DepthService {
	s.limit = &limit
	return s
}

// Do send request
func (s *DepthService) Do(ctx context.Context, opts ...RequestOption) (res *futures.DepthResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/depth",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(futures.DepthResponse)
	err = json.Unmarshal(data, res)
	return res, err
}

// RecentTradesListService get recent trades
type RecentTradesListService struct {
	c      *Client
	symbol string
	limit  *int
}

// Symbol set symbol
func (s *RecentTradesListService) Symbol(symbol string) *RecentTradesListService {
	s.symbol = symbol
	return s
}

// Limit set limit
func (s *RecentTradesListService) Limit(limit int) *RecentTradesListService {
	s.limit = &limit
	return s
}

// Do send request
func (s *RecentTradesListService) Do(ctx context.Context, opts ...RequestOption) (res []*futures.Trade, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/trades",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*futures.Trade, 0)
	err = json.Unmarshal(data, &res)
	return res, err
}

// AggTradesService get aggregate trades
type AggTradesService struct {
	c         *Client
	symbol    string
	fromID    *int64
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *AggTradesService) Symbol(symbol string) *AggTradesService {
	s.symbol = symbol
	return s
}

// FromID set fromID
func (s *AggTradesService) FromID(fromID int64) *AggTradesService {
	s.fromID = &fromID
	return s
}

// StartTime set startTime
func (s *AggTradesService) StartTime(startTime int64) *AggTradesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *AggTradesService) EndTime(endTime int64) *AggTradesService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *AggTradesService) Limit(limit int) *AggTradesService {
	s.limit = &limit
	return s
}

// Do send request
func (s *AggTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*futures.AggTrade, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/aggTrades",
	}
	r.setParam("symbol", s.symbol)
	if s.fromID != nil {
		r.setParam("fromId", *s.fromID)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*futures.AggTrade, 0)
	err = json.Unmarshal(data, &res)
	return res, err
}

// KlinesService get klines
type KlinesService struct {
	c         *Client
	symbol    string
	interval  futures.KlineInterval
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *KlinesService) Symbol(symbol string) *KlinesService {
	s.symbol = symbol
	return s
}

// Interval set interval
func (s *KlinesService) Interval(interval futures.KlineInterval) *KlinesService {
	s.interval = interval
	return s
}

// StartTime set startTime
func (s *KlinesService) StartTime(startTime int64) *KlinesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *KlinesService) EndTime(endTime int64) *KlinesService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *KlinesService) Limit(limit int) *KlinesService {
	s.limit = &limit
	return s
}

// Do send request
func (s *KlinesService) Do(ctx context.Context, opts ...RequestOption) (res []*futures.Kline, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/klines",
	}
	r.setParam("symbol", s.symbol)
	r.setParam("interval", s.interval)
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*futures.Kline, 0)
	err = json.Unmarshal(data, &res)
	return res, err
}

// TickerPriceService get ticker price
type TickerPriceService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *TickerPriceService) Symbol(symbol string) *TickerPriceService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *TickerPriceService) Do(ctx context.Context, opts ...RequestOption) (res []*futures.TickerPrice, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/ticker/price",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*futures.TickerPrice, 0)
	err = json.Unmarshal(data, &res)
	return res, err
}

// BookTickerService get book ticker
type BookTickerService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *BookTickerService) Symbol(symbol string) *BookTickerService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *BookTickerService) Do(ctx context.Context, opts ...RequestOption) (res []*futures.BookTicker, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/ticker/bookTicker",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*futures.BookTicker, 0)
	err = json.Unmarshal(data, &res)
	return res, err
}