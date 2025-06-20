package main

import (
	"context"
	"fmt"
	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/p2p/pb"
	"golang-chain/pkg/wallet"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("🚀 Sending signed tx to leader...")

	// Tạo ví thật sự
	w, _ := wallet.NewWallet()
	pubBytes, _ := wallet.EncodePublicKey(w.PublicKey)

	tx := blockchain.NewTransaction(pubBytes, []byte("bob"), 10)

	// Ký bằng private key thật
	err := tx.Sign(w.PrivateKey)
	if err != nil {
		log.Fatalln("❌ Failed to sign tx:", err)
	}

	// Gửi qua gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)
	resp, err := client.SendTransaction(context.Background(), &pb.Transaction{
		Sender:    tx.Sender,
		Receiver:  tx.Receiver,
		Amount:    tx.Amount,
		Timestamp: tx.Timestamp,
		Signature: tx.Signature,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("📨", resp.Message)
}
