package futures

import (
	"context"
	"net/http"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
)

// PingService ping server
type PingService struct {
	C *aster.BaseClient
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...aster.RequestOption) (err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/ping", aster.SecTypeNone)
	_, err = s.C.CallAPI(ctx, r, opts...)
	return err
}

// ServerTimeService get server time
type ServerTimeService struct {
	C *aster.BaseClient
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...aster.RequestOption) (serverTime int64, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/time", aster.SecTypeNone)
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	j, err := aster.NewJSON(data)
	if err != nil {
		return 0, err
	}
	serverTime = (*j).Get("serverTime").ToInt64()
	return serverTime, nil
}

// ExchangeInfoService get exchange info
type ExchangeInfoService struct {
	C *aster.BaseClient
}

// Do send request
func (s *ExchangeInfoService) Do(ctx context.Context, opts ...aster.RequestOption) (res *ExchangeInfo, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/exchangeInfo", aster.SecTypeNone)
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(ExchangeInfo)
	err = aster.JSON.Unmarshal(data, res)
	return res, err
}

// ExchangeInfo exchange info
type ExchangeInfo struct {
	Timezone        string             `json:"timezone"`
	ServerTime      int64              `json:"serverTime"`
	FuturesType     string             `json:"futuresType"`
	RateLimits      []common.RateLimit `json:"rateLimits"`
	ExchangeFilters []interface{}      `json:"exchangeFilters"`
	Assets          []interface{}      `json:"assets"`
	Symbols         []Symbol           `json:"symbols"`
}

// DepthService get order book
type DepthService struct {
	C      *aster.BaseClient
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
func (s *DepthService) Do(ctx context.Context, opts ...aster.RequestOption) (res *DepthResponse, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/depth", aster.SecTypeNone)
	r.SetParam("symbol", s.symbol)
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(DepthResponse)
	err = aster.JSON.Unmarshal(data, res)
	return res, err
}

// DepthResponse depth response
type DepthResponse struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	EventTime    int64      `json:"E"`
	TransactionTime int64   `json:"T"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// RecentTradesListService list recent trades
type RecentTradesListService struct {
	C      *aster.BaseClient
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
func (s *RecentTradesListService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*Trade, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/trades", aster.SecTypeNone)
	r.SetParam("symbol", s.symbol)
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*Trade, 0)
	err = aster.JSON.Unmarshal(data, &res)
	return res, err
}

// Trade represents a trade
type Trade struct {
	ID           int64  `json:"id"`
	Price        string `json:"price"`
	Quantity     string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
}

// AggTradesService list aggregate trades
type AggTradesService struct {
	C         *aster.BaseClient
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
func (s *AggTradesService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*AggTrade, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/aggTrades", aster.SecTypeNone)
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
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*AggTrade, 0)
	err = aster.JSON.Unmarshal(data, &res)
	return res, err
}

// AggTrade represents aggregate trade
type AggTrade struct {
	AggTradeID   int64  `json:"a"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	FirstTradeID int64  `json:"f"`
	LastTradeID  int64  `json:"l"`
	Time         int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
}

// KlinesService list klines
type KlinesService struct {
	C         *aster.BaseClient
	symbol    string
	interval  common.Interval
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
func (s *KlinesService) Interval(interval common.Interval) *KlinesService {
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
func (s *KlinesService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*Kline, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/klines", aster.SecTypeNone)
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
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Klines are returned as arrays
	var klineData [][]interface{}
	err = aster.JSON.Unmarshal(data, &klineData)
	if err != nil {
		return nil, err
	}
	
	res = make([]*Kline, 0, len(klineData))
	for _, k := range klineData {
		if len(k) < 12 {
			continue
		}
		res = append(res, &Kline{
			OpenTime:                 int64(k[0].(float64)),
			Open:                     k[1].(string),
			High:                     k[2].(string),
			Low:                      k[3].(string),
			Close:                    k[4].(string),
			Volume:                   k[5].(string),
			CloseTime:                int64(k[6].(float64)),
			QuoteVolume:              k[7].(string),
			TradeNum:                 int64(k[8].(float64)),
			TakerBuyBaseAssetVolume:  k[9].(string),
			TakerBuyQuoteAssetVolume: k[10].(string),
		})
	}
	
	return res, nil
}

// Kline represents a kline
type Kline struct {
	OpenTime                 int64  `json:"openTime"`
	Open                     string `json:"open"`
	High                     string `json:"high"`
	Low                      string `json:"low"`
	Close                    string `json:"close"`
	Volume                   string `json:"volume"`
	CloseTime                int64  `json:"closeTime"`
	QuoteVolume              string `json:"quoteVolume"`
	TradeNum                 int64  `json:"tradeNum"`
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"`
}

// ContinuousKlinesService list continuous klines
type ContinuousKlinesService struct {
	C            *aster.BaseClient
	pair         string
	contractType ContractType
	interval     common.Interval
	startTime    *int64
	endTime      *int64
	limit        *int
}

// Pair set pair
func (s *ContinuousKlinesService) Pair(pair string) *ContinuousKlinesService {
	s.pair = pair
	return s
}

// ContractType set contractType
func (s *ContinuousKlinesService) ContractType(contractType ContractType) *ContinuousKlinesService {
	s.contractType = contractType
	return s
}

// Interval set interval
func (s *ContinuousKlinesService) Interval(interval common.Interval) *ContinuousKlinesService {
	s.interval = interval
	return s
}

// StartTime set startTime
func (s *ContinuousKlinesService) StartTime(startTime int64) *ContinuousKlinesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *ContinuousKlinesService) EndTime(endTime int64) *ContinuousKlinesService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *ContinuousKlinesService) Limit(limit int) *ContinuousKlinesService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ContinuousKlinesService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*Kline, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/continuousKlines", aster.SecTypeNone)
	r.SetParam("pair", s.pair)
	r.SetParam("contractType", s.contractType)
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
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Parse klines same way as regular klines
	var klineData [][]interface{}
	err = aster.JSON.Unmarshal(data, &klineData)
	if err != nil {
		return nil, err
	}
	
	res = make([]*Kline, 0, len(klineData))
	for _, k := range klineData {
		if len(k) < 12 {
			continue
		}
		res = append(res, &Kline{
			OpenTime:                 int64(k[0].(float64)),
			Open:                     k[1].(string),
			High:                     k[2].(string),
			Low:                      k[3].(string),
			Close:                    k[4].(string),
			Volume:                   k[5].(string),
			CloseTime:                int64(k[6].(float64)),
			QuoteVolume:              k[7].(string),
			TradeNum:                 int64(k[8].(float64)),
			TakerBuyBaseAssetVolume:  k[9].(string),
			TakerBuyQuoteAssetVolume: k[10].(string),
		})
	}
	
	return res, nil
}

// MarkPriceService get mark price
type MarkPriceService struct {
	C      *aster.BaseClient
	symbol *string
}

// Symbol set symbol
func (s *MarkPriceService) Symbol(symbol string) *MarkPriceService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *MarkPriceService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*MarkPrice, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/premiumIndex", aster.SecTypeNone)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Handle both single and array response
	if s.symbol != nil {
		var markPrice MarkPrice
		err = aster.JSON.Unmarshal(data, &markPrice)
		if err != nil {
			return nil, err
		}
		res = []*MarkPrice{&markPrice}
	} else {
		res = make([]*MarkPrice, 0)
		err = aster.JSON.Unmarshal(data, &res)
	}
	
	return res, err
}

// FundingRateService get funding rate history
type FundingRateService struct {
	C         *aster.BaseClient
	symbol    string
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *FundingRateService) Symbol(symbol string) *FundingRateService {
	s.symbol = symbol
	return s
}

// StartTime set startTime
func (s *FundingRateService) StartTime(startTime int64) *FundingRateService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *FundingRateService) EndTime(endTime int64) *FundingRateService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *FundingRateService) Limit(limit int) *FundingRateService {
	s.limit = &limit
	return s
}

// Do send request
func (s *FundingRateService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*FundingRate, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/fundingRate", aster.SecTypeNone)
	r.SetParam("symbol", s.symbol)
	if s.startTime != nil {
		r.SetParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*FundingRate, 0)
	err = aster.JSON.Unmarshal(data, &res)
	return res, err
}

// ListPriceChangeStatsService show price change stats
type ListPriceChangeStatsService struct {
	C      *aster.BaseClient
	symbol *string
}

// Symbol set symbol
func (s *ListPriceChangeStatsService) Symbol(symbol string) *ListPriceChangeStatsService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListPriceChangeStatsService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*PriceChangeStats, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/ticker/24hr", aster.SecTypeNone)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Handle both single and array response
	if s.symbol != nil {
		var stats PriceChangeStats
		err = aster.JSON.Unmarshal(data, &stats)
		if err != nil {
			return nil, err
		}
		res = []*PriceChangeStats{&stats}
	} else {
		res = make([]*PriceChangeStats, 0)
		err = aster.JSON.Unmarshal(data, &res)
	}
	
	return res, err
}

// ListPricesService list all prices
type ListPricesService struct {
	C      *aster.BaseClient
	symbol *string
}

// Symbol set symbol
func (s *ListPricesService) Symbol(symbol string) *ListPricesService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListPricesService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*SymbolPrice, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/ticker/price", aster.SecTypeNone)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Handle both single and array response
	if s.symbol != nil {
		var price SymbolPrice
		err = aster.JSON.Unmarshal(data, &price)
		if err != nil {
			return nil, err
		}
		res = []*SymbolPrice{&price}
	} else {
		res = make([]*SymbolPrice, 0)
		err = aster.JSON.Unmarshal(data, &res)
	}
	
	return res, err
}

// SymbolPrice represents price of a symbol
type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

// ListBookTickersService list best price/qty on the order book
type ListBookTickersService struct {
	C      *aster.BaseClient
	symbol *string
}

// Symbol set symbol
func (s *ListBookTickersService) Symbol(symbol string) *ListBookTickersService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListBookTickersService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*BookTicker, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/ticker/bookTicker", aster.SecTypeNone)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	
	// Handle both single and array response
	if s.symbol != nil {
		var ticker BookTicker
		err = aster.JSON.Unmarshal(data, &ticker)
		if err != nil {
			return nil, err
		}
		res = []*BookTicker{&ticker}
	} else {
		res = make([]*BookTicker, 0)
		err = aster.JSON.Unmarshal(data, &res)
	}
	
	return res, err
}

// BookTicker represents book ticker
type BookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
	Time     int64  `json:"time"`
}