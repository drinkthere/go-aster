package futures

import (
	"context"
	"net/http"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
)

// CreateOrderService create futures order
type CreateOrderService struct {
	C                *aster.BaseClient
	symbol           string
	side             common.SideType
	positionSide     *PositionSideType
	orderType        common.OrderType
	timeInForce      *common.TimeInForceType
	quantity         string
	reduceOnly       *bool
	price            *string
	newClientOrderID *string
	stopPrice        *string
	closePosition    *bool
	activationPrice  *string
	callbackRate     *string
	workingType      *WorkingType
	priceProtect     *bool
	newOrderRespType *common.NewOrderRespType
}

// Symbol set symbol
func (s *CreateOrderService) Symbol(symbol string) *CreateOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CreateOrderService) Side(side common.SideType) *CreateOrderService {
	s.side = side
	return s
}

// PositionSide set positionSide
func (s *CreateOrderService) PositionSide(positionSide PositionSideType) *CreateOrderService {
	s.positionSide = &positionSide
	return s
}

// Type set type
func (s *CreateOrderService) Type(orderType common.OrderType) *CreateOrderService {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *CreateOrderService) TimeInForce(timeInForce common.TimeInForceType) *CreateOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CreateOrderService) Quantity(quantity string) *CreateOrderService {
	s.quantity = quantity
	return s
}

// ReduceOnly set reduceOnly
func (s *CreateOrderService) ReduceOnly(reduceOnly bool) *CreateOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// Price set price
func (s *CreateOrderService) Price(price string) *CreateOrderService {
	s.price = &price
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CreateOrderService) NewClientOrderID(newClientOrderID string) *CreateOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// StopPrice set stopPrice
func (s *CreateOrderService) StopPrice(stopPrice string) *CreateOrderService {
	s.stopPrice = &stopPrice
	return s
}

// ClosePosition set closePosition
func (s *CreateOrderService) ClosePosition(closePosition bool) *CreateOrderService {
	s.closePosition = &closePosition
	return s
}

// ActivationPrice set activationPrice
func (s *CreateOrderService) ActivationPrice(activationPrice string) *CreateOrderService {
	s.activationPrice = &activationPrice
	return s
}

// CallbackRate set callbackRate
func (s *CreateOrderService) CallbackRate(callbackRate string) *CreateOrderService {
	s.callbackRate = &callbackRate
	return s
}

// WorkingType set workingType
func (s *CreateOrderService) WorkingType(workingType WorkingType) *CreateOrderService {
	s.workingType = &workingType
	return s
}

// PriceProtect set priceProtect
func (s *CreateOrderService) PriceProtect(priceProtect bool) *CreateOrderService {
	s.priceProtect = &priceProtect
	return s
}

// NewOrderRespType set newOrderRespType
func (s *CreateOrderService) NewOrderRespType(newOrderRespType common.NewOrderRespType) *CreateOrderService {
	s.newOrderRespType = &newOrderRespType
	return s
}

// Do send request
func (s *CreateOrderService) Do(ctx context.Context, opts ...aster.RequestOption) (res *Order, err error) {
	r := aster.NewRequest(http.MethodPost, "/fapi/v1/order", aster.SecTypeSigned)
	m := aster.Params{
		"symbol":    s.symbol,
		"side":      s.side,
		"type":      s.orderType,
		"quantity":  s.quantity,
	}
	if s.positionSide != nil {
		m["positionSide"] = *s.positionSide
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.reduceOnly != nil {
		m["reduceOnly"] = *s.reduceOnly
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
	if s.closePosition != nil {
		m["closePosition"] = *s.closePosition
	}
	if s.activationPrice != nil {
		m["activationPrice"] = *s.activationPrice
	}
	if s.callbackRate != nil {
		m["callbackRate"] = *s.callbackRate
	}
	if s.workingType != nil {
		m["workingType"] = *s.workingType
	}
	if s.priceProtect != nil {
		m["priceProtect"] = *s.priceProtect
	}
	if s.newOrderRespType != nil {
		m["newOrderRespType"] = *s.newOrderRespType
	}
	r.SetFormParams(m)
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Order)
	err = aster.JSON.Unmarshal(data, res)
	return res, err
}

// GetOrderService get order
type GetOrderService struct {
	C                 *aster.BaseClient
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Symbol set symbol
func (s *GetOrderService) Symbol(symbol string) *GetOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *GetOrderService) OrderID(orderID int64) *GetOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *GetOrderService) OrigClientOrderID(origClientOrderID string) *GetOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *GetOrderService) Do(ctx context.Context, opts ...aster.RequestOption) (res *Order, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/order", aster.SecTypeSigned)
	r.SetParam("symbol", s.symbol)
	if s.orderID != nil {
		r.SetParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetParam("origClientOrderId", *s.origClientOrderID)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Order)
	err = aster.JSON.Unmarshal(data, res)
	return res, err
}

// CancelOrderService cancel order
type CancelOrderService struct {
	C                 *aster.BaseClient
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Symbol set symbol
func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *CancelOrderService) OrderID(orderID int64) *CancelOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelOrderService) OrigClientOrderID(origClientOrderID string) *CancelOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *CancelOrderService) Do(ctx context.Context, opts ...aster.RequestOption) (res *Order, err error) {
	r := aster.NewRequest(http.MethodDelete, "/fapi/v1/order", aster.SecTypeSigned)
	r.SetParam("symbol", s.symbol)
	if s.orderID != nil {
		r.SetParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetParam("origClientOrderId", *s.origClientOrderID)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Order)
	err = aster.JSON.Unmarshal(data, res)
	return res, err
}

// CancelAllOpenOrdersService cancel all open orders
type CancelAllOpenOrdersService struct {
	C      *aster.BaseClient
	symbol string
}

// Symbol set symbol
func (s *CancelAllOpenOrdersService) Symbol(symbol string) *CancelAllOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CancelAllOpenOrdersService) Do(ctx context.Context, opts ...aster.RequestOption) error {
	r := aster.NewRequest(http.MethodDelete, "/fapi/v1/allOpenOrders", aster.SecTypeSigned)
	r.SetParam("symbol", s.symbol)
	_, err := s.C.CallAPI(ctx, r, opts...)
	return err
}

// ListOpenOrdersService list open orders
type ListOpenOrdersService struct {
	C      *aster.BaseClient
	symbol *string
}

// Symbol set symbol
func (s *ListOpenOrdersService) Symbol(symbol string) *ListOpenOrdersService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListOpenOrdersService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*Order, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/openOrders", aster.SecTypeSigned)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*Order, 0)
	err = aster.JSON.Unmarshal(data, &res)
	return res, err
}

// ListOrdersService list all orders
type ListOrdersService struct {
	C         *aster.BaseClient
	symbol    string
	orderID   *int64
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *ListOrdersService) Symbol(symbol string) *ListOrdersService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *ListOrdersService) OrderID(orderID int64) *ListOrdersService {
	s.orderID = &orderID
	return s
}

// StartTime set startTime
func (s *ListOrdersService) StartTime(startTime int64) *ListOrdersService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *ListOrdersService) EndTime(endTime int64) *ListOrdersService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *ListOrdersService) Limit(limit int) *ListOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListOrdersService) Do(ctx context.Context, opts ...aster.RequestOption) (res []*Order, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/allOrders", aster.SecTypeSigned)
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
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*Order, 0)
	err = aster.JSON.Unmarshal(data, &res)
	return res, err
}