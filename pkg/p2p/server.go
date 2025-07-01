package p2p

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/p2p/pb"
	"golang-chain/pkg/storage"
	"golang-chain/pkg/wallet"

	"golang-chain/pkg/consensus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type NodeServer struct {
	pb.UnimplementedNodeServiceServer
	DBPath     string
	NodeID     string
	DB         *storage.DB
	State      *NodeState
	Priority   int
	LeaderID   string
	Priorities map[string]int // 🆕 lưu priority của các node khác
}

func (s *NodeServer) SendTransaction(ctx context.Context, tx *pb.Transaction) (*pb.TxResponse, error) {
	if *s.State != StateLeader {
		return &pb.TxResponse{
			Status:  "error",
			Message: "Only the leader can accept transactions",
		}, nil
	}

	fromName := wallet.ResolveSenderName(tx.Sender)
	toName := wallet.ResolveSenderName(tx.Receiver)

	log.Printf("Received transaction from %s to %s (%.2f coins)", fromName, toName, tx.Amount)

	// 🔍 Kiểm tra số dư trước
	balance, err := s.DB.GetBalance(fromName)
	if err != nil {
		return &pb.TxResponse{
			Status:  "error",
			Message: fmt.Sprintf("❌ Failed to get balance: %v", err),
		}, nil
	}

	amount := big.NewFloat(tx.Amount)
	if balance.Cmp(amount) < 0 {
		return &pb.TxResponse{
			Status:  "fail",
			Message: fmt.Sprintf("❌ Insufficient balance. You have %s, trying to send %.2f", balance.Text('f', 2), tx.Amount),
		}, nil
	}

	t := &blockchain.Transaction{
		Sender:    tx.Sender,
		Receiver:  tx.Receiver,
		Amount:    tx.Amount,
		Timestamp: tx.Timestamp,
		Signature: tx.Signature,
	}

	blockchain.AddPendingTx(t)
	log.Printf("📥 Transaction added to pending pool.")

	return &pb.TxResponse{
		Status:  "ok",
		Message: "The transaction has been sent and is pending.",
	}, nil
}

func (s *NodeServer) Ping(ctx context.Context, e *pb.Empty) (*pb.TxResponse, error) {
	return &pb.TxResponse{
		Status:  "pong",
		Message: string(*s.State),
	}, nil
}

func StartGRPCServer(port, dbPath, nodeID string, db *storage.DB, state *NodeState) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	server := &NodeServer{
		DBPath:     dbPath,
		NodeID:     nodeID,
		DB:         db,
		State:      state,
		Priorities: make(map[string]int),
		Priority:   generatePriority(), // 🆕 Khởi tạo priority ngẫu nhiên
		LeaderID:   nodeID,             // ban đầu assume mình là leader
	}

	log.Printf("🎲 My priority is %d", server.Priority)

	grpcServer := grpc.NewServer()
	pb.RegisterNodeServiceServer(grpcServer, server)

	go func() {
		log.Printf("gRPC server listening on port %s", port)
		log.Println("🔎 Current state:", *state)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 🗳️ Delay election để gRPC sẵn sàng
	time.Sleep(2 * time.Second)
	StartElection(server, peersFromEnv())
}

// Follower xử lý block do Leader đề xuất để vote
func (s *NodeServer) ProposeBlock(ctx context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	if *s.State != StateFollower {
		log.Println("⚠️ Vote rejected: I am not a follower.")
		return &pb.VoteResponse{
			NodeId:   s.NodeID,
			Approved: false,
		}, nil
	}

	block := req.Block
	log.Printf("[Follower] Received proposed block: %s", block.CurrentBlockHash)

	latestBlock, err := s.DB.GetLatestBlock()
	if err != nil {
		log.Println("⚠️ No latest block found, assuming fresh node")
		if block.Height == 0 {
			log.Println("✅ Accepting genesis block proposal.")
			return &pb.VoteResponse{
				NodeId:   s.NodeID,
				Approved: true,
			}, nil
		}
		log.Println("❌ Rejected: Genesis block must have height 0")
		return &pb.VoteResponse{
			NodeId:   s.NodeID,
			Approved: false,
		}, nil
	}

	newBlock := convertPbBlock(block)
	isValid := consensus.VerifyBlock(newBlock, latestBlock)

	return &pb.VoteResponse{
		NodeId:   s.NodeID,
		Approved: isValid,
	}, nil
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
		return nil, err
	}
	return &pb.BlockResponse{Block: ConvertBlockToPb(block)}, nil
}

func DetectLeader(peers []string) string {
	for _, peer := range peers {
		conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
		if err != nil {
			log.Println("❌ Failed to connect to", peer)
			continue
		}
		defer conn.Close()

		client := pb.NewNodeServiceClient(conn)
		resp, err := client.Ping(context.Background(), &pb.Empty{})
		if err != nil {
			log.Println("❌ Ping failed to", peer)
			continue
		}

		if resp.Message == string(StateLeader) {
			log.Println("👑 Leader detected at", peer)
			return peer
		}
	}
	return ""
}

func (s *NodeServer) GetBalance(ctx context.Context, req *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	bal, err := s.DB.GetBalance(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get balance: %v", err)
	}
	return &pb.BalanceResponse{
		Balance: bal.Text('f', 2),
	}, nil
}

var priorityMap = make(map[string]int)

func (s *NodeServer) ExchangePriority(ctx context.Context, req *pb.PriorityRequest) (*pb.PriorityResponse, error) {
	log.Printf("🤝 Received priority %d from %s", req.Priority, req.NodeId)
	s.Priorities[req.NodeId] = int(req.Priority)
	return &pb.PriorityResponse{}, nil
}

func NewNodeServer(port, dbPath, nodeID string, db *storage.DB, state *NodeState) *NodeServer {
	return &NodeServer{
		DBPath:     dbPath,
		NodeID:     nodeID,
		DB:         db,
		State:      state,
		Priority:   generatePriority(),
		LeaderID:   nodeID, // mặc định assume mình
		Priorities: make(map[string]int),
	}
}

func (s *NodeServer) StartGRPC() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterNodeServiceServer(grpcServer, s)

	log.Println("🚀 gRPC server started on port", os.Getenv("PORT"))
	// log.Println("🔎 Current state:", *s.State)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
