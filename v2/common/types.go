package common

// Common types shared between spot and futures

// APIError represents an error response from the API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// Error returns the error message
func (e APIError) Error() string {
	return e.Message
}

// IsAPIError checks if error is APIError
func IsAPIError(e error) bool {
	_, ok := e.(*APIError)
	return ok
}

// Symbol filter types
type SymbolFilterType string

const (
	SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"
	SymbolFilterTypePriceFilter      SymbolFilterType = "PRICE_FILTER"
	SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"
	SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"
	SymbolFilterTypeIcebergParts     SymbolFilterType = "ICEBERG_PARTS"
	SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"
	SymbolFilterTypeMaxNumOrders     SymbolFilterType = "MAX_NUM_ORDERS"
	SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
)

// Order types
type OrderType string

const (
	OrderTypeLimit           OrderType = "LIMIT"
	OrderTypeMarket          OrderType = "MARKET"
	OrderTypeStopLoss        OrderType = "STOP_LOSS"
	OrderTypeStopLossLimit   OrderType = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
	OrderTypeLimitMaker      OrderType = "LIMIT_MAKER"
)

// Order side
type SideType string

const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"
)

// Order status
type OrderStatusType string

const (
	OrderStatusTypeNew             OrderStatusType = "NEW"
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
	OrderStatusTypeFilled          OrderStatusType = "FILLED"
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
	OrderStatusTypePendingCancel   OrderStatusType = "PENDING_CANCEL"
	OrderStatusTypeRejected        OrderStatusType = "REJECTED"
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"
)

// Time in force
type TimeInForceType string

const (
	TimeInForceTypeGTC TimeInForceType = "GTC"
	TimeInForceTypeIOC TimeInForceType = "IOC"
	TimeInForceTypeFOK TimeInForceType = "FOK"
)

// NewOrderRespType response type
type NewOrderRespType string

const (
	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"
	NewOrderRespTypeFULL   NewOrderRespType = "FULL"
)

// Rate limit intervals
type RateLimitInterval string

const (
	RateLimitIntervalSecond RateLimitInterval = "SECOND"
	RateLimitIntervalMinute RateLimitInterval = "MINUTE"
	RateLimitIntervalDay    RateLimitInterval = "DAY"
)

// Rate limit types
type RateLimitType string

const (
	RateLimitTypeRequestWeight RateLimitType = "REQUEST_WEIGHT"
	RateLimitTypeOrders        RateLimitType = "ORDERS"
	RateLimitTypeRawRequests   RateLimitType = "RAW_REQUESTS"
)

// RateLimit struct
type RateLimit struct {
	RateLimitType     RateLimitType     `json:"rateLimitType"`
	Interval          RateLimitInterval `json:"interval"`
	IntervalNum       int               `json:"intervalNum"`
	Limit             int               `json:"limit"`
}

// KlineInterval represents kline intervals
type Interval string

const (
	Interval1m  Interval = "1m"
	Interval3m  Interval = "3m"
	Interval5m  Interval = "5m"
	Interval15m Interval = "15m"
	Interval30m Interval = "30m"
	Interval1h  Interval = "1h"
	Interval2h  Interval = "2h"
	Interval4h  Interval = "4h"
	Interval6h  Interval = "6h"
	Interval8h  Interval = "8h"
	Interval12h Interval = "12h"
	Interval1d  Interval = "1d"
	Interval3d  Interval = "3d"
	Interval1w  Interval = "1w"
	Interval1M  Interval = "1M"
)