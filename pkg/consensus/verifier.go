package consensus

import (
	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/wallet"
	"log"
)

// Xác minh tính hợp lệ của block
func VerifyBlock(block, prevBlock *blockchain.Block) bool {
	// 1. Check Merkle Root
	expectedMerkle := blockchain.CalculateMerkleRoot(block.Transactions)
	if block.MerkleRoot != expectedMerkle {
		return false
	}

	// 2. Check block hash
	expectedHash := blockchain.HashBlock(block)
	if block.CurrentBlockHash != expectedHash {
		return false
	}

	// ✅ 3. Nếu có prevBlock thì check PrevBlockHash
	if prevBlock != nil {
		if block.PrevBlockHash != prevBlock.CurrentBlockHash {
			log.Println("❌ Prev hash mismatch")
			return false
		}
		if block.Height != prevBlock.Height+1 {
			log.Println("❌ Invalid height:", block.Height, "vs", prevBlock.Height+1)
			return false
		}
	} else {
		// Nếu là Genesis block, height phải là 0
		if block.Height != 0 {
			log.Println("❌ Genesis block must have height 0")
			return false
		}
	}

	// 4. Check chữ ký từng transaction
	for _, tx := range block.Transactions {
		pubKey, err := wallet.DecodePublicKey(tx.Sender)
		if err != nil {
			return false
		}
		valid, err := tx.Verify(pubKey)
		if err != nil || !valid {
			return false
		}
	}

	return true
}
