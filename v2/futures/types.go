package futures

import (
	"github.com/drinkthere/go-aster/v2/common"
)

// PositionSideType position side type
type PositionSideType string

const (
	PositionSideTypeBoth  PositionSideType = "BOTH"
	PositionSideTypeLong  PositionSideType = "LONG"
	PositionSideTypeShort PositionSideType = "SHORT"
)

// MarginType margin type
type MarginType string

const (
	MarginTypeCross    MarginType = "CROSS"
	MarginTypeIsolated MarginType = "ISOLATED"
)

// ContractType contract type
type ContractType string

const (
	ContractTypePerpetual     ContractType = "PERPETUAL"
	ContractTypeCurrentMonth  ContractType = "CURRENT_MONTH"
	ContractTypeNextMonth     ContractType = "NEXT_MONTH"
	ContractTypeCurrentQuarter ContractType = "CURRENT_QUARTER"
	ContractTypeNextQuarter   ContractType = "NEXT_QUARTER"
)

// OrderExecutionType order execution type
type OrderExecutionType string

const (
	OrderExecutionTypeNew              OrderExecutionType = "NEW"
	OrderExecutionTypePartialFill      OrderExecutionType = "PARTIAL_FILL"
	OrderExecutionTypeFill             OrderExecutionType = "FILL"
	OrderExecutionTypeCanceled         OrderExecutionType = "CANCELED"
	OrderExecutionTypeCalculated       OrderExecutionType = "CALCULATED"
	OrderExecutionTypeExpired          OrderExecutionType = "EXPIRED"
	OrderExecutionTypeTrade            OrderExecutionType = "TRADE"
)

// WorkingType working type
type WorkingType string

const (
	WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"
	WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE"
)

// SymbolType symbol type
type SymbolType string

const (
	SymbolTypeFuture SymbolType = "FUTURE"
)

// SymbolStatusType symbol status type
type SymbolStatusType string

const (
	SymbolStatusTypePreTrading   SymbolStatusType = "PRE_TRADING"
	SymbolStatusTypeTrading      SymbolStatusType = "TRADING"
	SymbolStatusTypePostTrading  SymbolStatusType = "POST_TRADING"
	SymbolStatusTypeEndOfDay     SymbolStatusType = "END_OF_DAY"
	SymbolStatusTypeHalt         SymbolStatusType = "HALT"
	SymbolStatusTypeAuctionMatch SymbolStatusType = "AUCTION_MATCH"
	SymbolStatusTypeBreak        SymbolStatusType = "BREAK"
)

// Symbol market symbol
type Symbol struct {
	Symbol                string                     `json:"symbol"`
	Pair                  string                     `json:"pair"`
	ContractType          ContractType               `json:"contractType"`
	DeliveryDate          int64                      `json:"deliveryDate"`
	OnboardDate           int64                      `json:"onboardDate"`
	Status                SymbolStatusType           `json:"status"`
	MaintMarginPercent    string                     `json:"maintMarginPercent"`
	RequiredMarginPercent string                     `json:"requiredMarginPercent"`
	BaseAsset             string                     `json:"baseAsset"`
	QuoteAsset            string                     `json:"quoteAsset"`
	MarginAsset           string                     `json:"marginAsset"`
	PricePrecision        int                        `json:"pricePrecision"`
	QuantityPrecision     int                        `json:"quantityPrecision"`
	BaseAssetPrecision    int                        `json:"baseAssetPrecision"`
	QuotePrecision        int                        `json:"quotePrecision"`
	UnderlyingType        string                     `json:"underlyingType"`
	UnderlyingSubType     []string                   `json:"underlyingSubType"`
	SettlePlan            int64                      `json:"settlePlan"`
	TriggerProtect        string                     `json:"triggerProtect"`
	Filters               []interface{}              `json:"filters"`
	OrderTypes            []common.OrderType         `json:"orderTypes"`
	TimeInForce           []common.TimeInForceType   `json:"timeInForce"`
	LiquidationFee        string                     `json:"liquidationFee"`
	MarketTakeBound       string                     `json:"marketTakeBound"`
}

// Balance account balance
type Balance struct {
	AccountAlias       string `json:"accountAlias"`
	Asset              string `json:"asset"`
	Balance            string `json:"balance"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPnl         string `json:"crossUnPnl"`
	AvailableBalance   string `json:"availableBalance"`
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
	MarginAvailable    bool   `json:"marginAvailable"`
	UpdateTime         int64  `json:"updateTime"`
}

// PositionRisk position risk
type PositionRisk struct {
	Symbol                 string           `json:"symbol"`
	PositionSide           PositionSideType `json:"positionSide"`
	PositionAmt            string           `json:"positionAmt"`
	EntryPrice             string           `json:"entryPrice"`
	MarkPrice              string           `json:"markPrice"`
	UnRealizedProfit       string           `json:"unRealizedProfit"`
	LiquidationPrice       string           `json:"liquidationPrice"`
	Leverage               string           `json:"leverage"`
	MaxNotionalValue       string           `json:"maxNotionalValue"`
	MarginType             MarginType       `json:"marginType"`
	IsolatedMargin         string           `json:"isolatedMargin"`
	IsAutoAddMargin        string           `json:"isAutoAddMargin"`
	PositionInitialMargin  string           `json:"positionInitialMargin"`
	PositionMaintMargin    string           `json:"positionMaintMargin"`
	OpenOrderInitialMargin string           `json:"openOrderInitialMargin"`
	IsolatedWallet         string           `json:"isolatedWallet"`
	UpdateTime             int64            `json:"updateTime"`
}

// Order futures order info
type Order struct {
	AvgPrice         string                   `json:"avgPrice"`
	ClientOrderID    string                   `json:"clientOrderId"`
	CumQuote         string                   `json:"cumQuote"`
	ExecutedQty      string                   `json:"executedQty"`
	OrderID          int64                    `json:"orderId"`
	OrigQty          string                   `json:"origQty"`
	OrigType         common.OrderType         `json:"origType"`
	Price            string                   `json:"price"`
	ReduceOnly       bool                     `json:"reduceOnly"`
	Side             common.SideType          `json:"side"`
	PositionSide     PositionSideType         `json:"positionSide"`
	Status           common.OrderStatusType   `json:"status"`
	StopPrice        string                   `json:"stopPrice"`
	ClosePosition    bool                     `json:"closePosition"`
	Symbol           string                   `json:"symbol"`
	Time             int64                    `json:"time"`
	TimeInForce      common.TimeInForceType   `json:"timeInForce"`
	Type             common.OrderType         `json:"type"`
	ActivatePrice    string                   `json:"activatePrice"`
	PriceRate        string                   `json:"priceRate"`
	UpdateTime       int64                    `json:"updateTime"`
	WorkingType      WorkingType              `json:"workingType"`
	PriceProtect     bool                     `json:"priceProtect"`
}

// Account futures account info
type Account struct {
	Assets                      []Balance      `json:"assets"`
	Positions                   []PositionRisk `json:"positions"`
	CanTrade                    bool           `json:"canTrade"`
	CanWithdraw                 bool           `json:"canWithdraw"`
	CanDeposit                  bool           `json:"canDeposit"`
	UpdateTime                  int64          `json:"updateTime"`
	TotalInitialMargin          string         `json:"totalInitialMargin"`
	TotalMaintMargin            string         `json:"totalMaintMargin"`
	TotalWalletBalance          string         `json:"totalWalletBalance"`
	TotalUnrealizedProfit       string         `json:"totalUnrealizedProfit"`
	TotalMarginBalance          string         `json:"totalMarginBalance"`
	TotalPositionInitialMargin  string         `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin string         `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     string         `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             string         `json:"totalCrossUnPnl"`
	AvailableBalance            string         `json:"availableBalance"`
	MaxWithdrawAmount           string         `json:"maxWithdrawAmount"`
}

// MarkPrice mark price and funding rate
type MarkPrice struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	InterestRate         string `json:"interestRate"`
	Time                 int64  `json:"time"`
}

// FundingRate funding rate
type FundingRate struct {
	Symbol      string `json:"symbol"`
	FundingRate string `json:"fundingRate"`
	FundingTime int64  `json:"fundingTime"`
}

// PriceChangeStats price change stats
type PriceChangeStats struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
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