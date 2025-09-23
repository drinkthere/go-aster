package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/crypto/sha3"
)

// SignatureType represents different signature methods
type SignatureType int

const (
	SignatureTypeHMAC SignatureType = iota
	SignatureTypeWeb3
)

// HMACSignature creates HMAC SHA256 signature
func HMACSignature(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// Web3Signature creates Web3 style signature
func Web3Signature(params map[string]string, privateKey string) (string, error) {
	// Convert params to sorted string
	messageStr := ParamsToSortedString(params)
	
	// Create message hash
	messageHash := CreateKeccakHash(messageStr)
	
	// Sign the message
	signature, err := SignWithPrivateKey(privateKey, messageHash)
	if err != nil {
		return "", err
	}
	
	return signature, nil
}

// ParamsToSortedString converts params map to sorted query string
func ParamsToSortedString(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	var parts []string
	for _, k := range keys {
		if v := params[k]; v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
	}
	
	return strings.Join(parts, "&")
}

// CreateKeccakHash creates Keccak256 hash
func CreateKeccakHash(message string) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(message))
	return hasher.Sum(nil)
}

// SignWithPrivateKey signs message with private key
func SignWithPrivateKey(privateKeyHex string, messageHash []byte) (string, error) {
	// Remove 0x prefix if present
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	
	// Convert hex string to private key
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %v", err)
	}
	
	// Create ECDSA private key
	privateKey, err := HexToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create private key: %v", err)
	}
	
	// Sign the message
	signature, err := SignHash(messageHash, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}
	
	return "0x" + hex.EncodeToString(signature), nil
}

// ParseParamsFromURL parses query and body parameters
func ParseParamsFromURL(queryString, bodyString string) map[string]string {
	params := make(map[string]string)
	
	// Parse query parameters
	if queryString != "" {
		queryParams, _ := url.ParseQuery(queryString)
		for key, values := range queryParams {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
	}
	
	// Parse body parameters
	if bodyString != "" {
		bodyParams, _ := url.ParseQuery(bodyString)
		for key, values := range bodyParams {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
	}
	
	return params
}