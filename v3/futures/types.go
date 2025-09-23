package futures

import "time"

// Symbol filter types
type SymbolFilterType string

const (
	SymbolFilterTypePriceFilter       SymbolFilterType = "PRICE_FILTER"
	SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"
	SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"
	SymbolFilterTypeMaxNumOrders     SymbolFilterType = "MAX_NUM_ORDERS"
	SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
	SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"
	SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"
)

// Order types
type OrderType string

const (
	OrderTypeLimit           OrderType = "LIMIT"
	OrderTypeMarket          OrderType = "MARKET"
	OrderTypeStop            OrderType = "STOP"
	OrderTypeStopMarket      OrderType = "STOP_MARKET"
	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket OrderType = "TAKE_PROFIT_MARKET"
	OrderTypeTrailingStopMarket OrderType = "TRAILING_STOP_MARKET"
)

// Order side
type SideType string

const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"
)

// Position side
type PositionSideType string

const (
	PositionSideTypeBoth  PositionSideType = "BOTH"
	PositionSideTypeLong  PositionSideType = "LONG"
	PositionSideTypeShort PositionSideType = "SHORT"
)

// Time in force
type TimeInForceType string

const (
	TimeInForceTypeGTC TimeInForceType = "GTC" // Good Till Cancel
	TimeInForceTypeIOC TimeInForceType = "IOC" // Immediate or Cancel
	TimeInForceTypeFOK TimeInForceType = "FOK" // Fill or Kill
	TimeInForceTypeGTX TimeInForceType = "GTX" // Good Till Crossing
)

// Order status
type OrderStatusType string

const (
	OrderStatusTypeNew             OrderStatusType = "NEW"
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
	OrderStatusTypeFilled          OrderStatusType = "FILLED"
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
	OrderStatusTypeRejected        OrderStatusType = "REJECTED"
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"
)

// Working type
type WorkingType string

const (
	WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"
	WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE"
)

// Response type
type NewOrderRespType string

const (
	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"
)

// Kline intervals
type KlineInterval string

const (
	KlineInterval1m  KlineInterval = "1m"
	KlineInterval3m  KlineInterval = "3m"
	KlineInterval5m  KlineInterval = "5m"
	KlineInterval15m KlineInterval = "15m"
	KlineInterval30m KlineInterval = "30m"
	KlineInterval1h  KlineInterval = "1h"
	KlineInterval2h  KlineInterval = "2h"
	KlineInterval4h  KlineInterval = "4h"
	KlineInterval6h  KlineInterval = "6h"
	KlineInterval8h  KlineInterval = "8h"
	KlineInterval12h KlineInterval = "12h"
	KlineInterval1d  KlineInterval = "1d"
	KlineInterval3d  KlineInterval = "3d"
	KlineInterval1w  KlineInterval = "1w"
	KlineInterval1M  KlineInterval = "1M"
)

// Rate limit types
type RateLimitType string

const (
	RateLimitTypeRequestWeight RateLimitType = "REQUEST_WEIGHT"
	RateLimitTypeOrders        RateLimitType = "ORDERS"
)

// Rate limit intervals
type RateLimitInterval string

const (
	RateLimitIntervalSecond RateLimitInterval = "SECOND"
	RateLimitIntervalMinute RateLimitInterval = "MINUTE"
	RateLimitIntervalDay    RateLimitInterval = "DAY"
)

// Common structures
type RateLimit struct {
	RateLimitType     RateLimitType     `json:"rateLimitType"`
	Interval          RateLimitInterval `json:"interval"`
	IntervalNum       int64             `json:"intervalNum"`
	Limit             int64             `json:"limit"`
}

type Symbol struct {
	Symbol                string               `json:"symbol"`
	Pair                  string               `json:"pair"`
	ContractType          string               `json:"contractType"`
	DeliveryDate          int64                `json:"deliveryDate"`
	OnboardDate           int64                `json:"onboardDate"`
	Status                string               `json:"status"`
	MaintMarginPercent    string               `json:"maintMarginPercent"`
	RequiredMarginPercent string               `json:"requiredMarginPercent"`
	BaseAsset             string               `json:"baseAsset"`
	QuoteAsset            string               `json:"quoteAsset"`
	MarginAsset           string               `json:"marginAsset"`
	PricePrecision        int                  `json:"pricePrecision"`
	QuantityPrecision     int                  `json:"quantityPrecision"`
	BaseAssetPrecision    int                  `json:"baseAssetPrecision"`
	QuotePrecision        int                  `json:"quotePrecision"`
	UnderlyingType        string               `json:"underlyingType"`
	UnderlyingSubType     []string             `json:"underlyingSubType"`
	SettlePlan            int64                `json:"settlePlan"`
	TriggerProtect        string               `json:"triggerProtect"`
	Filters               []map[string]interface{} `json:"filters"`
	OrderType             []OrderType          `json:"orderType"`
	TimeInForce           []TimeInForceType    `json:"timeInForce"`
}

type DepthEntry struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type DepthResponse struct {
	LastUpdateID int64        `json:"lastUpdateId"`
	EventTime    int64        `json:"E"`
	Time         int64        `json:"T"`
	Bids         []DepthEntry `json:"bids"`
	Asks         []DepthEntry `json:"asks"`
}

type Trade struct {
	ID           int64  `json:"id"`
	Price        string `json:"price"`
	Quantity     string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
}

type AggTrade struct {
	AggregateTradeID int64  `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeID     int64  `json:"f"`
	LastTradeID      int64  `json:"l"`
	Time             int64  `json:"T"`
	IsBuyerMaker     bool   `json:"m"`
}

type Kline struct {
	OpenTime                 int64  `json:"t"`
	Open                     string `json:"o"`
	High                     string `json:"h"`
	Low                      string `json:"l"`
	Close                    string `json:"c"`
	Volume                   string `json:"v"`
	CloseTime                int64  `json:"T"`
	QuoteVolume              string `json:"q"`
	TradeNum                 int64  `json:"n"`
	TakerBuyBaseAssetVolume  string `json:"V"`
	TakerBuyQuoteAssetVolume string `json:"Q"`
}

type TickerPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

type BookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
	Time     int64  `json:"time"`
}

type Order struct {
	Symbol           string           `json:"symbol"`
	OrderID          int64            `json:"orderId"`
	ClientOrderID    string           `json:"clientOrderId"`
	Price            string           `json:"price"`
	ReduceOnly       bool             `json:"reduceOnly"`
	OrigQty          string           `json:"origQty"`
	ExecutedQty      string           `json:"executedQty"`
	CumQty           string           `json:"cumQty"`
	CumQuote         string           `json:"cumQuote"`
	Status           OrderStatusType  `json:"status"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	Side             SideType         `json:"side"`
	StopPrice        string           `json:"stopPrice"`
	WorkingType      WorkingType      `json:"workingType"`
	OrigType         OrderType        `json:"origType"`
	PositionSide     PositionSideType `json:"positionSide"`
	ClosePosition    bool             `json:"closePosition"`
	ActivatePrice    string           `json:"activatePrice"`
	PriceRate        string           `json:"priceRate"`
	UpdateTime       int64            `json:"updateTime"`
	Time             int64            `json:"time"`
	AvgPrice         string           `json:"avgPrice"`
}

type Balance struct {
	Asset               string `json:"asset"`
	Balance             string `json:"balance"`
	CrossWalletBalance  string `json:"crossWalletBalance"`
	CrossUnPnl          string `json:"crossUnPnl"`
	AvailableBalance    string `json:"availableBalance"`
	MaxWithdrawAmount   string `json:"maxWithdrawAmount"`
}

type Position struct {
	Symbol                 string           `json:"symbol"`
	PositionSide           PositionSideType `json:"positionSide"`
	PositionAmt            string           `json:"positionAmt"`
	MarginType             string           `json:"marginType"`
	IsolatedWallet         string           `json:"isolatedWallet"`
	MarkPrice              string           `json:"markPrice"`
	UnrealizedProfit       string           `json:"unRealizedProfit"`
	MaintMargin            string           `json:"maintMargin"`
	InitialMargin          string           `json:"initialMargin"`
	PositionInitialMargin  string           `json:"positionInitialMargin"`
	OpenOrderInitialMargin string           `json:"openOrderInitialMargin"`
	Leverage               string           `json:"leverage"`
	Isolated               bool             `json:"isolated"`
	EntryPrice             string           `json:"entryPrice"`
	MaxNotional            string           `json:"maxNotional"`
	UpdateTime             int64            `json:"updateTime"`
}

type Account struct {
	Assets                      []Balance  `json:"assets"`
	Positions                   []Position `json:"positions"`
	CanTrade                    bool       `json:"canTrade"`
	CanWithdraw                 bool       `json:"canWithdraw"`
	CanDeposit                  bool       `json:"canDeposit"`
	UpdateTime                  int64      `json:"updateTime"`
	TotalInitialMargin          string     `json:"totalInitialMargin"`
	TotalMaintMargin            string     `json:"totalMaintMargin"`
	TotalWalletBalance          string     `json:"totalWalletBalance"`
	TotalUnrealizedProfit       string     `json:"totalUnrealizedProfit"`
	TotalMarginBalance          string     `json:"totalMarginBalance"`
	TotalPositionInitialMargin  string     `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin string     `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     string     `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             string     `json:"totalCrossUnPnl"`
	AvailableBalance            string     `json:"availableBalance"`
	MaxWithdrawAmount           string     `json:"maxWithdrawAmount"`
}