package consensus

import (
	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/wallet"
	"log"
)

// VerifyBlock checks whether a proposed block is valid before accepting it.
// It performs Merkle root verification, hash validation, previous block linkage, height consistency,
// and signature checks on all transactions.
func VerifyBlock(block, prevBlock *blockchain.Block) bool {
	// 1. Recompute and compare Merkle root to ensure integrity of transactions
	expectedMerkle := blockchain.CalculateMerkleRoot(block.Transactions)
	if block.MerkleRoot != expectedMerkle {
		return false
	}

	// 2. Recompute and compare block hash to detect tampering
	expectedHash := blockchain.HashBlock(block)
	if block.CurrentBlockHash != expectedHash {
		return false
	}

	// 3. If a previous block is provided, verify linkage and height
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
		// If this is a genesis block, its height must be 0
		if block.Height != 0 {
			log.Println("❌ Genesis block must have height 0")
			return false
		}
	}

	// 4. Verify digital signatures of all transactions in the block
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
