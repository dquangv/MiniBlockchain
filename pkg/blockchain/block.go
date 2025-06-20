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

func HashBlock(b *Block) string {
	copyBlock := *b
	copyBlock.CurrentBlockHash = ""
	data, _ := json.Marshal(copyBlock)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

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
