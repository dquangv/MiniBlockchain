package storage

import (
	"encoding/json"
	"fmt"
	"golang-chain/pkg/blockchain"

	"github.com/syndtr/goleveldb/leveldb"
)

// DB wraps the LevelDB instance for blockchain storage access
type DB struct {
	db *leveldb.DB
}

// NewDB opens or creates a LevelDB database at the given path
func NewDB(path string) (*DB, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db: ldb}, nil
}

// SaveBlock stores a block in the database
// It saves the block under three keys:
// - the block hash (for lookup by hash),
// - the block height (for sequential access),
// - and updates the "latest" pointer to this block
func (d *DB) SaveBlock(block *blockchain.Block) error {
	// Serialize the block to JSON
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	// Save block by hash
	err = d.db.Put([]byte(block.CurrentBlockHash), data, nil)
	if err != nil {
		return err
	}

	// Save block by height
	heightKey := []byte(fmt.Sprintf("height_%d", block.Height))
	if err := d.db.Put(heightKey, data, nil); err != nil {
		return err
	}

	// Update latest block pointer
	return d.db.Put([]byte("latest"), []byte(block.CurrentBlockHash), nil)
}

// GetBlock fetches a block by its hash
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

// GetLatestBlock returns the latest block based on the stored "latest" pointer
func (db *DB) GetLatestBlock() (*blockchain.Block, error) {
	hashBytes, err := db.db.Get([]byte("latest"), nil)
	if err != nil {
		return nil, err
	}
	return db.GetBlock(hashBytes)
}

// GetBlockByHeight fetches a block by its height in the chain
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
