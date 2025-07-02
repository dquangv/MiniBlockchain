package p2p

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang-chain/pkg/p2p/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func generatePriority() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000)
}

func StartElection(server *NodeServer, peers []string) {
	myID := server.NodeID
	myPriority := server.Priority

	// Reset lại priorities và leader
	server.Priorities = map[string]int{
		myID: myPriority,
	}
	CurrentLeader = ""
	log.Printf("🎲 My priority is %d", myPriority)

	aliveNodes := map[string]bool{myID: true} // node hiện tại là alive

	// Gửi priority cho các node khác
	for _, peer := range peers {
		if strings.Contains(peer, myID) {
			continue
		}

		conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("❌ Cannot connect to peer %s", peer)
			continue
		}
		defer conn.Close()

		client := pb.NewNodeServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err = client.ExchangePriority(ctx, &pb.PriorityRequest{
			NodeId:   myID,
			Priority: int32(myPriority),
		})
		if err != nil {
			log.Printf("⚠️ Failed to send priority to %s", peer)
			continue
		}

		peerID := extractNodeID(peer)
		aliveNodes[peerID] = true
		log.Printf("✅ Got response from %s", peer)
	}

	// Đợi các node khác phản hồi xong
	time.Sleep(1 * time.Second)

	// Chọn leader từ các node còn sống
	highest := myPriority
	leader := myID

	for id, p := range server.Priorities {
		if !aliveNodes[id] {
			continue // Loại node chết ra
		}
		if p > highest || (p == highest && id > leader) {
			highest = p
			leader = id
		}
	}

	server.LeaderID = leader
	CurrentLeader = peerAddressByID(leader, peers)

	if leader == myID {
		*server.State = StateLeader
		log.Println("👑 Elected as leader after full priority comparison")
		StartLeaderLoop(server.DB, peers)
	} else {
		*server.State = StateFollower
		log.Printf("🤖 I am a follower. Leader is %s", leader)
	}
}

func peersFromEnv() []string {
	raw := os.Getenv("PEERS")
	return strings.Split(raw, ",")
}

func peerAddressByID(nodeID string, peers []string) string {
	for _, p := range peers {
		if strings.Contains(p, nodeID) {
			return p
		}
	}
	return ""
}

func extractNodeID(addr string) string {
	if idx := strings.Index(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}
