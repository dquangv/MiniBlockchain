package storage

import (
	"encoding/json"
	"math/big"
)

// Lưu số dư dưới dạng string để tránh mất độ chính xác khi marshal float
func (d *DB) SetBalance(address string, amount *big.Float) error {
	bytes, err := json.Marshal(amount.Text('f', 8)) // lưu dạng chuỗi
	if err != nil {
		return err
	}
	return d.db.Put([]byte("balance_"+address), bytes, nil)
}

func (d *DB) GetBalance(address string) (*big.Float, error) {
	data, err := d.db.Get([]byte("balance_"+address), nil)
	if err != nil {
		return big.NewFloat(0), nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return nil, err
	}

	val, _, err := big.ParseFloat(str, 10, 256, big.ToNearestEven)
	if err != nil {
		return nil, err
	}
	return val, nil
}
