package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

// Custom HTTP client to intercept requests
type debugTransport struct {
	Transport http.RoundTripper
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf("\n=== DEBUG HTTP Request ===\n")
	fmt.Printf("Method: %s\n", req.Method)
	fmt.Printf("URL: %s\n", req.URL.String())
	fmt.Printf("Headers:\n")
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	
	// Parse query to see parameter order
	fmt.Printf("\nQuery Parameters (parsed):\n")
	query := req.URL.Query()
	for key, values := range query {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	
	// Show raw query string
	fmt.Printf("\nRaw Query String: %s\n", req.URL.RawQuery)
	
	return d.Transport.RoundTrip(req)
}

func main() {
	apiKey := "1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993"
	secretKey := "f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c"

	// Create custom HTTP client with debug transport
	httpClient := &http.Client{
		Transport: &debugTransport{
			Transport: http.DefaultTransport,
		},
	}

	// Create futures client with custom HTTP client
	client := aster.NewFuturesClient(apiKey, secretKey,
		aster.WithHTTPClient(httpClient),
		aster.WithDebug(true))

	fmt.Println("=== Testing Aster Futures SDK Signature ===")

	// Test account service
	fmt.Println("\nTesting GetAccount...")
	accountService := &futures.GetAccountService{C: client}
	_, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Let's also check how URL encoding works
	fmt.Println("\n=== URL Encoding Test ===")
	v := url.Values{}
	v.Set("timestamp", "1234567890")
	v.Set("signature", "test")
	fmt.Printf("url.Values.Encode(): %s\n", v.Encode())
}