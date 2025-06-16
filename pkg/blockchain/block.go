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
}

// Tạo Merkle Root từ danh sách giao dịch
func calculateMerkleRoot(txs []*Transaction) string {
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
			right = left // copy node cuối nếu lẻ
		}

		hash := sha256.Sum256(append(left, right...))
		newLevel = append(newLevel, hash[:])
	}

	return buildMerkleRoot(newLevel)
}

func hashBlock(b *Block) string {
	copyBlock := *b
	copyBlock.CurrentBlockHash = "" // bỏ self-hash trước khi hash
	data, _ := json.Marshal(copyBlock)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func NewBlock(txs []*Transaction, prevHash string) *Block {
	block := &Block{
		Transactions:  txs,
		PrevBlockHash: prevHash,
	}
	block.MerkleRoot = calculateMerkleRoot(txs)
	block.CurrentBlockHash = hashBlock(block)
	return block
}
