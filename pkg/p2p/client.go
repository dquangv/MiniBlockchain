package p2p

import (
	"context"
	"golang-chain/pkg/p2p/pb"
	"golang-chain/pkg/storage"
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

func SyncBlocksFromPeer(peer string, localLatestHash string) {
	log.Println("🛠 SyncBlocksFromPeer CALLED: peer =", peer, ", localHash =", localLatestHash)
	conn, err := grpc.Dial(peer, grpc.WithInsecure())
	if err != nil {
		log.Println("Connect error:", err)
		return
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)

	// Lấy block mới nhất từ peer
	latestResp, err := client.GetLatestBlock(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println("Failed to get latest block from peer:", err)
		return
	}

	db, _ := storage.NewDB("blockdata")
	defer db.Close()

	remoteBlock := latestResp.Block
	current := remoteBlock

	for current != nil && current.CurrentBlockHash != localLatestHash {
		blk := convertPbBlock(current)
		err := db.SaveBlock(blk)
		if err != nil {
			log.Println("Failed to save block:", err)
			return
		}

		// Lấy block trước
		resp, err := client.GetBlock(context.Background(), &pb.BlockRequest{Hash: current.PrevBlockHash})
		if err != nil {
			break
		}
		current = resp.Block
	}

	log.Println("✅ Synced all missing blocks from peer!")
}
