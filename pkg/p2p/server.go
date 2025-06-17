package p2p

import (
	"context"
	"fmt"
	"log"
	"net"

	"golang-chain/pkg/p2p/pb"

	"google.golang.org/grpc"
)

type NodeServer struct {
	pb.UnimplementedNodeServiceServer
}

func (s *NodeServer) SendTransaction(ctx context.Context, tx *pb.Transaction) (*pb.TxResponse, error) {
	log.Printf("Received transaction from %s to %s (%.2f)", tx.Sender, tx.Receiver, tx.Amount)
	return &pb.TxResponse{
		Status:  "ok",
		Message: "tx received",
	}, nil
}

func (s *NodeServer) Ping(ctx context.Context, e *pb.Empty) (*pb.TxResponse, error) {
	return &pb.TxResponse{
		Status:  "pong",
		Message: "I'm alive",
	}, nil
}

func StartGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNodeServiceServer(grpcServer, &NodeServer{})

	fmt.Println("gRPC server listening on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *NodeServer) ProposeBlock(ctx context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	block := req.Block
	log.Printf("[Follower] Received proposed block: %s", block.CurrentBlockHash)

	// Giả sử block hợp lệ (chưa verify kỹ, sẽ bổ sung sau)
	vote := &pb.VoteResponse{
		NodeId:   "follower-1",
		Approved: true,
	}
	return vote, nil
}

func (s *NodeServer) CommitBlock(ctx context.Context, block *pb.Block) (*pb.TxResponse, error) {
	log.Printf("[Follower] Committing block: %s", block.CurrentBlockHash)

	// TODO: Convert pb.Block → blockchain.Block → Save vào LevelDB

	return &pb.TxResponse{
		Status:  "success",
		Message: "block committed",
	}, nil
}
