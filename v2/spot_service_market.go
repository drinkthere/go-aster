package aster

import (
	"context"
	"net/http"

	"github.com/drinkthere/go-aster/v2/common"
)

// SpotPingService ping server
type SpotPingService struct {
	c *BaseClient
}

// Do send request
func (s *SpotPingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := newRequest(http.MethodGet, "/api/v3/ping", secTypeNone)
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// SpotServerTimeService get server time
type SpotServerTimeService struct {
	c *BaseClient
}

// Do send request
func (s *SpotServerTimeService) Do(ctx context.Context, opts ...RequestOption) (serverTime int64, err error) {
	r := newRequest(http.MethodGet, "/api/v3/time", secTypeNone)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	j, err := newJSON(data)
	if err != nil {
		return 0, err
	}
	serverTime = (*j).Get("serverTime").ToInt64()
	return serverTime, nil
}

// SpotExchangeInfoService get exchange info
type SpotExchangeInfoService struct {
	c       *BaseClient
	symbol  *string
	symbols []string
}

// Symbol set symbol
func (s *SpotExchangeInfoService) Symbol(symbol string) *SpotExchangeInfoService {
	s.symbol = &symbol
	return s
}

// Symbols set symbols
func (s *SpotExchangeInfoService) Symbols(symbols ...string) *SpotExchangeInfoService {
	s.symbols = symbols
	return s
}

// Do send request
func (s *SpotExchangeInfoService) Do(ctx context.Context, opts ...RequestOption) (res *SpotExchangeInfo, err error) {
	r := newRequest(http.MethodGet, "/api/v3/exchangeInfo", secTypeNone)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	if len(s.symbols) > 0 {
		r.SetParam("symbols", s.symbols)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(SpotExchangeInfo)
	err = JSON.Unmarshal(data, res)
	return res, err
}

// SpotExchangeInfo represents exchange info
type SpotExchangeInfo struct {
	Timezone        string              `json:"timezone"`
	ServerTime      int64               `json:"serverTime"`
	RateLimits      []common.RateLimit  `json:"rateLimits"`
	ExchangeFilters []interface{}       `json:"exchangeFilters"`
	Symbols         []SpotSymbol        `json:"symbols"`
}

// SpotSymbol represents a symbol in spot
type SpotSymbol struct {
	Symbol                     string                   `json:"symbol"`
	Status                     string                   `json:"status"`
	BaseAsset                  string                   `json:"baseAsset"`
	BaseAssetPrecision         int                      `json:"baseAssetPrecision"`
	QuoteAsset                 string                   `json:"quoteAsset"`
	QuotePrecision             int                      `json:"quotePrecision"`
	QuoteAssetPrecision        int                      `json:"quoteAssetPrecision"`
	OrderTypes                 []common.OrderType       `json:"orderTypes"`
	IcebergAllowed             bool                     `json:"icebergAllowed"`
	OcoAllowed                 bool                     `json:"ocoAllowed"`
	IsSpotTradingAllowed       bool                     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool                     `json:"isMarginTradingAllowed"`
	Filters                    []interface{}            `json:"filters"`
	Permissions                []string                 `json:"permissions"`
}

// SpotDepthService get order book
type SpotDepthService struct {
	c      *BaseClient
	symbol string
	limit  *int
}

// Symbol set symbol
func (s *SpotDepthService) Symbol(symbol string) *SpotDepthService {
	s.symbol = symbol
	return s
}

// Limit set limit
func (s *SpotDepthService) Limit(limit int) *SpotDepthService {
	s.limit = &limit
	return s
}

// Do send request
func (s *SpotDepthService) Do(ctx context.Context, opts ...RequestOption) (res *SpotDepthResponse, err error) {
	r := newRequest(http.MethodGet, "/api/v3/depth", secTypeNone)
	r.SetParam("symbol", s.symbol)
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(SpotDepthResponse)
	err = JSON.Unmarshal(data, res)
	return res, err
}

// SpotDepthResponse represents depth info
type SpotDepthResponse struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// SpotRecentTradesListService list recent trades
type SpotRecentTradesListService struct {
	c      *BaseClient
	symbol string
	limit  *int
}

// Symbol set symbol
func (s *SpotRecentTradesListService) Symbol(symbol string) *SpotRecentTradesListService {
	s.symbol = symbol
	return s
}

// Limit set limit
func (s *SpotRecentTradesListService) Limit(limit int) *SpotRecentTradesListService {
	s.limit = &limit
	return s
}

// Do send request
func (s *SpotRecentTradesListService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotMarketTrade, err error) {
	r := newRequest(http.MethodGet, "/api/v3/trades", secTypeNone)
	r.SetParam("symbol", s.symbol)
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*SpotMarketTrade, 0)
	err = JSON.Unmarshal(data, &res)
	return res, err
}

// SpotMarketTrade represents a market trade
type SpotMarketTrade struct {
	ID           int64  `json:"id"`
	Price        string `json:"price"`
	Quantity     string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

// SpotAggTradesService list aggregate trades
type SpotAggTradesService struct {
	c         *BaseClient
	symbol    string
	fromID    *int64
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *SpotAggTradesService) Symbol(symbol string) *SpotAggTradesService {
	s.symbol = symbol
	return s
}

// FromID set fromID
func (s *SpotAggTradesService) FromID(fromID int64) *SpotAggTradesService {
	s.fromID = &fromID
	return s
}

// StartTime set startTime
func (s *SpotAggTradesService) StartTime(startTime int64) *SpotAggTradesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *SpotAggTradesService) EndTime(endTime int64) *SpotAggTradesService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *SpotAggTradesService) Limit(limit int) *SpotAggTradesService {
	s.limit = &limit
	return s
}

// Do send request
func (s *SpotAggTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotAggTrade, err error) {
	r := newRequest(http.MethodGet, "/api/v3/aggTrades", secTypeNone)
	r.SetParam("symbol", s.symbol)
	if s.fromID != nil {
		r.SetParam("fromId", *s.fromID)
	}
	if s.startTime != nil {
		r.SetParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*SpotAggTrade, 0)
	err = JSON.Unmarshal(data, &res)
	return res, err
}

// SpotAggTrade represents aggregated trade
type SpotAggTrade struct {
	TradeID      int64  `json:"a"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	FirstTradeID int64  `json:"f"`
	LastTradeID  int64  `json:"l"`
	Time         int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
	IsBestMatch  bool   `json:"M"`
}

// SpotKlinesService list klines
type SpotKlinesService struct {
	c         *BaseClient
	symbol    string
	interval  common.Interval
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *SpotKlinesService) Symbol(symbol string) *SpotKlinesService {
	s.symbol = symbol
	return s
}

// Interval set interval
func (s *SpotKlinesService) Interval(interval common.Interval) *SpotKlinesService {
	s.interval = interval
	return s
}

// StartTime set startTime
func (s *SpotKlinesService) StartTime(startTime int64) *SpotKlinesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *SpotKlinesService) EndTime(endTime int64) *SpotKlinesService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *SpotKlinesService) Limit(limit int) *SpotKlinesService {
	s.limit = &limit
	return s
}

// Do send request
func (s *SpotKlinesService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotKline, err error) {
	r := newRequest(http.MethodGet, "/api/v3/klines", secTypeNone)
	r.SetParam("symbol", s.symbol)
	r.SetParam("interval", s.interval)
	if s.startTime != nil {
		r.SetParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Klines are returned as arrays, need to unmarshal properly
	var klineData [][]interface{}
	err = JSON.Unmarshal(data, &klineData)
	if err != nil {
		return nil, err
	}
	
	res = make([]*SpotKline, 0, len(klineData))
	for _, k := range klineData {
		if len(k) < 12 {
			continue
		}
		res = append(res, &SpotKline{
			OpenTime:                 int64(k[0].(float64)),
			Open:                     k[1].(string),
			High:                     k[2].(string),
			Low:                      k[3].(string),
			Close:                    k[4].(string),
			Volume:                   k[5].(string),
			CloseTime:                int64(k[6].(float64)),
			QuoteAssetVolume:         k[7].(string),
			TradeNum:                 int64(k[8].(float64)),
			TakerBuyBaseAssetVolume:  k[9].(string),
			TakerBuyQuoteAssetVolume: k[10].(string),
		})
	}
	
	return res, nil
}

// SpotKline represents a kline
type SpotKline struct {
	OpenTime                 int64  `json:"openTime"`
	Open                     string `json:"open"`
	High                     string `json:"high"`
	Low                      string `json:"low"`
	Close                    string `json:"close"`
	Volume                   string `json:"volume"`
	CloseTime                int64  `json:"closeTime"`
	QuoteAssetVolume         string `json:"quoteAssetVolume"`
	TradeNum                 int64  `json:"tradeNum"`
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"`
}

// SpotListPricesService list latest price for all symbols
type SpotListPricesService struct {
	c      *BaseClient
	symbol *string
}

// Symbol set symbol
func (s *SpotListPricesService) Symbol(symbol string) *SpotListPricesService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *SpotListPricesService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotSymbolPrice, err error) {
	r := newRequest(http.MethodGet, "/api/v3/ticker/price", secTypeNone)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Handle both single and array response
	if s.symbol != nil {
		var price SpotSymbolPrice
		err = JSON.Unmarshal(data, &price)
		if err != nil {
			return nil, err
		}
		res = []*SpotSymbolPrice{&price}
	} else {
		res = make([]*SpotSymbolPrice, 0)
		err = JSON.Unmarshal(data, &res)
	}
	
	return res, err
}

// SpotSymbolPrice represents price of a symbol
type SpotSymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// SpotListBookTickersService list best price/qty on the order book
type SpotListBookTickersService struct {
	c      *BaseClient
	symbol *string
}

// Symbol set symbol
func (s *SpotListBookTickersService) Symbol(symbol string) *SpotListBookTickersService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *SpotListBookTickersService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotBookTicker, err error) {
	r := newRequest(http.MethodGet, "/api/v3/ticker/bookTicker", secTypeNone)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Handle both single and array response
	if s.symbol != nil {
		var ticker SpotBookTicker
		err = JSON.Unmarshal(data, &ticker)
		if err != nil {
			return nil, err
		}
		res = []*SpotBookTicker{&ticker}
	} else {
		res = make([]*SpotBookTicker, 0)
		err = JSON.Unmarshal(data, &res)
	}
	
	return res, err
}

// SpotBookTicker represents book ticker
type SpotBookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}

// SpotListPriceChangeStatsService show price change stats
type SpotListPriceChangeStatsService struct {
	c      *BaseClient
	symbol *string
}

// Symbol set symbol
func (s *SpotListPriceChangeStatsService) Symbol(symbol string) *SpotListPriceChangeStatsService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *SpotListPriceChangeStatsService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotPriceChangeStats, err error) {
	r := newRequest(http.MethodGet, "/api/v3/ticker/24hr", secTypeNone)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Handle both single and array response
	if s.symbol != nil {
		var stats SpotPriceChangeStats
		err = JSON.Unmarshal(data, &stats)
		if err != nil {
			return nil, err
		}
		res = []*SpotPriceChangeStats{&stats}
	} else {
		res = make([]*SpotPriceChangeStats, 0)
		err = JSON.Unmarshal(data, &res)
	}
	
	return res, err
}

// SpotPriceChangeStats represents price change stats
type SpotPriceChangeStats struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstID            int64  `json:"firstId"`
	LastID             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}