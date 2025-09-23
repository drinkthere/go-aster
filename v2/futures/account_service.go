package futures

import (
	"context"
	"net/http"

	"github.com/drinkthere/go-aster/v2"
)

// GetAccountService get account info
type GetAccountService struct {
	C *aster.BaseClient
}

// Do send request
func (s *GetAccountService) Do(ctx context.Context, opts ...aster.RequestOption) (res *Account, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v2/account", aster.SecTypeSigned)
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Account)
	err = aster.JSON.Unmarshal(data, res)
	return res, err
}

// GetBalanceService get balance info
type GetBalanceService struct {
	C *aster.BaseClient
}

// Do send request
func (s *GetBalanceService) Do(ctx context.Context, opts ...aster.RequestOption) (res []Balance, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v2/balance", aster.SecTypeSigned)
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]Balance, 0)
	err = aster.JSON.Unmarshal(data, &res)
	return res, err
}

// GetPositionRiskService get position risk
type GetPositionRiskService struct {
	C      *aster.BaseClient
	symbol *string
}

// Symbol set symbol
func (s *GetPositionRiskService) Symbol(symbol string) *GetPositionRiskService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *GetPositionRiskService) Do(ctx context.Context, opts ...aster.RequestOption) (res []PositionRisk, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v2/positionRisk", aster.SecTypeSigned)
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]PositionRisk, 0)
	err = aster.JSON.Unmarshal(data, &res)
	return res, err
}

// ChangeLeverageService change leverage
type ChangeLeverageService struct {
	C        *aster.BaseClient
	symbol   string
	leverage int
}

// Symbol set symbol
func (s *ChangeLeverageService) Symbol(symbol string) *ChangeLeverageService {
	s.symbol = symbol
	return s
}

// Leverage set leverage
func (s *ChangeLeverageService) Leverage(leverage int) *ChangeLeverageService {
	s.leverage = leverage
	return s
}

// Do send request
func (s *ChangeLeverageService) Do(ctx context.Context, opts ...aster.RequestOption) (res *SymbolLeverage, err error) {
	r := aster.NewRequest(http.MethodPost, "/fapi/v1/leverage", aster.SecTypeSigned)
	r.SetFormParam("symbol", s.symbol)
	r.SetFormParam("leverage", s.leverage)
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(SymbolLeverage)
	err = aster.JSON.Unmarshal(data, res)
	return res, err
}

// SymbolLeverage symbol leverage response
type SymbolLeverage struct {
	Symbol           string `json:"symbol"`
	Leverage         int    `json:"leverage"`
	MaxNotionalValue string `json:"maxNotionalValue"`
}

// ChangeMarginTypeService change margin type
type ChangeMarginTypeService struct {
	C          *aster.BaseClient
	symbol     string
	marginType MarginType
}

// Symbol set symbol
func (s *ChangeMarginTypeService) Symbol(symbol string) *ChangeMarginTypeService {
	s.symbol = symbol
	return s
}

// MarginType set margin type
func (s *ChangeMarginTypeService) MarginType(marginType MarginType) *ChangeMarginTypeService {
	s.marginType = marginType
	return s
}

// Do send request
func (s *ChangeMarginTypeService) Do(ctx context.Context, opts ...aster.RequestOption) error {
	r := aster.NewRequest(http.MethodPost, "/fapi/v1/marginType", aster.SecTypeSigned)
	r.SetFormParam("symbol", s.symbol)
	r.SetFormParam("marginType", s.marginType)
	_, err := s.C.CallAPI(ctx, r, opts...)
	return err
}

// UpdatePositionMarginService update position margin
type UpdatePositionMarginService struct {
	C            *aster.BaseClient
	symbol       string
	positionSide *PositionSideType
	amount       string
	actionType   int
}

// Symbol set symbol
func (s *UpdatePositionMarginService) Symbol(symbol string) *UpdatePositionMarginService {
	s.symbol = symbol
	return s
}

// PositionSide set position side
func (s *UpdatePositionMarginService) PositionSide(positionSide PositionSideType) *UpdatePositionMarginService {
	s.positionSide = &positionSide
	return s
}

// Amount set amount
func (s *UpdatePositionMarginService) Amount(amount string) *UpdatePositionMarginService {
	s.amount = amount
	return s
}

// Type set action type (1: Add position margin, 2: Reduce position margin)
func (s *UpdatePositionMarginService) Type(actionType int) *UpdatePositionMarginService {
	s.actionType = actionType
	return s
}

// Do send request
func (s *UpdatePositionMarginService) Do(ctx context.Context, opts ...aster.RequestOption) error {
	r := aster.NewRequest(http.MethodPost, "/fapi/v1/positionMargin", aster.SecTypeSigned)
	r.SetFormParam("symbol", s.symbol)
	if s.positionSide != nil {
		r.SetFormParam("positionSide", *s.positionSide)
	}
	r.SetFormParam("amount", s.amount)
	r.SetFormParam("type", s.actionType)
	_, err := s.C.CallAPI(ctx, r, opts...)
	return err
}

// CommissionRateService get commission rate
type CommissionRateService struct {
	C      *aster.BaseClient
	symbol string
}

// Symbol set symbol
func (s *CommissionRateService) Symbol(symbol string) *CommissionRateService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CommissionRateService) Do(ctx context.Context, opts ...aster.RequestOption) (res *CommissionRate, err error) {
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/commissionRate", aster.SecTypeSigned)
	if s.symbol != "" {
		r.SetParam("symbol", s.symbol)
	}
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CommissionRate)
	err = aster.JSON.Unmarshal(data, res)
	return res, err
}

// CommissionRate represents commission rate
type CommissionRate struct {
	Symbol              string `json:"symbol"`
	MakerCommissionRate string `json:"makerCommissionRate"`
	TakerCommissionRate string `json:"takerCommissionRate"`
}