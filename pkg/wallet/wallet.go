package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
)

// Wallet stores a pair of ECDSA private/public keys
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// NewWallet generates a new ECDSA key pair and returns a wallet instance
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

// EncodePublicKey converts an ECDSA public key to a PEM-encoded []byte
func EncodePublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	return pem.EncodeToMemory(block), nil
}

// DecodePublicKey converts a PEM-encoded []byte back into an ECDSA public key
func DecodePublicKey(data []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM block for public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pk, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not a valid ECDSA public key")
	}
	return pk, nil
}

// EncodePrivateKey converts an ECDSA private key to a PEM-encoded []byte
func EncodePrivateKey(priv *ecdsa.PrivateKey) ([]byte, error) {
	der, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	}
	return pem.EncodeToMemory(block), nil
}

// DecodePrivateKey converts a PEM-encoded []byte back into an ECDSA private key
func DecodePrivateKey(data []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("invalid PEM block for private key")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

// PublicKeyToAddress generates a wallet address from an ECDSA public key
// by hashing its X and Y coordinates with SHA256
func PublicKeyToAddress(pub *ecdsa.PublicKey) string {
	pubBytes := append(pub.X.Bytes(), pub.Y.Bytes()...)
	hash := sha256.Sum256(pubBytes)
	return hex.EncodeToString(hash[:])
}

// LoadWallet reads a wallet from a JSON file in the "wallets/" folder
// The file should contain base64 PEM strings for the public and private key
func LoadWallet(name string) (*Wallet, error) {
	path := filepath.Join("wallets", name+"_wallet.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var keyData map[string]string
	if err := json.Unmarshal(data, &keyData); err != nil {
		return nil, err
	}

	privKey, err := DecodePrivateKey([]byte(keyData["privateKey"]))
	if err != nil {
		return nil, err
	}
	pubKey, err := DecodePublicKey([]byte(keyData["publicKey"]))
	if err != nil {
		return nil, err
	}

	if privKey == nil || pubKey == nil {
		return nil, errors.New("invalid keys")
	}

	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}, nil
}

// ResolveSenderName attempts to match a given public key to a wallet name by
// scanning through all JSON wallet files in the "wallets/" directory
func ResolveSenderName(pub []byte) string {
	files, err := os.ReadDir("wallets")
	if err != nil {
		return "Unknown"
	}

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join("wallets", f.Name()))
		if err != nil {
			continue
		}

		var keyMap map[string]string
		if err := json.Unmarshal(data, &keyMap); err != nil {
			continue
		}

		if keyMap["publicKey"] == string(pub) {
			// Strip off "_wallet.json" to get the user-friendly wallet name
			return f.Name()[:len(f.Name())-12]
		}
	}

	return "Unknown"
}

func WalletExists(name string) bool {
	path := filepath.Join("wallets", name+"_wallet.json")
	_, err := os.Stat(path)
	return err == nil
}
