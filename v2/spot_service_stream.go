package aster

import (
	"context"
	"net/http"
)

// StartSpotUserStreamService create listen key for user stream
type StartSpotUserStreamService struct {
	c *BaseClient
}

// Do send request
func (s *StartSpotUserStreamService) Do(ctx context.Context, opts ...RequestOption) (listenKey string, err error) {
	r := newRequest(http.MethodPost, "/api/v3/userDataStream", secTypeAPIKey)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return "", err
	}
	j, err := newJSON(data)
	if err != nil {
		return "", err
	}
	listenKey = (*j).Get("listenKey").ToString()
	return listenKey, nil
}

// KeepaliveSpotUserStreamService update listen key
type KeepaliveSpotUserStreamService struct {
	c         *BaseClient
	listenKey string
}

// ListenKey set listen key
func (s *KeepaliveSpotUserStreamService) ListenKey(listenKey string) *KeepaliveSpotUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *KeepaliveSpotUserStreamService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := newRequest(http.MethodPut, "/api/v3/userDataStream", secTypeAPIKey)
	r.setFormParam("listenKey", s.listenKey)
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// CloseSpotUserStreamService close user stream
type CloseSpotUserStreamService struct {
	c         *BaseClient
	listenKey string
}

// ListenKey set listen key
func (s *CloseSpotUserStreamService) ListenKey(listenKey string) *CloseSpotUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *CloseSpotUserStreamService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := newRequest(http.MethodDelete, "/api/v3/userDataStream", secTypeAPIKey)
	r.setFormParam("listenKey", s.listenKey)
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}