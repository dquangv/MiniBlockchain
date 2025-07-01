package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"golang-chain/pkg/p2p/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	name := flag.String("name", "", "Wallet name")
	node := flag.String("node", "localhost:50051", "Node address (host:port)")
	flag.Parse()

	if *name == "" {
		log.Fatalln("⚠️  Usage: ./balance --name Alice")
	}

	conn, err := grpc.Dial(*node, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to node: %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)

	resp, err := client.GetBalance(context.Background(), &pb.BalanceRequest{Name: *name})
	if err != nil {
		log.Fatalf("❌ Failed to get balance: %v", err)
	}

	fmt.Printf("💰 Balance of %s: %s coins\n", *name, resp.Balance)
}
