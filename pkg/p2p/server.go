package p2p

import (
	"context"
	"fmt"
	"log"
	"net"

	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/p2p/pb"
	"golang-chain/pkg/storage"

	"golang-chain/pkg/consensus"

	"google.golang.org/grpc"
)

type NodeServer struct {
	pb.UnimplementedNodeServiceServer
	DBPath string
	NodeID string
	DB     *storage.DB
	State  NodeState
}

func (s *NodeServer) SendTransaction(ctx context.Context, tx *pb.Transaction) (*pb.TxResponse, error) {
	if s.State != StateLeader {
		return &pb.TxResponse{
			Status:  "error",
			Message: "Only the leader can accept transactions",
		}, nil
	}

	log.Printf("Received transaction from %s to %s (%.2f)", tx.Sender, tx.Receiver, tx.Amount)

	// Convert pb.Transaction → blockchain.Transaction
	t := &blockchain.Transaction{
		Sender:    tx.Sender,
		Receiver:  tx.Receiver,
		Amount:    tx.Amount,
		Timestamp: tx.Timestamp,
		Signature: tx.Signature,
	}

	blockchain.AddPendingTx(t) // 🆕 Gửi vào hàng chờ

	return &pb.TxResponse{
		Status:  "ok",
		Message: "tx received and pending",
	}, nil
}

func (s *NodeServer) Ping(ctx context.Context, e *pb.Empty) (*pb.TxResponse, error) {
	return &pb.TxResponse{
		Status:  "pong",
		Message: "I'm alive",
	}, nil
}

func StartGRPCServer(port, dbPath, nodeID string, db *storage.DB, state NodeState) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	server := &NodeServer{
		DBPath: dbPath,
		NodeID: nodeID,
		DB:     db, // 🆕 xài lại db đã mở
		State:  state,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNodeServiceServer(grpcServer, server)

	fmt.Println("gRPC server listening on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Follower xử lý block do Leader đề xuất để vote
func (s *NodeServer) ProposeBlock(ctx context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	block := req.Block
	log.Printf("[Follower] Received proposed block: %s", block.CurrentBlockHash)

	latestBlock, err := s.DB.GetLatestBlock()
	if err != nil {
		log.Println("No latest block:", err)
		latestBlock = nil
	}

	newBlock := convertPbBlock(block)
	isValid := consensus.VerifyBlock(newBlock, latestBlock)

	vote := &pb.VoteResponse{
		NodeId:   s.NodeID,
		Approved: isValid,
	}

	return vote, nil
}

func convertPbBlock(pbBlock *pb.Block) *blockchain.Block {
	var txs []*blockchain.Transaction
	for _, tx := range pbBlock.Transactions {
		txs = append(txs, &blockchain.Transaction{
			Sender:    append([]byte(nil), tx.Sender...),
			Receiver:  append([]byte(nil), tx.Receiver...),
			Amount:    tx.Amount,
			Timestamp: tx.Timestamp,
			Signature: append([]byte(nil), tx.Signature...),
		})
	}

	block := &blockchain.Block{
		Transactions:     txs,
		PrevBlockHash:    pbBlock.PrevBlockHash,
		CurrentBlockHash: pbBlock.CurrentBlockHash,
		Height:           pbBlock.Height,
	}

	block.MerkleRoot = blockchain.CalculateMerkleRoot(txs)

	return block
}

func (s *NodeServer) CommitBlock(ctx context.Context, pbBlock *pb.Block) (*pb.TxResponse, error) {
	block := convertPbBlock(pbBlock)

	err := s.DB.SaveBlock(block)
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
	latestBlock, err := s.DB.GetLatestBlock()
	if err != nil {
		return nil, err
	}

	return &pb.BlockResponse{Block: ConvertBlockToPb(latestBlock)}, nil
}

func (s *NodeServer) GetBlock(ctx context.Context, req *pb.BlockRequest) (*pb.BlockResponse, error) {
	blk, err := s.DB.GetBlock([]byte(req.Hash))
	if err != nil {
		return nil, err
	}

	return &pb.BlockResponse{Block: ConvertBlockToPb(blk)}, nil
}

func ConvertBlockToPb(block *blockchain.Block) *pb.Block {
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
		Height:           block.Height,
	}
}

func (s *NodeServer) GetBlockByHeight(ctx context.Context, req *pb.HeightRequest) (*pb.BlockResponse, error) {
	block, err := s.DB.GetBlockByHeight(req.Height)
	if err != nil {
		log.Println("❌ GetBlockByHeight error:", err)
		return nil, err
	}
	return &pb.BlockResponse{Block: ConvertBlockToPb(block)}, nil
}
