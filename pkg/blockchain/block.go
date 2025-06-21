package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	Transactions     []*Transaction
	MerkleRoot       string
	PrevBlockHash    string
	CurrentBlockHash string
	Height           int64
}

// CalculateMerkleRoot computes the Merkle root from all transactions in the block.
func CalculateMerkleRoot(txs []*Transaction) string {
	var txHashes [][]byte
	for _, tx := range txs {
		hash, err := tx.Hash()
		if err != nil {
			panic(err)
		}
		txHashes = append(txHashes, hash)
	}
	return hex.EncodeToString(buildMerkleRoot(txHashes))
}

// buildMerkleRoot recursively builds the Merkle tree and returns the root hash.
// If there's an odd number of nodes, the last one is duplicated to balance the tree.
func buildMerkleRoot(leaves [][]byte) []byte {
	if len(leaves) == 0 {
		return []byte{}
	}
	if len(leaves) == 1 {
		return leaves[0]
	}

	var newLevel [][]byte
	for i := 0; i < len(leaves); i += 2 {
		left := leaves[i]
		var right []byte
		if i+1 < len(leaves) {
			right = leaves[i+1]
		} else {
			right = left
		}
		hash := sha256.Sum256(append(left, right...))
		newLevel = append(newLevel, hash[:])
	}

	return buildMerkleRoot(newLevel)
}

// HashBlock computes a SHA-256 hash of the block's contents,
// excluding its own current hash to avoid circular dependency.
func HashBlock(b *Block) string {
	copyBlock := *b
	copyBlock.CurrentBlockHash = ""
	data, _ := json.Marshal(copyBlock)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// NewBlock creates a new block with the given transactions,
// links it to the previous block via PrevBlockHash, and calculates the Merkle root and hash.
func NewBlock(txs []*Transaction, prevHash string, height int64) *Block {
	block := &Block{
		Transactions:  txs,
		PrevBlockHash: prevHash,
		MerkleRoot:    CalculateMerkleRoot(txs),
		Height:        height,
	}
	block.CurrentBlockHash = HashBlock(block)
	return block
}

// CreateGenesisBlock initializes the first block of the blockchain.
// It has no transactions and no previous hash, and is always at height 0.
func CreateGenesisBlock() *Block {
	block := &Block{
		Height:        0,
		PrevBlockHash: "",
		Transactions:  []*Transaction{},
		MerkleRoot:    "",
	}
	block.CurrentBlockHash = HashBlock(block)
	return block
}
