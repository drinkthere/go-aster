package aster

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/crypto/sha3"
	"github.com/drinkthere/go-aster/common"
)

func (c *Client) signRequest(queryString, bodyString string, r *request) (string, error) {
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
	
	// Convert params to sorted string
	messageStr := paramsToSortedString(params)
	
	// Create message hash
	messageHash := createMessageHash(messageStr)
	
	// Sign the message
	signature, err := signMessage(c.PrivateKey, messageHash)
	if err != nil {
		return "", err
	}
	
	return signature, nil
}

func paramsToSortedString(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	
	return strings.Join(parts, "&")
}

func createMessageHash(message string) []byte {
	// Create Keccak256 hash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(message))
	return hasher.Sum(nil)
}

func signMessage(privateKeyHex string, messageHash []byte) (string, error) {
	// Remove 0x prefix if present
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	
	// Convert hex string to private key
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %v", err)
	}
	
	// Create ECDSA private key
	privateKey, err := common.HexToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create private key: %v", err)
	}
	
	// Sign the message
	signature, err := common.SignHash(messageHash, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}
	
	return "0x" + hex.EncodeToString(signature), nil
}