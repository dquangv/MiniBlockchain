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
	server.Priorities[myID] = myPriority
	log.Printf("🎲 My priority is %d", myPriority)

	// Gửi priority cho tất cả peer
	for _, peer := range peers {
		if strings.Contains(peer, myID) {
			continue
		}

		conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("❌ Failed to connect to %s: %v", peer, err)
			continue
		}
		defer conn.Close()

		client := pb.NewNodeServiceClient(conn)
		_, err = client.ExchangePriority(context.Background(), &pb.PriorityRequest{
			NodeId:   myID,
			Priority: int32(myPriority),
		})
		if err != nil {
			log.Printf("⚠️ Error exchanging with %s: %v", peer, err)
			continue
		}
	}

	// 🕒 Đợi tất cả node gửi xong (ví dụ 2 giây)
	time.Sleep(2 * time.Second)

	// 🧠 Lúc này tất cả priority đã được lưu → mới bắt đầu chọn leader
	highest := myPriority
	leader := myID

	for id, p := range server.Priorities {
		if p > highest || (p == highest && id > leader) {
			highest = p
			leader = id
		}
	}

	server.LeaderID = leader
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
