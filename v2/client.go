package aster

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/drinkthere/go-aster/v2/common"
	"github.com/json-iterator/go"
)

// newJSON parses response data to JSON
func newJSON(data []byte) (*jsoniter.Any, error) {
	j := jsoniter.Get(data)
	return &j, nil
}

const (
	// Spot API endpoints
	baseSpotAPIURL     = "https://sapi.asterdex.com"
	baseSpotAPITestURL = "https://testnet-sapi.asterdex.com"

	// Futures API endpoints
	baseFuturesAPIURL         = "https://fapi.asterdex.com"
	baseFuturesAPIIntranetURL = "https://fapi3.asterdex.com"
	baseFuturesAPITestURL     = "https://testnet.asterdex.com"

	// WebSocket endpoints
	baseWsMainURL    = "wss://stream.asterdex.com:9443"
	baseWsTestURL    = "wss://testnet.asterdex.com"
	baseWsFuturesURL = "wss://fstream.asterdex.com"
)

// doFunc represents the function to do HTTP request
type doFunc func(req *http.Request) (*http.Response, error)

// BaseClient is the base client for all APIs
type BaseClient struct {
	APIKey       string
	SecretKey    string
	BaseURL      string
	UserAgent    string
	HTTPClient   *http.Client
	Debug        bool
	Logger       *log.Logger
	TimeOffset   int64
	do           doFunc
	LocalAddress string // Local IP address for outbound connections

	// For futures API with Web3 signature
	UserAddress   string
	SignerAddress string
	PrivateKey    string
	SignatureType common.SignatureType
}

// ClientOption is a function to set options for client
type ClientOption func(*BaseClient)

// WithAPIKey sets API key
func WithAPIKey(key string) ClientOption {
	return func(c *BaseClient) {
		c.APIKey = key
	}
}

// WithSecretKey sets secret key
func WithSecretKey(secret string) ClientOption {
	return func(c *BaseClient) {
		c.SecretKey = secret
	}
}

// WithUserAddress sets user address (for futures)
func WithUserAddress(address string) ClientOption {
	return func(c *BaseClient) {
		c.UserAddress = address
	}
}

// WithSignerAddress sets signer address (for futures)
func WithSignerAddress(address string) ClientOption {
	return func(c *BaseClient) {
		c.SignerAddress = address
	}
}

// WithPrivateKey sets private key (for futures)
func WithPrivateKey(key string) ClientOption {
	return func(c *BaseClient) {
		c.PrivateKey = key
	}
}

// WithBaseURL sets base URL
func WithBaseURL(url string) ClientOption {
	return func(c *BaseClient) {
		c.BaseURL = url
	}
}

// WithHTTPClient sets HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *BaseClient) {
		c.HTTPClient = client
	}
}

// WithDebug sets debug mode
func WithDebug(debug bool) ClientOption {
	return func(c *BaseClient) {
		c.Debug = debug
	}
}

// WithLocalAddress sets local IP address for outbound connections
func WithLocalAddress(localAddr string) ClientOption {
	return func(c *BaseClient) {
		c.LocalAddress = localAddr
	}
}

// NewBaseClient creates a new base client
func NewBaseClient(opts ...ClientOption) *BaseClient {
	c := &BaseClient{
		UserAgent: "go-aster/2.0",
		Logger:    log.New(os.Stderr, "[ASTER] ", log.LstdFlags),
	}

	// Apply options first
	for _, opt := range opts {
		opt(c)
	}

	// Create HTTP client with optional local address binding
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// If LocalAddress is specified, configure the dialer
	if c.LocalAddress != "" {
		localAddr, err := net.ResolveTCPAddr("tcp", c.LocalAddress+":0")
		if err == nil {
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := &net.Dialer{
					LocalAddr: localAddr,
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}
				return dialer.DialContext(ctx, network, addr)
			}
		}
	}

	c.HTTPClient = &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return c
}

func (c *BaseClient) parseRequest(r *request, opts ...RequestOption) (err error) {
	// Apply request options
	for _, opt := range opts {
		opt(r)
	}
	if err = r.validate(); err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)
	if r.recvWindow > 0 {
		r.SetParam("recvWindow", r.recvWindow)
	}

	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}

	// Set headers
	header.Set("User-Agent", c.UserAgent)
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}

	// Handle authentication
	if r.secType == secTypeAPIKey || r.secType == secTypeSigned {
		if c.APIKey == "" {
			return fmt.Errorf("API key is required")
		}
		header.Set("X-MBX-APIKEY", c.APIKey)
		if c.Debug {
			c.Logger.Printf("Setting API Key header: X-MBX-APIKEY = %s...", c.APIKey[:20])
		}
	}

	// Handle signature
	if r.secType == secTypeSigned {
		timestamp := time.Now().UnixNano() / 1e6
		r.SetParam("timestamp", timestamp)

		// For futures with Web3 signature
		if c.SignatureType == common.SignatureTypeWeb3 {
			if c.UserAddress == "" || c.SignerAddress == "" || c.PrivateKey == "" {
				return fmt.Errorf("user address, signer address and private key are required for Web3 signature")
			}
			r.SetParam("userAddress", c.UserAddress)
			r.SetParam("signerAddress", c.SignerAddress)
			r.SetParam("nonce", time.Now().UnixMicro())

			// Get updated query string after adding auth params
			queryString = r.query.Encode()

			// Parse all params
			params := common.ParseParamsFromURL(queryString, bodyString)

			// Create Web3 signature
			signature, err := common.Web3Signature(params, c.PrivateKey)
			if err != nil {
				return err
			}
			r.SetParam("signature", signature)
			// Update query string with signature
			queryString = r.query.Encode()
		} else {
			// HMAC signature for spot
			if c.SecretKey == "" {
				return fmt.Errorf("secret key is required for HMAC signature")
			}

			// Build query string without signature first
			queryString = r.query.Encode()
			payload := queryString + bodyString
			signature := common.HMACSignature(payload, c.SecretKey)

			// Now append signature to the query string manually to ensure it's last
			if queryString != "" {
				queryString = queryString + "&signature=" + signature
			} else {
				queryString = "signature=" + signature
			}
		}
	}

	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	if c.Debug {
		c.Logger.Printf("Full URL: %s, Body: %s", fullURL, bodyString)
	}

	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *BaseClient) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
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
		c.Logger.Printf("Request: %#v", req)
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
		c.Logger.Printf("Response Status: %d, Body: %s", res.StatusCode, string(data))
	}

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := JSON.Unmarshal(data, apiErr)
		if e != nil {
			if c.Debug {
				c.Logger.Printf("Failed to unmarshal error: %s", e)
			}
			return nil, fmt.Errorf("request failed with status %d: %s", res.StatusCode, string(data))
		}
		return nil, apiErr
	}
	return data, nil
}

// SetServerTimeOffset sets time offset
func (c *BaseClient) SetServerTimeOffset(offset int64) *sync.Map {
	c.TimeOffset = offset
	return nil
}

// NewSpotClient creates a spot trading client (using HMAC signature)
func NewSpotClient(apiKey, secretKey string, opts ...ClientOption) *BaseClient {
	// Default options for spot
	defaultOpts := []ClientOption{
		WithAPIKey(apiKey),
		WithSecretKey(secretKey),
		WithBaseURL(baseSpotAPIURL),
	}

	// Append user options
	defaultOpts = append(defaultOpts, opts...)

	client := NewBaseClient(defaultOpts...)
	client.SignatureType = common.SignatureTypeHMAC
	return client
}

// makeFuturesClient creates a futures trading client (using HMAC signature)
func makeFuturesClient(apiKey, secretKey string, useIntranet bool, opts ...ClientOption) *BaseClient {
	// Default options for futures
	defaultOpts := []ClientOption{
		WithAPIKey(apiKey),
		WithSecretKey(secretKey),
	}
	if useIntranet {
		defaultOpts = append(defaultOpts, WithBaseURL(baseFuturesAPIIntranetURL))
	} else {
		defaultOpts = append(defaultOpts, WithBaseURL(baseFuturesAPIURL))
	}

	// Append user options
	defaultOpts = append(defaultOpts, opts...)

	client := NewBaseClient(defaultOpts...)
	client.SignatureType = common.SignatureTypeHMAC
	return client
}

// NewFuturesClient creates a futures trading client (using HMAC signature)
func NewFuturesClient(apiKey, secretKey string, opts ...ClientOption) *BaseClient {

	return makeFuturesClient(apiKey, secretKey, false, opts...)
}

// NewFuturesIntranetClient creates a futures trading intranet client (using HMAC signature)
func NewFuturesIntranetClient(apiKey, secretKey string, opts ...ClientOption) *BaseClient {

	return makeFuturesClient(apiKey, secretKey, true, opts...)
}

// makeFuturesClientWithWeb3 creates a futures trading client (using Web3 signature)
func makeFuturesClientWithWeb3(userAddress, signerAddress, privateKey string, useIntranet bool, opts ...ClientOption) *BaseClient {
	// Default options for futures
	defaultOpts := []ClientOption{
		WithUserAddress(userAddress),
		WithSignerAddress(signerAddress),
		WithPrivateKey(privateKey),
	}
	if useIntranet {
		defaultOpts = append(defaultOpts, WithBaseURL(baseFuturesAPIIntranetURL))
	} else {
		defaultOpts = append(defaultOpts, WithBaseURL(baseFuturesAPIURL))
	}

	// Append user options
	defaultOpts = append(defaultOpts, opts...)

	client := NewBaseClient(defaultOpts...)
	client.SignatureType = common.SignatureTypeWeb3
	return client
}

// NewFuturesClientWithWeb3 creates a futures trading client (using Web3 signature)
func NewFuturesClientWithWeb3(userAddress, signerAddress, privateKey string, opts ...ClientOption) *BaseClient {

	return makeFuturesClientWithWeb3(userAddress, signerAddress, privateKey, false, opts...)
}

// NewFuturesIntranetClientWithWeb3 creates a futures trading client (using Web3 signature)
func NewFuturesIntranetClientWithWeb3(userAddress, signerAddress, privateKey string, opts ...ClientOption) *BaseClient {

	return makeFuturesClientWithWeb3(userAddress, signerAddress, privateKey, true, opts...)
}
