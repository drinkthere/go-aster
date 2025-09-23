package aster

import (
	"context"
	"net/http"

	"github.com/drinkthere/go-aster/v2/common"
)

// CreateSpotOrderService create spot order
type CreateSpotOrderService struct {
	c                *BaseClient
	symbol           string
	side             common.SideType
	orderType        common.OrderType
	timeInForce      *common.TimeInForceType
	quantity         *string
	quoteOrderQty    *string
	price            *string
	newClientOrderID *string
	stopPrice        *string
	icebergQty       *string
	newOrderRespType *common.NewOrderRespType
}

// Symbol set symbol
func (s *CreateSpotOrderService) Symbol(symbol string) *CreateSpotOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CreateSpotOrderService) Side(side common.SideType) *CreateSpotOrderService {
	s.side = side
	return s
}

// Type set order type
func (s *CreateSpotOrderService) Type(orderType common.OrderType) *CreateSpotOrderService {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *CreateSpotOrderService) TimeInForce(timeInForce common.TimeInForceType) *CreateSpotOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CreateSpotOrderService) Quantity(quantity string) *CreateSpotOrderService {
	s.quantity = &quantity
	return s
}

// QuoteOrderQty set quoteOrderQty
func (s *CreateSpotOrderService) QuoteOrderQty(quoteOrderQty string) *CreateSpotOrderService {
	s.quoteOrderQty = &quoteOrderQty
	return s
}

// Price set price
func (s *CreateSpotOrderService) Price(price string) *CreateSpotOrderService {
	s.price = &price
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CreateSpotOrderService) NewClientOrderID(newClientOrderID string) *CreateSpotOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// StopPrice set stopPrice
func (s *CreateSpotOrderService) StopPrice(stopPrice string) *CreateSpotOrderService {
	s.stopPrice = &stopPrice
	return s
}

// IcebergQty set icebergQty
func (s *CreateSpotOrderService) IcebergQty(icebergQty string) *CreateSpotOrderService {
	s.icebergQty = &icebergQty
	return s
}

// NewOrderRespType set newOrderRespType
func (s *CreateSpotOrderService) NewOrderRespType(newOrderRespType common.NewOrderRespType) *CreateSpotOrderService {
	s.newOrderRespType = &newOrderRespType
	return s
}

// Do send request
func (s *CreateSpotOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateSpotOrderResponse, err error) {
	r := newRequest(http.MethodPost, "/api/v3/order", secTypeSigned)
	m := params{
		"symbol": s.symbol,
		"side":   s.side,
		"type":   s.orderType,
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.quantity != nil {
		m["quantity"] = *s.quantity
	}
	if s.quoteOrderQty != nil {
		m["quoteOrderQty"] = *s.quoteOrderQty
	}
	if s.price != nil {
		m["price"] = *s.price
	}
	if s.newClientOrderID != nil {
		m["newClientOrderId"] = *s.newClientOrderID
	}
	if s.stopPrice != nil {
		m["stopPrice"] = *s.stopPrice
	}
	if s.icebergQty != nil {
		m["icebergQty"] = *s.icebergQty
	}
	if s.newOrderRespType != nil {
		m["newOrderRespType"] = *s.newOrderRespType
	}
	r.setFormParams(m)
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateSpotOrderResponse)
	err = JSON.Unmarshal(data, res)
	return res, err
}

// CreateSpotOrderResponse represents a response to create spot order
type CreateSpotOrderResponse struct {
	Symbol                  string                  `json:"symbol"`
	OrderID                 int64                   `json:"orderId"`
	OrderListID             int64                   `json:"orderListId"`
	ClientOrderID           string                  `json:"clientOrderId"`
	TransactTime            int64                   `json:"transactTime"`
	Price                   string                  `json:"price"`
	OrigQty                 string                  `json:"origQty"`
	ExecutedQty             string                  `json:"executedQty"`
	CumulativeQuoteQty      string                  `json:"cumulativeQuoteQty"`
	Status                  common.OrderStatusType  `json:"status"`
	TimeInForce             common.TimeInForceType  `json:"timeInForce"`
	Type                    common.OrderType        `json:"type"`
	Side                    common.SideType         `json:"side"`
	Fills                   []SpotFill              `json:"fills"`
}

// SpotFill represents a fill in spot order
type SpotFill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	TradeId         int64  `json:"tradeId"`
}

// GetSpotOrderService get spot order
type GetSpotOrderService struct {
	c                 *BaseClient
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Symbol set symbol
func (s *GetSpotOrderService) Symbol(symbol string) *GetSpotOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *GetSpotOrderService) OrderID(orderID int64) *GetSpotOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *GetSpotOrderService) OrigClientOrderID(origClientOrderID string) *GetSpotOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *GetSpotOrderService) Do(ctx context.Context, opts ...RequestOption) (res *SpotOrder, err error) {
	r := newRequest(http.MethodGet, "/api/v3/order", secTypeSigned)
	r.SetParam("symbol", s.symbol)
	if s.orderID != nil {
		r.SetParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetParam("origClientOrderId", *s.origClientOrderID)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(SpotOrder)
	err = JSON.Unmarshal(data, res)
	return res, err
}

// CancelSpotOrderService cancel spot order
type CancelSpotOrderService struct {
	c                 *BaseClient
	symbol            string
	orderID           *int64
	origClientOrderID *string
	newClientOrderID  *string
}

// Symbol set symbol
func (s *CancelSpotOrderService) Symbol(symbol string) *CancelSpotOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *CancelSpotOrderService) OrderID(orderID int64) *CancelSpotOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelSpotOrderService) OrigClientOrderID(origClientOrderID string) *CancelSpotOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CancelSpotOrderService) NewClientOrderID(newClientOrderID string) *CancelSpotOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// Do send request
func (s *CancelSpotOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelSpotOrderResponse, err error) {
	r := newRequest(http.MethodDelete, "/api/v3/order", secTypeSigned)
	r.SetParam("symbol", s.symbol)
	if s.orderID != nil {
		r.SetParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetParam("origClientOrderId", *s.origClientOrderID)
	}
	if s.newClientOrderID != nil {
		r.SetParam("newClientOrderId", *s.newClientOrderID)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CancelSpotOrderResponse)
	err = JSON.Unmarshal(data, res)
	return res, err
}

// CancelSpotOrderResponse represents response of canceling spot order
type CancelSpotOrderResponse struct {
	Symbol                  string                  `json:"symbol"`
	OrigClientOrderID       string                  `json:"origClientOrderId"`
	OrderID                 int64                   `json:"orderId"`
	OrderListID             int64                   `json:"orderListId"`
	ClientOrderID           string                  `json:"clientOrderId"`
	Price                   string                  `json:"price"`
	OrigQty                 string                  `json:"origQty"`
	ExecutedQty             string                  `json:"executedQty"`
	CumulativeQuoteQty      string                  `json:"cumulativeQuoteQty"`
	Status                  common.OrderStatusType  `json:"status"`
	TimeInForce             common.TimeInForceType  `json:"timeInForce"`
	Type                    common.OrderType        `json:"type"`
	Side                    common.SideType         `json:"side"`
}

// CancelOpenSpotOrdersService cancel all open orders
type CancelOpenSpotOrdersService struct {
	c      *BaseClient
	symbol string
}

// Symbol set symbol
func (s *CancelOpenSpotOrdersService) Symbol(symbol string) *CancelOpenSpotOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CancelOpenSpotOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*CancelSpotOrderResponse, err error) {
	r := newRequest(http.MethodDelete, "/api/v3/openOrders", secTypeSigned)
	r.SetParam("symbol", s.symbol)
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*CancelSpotOrderResponse, 0)
	err = JSON.Unmarshal(data, &res)
	return res, err
}

// ListOpenSpotOrdersService list open orders
type ListOpenSpotOrdersService struct {
	c      *BaseClient
	symbol *string
}

// Symbol set symbol
func (s *ListOpenSpotOrdersService) Symbol(symbol string) *ListOpenSpotOrdersService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListOpenSpotOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotOrder, err error) {
	r := newRequest(http.MethodGet, "/api/v3/openOrders", secTypeSigned)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*SpotOrder, 0)
	err = JSON.Unmarshal(data, &res)
	return res, err
}

// ListSpotOrdersService list all orders
type ListSpotOrdersService struct {
	c         *BaseClient
	symbol    string
	orderID   *int64
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *ListSpotOrdersService) Symbol(symbol string) *ListSpotOrdersService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *ListSpotOrdersService) OrderID(orderID int64) *ListSpotOrdersService {
	s.orderID = &orderID
	return s
}

// StartTime set startTime
func (s *ListSpotOrdersService) StartTime(startTime int64) *ListSpotOrdersService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *ListSpotOrdersService) EndTime(endTime int64) *ListSpotOrdersService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *ListSpotOrdersService) Limit(limit int) *ListSpotOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListSpotOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*SpotOrder, err error) {
	r := newRequest(http.MethodGet, "/api/v3/allOrders", secTypeSigned)
	r.SetParam("symbol", s.symbol)
	if s.orderID != nil {
		r.SetParam("orderId", *s.orderID)
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
	res = make([]*SpotOrder, 0)
	err = JSON.Unmarshal(data, &res)
	return res, err
}

// SpotOrder represents spot order info
type SpotOrder struct {
	Symbol              string                  `json:"symbol"`
	OrderID             int64                   `json:"orderId"`
	OrderListID         int64                   `json:"orderListId"`
	ClientOrderID       string                  `json:"clientOrderId"`
	Price               string                  `json:"price"`
	OrigQty             string                  `json:"origQty"`
	ExecutedQty         string                  `json:"executedQty"`
	CumulativeQuoteQty  string                  `json:"cumulativeQuoteQty"`
	Status              common.OrderStatusType  `json:"status"`
	TimeInForce         common.TimeInForceType  `json:"timeInForce"`
	Type                common.OrderType        `json:"type"`
	Side                common.SideType         `json:"side"`
	StopPrice           string                  `json:"stopPrice"`
	IcebergQty          string                  `json:"icebergQty"`
	Time                int64                   `json:"time"`
	UpdateTime          int64                   `json:"updateTime"`
	IsWorking           bool                    `json:"isWorking"`
	OrigQuoteOrderQty   string                  `json:"origQuoteOrderQty"`
}