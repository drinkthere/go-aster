package aster

import (
	"context"
	"net/http"

	"github.com/your-org/go-aster/v3/futures"
)

// GetAccountService get account
type GetAccountService struct {
	c *Client
}

// Do send request
func (s *GetAccountService) Do(ctx context.Context, opts ...RequestOption) (res *futures.Account, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v2/account",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(futures.Account)
	err = json.Unmarshal(data, res)
	return res, err
}

// GetBalanceService get balance
type GetBalanceService struct {
	c *Client
}

// Do send request
func (s *GetBalanceService) Do(ctx context.Context, opts ...RequestOption) (res []*futures.Balance, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v2/balance",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*futures.Balance, 0)
	err = json.Unmarshal(data, &res)
	return res, err
}

// GetPositionsService get positions
type GetPositionsService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *GetPositionsService) Symbol(symbol string) *GetPositionsService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *GetPositionsService) Do(ctx context.Context, opts ...RequestOption) (res []*futures.Position, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v2/positionRisk",
		secType:  secTypeSigned,
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*futures.Position, 0)
	err = json.Unmarshal(data, &res)
	return res, err
}