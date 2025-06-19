package main

import (
	"context"
	"fmt"
	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/p2p"
	"golang-chain/pkg/p2p/pb"
	"golang-chain/pkg/wallet"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("🚀 Leader node is proposing a new block...")

	// Tạo ví mới cho Alice
	aliceWallet, _ := wallet.NewWallet()
	encodedSender, _ := wallet.EncodePublicKey(aliceWallet.PublicKey)
	fmt.Println("public key: ", aliceWallet.PublicKey)

	// Tạo transaction từ Alice gửi Bob
	tx := blockchain.NewTransaction(encodedSender, []byte("bob"), 10)
	err := tx.Sign(aliceWallet.PrivateKey)
	if err != nil {
		log.Fatalln("❌ Failed to sign transaction:", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("❌ Can't connect to leader:", err)
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)
	resp, err := client.GetLatestBlock(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Println("⚠️  No latest block found:", err)
	}

	prevHash := ""
	if err != nil || resp == nil || resp.Block == nil {
		fmt.Println("⚠️  No latest block found:", err)
	} else {
		prevHash = resp.Block.CurrentBlockHash
	}

	log.Println("prevHash =", prevHash)

	// 4. Tạo block chứa transaction
	block := blockchain.NewBlock([]*blockchain.Transaction{tx}, prevHash)

	// Convert sang pb.Block để gửi qua gRPC
	pbBlock := p2p.ConvertBlockToPb(block)

	// Gửi block cho follower vote
	peers := []string{"localhost:50052", "localhost:50053"}
	votes := p2p.SendBlockForVote(peers, pbBlock)

	// Tính số phiếu đồng ý
	approveCount := 0
	for _, vote := range votes {
		if vote.Approved {
			approveCount++
		}
	}

	// Nếu đủ thì commit
	if approveCount >= 2 {
		p2p.BroadcastCommit(peers, pbBlock)
		fmt.Println("✅ Enough votes → commit block")
	} else {
		fmt.Println("❌ Not enough votes → discard block")
	}
}
