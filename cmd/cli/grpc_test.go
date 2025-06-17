// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"
// 	"golang-chain/pkg/p2p/pb"
// 	"google.golang.org/grpc"
// )

// func main() {
// 	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()

// 	client := pb.NewNodeServiceClient(conn)

// 	// Gửi Ping thử
// 	resp, err := client.Ping(context.Background(), &pb.Empty{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Ping response:", resp)

// 	// Gửi Transaction thử
// 	tx := &pb.Transaction{
// 		Sender:    "alice",
// 		Receiver:  "bob",
// 		Amount:    10.5,
// 		Timestamp: time.Now().Unix(),
// 		Signature: []byte("dummy-sig"),
// 	}

// 	res, err := client.SendTransaction(context.Background(), tx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("SendTransaction response:", res)
// }
