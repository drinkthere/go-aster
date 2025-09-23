package aster

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// secType represents security type
type secType int

const (
	secTypeNone secType = iota
	secTypeAPIKey
	secTypeSigned // private request
)

// params represents parameters map
type params map[string]interface{}

// request represents an API request
type request struct {
	method     string
	endpoint   string
	query      url.Values
	form       url.Values
	recvWindow int64
	secType    secType
	header     http.Header
	body       io.Reader
	fullURL    string
}

// RequestOption represents optional parameters for requests
type RequestOption func(*request)

// WithRecvWindow sets receive window
func WithRecvWindow(recvWindow int64) RequestOption {
	return func(r *request) {
		r.recvWindow = recvWindow
	}
}

// WithHeader sets a header
func WithHeader(key, value string, replace bool) RequestOption {
	return func(r *request) {
		if r.header == nil {
			r.header = http.Header{}
		}
		if replace {
			r.header.Set(key, value)
		} else {
			r.header.Add(key, value)
		}
	}
}

// WithHeaders sets headers
func WithHeaders(header http.Header) RequestOption {
	return func(r *request) {
		r.header = header.Clone()
	}
}

// newRequest creates a new request
func newRequest(method, endpoint string, secType secType) *request {
	return &request{
		method:   method,
		endpoint: endpoint,
		secType:  secType,
		query:    url.Values{},
		form:     url.Values{},
		header:   http.Header{},
	}
}

// setParam sets a query parameter
func (r *request) setParam(key string, value interface{}) *request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Set(key, fmt.Sprintf("%v", value))
	return r
}

// setParams sets multiple query parameters
func (r *request) setParams(m params) *request {
	for k, v := range m {
		r.SetParam(k, v)
	}
	return r
}

// setFormParam sets a form parameter
func (r *request) setFormParam(key string, value interface{}) *request {
	if r.form == nil {
		r.form = url.Values{}
	}
	r.form.Set(key, fmt.Sprintf("%v", value))
	return r
}

// setFormParams sets multiple form parameters
func (r *request) setFormParams(m params) *request {
	for k, v := range m {
		r.setFormParam(k, v)
	}
	return r
}

// validate validates the request
func (r *request) validate() (err error) {
	if r.query == nil {
		r.query = url.Values{}
	}
	if r.form == nil {
		r.form = url.Values{}
	}
	return nil
}