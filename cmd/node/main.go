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
	// log.Println("🔁 Sync complete. Now acting as Follower.")

	// ✅ Tạo server instance (quan trọng để giữ priority & state)
	server := p2p.NewNodeServer(port, dbPath, nodeID, db, &state)

	// 🚀 Start gRPC server
	go server.StartGRPC()

	// Bắt đầu monitor leader
	p2p.MonitorLeader(server, peers)

	// 🗳️ Bắt đầu bầu cử sau khi server sẵn sàng
	time.Sleep(2 * time.Second)
	p2p.StartElection(server, peers)

	select {} // giữ cho chương trình chạy hoài
}
