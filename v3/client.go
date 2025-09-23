package aster

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/drinkthere/go-aster/common"
	jsoniter "github.com/json-iterator/go"
)

const (
	baseAPIMainURL = "https://fapi.asterdex.com"
	baseAPITestURL = "https://testnet.asterdex.com"
)

var (
	useTestnet = false
)

func UseTestnet() {
	useTestnet = true
}

type doFunc func(req *http.Request) (*http.Response, error)

type Client struct {
	BaseURL      string
	UserAddress  string
	SignerAddress string
	PrivateKey   string
	HTTPClient   *http.Client
	Debug        bool
	Logger       *log.Logger
	TimeOffset   int64
	do           doFunc
}

func NewClient(userAddress, signerAddress, privateKey string) *Client {
	return &Client{
		BaseURL:      getAPIEndpoint(),
		UserAddress:  userAddress,
		SignerAddress: signerAddress,
		PrivateKey:   privateKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		Logger: log.New(os.Stderr, "[Aster-golang] ", log.LstdFlags),
	}
}

func (c *Client) init() {
	c.do = c.HTTPClient.Do
}

func getAPIEndpoint() string {
	if useTestnet {
		return baseAPITestURL
	}
	return baseAPIMainURL
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	for _, opt := range opts {
		opt(r)
	}
	if err = r.validate(); err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow)
	}
	if r.secType == secTypeSigned {
		r.setParam("userAddress", c.UserAddress)
		r.setParam("signerAddress", c.SignerAddress)
		nonce := time.Now().UnixMicro()
		r.setParam("nonce", nonce)
	}
	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.secType == secTypeSigned {
		signature, err := c.signRequest(queryString, bodyString, r)
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Set(signatureKey, signature)
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	if c.Debug {
		c.Logger.Printf("full url: %s, body: %s", fullURL, bodyString)
	}

	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	if c.Debug {
		c.Logger.Printf("request: %#v", req)
	}
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	if c.Debug {
		c.Logger.Printf("response: %#v", res)
		c.Logger.Printf("response body: %s", string(data))
		c.Logger.Printf("response status code: %d", res.StatusCode)
	}

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := jsoniter.Unmarshal(data, apiErr)
		if e != nil {
			if c.Debug {
				c.Logger.Printf("failed to unmarshal json: %s", e)
			}
		}
		return nil, apiErr
	}
	return data, nil
}

func (c *Client) SetServerTimeOffset(offset int64) {
	c.TimeOffset = offset
}