package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/" + os.Getenv("NODE_ID")
	}

	isLeader := os.Getenv("IS_LEADER") == "true"
	peerList := os.Getenv("PEERS")
	var peers []string
	if peerList != "" {
		peers = strings.Split(peerList, ",")
	}

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = fmt.Sprintf("follower-%s", port)
	}

	// Mở DB
	db, err := storage.NewDB(dbPath)
	if err != nil {
		log.Fatalln("❌ Failed to open DB:", err)
	}
	defer db.Close()

	state := p2p.StateSyncing
	if isLeader {
		// Tạo Genesis block nếu chưa có block nào
		_, err = db.GetLatestBlock()
		if err != nil {
			log.Println("📦 No blocks found. Creating genesis block...")
			genesis := blockchain.CreateGenesisBlock()
			if err := db.SaveBlock(genesis); err != nil {
				log.Fatalln("❌ Failed to create genesis block:", err)
			}
			log.Println("✅ Genesis block created.")
		}

		log.Println("🧠 This node is the Leader.")
		state = p2p.StateLeader
		go p2p.StartLeaderLoop(db, peers)
	} else {
		log.Println("🔄 This node is Syncing...")
		if len(peers) > 0 {
			p2p.SyncFromPeerByHeight(peers[0], db)
		} else {
			log.Println("⚠️ No peers found to sync from.")
		}

		state = p2p.StateFollower
		log.Println("🔁 Sync complete. Now acting as Follower.")
	}

	p2p.StartGRPCServer(port, dbPath, nodeID, db, &state)
}
