package aster

import (
	"github.com/drinkthere/go-aster/v2/common"
)

// Base WebSocket endpoints
const (
	baseWsMainnetURL        = "wss://sstream.asterdex.com"
	baseWsFuturesMainnetURL = "wss://fstream.asterdex.com"
	baseWsTestnetURL        = "wss://testnet.asterdex.com"
	
	combinedBaseURL        = "wss://sstream.asterdex.com/stream"
	combinedFuturesBaseURL = "wss://fstream.asterdex.com/stream"
)

// getWsEndpoint returns the websocket endpoint
func getWsEndpoint(isFutures, isTestnet bool) string {
	if isTestnet {
		return baseWsTestnetURL
	}
	if isFutures {
		return baseWsFuturesMainnetURL
	}
	return baseWsMainnetURL
}

// Depth handlers
type WsDepthHandler func(event *WsDepthEvent)

// Bid represents a bid in the order book
type Bid struct {
	Price    string
	Quantity string
}

// Ask represents an ask in the order book
type Ask struct {
	Price    string
	Quantity string
}

// WsDepthEvent define websocket depth event
type WsDepthEvent struct {
	Event         string `json:"e"`
	Time          int64  `json:"E"`
	TransactionTime int64  `json:"T"`
	Symbol        string `json:"s"`
	FirstUpdateID int64  `json:"U"`
	LastUpdateID  int64  `json:"u"`
	Bids          []Bid  `json:"b"`
	Asks          []Ask  `json:"a"`
}

// UnmarshalJSON custom unmarshal for Bid
func (b *Bid) UnmarshalJSON(data []byte) error {
	var items []string
	if err := JSON.Unmarshal(data, &items); err != nil {
		return err
	}
	if len(items) >= 2 {
		b.Price = items[0]
		b.Quantity = items[1]
	}
	return nil
}

// UnmarshalJSON custom unmarshal for Ask
func (a *Ask) UnmarshalJSON(data []byte) error {
	var items []string
	if err := JSON.Unmarshal(data, &items); err != nil {
		return err
	}
	if len(items) >= 2 {
		a.Price = items[0]
		a.Quantity = items[1]
	}
	return nil
}

// Kline handlers
type WsSpotKlineHandler func(event *WsSpotKlineEvent)
type WsFuturesKlineHandler func(event *WsFuturesKlineEvent)

// WsSpotKlineEvent define websocket kline event for spot
type WsSpotKlineEvent struct {
	Event  string      `json:"e"`
	Time   int64       `json:"E"`
	Symbol string      `json:"s"`
	Kline  WsSpotKline `json:"k"`
}

// WsSpotKline define websocket kline
type WsSpotKline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Symbol               string `json:"s"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
	Ignore               string `json:"B"`
}

// WsFuturesKlineEvent define websocket kline event for futures
type WsFuturesKlineEvent struct {
	Event  string          `json:"e"`
	Time   int64           `json:"E"`
	Symbol string          `json:"s"`
	Kline  WsFuturesKline  `json:"k"`
}

// WsFuturesKline define websocket kline for futures
type WsFuturesKline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Symbol               string `json:"s"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
}

// Aggregate trade handlers
type WsSpotAggTradeHandler func(event *WsSpotAggTradeEvent)
type WsFuturesAggTradeHandler func(event *WsFuturesAggTradeEvent)

// WsSpotAggTradeEvent define websocket aggregate trade event for spot
type WsSpotAggTradeEvent struct {
	Event                 string `json:"e"`
	Time                  int64  `json:"E"`
	Symbol                string `json:"s"`
	AggregateTradeID      int64  `json:"a"`
	Price                 string `json:"p"`
	Quantity              string `json:"q"`
	FirstBreakdownTradeID int64  `json:"f"`
	LastBreakdownTradeID  int64  `json:"l"`
	TradeTime             int64  `json:"T"`
	IsBuyerMaker          bool   `json:"m"`
	Ignore                bool   `json:"M"`
}

// WsFuturesAggTradeEvent define websocket aggregate trade event for futures
type WsFuturesAggTradeEvent struct {
	Event                 string `json:"e"`
	Time                  int64  `json:"E"`
	Symbol                string `json:"s"`
	AggregateTradeID      int64  `json:"a"`
	Price                 string `json:"p"`
	Quantity              string `json:"q"`
	FirstBreakdownTradeID int64  `json:"f"`
	LastBreakdownTradeID  int64  `json:"l"`
	TradeTime             int64  `json:"T"`
	IsBuyerMaker          bool   `json:"m"`
}

// Book ticker handler
type WsBookTickerHandler func(event *WsBookTickerEvent)

// WsBookTickerEvent define websocket best price/qty event
type WsBookTickerEvent struct {
	Event           string `json:"e"`
	UpdateID        int64  `json:"u"`
	Time            int64  `json:"E"`
	TransactionTime int64  `json:"T"`
	Symbol          string `json:"s"`
	BestBidPrice    string `json:"b"`
	BestBidQty      string `json:"B"`
	BestAskPrice    string `json:"a"`
	BestAskQty      string `json:"A"`
}

// Market statistics handler for spot
type WsSpotAllMarketsStatHandler func(event WsSpotAllMarketsStatEvent)

// WsSpotAllMarketsStatEvent define websocket market statistics event for spot
type WsSpotAllMarketsStatEvent []*WsSpotMarketStatEvent

// WsSpotMarketStatEvent define websocket market statistics event
type WsSpotMarketStatEvent struct {
	Event              string `json:"e"`
	Time               int64  `json:"E"`
	Symbol             string `json:"s"`
	PriceChange        string `json:"p"`
	PriceChangePercent string `json:"P"`
	WeightedAvgPrice   string `json:"w"`
	PrevClosePrice     string `json:"x"`
	LastPrice          string `json:"c"`
	CloseQty           string `json:"Q"`
	BidPrice           string `json:"b"`
	BidQty             string `json:"B"`
	AskPrice           string `json:"a"`
	AskQty             string `json:"A"`
	OpenPrice          string `json:"o"`
	HighPrice          string `json:"h"`
	LowPrice           string `json:"l"`
	BaseVolume         string `json:"v"`
	QuoteVolume        string `json:"q"`
	OpenTime           int64  `json:"O"`
	CloseTime          int64  `json:"C"`
	FirstID            int64  `json:"F"`
	LastID             int64  `json:"L"`
	Count              int64  `json:"n"`
}

// Mark price handler for futures
type WsFuturesMarkPriceHandler func(event *WsFuturesMarkPriceEvent)

// WsFuturesMarkPriceEvent define websocket mark price event for futures
type WsFuturesMarkPriceEvent struct {
	Event                string `json:"e"`
	Time                 int64  `json:"E"`
	Symbol               string `json:"s"`
	MarkPrice            string `json:"p"`
	IndexPrice           string `json:"i"`
	EstimatedSettlePrice string `json:"P"`
	FundingRate          string `json:"r"`
	NextFundingTime      int64  `json:"T"`
}

// User data handlers
type WsSpotUserDataHandler func(event *WsSpotUserDataEvent)
type WsFuturesUserDataHandler func(event *WsFuturesUserDataEvent)

// WsSpotUserDataEvent represents a spot user data event
type WsSpotUserDataEvent struct {
	Event            string                 `json:"e"`
	Time             int64                  `json:"E"`
	
	// Account update event
	AccountUpdate    *WsSpotAccountUpdate   `json:"accountUpdate,omitempty"`
	
	// Order update event
	OrderUpdate      *WsSpotOrderUpdate     `json:"orderUpdate,omitempty"`
}

// WsFuturesUserDataEvent represents a futures user data event
type WsFuturesUserDataEvent struct {
	Event            string                   `json:"e"`
	Time             int64                    `json:"E"`
	TransactionTime  int64                    `json:"T"`
	
	// Account update event
	AccountUpdate    *WsFuturesAccountUpdate  `json:"a,omitempty"`
	
	// Order update event
	OrderUpdate      *WsFuturesOrderUpdate    `json:"o,omitempty"`
	
	// Account config update event
	AccountConfigUpdate *WsFuturesAccountConfigUpdate `json:"ac,omitempty"`
	
	// Margin call event
	MarginCall       *WsFuturesMarginCall     `json:"margin_call,omitempty"`
}

// WsSpotAccountUpdate represents spot account update
type WsSpotAccountUpdate struct {
	Balances []WsSpotBalance `json:"B"`
}

// WsSpotBalance represents spot balance
type WsSpotBalance struct {
	Asset  string `json:"a"`
	Free   string `json:"f"`
	Locked string `json:"l"`
}

// WsSpotOrderUpdate represents spot order update
type WsSpotOrderUpdate struct {
	Symbol                   string                    `json:"s"`
	ClientOrderID            string                    `json:"c"`
	Side                     common.SideType           `json:"S"`
	OrderType                common.OrderType          `json:"o"`
	TimeInForce              common.TimeInForceType    `json:"f"`
	OriginalQuantity         string                    `json:"q"`
	OriginalPrice            string                    `json:"p"`
	AveragePrice             string                    `json:"ap"`
	StopPrice                string                    `json:"sp"`
	ExecutionType            string                    `json:"x"`
	OrderStatus              common.OrderStatusType    `json:"X"`
	OrderID                  int64                     `json:"i"`
	LastFilledQuantity       string                    `json:"l"`
	CumulativeFilledQuantity string                    `json:"z"`
	LastFilledPrice          string                    `json:"L"`
	Commission               string                    `json:"n"`
	CommissionAsset          string                    `json:"N"`
	TransactionTime          int64                     `json:"T"`
	TradeID                  int64                     `json:"t"`
	OrderCreatedTime         int64                     `json:"O"`
	CumulativeQuoteQty       string                    `json:"Z"`
	LastQuoteQty             string                    `json:"Y"`
}

// WsFuturesAccountUpdate represents futures account update
type WsFuturesAccountUpdate struct {
	Event           string                    `json:"e"`
	Time            int64                     `json:"E"`
	TransactionTime int64                     `json:"T"`
	UpdateData      WsFuturesAccountUpdateData `json:"a"`
}

// WsFuturesAccountUpdateData represents futures account update data
type WsFuturesAccountUpdateData struct {
	Reason    string                `json:"m"`
	Balances  []WsFuturesBalance    `json:"B"`
	Positions []WsFuturesPosition   `json:"P"`
}

// WsFuturesBalance represents futures balance
type WsFuturesBalance struct {
	Asset              string `json:"a"`
	WalletBalance      string `json:"wb"`
	CrossWalletBalance string `json:"cw"`
	BalanceChange      string `json:"bc"`
}

// WsFuturesPosition represents futures position
type WsFuturesPosition struct {
	Symbol         string `json:"s"`   // 交易对
	Side           string `json:"ps"`  // 持仓方向
	Amount         string `json:"pa"`  // 持仓数量
	MarginType     string `json:"mt"`  // 保证金模式
	IsolatedWallet string `json:"iw"`  // 若为逐仓，仓位保证金
	MarkPrice      string `json:"mp"`  // 标记价格
	UnrealizedPnL  string `json:"up"`  // 持仓未实现盈亏
	EntryPrice     string `json:"ep"`  // 持仓成本价
}

// WsFuturesOrderUpdate represents futures order update
type WsFuturesOrderUpdate struct {
	Symbol                   string                    `json:"s"`
	ClientOrderID            string                    `json:"c"`
	Side                     common.SideType           `json:"S"`
	OrderType                common.OrderType          `json:"o"`
	TimeInForce              common.TimeInForceType    `json:"f"`
	OriginalQuantity         string                    `json:"q"`
	OriginalPrice            string                    `json:"p"`
	AveragePrice             string                    `json:"ap"`
	StopPrice                string                    `json:"sp"`
	ExecutionType            string                    `json:"x"`
	OrderStatus              common.OrderStatusType    `json:"X"`
	OrderID                  int64                     `json:"i"`
	LastFilledQuantity       string                    `json:"l"`
	CumulativeFilledQuantity string                    `json:"z"`
	LastFilledPrice          string                    `json:"L"`
	CommissionAsset          string                    `json:"N"`
	Commission               string                    `json:"n"`
	OrderTradeTime           int64                     `json:"T"`
	TradeID                  int64                     `json:"t"`
	BidsNotional             string                    `json:"b"`
	AsksNotional             string                    `json:"a"`
	IsMaker                  bool                      `json:"m"`
	IsReduceOnly             bool                      `json:"R"`
	WorkingType              string                    `json:"wt"`
	OriginalType             common.OrderType          `json:"ot"`
	PositionSide             string                    `json:"ps"`
	IsClosingPosition        bool                      `json:"cp"`
	ActivationPrice          string                    `json:"AP"`
	CallbackRate             string                    `json:"cr"`
	RealizedProfit           string                    `json:"rp"`
}

// WsFuturesAccountConfigUpdate represents futures account configuration update
type WsFuturesAccountConfigUpdate struct {
	Symbol     string `json:"s"`
	Leverage   int    `json:"l"`
	MarginType string `json:"mt"`
}

// WsFuturesMarginCall represents futures margin call
type WsFuturesMarginCall struct {
	Positions []WsFuturesMarginCallPosition `json:"p"`
}

// WsFuturesMarginCallPosition represents futures margin call position
type WsFuturesMarginCallPosition struct {
	Symbol      string            `json:"s"`
	Side        string            `json:"ps"`
	Amount      string            `json:"pa"`
	MarginType  string            `json:"mt"`
	IsolatedWallet string         `json:"iw"`
	MarkPrice   string            `json:"mp"`
	UnrealizedPnL string          `json:"up"`
	MaintenanceMarginRequired string `json:"mm"`
}