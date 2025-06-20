package storage

import (
	"encoding/json"
	"fmt"
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

	err = d.db.Put([]byte(block.CurrentBlockHash), data, nil)
	if err != nil {
		return err
	}

	// Save by height (new)
	heightKey := []byte(fmt.Sprintf("height_%d", block.Height))
	if err := d.db.Put(heightKey, data, nil); err != nil {
		return err
	}

	// Save latest hash
	return d.db.Put([]byte("latest"), []byte(block.CurrentBlockHash), nil)
}

func (d *DB) GetBlock(hash []byte) (*blockchain.Block, error) {
	data, err := d.db.Get(hash, nil)
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

func (db *DB) GetLatestBlock() (*blockchain.Block, error) {
	hashBytes, err := db.db.Get([]byte("latest"), nil)
	if err != nil {
		return nil, err
	}
	return db.GetBlock(hashBytes)
}

func (d *DB) GetBlockByHeight(height int64) (*blockchain.Block, error) {
	data, err := d.db.Get([]byte(fmt.Sprintf("height_%d", height)), nil)
	if err != nil {
		return nil, err
	}
	var blk blockchain.Block
	if err := json.Unmarshal(data, &blk); err != nil {
		return nil, err
	}
	return &blk, nil
}
