package p2p

import (
	"context"
	"fmt"
	"log"
	"net"

	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/p2p/pb"
	"golang-chain/pkg/storage"

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

func convertPbBlock(pbBlock *pb.Block) *blockchain.Block {
	var txs []*blockchain.Transaction
	for _, tx := range pbBlock.Transactions {
		txs = append(txs, &blockchain.Transaction{
			Sender:    tx.Sender,
			Receiver:  tx.Receiver,
			Amount:    tx.Amount,
			Timestamp: tx.Timestamp,
			Signature: tx.Signature,
		})
	}

	return &blockchain.Block{
		Transactions:     txs,
		MerkleRoot:       pbBlock.MerkleRoot,
		PrevBlockHash:    pbBlock.PrevBlockHash,
		CurrentBlockHash: pbBlock.CurrentBlockHash,
	}
}

func (s *NodeServer) CommitBlock(ctx context.Context, pbBlock *pb.Block) (*pb.TxResponse, error) {
	block := convertPbBlock(pbBlock)
	db, err := storage.NewDB("blockdata")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.SaveBlock(block)
	if err != nil {
		return nil, err
	}

	log.Printf("[Follower] Block committed: %s", block.CurrentBlockHash)
	return &pb.TxResponse{
		Status:  "success",
		Message: "block saved",
	}, nil
}

func (s *NodeServer) GetLatestBlock(ctx context.Context, _ *pb.Empty) (*pb.BlockResponse, error) {
	db, _ := storage.NewDB("blockdata")
	defer db.Close()

	latestBlock, err := db.GetLatestBlock()
	if err != nil {
		return nil, err
	}

	return &pb.BlockResponse{Block: convertBlockToPb(latestBlock)}, nil
}

func (s *NodeServer) GetBlock(ctx context.Context, req *pb.BlockRequest) (*pb.BlockResponse, error) {
	db, err := storage.NewDB("node1_db") // ❗ hardcoded path tạm thời
	if err != nil {
		log.Println("GetBlock: failed to open DB:", err)
		return nil, err
	}
	defer db.Close()

	blk, err := db.GetBlock(req.Hash)
	if err != nil {
		return nil, err
	}

	return &pb.BlockResponse{Block: convertBlockToPb(blk)}, nil
}

func convertBlockToPb(block *blockchain.Block) *pb.Block {
	var txs []*pb.Transaction
	for _, tx := range block.Transactions {
		txs = append(txs, &pb.Transaction{
			Sender:    tx.Sender,
			Receiver:  tx.Receiver,
			Amount:    tx.Amount,
			Timestamp: tx.Timestamp,
			Signature: tx.Signature,
		})
	}

	return &pb.Block{
		Transactions:     txs,
		MerkleRoot:       block.MerkleRoot,
		PrevBlockHash:    block.PrevBlockHash,
		CurrentBlockHash: block.CurrentBlockHash,
	}
}
