package p2p

import (
	"context"
	"golang-chain/pkg/p2p/pb"
	"log"

	"google.golang.org/grpc"
)

func SendBlockForVote(peers []string, block *pb.Block) []*pb.VoteResponse {
	var votes []*pb.VoteResponse

	for _, addr := range peers {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
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

		log.Printf("[Leader] Peer %s voted %v\n", vote.NodeId, vote.Approved)
		votes = append(votes, vote)
	}

	return votes
}

func BroadcastCommit(peers []string, block *pb.Block) {
	for _, addr := range peers {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
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
