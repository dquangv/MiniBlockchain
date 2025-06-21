package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"time"
)

type Transaction struct {
	Sender    []byte
	Receiver  []byte
	Amount    float64
	Timestamp int64
	Signature []byte
}

func NewTransaction(sender, receiver []byte, amount float64) *Transaction {
	return &Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}
}

// Hash calculates the SHA-256 hash of the transaction data.
func (t *Transaction) Hash() ([]byte, error) {
	txMap := map[string]interface{}{
		"sender":    hex.EncodeToString(t.Sender),
		"receiver":  hex.EncodeToString(t.Receiver),
		"amount":    t.Amount,
		"timestamp": t.Timestamp,
	}

	jsonData, err := json.Marshal(txMap)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(jsonData)
	return hash[:], nil
}

// Sign signs the transaction hash using the sender's private key.
// It uses ECDSA and embeds the resulting signature in the transaction.
func (t *Transaction) Sign(priv *ecdsa.PrivateKey) error {
	hash, err := t.Hash()
	if err != nil {
		return err
	}
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash)
	if err != nil {
		return err
	}

	// Normalize the signature to fixed 64 bytes: 32 bytes for R + 32 bytes for S
	rBytes := r.FillBytes(make([]byte, 32))
	sBytes := s.FillBytes(make([]byte, 32))
	t.Signature = append(rBytes, sBytes...)
	return nil
}

// Verify checks whether the transaction's signature is valid
// using the sender's public key. It ensures authenticity and integrity.
func (t *Transaction) Verify(pub *ecdsa.PublicKey) (bool, error) {
	if len(t.Signature) != 64 {
		return false, errors.New("invalid signature length")
	}

	hash, err := t.Hash()
	if err != nil {
		return false, err
	}

	r := new(big.Int).SetBytes(t.Signature[:32])
	s := new(big.Int).SetBytes(t.Signature[32:])
	return ecdsa.Verify(pub, hash, r, s), nil
}
