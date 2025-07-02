package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/p2p"
	"golang-chain/pkg/storage"
)

func main() {
	fmt.Println("Hello from validator node!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = fmt.Sprintf("node-%s", port)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/" + nodeID
	}

	var peers []string
	if raw := os.Getenv("PEERS"); raw != "" {
		peers = strings.Split(raw, ",")
	}

	db, err := storage.NewDB(dbPath)
	if err != nil {
		log.Fatalln("❌ Failed to open DB:", err)
	}
	defer db.Close()

	if _, err := db.GetLatestBlock(); err != nil {
		log.Println("📦 No blocks found. Creating genesis block...")
		genesis := blockchain.CreateGenesisBlock()
		if err := db.SaveBlock(genesis); err != nil {
			log.Fatalln("❌ Failed to create genesis block:", err)
		}
		log.Println("✅ Genesis block created.")
	}

	log.Println("🔄 This node is Syncing...")
	state := p2p.StateSyncing
	if len(peers) > 0 {
		p2p.SyncFromPeerByHeight(peers[0], db)
		log.Println("🎉 Sync completed successfully.")
	} else {
		log.Println("⚠️ No peers found to sync from.")
	}
	state = p2p.StateFollower

	// ✅ Tạo NodeServer instance
	server := p2p.NewNodeServer(port, dbPath, nodeID, db, &state)

	// 🚀 Khởi động gRPC server
	go server.StartGRPC()

	// ⏱️ Đợi gRPC ổn định
	time.Sleep(2 * time.Second)

	// 🧠 Kiểm tra xem có leader nào đang online không
	if has, leaderAddr := p2p.HasLeader(peers); has {
		log.Printf("✅ Detected existing leader: %s. Joining as follower.", leaderAddr)
		server.LeaderID = p2p.ExtractNodeID(leaderAddr)
		p2p.CurrentLeader = leaderAddr
		*server.State = p2p.StateFollower
	} else {
		log.Println("🗳️ No leader detected. Starting election...")
		p2p.StartElection(server, peers)
	}

	// 🕵️ Theo dõi leader
	p2p.MonitorLeader(server, peers)

	select {} // giữ chương trình chạy hoài
}
