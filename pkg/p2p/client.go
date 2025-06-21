package p2p

import (
	"context"
	"golang-chain/pkg/p2p/pb"
	"golang-chain/pkg/storage"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendBlockForVote(peers []string, block *pb.Block) []*pb.VoteResponse {
	var votes []*pb.VoteResponse

	for _, addr := range peers {
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

func SyncBlocksFromPeer(peer string, localLatestHash string, db *storage.DB) {
	log.Println("🛠 SyncBlocksFromPeer CALLED: peer =", peer, ", localHash =", localLatestHash)

	conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	remoteBlock := latestResp.Block
	current := remoteBlock

	for current != nil && current.CurrentBlockHash != localLatestHash {
		blk := convertPbBlock(current)
		if err := db.SaveBlock(blk); err != nil {
			log.Println("Failed to save block:", err)
			return
		}

		resp, err := client.GetBlock(context.Background(), &pb.BlockRequest{Hash: current.PrevBlockHash})
		if err != nil {
			log.Println("Failed to get previous block:", err)
			break
		}
		current = resp.Block
	}

	log.Println("✅ Synced all missing blocks from peer!")
}

func SyncFromPeerByHeight(peer string, db *storage.DB) {
	log.Println("🌐 Syncing from peer:", peer)

	conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("❌ Connect error:", err)
		return
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)

	// 1. Lấy latest block của follower
	local, err := db.GetLatestBlock()
	start := int64(0)
	if err == nil && local != nil {
		start = local.Height + 1
		log.Printf("🔍 Local height: %d — starting sync from %d", local.Height, start)
	} else {
		log.Println("📭 No local block found — full sync from height 0")
	}

	// 2. Lấy latest block từ leader
	latestResp, err := client.GetLatestBlock(context.Background(), &pb.Empty{})
	if err != nil || latestResp.Block == nil {
		log.Println("❌ Cannot fetch latest block from peer")
		return
	}
	leaderHeight := latestResp.Block.Height
	log.Printf("🌐 Peer has block height: %d", leaderHeight)

	// 3. Lặp và sync từng block còn thiếu
	for h := start; h <= leaderHeight; h++ {
		resp, err := client.GetBlockByHeight(context.Background(), &pb.HeightRequest{Height: h})
		if err != nil || resp.Block == nil {
			log.Printf("❌ Failed to get block at height %d", h)
			break
		}

		block := convertPbBlock(resp.Block)
		err = db.SaveBlock(block)
		if err != nil {
			log.Printf("❌ Failed to save block at height %d: %v", h, err)
			break
		}

		log.Printf("✅ Synced block at height %d (hash: %s)", h, block.CurrentBlockHash)
	}
	log.Println("🎉 Sync completed successfully.")
}
