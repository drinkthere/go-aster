package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	secretKey := "f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c"
	timestamp := "1758619433599"

	// Test different parameter orders
	payload1 := fmt.Sprintf("timestamp=%s", timestamp)
	
	// Calculate signatures
	mac1 := hmac.New(sha256.New, []byte(secretKey))
	mac1.Write([]byte(payload1))
	sig1 := hex.EncodeToString(mac1.Sum(nil))
	
	fmt.Printf("Payload 1: %s\n", payload1)
	fmt.Printf("Signature 1: %s\n", sig1)
	
	// This is what the SDK seems to be doing - wrong!
	// It calculates signature on "timestamp=X" but sends "signature=Y&timestamp=X"
	fmt.Printf("\nThe SDK calculates signature on: timestamp=%s\n", timestamp)
	fmt.Printf("But sends URL with parameters in order: signature=%s&timestamp=%s\n", sig1, timestamp)
	fmt.Printf("This is incorrect!\n")
}