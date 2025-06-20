package p2p

import (
	"log"
	"time"

	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/storage"
)

func StartLeaderLoop(db *storage.DB, peers []string) {
	// 🧱 Tạo Genesis Block nếu DB trống
	latest, err := db.GetLatestBlock()
	if err != nil || latest == nil {
		log.Println("🧱 No blocks found. Creating genesis block...")

		genesis := blockchain.CreateGenesisBlock()
		if err := db.SaveBlock(genesis); err != nil {
			log.Fatalln("❌ Failed to save genesis block:", err)
		} else {
			log.Println("✅ Genesis block created.")
		}
	}

	// ⏱ Tạo block định kỳ
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pending := blockchain.GetAndClearPendingTxs()
		if len(pending) == 0 {
			continue
		}

		// Lấy block mới nhất từ local DB
		latest, _ := db.GetLatestBlock()

		// Tính prevHash và height
		prevHash := ""
		newHeight := int64(0)
		if latest != nil {
			prevHash = latest.CurrentBlockHash
			newHeight = latest.Height + 1
		}

		// Tạo block mới
		block := blockchain.NewBlock(pending, prevHash, newHeight)

		pbBlock := ConvertBlockToPb(block)

		// Gửi đến các follower để vote
		votes := SendBlockForVote(peers, pbBlock)
		approveCount := 0
		for _, v := range votes {
			if v.Approved {
				approveCount++
			}
		}

		if approveCount >= 2 {
			BroadcastCommit(peers, pbBlock)
			log.Println("✅ Committed block at height", block.Height, "with", len(pending), "txs")
			db.SaveBlock(block)
		} else {
			log.Println("❌ Not enough votes to commit block at height", block.Height)
		}
	}
}
