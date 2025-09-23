package aster

import (
	"context"
	"net/http"
)

// GetSpotAccountService get account info
type GetSpotAccountService struct {
	c *BaseClient
}

// Do send request
func (s *GetSpotAccountService) Do(ctx context.Context, opts ...RequestOption) (res *SpotAccount, err error) {
	r := newRequest(http.MethodGet, "/api/v3/account", secTypeSigned)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(SpotAccount)
	err = JSON.Unmarshal(data, res)
	return res, err
}

// SpotAccount represents spot account info
type SpotAccount struct {
	MakerCommission  int64          `json:"makerCommission"`
	TakerCommission  int64          `json:"takerCommission"`
	BuyerCommission  int64          `json:"buyerCommission"`
	SellerCommission int64          `json:"sellerCommission"`
	CanTrade         bool           `json:"canTrade"`
	CanWithdraw      bool           `json:"canWithdraw"`
	CanDeposit       bool           `json:"canDeposit"`
	UpdateTime       int64          `json:"updateTime"`
	AccountType      string         `json:"accountType"`
	Balances         []SpotBalance  `json:"balances"`
	Permissions      []string       `json:"permissions"`
}

// SpotBalance represents a balance in spot account
type SpotBalance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// ListSpotTradesService list trades
type ListSpotTradesService struct {
	c         *BaseClient
	symbol    string
	orderId   *int64
	startTime *int64
	endTime   *int64
	fromId    *int64
	limit     *int
}

// Symbol set symbol
func (s *ListSpotTradesService) Symbol(symbol string) *ListSpotTradesService {
	s.symbol = symbol
	return s
}

// OrderId set orderId
func (s *ListSpotTradesService) OrderId(orderId int64) *ListSpotTradesService {
	s.orderId = &orderId
	return s
}

// StartTime set startTime
func (s *ListSpotTradesService) StartTime(startTime int64) *ListSpotTradesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *ListSpotTradesService) EndTime(endTime int64) *ListSpotTradesService {
	s.endTime = &endTime
	return s
}

// FromId set fromId
func (s *ListSpotTradesService) FromId(fromId int64) *ListSpotTradesService {
	s.fromId = &fromId
	return s
}

// Limit set limit
func (s *ListSpotTradesService) Limit(limit int) *ListSpotTradesService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListSpotTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotTrade, err error) {
	r := newRequest(http.MethodGet, "/api/v3/myTrades", secTypeSigned)
	r.SetParam("symbol", s.symbol)
	if s.orderId != nil {
		r.SetParam("orderId", *s.orderId)
	}
	if s.startTime != nil {
		r.SetParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetParam("endTime", *s.endTime)
	}
	if s.fromId != nil {
		r.SetParam("fromId", *s.fromId)
	}
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*SpotTrade, 0)
	err = JSON.Unmarshal(data, &res)
	return res, err
}

// SpotTrade represents a trade in spot
type SpotTrade struct {
	Symbol          string `json:"symbol"`
	Id              int64  `json:"id"`
	OrderId         int64  `json:"orderId"`
	OrderListId     int64  `json:"orderListId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}