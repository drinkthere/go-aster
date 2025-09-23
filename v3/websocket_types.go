package aster

// WsDepthHandler websocket depth handler
type WsDepthHandler func(event *WsDepthEvent)

// WsDepthEvent websocket depth event
type WsDepthEvent struct {
	Event         string     `json:"e"`
	Time          int64      `json:"E"`
	Symbol        string     `json:"s"`
	FirstUpdateID int64      `json:"U"`
	LastUpdateID  int64      `json:"u"`
	Bids          [][]string `json:"b"`
	Asks          [][]string `json:"a"`
}

// WsAggTradeHandler websocket aggregate trade handler
type WsAggTradeHandler func(event *WsAggTradeEvent)

// WsAggTradeEvent websocket aggregate trade event
type WsAggTradeEvent struct {
	Event            string `json:"e"`
	Time             int64  `json:"E"`
	Symbol           string `json:"s"`
	AggregateTradeID int64  `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeID     int64  `json:"f"`
	LastTradeID      int64  `json:"l"`
	TradeTime        int64  `json:"T"`
	IsBuyerMaker     bool   `json:"m"`
}

// WsKlineHandler websocket kline handler
type WsKlineHandler func(event *WsKlineEvent)

// WsKlineEvent websocket kline event
type WsKlineEvent struct {
	Event  string  `json:"e"`
	Time   int64   `json:"E"`
	Symbol string  `json:"s"`
	Kline  WsKline `json:"k"`
}

// WsKline websocket kline
type WsKline struct {
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

// WsBookTickerHandler websocket book ticker handler
type WsBookTickerHandler func(event *WsBookTickerEvent)

// WsBookTickerEvent websocket book ticker event
type WsBookTickerEvent struct {
	Event         string `json:"e"`          // Event type
	UpdateID      int64  `json:"u"`          // Order book updateId
	EventTime     int64  `json:"E"`          // Event time
	TransactionTime int64 `json:"T"`        // Transaction time
	Symbol        string `json:"s"`          // Symbol
	BestBidPrice  string `json:"b"`          // Best bid price
	BestBidQty    string `json:"B"`          // Best bid quantity
	BestAskPrice  string `json:"a"`          // Best ask price
	BestAskQty    string `json:"A"`          // Best ask quantity
}

// WsUserDataEvent websocket user data event
type WsUserDataEvent struct {
	Event            string            `json:"e"`
	Time             int64             `json:"E"`
	TransactionTime  int64             `json:"T"`
	AccountUpdate    *WsAccountUpdate  `json:"a,omitempty"`
	OrderUpdate      *WsOrderUpdate    `json:"o,omitempty"`
	AccountConfigUpdate *WsAccountConfigUpdate `json:"ac,omitempty"`
}

// WsAccountUpdate account update
type WsAccountUpdate struct {
	Reason    string       `json:"m"`    // Event reason type
	Balances  []WsBalance  `json:"B"`    // Balances
	Positions []WsPosition `json:"P"`    // Positions
}

// WsBalance balance
type WsBalance struct {
	Asset              string `json:"a"`   // Asset
	Balance            string `json:"wb"`  // Wallet balance
	CrossWalletBalance string `json:"cw"`  // Cross wallet balance
}

// WsPosition position
type WsPosition struct {
	Symbol              string `json:"s"`   // Symbol
	PositionAmount      string `json:"pa"`  // Position amount
	EntryPrice          string `json:"ep"`  // Entry price
	AccumulatedRealized string `json:"cr"`  // (Pre-fee) accumulated realized
	UnrealizedPnL       string `json:"up"`  // Unrealized PnL
	MarginType          string `json:"mt"`  // Margin type
	IsolatedWallet      string `json:"iw"`  // Isolated wallet (if isolated position)
	PositionSide        string `json:"ps"`  // Position side
}

// WsOrderUpdate order update
type WsOrderUpdate struct {
	Symbol               string `json:"s"`      // Symbol
	ClientOrderID        string `json:"c"`      // Client order ID
	Side                 string `json:"S"`      // Side
	Type                 string `json:"o"`      // Order type
	TimeInForce          string `json:"f"`      // Time in force
	OriginalQty          string `json:"q"`      // Original quantity
	Price                string `json:"p"`      // Price
	AveragePrice         string `json:"ap"`     // Average price
	StopPrice            string `json:"sp"`     // Stop price
	ExecutionType        string `json:"x"`      // Execution type
	Status               string `json:"X"`      // Order status
	OrderID              int64  `json:"i"`      // Order ID
	LastFilledQty        string `json:"l"`      // Order last filled quantity
	FilledAccumulatedQty string `json:"z"`      // Order filled accumulated quantity
	LastFilledPrice      string `json:"L"`      // Last filled price
	CommissionAsset      string `json:"N"`      // Commission asset
	Commission           string `json:"n"`      // Commission
	OrderTradeTime       int64  `json:"T"`      // Order trade time
	TradeID              int64  `json:"t"`      // Trade ID
	BidsNotional         string `json:"b"`      // Bids notional
	AsksNotional         string `json:"a"`      // Asks notional
	IsMaker              bool   `json:"m"`      // Is this trade the maker side?
	IsReduceOnly         bool   `json:"R"`      // Is this reduce only
	WorkingType          string `json:"wt"`     // Stop price working type
	OriginalType         string `json:"ot"`     // Original order type
	PositionSide         string `json:"ps"`     // Position side
	ClosePosition        bool   `json:"cp"`     // If close all, pushed with conditional order
	ActivationPrice      string `json:"AP"`     // Activation price, only pushed with TRAILING_STOP_MARKET order
	CallbackRate         string `json:"cr"`     // Callback rate, only pushed with TRAILING_STOP_MARKET order
	RealizedProfit       string `json:"rp"`     // Realized profit of the trade
}

// WsAccountConfigUpdate account configuration update
type WsAccountConfigUpdate struct {
	Symbol     string `json:"s"`    // Symbol
	Leverage   int64  `json:"l"`    // Leverage
	MarginType string `json:"mt"`   // Margin type
}

// WsUserDataHandler handles all user data websocket events
type WsUserDataHandler func(event *WsUserDataEvent)