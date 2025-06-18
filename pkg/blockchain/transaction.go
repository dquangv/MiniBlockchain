package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
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

// Hash nội dung giao dịch (không gồm chữ ký)
func (t *Transaction) Hash() ([]byte, error) {
	txCopy := *t
	txCopy.Signature = nil
	jsonData, err := json.Marshal(txCopy)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(jsonData)
	return hash[:], nil
}

// Ký giao dịch bằng private key
func (t *Transaction) Sign(priv *ecdsa.PrivateKey) error {
	hash, err := t.Hash()
	if err != nil {
		return err
	}
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash)
	if err != nil {
		return err
	}

	rBytes := r.FillBytes(make([]byte, 32)) // chuẩn hóa thành 32 byte
	sBytes := s.FillBytes(make([]byte, 32))
	t.Signature = append(rBytes, sBytes...)
	return nil
}

// Xác thực chữ ký giao dịch
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
