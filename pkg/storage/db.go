package storage

import (
	"encoding/json"
	"golang-chain/pkg/blockchain"

	"github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
	db *leveldb.DB
}

func NewDB(path string) (*DB, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db: ldb}, nil
}

func (d *DB) SaveBlock(block *blockchain.Block) error {
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}
	return d.db.Put([]byte(block.CurrentBlockHash), data, nil)
}

func (d *DB) GetBlock(hash string) (*blockchain.Block, error) {
	data, err := d.db.Get([]byte(hash), nil)
	if err != nil {
		return nil, err
	}
	var block blockchain.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func (d *DB) Close() {
	d.db.Close()
}
