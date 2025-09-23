package futures

import (
	"context"
	"net/http"

	"github.com/drinkthere/go-aster/v2"
)

// StartUserStreamService create listen key for user stream
type StartUserStreamService struct {
	C *aster.BaseClient
}

// Do send request
func (s *StartUserStreamService) Do(ctx context.Context, opts ...aster.RequestOption) (listenKey string, err error) {
	r := aster.NewRequest(http.MethodPost, "/fapi/v1/listenKey", aster.SecTypeSigned)
	data, err := s.C.CallAPI(ctx, r, opts...)
	if err != nil {
		return "", err
	}
	j, err := aster.NewJSON(data)
	if err != nil {
		return "", err
	}
	listenKey = (*j).Get("listenKey").ToString()
	return listenKey, nil
}

// KeepaliveUserStreamService update listen key
type KeepaliveUserStreamService struct {
	C         *aster.BaseClient
	listenKey string
}

// ListenKey set listen key
func (s *KeepaliveUserStreamService) ListenKey(listenKey string) *KeepaliveUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *KeepaliveUserStreamService) Do(ctx context.Context, opts ...aster.RequestOption) (err error) {
	r := aster.NewRequest(http.MethodPut, "/fapi/v1/listenKey", aster.SecTypeSigned)
	r.SetFormParam("listenKey", s.listenKey)
	_, err = s.C.CallAPI(ctx, r, opts...)
	return err
}

// CloseUserStreamService close user stream
type CloseUserStreamService struct {
	C         *aster.BaseClient
	listenKey string
}

// ListenKey set listen key
func (s *CloseUserStreamService) ListenKey(listenKey string) *CloseUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *CloseUserStreamService) Do(ctx context.Context, opts ...aster.RequestOption) (err error) {
	r := aster.NewRequest(http.MethodDelete, "/fapi/v1/listenKey", aster.SecTypeSigned)
	r.SetFormParam("listenKey", s.listenKey)
	_, err = s.C.CallAPI(ctx, r, opts...)
	return err
}