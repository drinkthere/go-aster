package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"math/big"
)

// HexToECDSA converts hex string to ECDSA private key
func HexToECDSA(hexkey []byte) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	
	if 8*len(hexkey) != priv.Params().BitSize {
		return nil, errors.New("invalid length, need 256 bits")
	}
	
	priv.D = new(big.Int).SetBytes(hexkey)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(hexkey)
	
	return priv, nil
}

// SignHash signs the given hash with the private key
func SignHash(hash []byte, priv *ecdsa.PrivateKey) ([]byte, error) {
	r, s, err := ecdsa.Sign(nil, priv, hash)
	if err != nil {
		return nil, err
	}
	
	// Serialize signature
	sig := make([]byte, 65)
	copy(sig[0:32], r.Bytes())
	copy(sig[32:64], s.Bytes())
	sig[64] = 0 // Recovery ID
	
	return sig, nil
}