package p2p

import (
	"log"
	"time"

	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/storage"
)

func StartLeaderLoop(db *storage.DB, peers []string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pending := blockchain.GetAndClearPendingTxs()
		if len(pending) == 0 {
			continue
		}

		latest, _ := db.GetLatestBlock()
		prevHash := ""
		if latest != nil {
			prevHash = latest.CurrentBlockHash
		}

		block := blockchain.NewBlock(pending, prevHash)
		pbBlock := ConvertBlockToPb(block)

		votes := SendBlockForVote(peers, pbBlock)
		approveCount := 0
		for _, v := range votes {
			if v.Approved {
				approveCount++
			}
		}

		if approveCount >= 2 {
			BroadcastCommit(peers, pbBlock)
			log.Println("✅ Committed block with", len(pending), "txs")
			db.SaveBlock(block)
		} else {
			log.Println("❌ Not enough votes to commit block.")
		}
	}
}
