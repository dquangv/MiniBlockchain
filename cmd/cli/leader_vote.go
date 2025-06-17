package main

import (
	"fmt"
	"golang-chain/pkg/p2p"
	"golang-chain/pkg/p2p/pb"
)

func main() {
	// Giả lập block (thực tế nên convert từ blockchain.Block → pb.Block)
	block := &pb.Block{
		MerkleRoot:       "abc123",
		PrevBlockHash:    "def456",
		CurrentBlockHash: "xyz789",
	}

	peers := []string{"localhost:50052", "localhost:50053"}
	votes := p2p.SendBlockForVote(peers, block)

	approveCount := 0
	for _, vote := range votes {
		if vote.Approved {
			approveCount++
		}
	}

	if approveCount >= 2 {
		fmt.Println("Enough votes → commit block")
		p2p.BroadcastCommit(peers, block)
	} else {
		fmt.Println("Not enough votes → discard block")
	}
}
