package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
)

// Wallet chứa cặp khóa
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// Tạo ví mới
func NewWallet() (*Wallet, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: priv,
		PublicKey:  &priv.PublicKey,
	}, nil
}

// Encode public key thành []byte để gắn vào Sender
func EncodePublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

// Decode []byte trở lại *ecdsa.PublicKey
func DecodePublicKey(data []byte) (*ecdsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	pk, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not a valid ECDSA public key")
	}
	return pk, nil
}

func PublicKeyToAddress(pub *ecdsa.PublicKey) string {
	pubBytes := append(pub.X.Bytes(), pub.Y.Bytes()...)
	hash := sha256.Sum256(pubBytes)
	return hex.EncodeToString(hash[:])
}
