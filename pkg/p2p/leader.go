package p2p

import (
	"log"
	"math/big"
	"time"

	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/storage"
	"golang-chain/pkg/wallet"
)

// StartLeaderLoop runs on the leader node and periodically checks for pending transactions.
// If any exist, it creates a new block, proposes it to followers, and commits it if enough votes are received.
func StartLeaderLoop(db *storage.DB, peers []string) {
	// ⏱ Create a ticker that fires every 5 seconds (can be adjusted)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("⏳ Tick! Checking for pending transactions...")

		// 1. Retrieve and clear the list of pending transactions
		pending := blockchain.GetAndClearPendingTxs()
		if len(pending) == 0 {
			log.Println("🔍 No pending transactions. Skipping block creation.")
			continue
		}

		log.Printf("📨 Found %d pending transaction(s). Creating new block...", len(pending))

		// 2. Get the latest block from the local DB to determine previous hash and height
		latest, _ := db.GetLatestBlock()

		prevHash := ""
		newHeight := int64(0)
		if latest != nil {
			prevHash = latest.CurrentBlockHash
			newHeight = latest.Height + 1
		}

		// 3. Create a new block with the pending transactions
		block := blockchain.NewBlock(pending, prevHash, newHeight)
		pbBlock := ConvertBlockToPb(block)

		// 4. Propose the block to follower nodes and collect votes
		votes := SendBlockForVote(peers, pbBlock)
		approveCount := 0
		for _, v := range votes {
			if v.Approved {
				approveCount++
			}
		}

		// 5. Commit the block if majority votes are received (>=2 out of 3 nodes)
		if approveCount >= 2 {
			BroadcastCommit(peers, pbBlock) // Notify followers to commit
			log.Println("✅ Committed block at height", block.Height, "with", len(pending), "txs")
			db.SaveBlock(block) // Save block locally

			// ✅ Update balances
			for _, tx := range block.Transactions {
				sender := wallet.ResolveSenderName(tx.Sender)
				receiver := string(tx.Receiver)

				fromBal, _ := db.GetBalance(sender)
				toBal, _ := db.GetBalance(receiver)

				fromBal.Sub(fromBal, big.NewFloat(tx.Amount))
				toBal.Add(toBal, big.NewFloat(tx.Amount))

				db.SetBalance(sender, fromBal)
				db.SetBalance(receiver, toBal)
			}
		} else {
			log.Println("❌ Not enough votes to commit block at height", block.Height)
		}
	}

}
