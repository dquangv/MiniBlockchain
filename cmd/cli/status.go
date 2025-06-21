package main

import (
	"context"
	"flag"
	"fmt"
	"golang-chain/pkg/p2p/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	node := flag.String("node", "localhost:50051", "Địa chỉ node validator")
	flag.Parse()

	conn, err := grpc.Dial(*node, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("❌ Không thể kết nối node:", err)
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)
	resp, err := client.GetLatestBlock(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalln("❌ Không lấy được block:", err)
	}

	block := resp.Block
	fmt.Println("📦 The latest block:")
	fmt.Println("👉 Height:       ", block.Height)
	fmt.Println("👉 Hash:         ", block.CurrentBlockHash)
	fmt.Println("👉 Prev Hash:    ", block.PrevBlockHash)
	fmt.Println("👉 Tx count:     ", len(block.Transactions))
}
