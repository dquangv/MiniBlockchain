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

	log.Printf("🗳️ No leader detected. Starting election...")
	log.Printf("🎲 My priority is %d", myPriority)

	// 🧠 Xác định các node online
	alivePeers := getAlivePeers(myID, peers)
	expected := len(alivePeers)
	log.Printf("📡 Alive peers: %v", alivePeers)

	// 🧹 Không reset hoàn toàn map để không mất dữ liệu cũ
	server.Mutex.Lock()
	if server.Priorities == nil {
		server.Priorities = make(map[string]int)
	}
	server.Priorities[myID] = myPriority
	server.Mutex.Unlock()

	// 📤 Gửi priority cho các node còn sống
	for _, peer := range peers {
		if strings.Contains(peer, myID) {
			continue
		}
		go func(peer string) {
			conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("⚠️ Failed to connect to %s: %v", peer, err)
				return
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
				log.Printf("⚠️ Failed to send priority to %s: %v", peer, err)
			}
		}(peer)
	}

	// ⏳ Chờ đến khi đủ priority từ các node còn sống
	log.Printf("🕒 Waiting for %d priorities...", expected)
	start := time.Now()
	for {
		server.Mutex.Lock()
		received := 0
		for _, peer := range alivePeers {
			if _, ok := server.Priorities[peer]; ok {
				received++
			}
		}
		server.Mutex.Unlock()

		if received >= expected {
			break
		}

		if time.Since(start) > 5*time.Second {
			log.Printf("🚫 Timeout. Collected %d/%d priorities. Skipping election.", received, expected)
			return
		}

		time.Sleep(500 * time.Millisecond)
	}

	// ✅ Tiến hành bầu chọn leader
	server.Mutex.Lock()

	// 🧼 Cleanup: chỉ giữ priority của các alivePeers
	filtered := make(map[string]int)
	for _, peerID := range alivePeers {
		if val, ok := server.Priorities[peerID]; ok {
			filtered[peerID] = val
		}
	}
	server.Priorities = filtered

	defer server.Mutex.Unlock()

	highest := myPriority
	leader := myID

	for id, p := range server.Priorities {
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

func ExtractNodeID(addr string) string {
	if idx := strings.Index(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

func getAlivePeers(myID string, peers []string) []string {
	alive := []string{myID}
	for _, peer := range peers {
		if strings.Contains(peer, myID) {
			continue
		}
		conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		defer conn.Close()

		client := pb.NewNodeServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err = client.Ping(ctx, &pb.Empty{})
		if err == nil {
			alive = append(alive, ExtractNodeID(peer))
		}
	}
	return alive
}

func HasLeader(peers []string) (bool, string) {
	for _, addr := range peers {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		defer conn.Close()

		client := pb.NewNodeServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		resp, err := client.Ping(ctx, &pb.Empty{})
		if err == nil && resp.Message == string(StateLeader) {
			return true, addr
		}
	}
	return false, ""
}
