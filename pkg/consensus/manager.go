package consensus

import (
	"context"
	"golang-chain/pkg/p2p/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConsensusManager struct {
	Peers []string
}

// Gửi block đi để vote
func (cm *ConsensusManager) ProposeBlock(block *pb.Block) []*pb.VoteResponse {
	var votes []*pb.VoteResponse

	for _, addr := range cm.Peers {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("Can't connect to peer", addr)
			continue
		}
		defer conn.Close()

		client := pb.NewNodeServiceClient(conn)
		vote, err := client.ProposeBlock(context.Background(), &pb.VoteRequest{Block: block})
		if err != nil {
			log.Println("Peer failed to vote:", err)
			continue
		}

		log.Printf("[Leader] Peer %s voted %v", vote.NodeId, vote.Approved)
		votes = append(votes, vote)
	}

	return votes
}

// Gửi block đi để commit sau khi có đủ vote
func (cm *ConsensusManager) BroadcastCommit(block *pb.Block) {
	for _, addr := range cm.Peers {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("Can't connect to peer", addr)
			continue
		}
		defer conn.Close()

		client := pb.NewNodeServiceClient(conn)
		_, err = client.CommitBlock(context.Background(), block)
		if err != nil {
			log.Println("Commit failed to", addr)
		} else {
			log.Println("Committed to", addr)
		}
	}
}
