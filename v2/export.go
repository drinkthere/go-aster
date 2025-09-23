package aster

import (
	"context"

	"github.com/json-iterator/go"
)

// Export some internal functions for sub-packages

var (
	// JSON is json-iterator instance
	JSON = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Export newRequest
var NewRequest = newRequest

// Export newJSON
var NewJSON = newJSON

// Export security types
const (
	SecTypeNone   = secTypeNone
	SecTypeAPIKey = secTypeAPIKey
	SecTypeSigned = secTypeSigned
)

// Export BaseClient methods
func (c *BaseClient) CallAPI(ctx context.Context, r *request, opts ...RequestOption) ([]byte, error) {
	return c.callAPI(ctx, r, opts...)
}

// Export request methods for packages
func (r *request) SetParam(key string, value interface{}) *request {
	return r.setParam(key, value)
}

func (r *request) SetFormParam(key string, value interface{}) *request {
	return r.setFormParam(key, value)
}

func (r *request) SetFormParams(m params) *request {
	return r.setFormParams(m)
}

// Export params type
type Params = params