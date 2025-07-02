# 🔗 Golang Blockchain Mini System

A lightweight blockchain simulation built with Golang. This system demonstrates coin transfers between two users (**Alice** and **Bob**) with a leader-follower consensus mechanism across 3 validator nodes—all containerized with Docker.

---

## 🧠 Architecture & Technologies

### 🛠 Tech Stack:
- **Golang**: Main programming language.
- **ECDSA** (`crypto/ecdsa`): For digital signatures.
- **LevelDB** (`github.com/syndtr/goleveldb/leveldb`): Embedded key-value store for blocks.
- **gRPC**: Communication between validator nodes.
- **Docker + Docker Compose**: Containerized environment for the full network.

### 🏗️ System Architecture:
- 3 validator nodes: node1, node2, node3 — each with a unique NODE_ID, communicating via gRPC.
- Leader is elected automatically based on randomly generated priority.
- Only the Leader receives transactions, creates blocks, and proposes them to other nodes for voting.
- A block is committed when ≥2 votes are received (Leader + 1 follower).
- When a node restarts, it automatically syncs missing blocks from peers to catch up.

---

## 🗂 Project Structure
```csharp
├── cmd/
│ ├── node/ # Main node startup
│ └── cli/ # CLI tools: wallet creation, send transaction, check status
├── pkg/
│ ├── blockchain/ # Block, Transaction, Merkle root logic
│ ├── consensus/ # Voting and commit verification
│ ├── p2p/ # gRPC communication
│ ├── storage/ # LevelDB database wrapper
│ └── wallet/ # ECDSA key pair and wallet logic
├── wallets/ # Alice & Bob wallet files (PEM encoded)
├── Dockerfile # Multi-stage Docker build
├── docker-compose.yml # Spin up the full validator network
```
## 🚀 Getting Started

### 1. Clone the repo and start containers
```bash
git clone https://github.com/dquangv/golang-chain.git
cd golang-chain
docker compose up --build
```
✨ Expected output
If everything runs correctly, you'll see logs like this (abbreviated):
```csharp
2025-07-02 17:32:46 node1  | Hello from validator node!
2025-07-02 17:32:46 node1  | 2025/07/02 10:32:46 🔄 This node is Syncing...
2025-07-02 17:32:46 node1  | 2025/07/02 10:32:46 🌐 Syncing from peer: node2:50051
2025-07-02 17:32:46 node1  | 2025/07/02 10:32:46 🔍 Local height: 14 — starting sync from 15
2025-07-02 17:32:48 node2  | ⏳ Waiting for node1:50051...
2025-07-02 17:32:48 node3  | ⏳ Waiting for node1:50051...
2025-07-02 17:32:49 node2  | ⏳ Waiting for node1:50051...
2025-07-02 17:32:49 node3  | ⏳ Waiting for node1:50051...
2025-07-02 17:32:50 node2  | ⏳ Waiting for node1:50051...
2025-07-02 17:32:50 node3  | ⏳ Waiting for node1:50051...
2025-07-02 17:32:50 node1  | 2025/07/02 10:32:50 ❌ Cannot fetch latest block from peer
2025-07-02 17:32:50 node1  | 2025/07/02 10:32:50 🎉 Sync completed successfully.
2025-07-02 17:32:50 node1  | 2025/07/02 10:32:50 🚀 gRPC server started on port 50051
2025-07-02 17:32:51 node2  | Hello from validator node!
2025-07-02 17:32:51 node2  | 2025/07/02 10:32:51 🔄 This node is Syncing...
2025-07-02 17:32:51 node2  | 2025/07/02 10:32:51 🌐 Syncing from peer: node1:50051
2025-07-02 17:32:51 node2  | 2025/07/02 10:32:51 🔍 Local height: 14 — starting sync from 15
2025-07-02 17:32:51 node2  | 2025/07/02 10:32:51 🌐 Peer has block height: 14
2025-07-02 17:32:51 node2  | 2025/07/02 10:32:51 🎉 Sync completed successfully.
2025-07-02 17:32:51 node2  | 2025/07/02 10:32:51 🎉 Sync completed successfully.
2025-07-02 17:32:51 node2  | 2025/07/02 10:32:51 🚀 gRPC server started on port 50051
2025-07-02 17:32:51 node3  | Hello from validator node!
2025-07-02 17:32:51 node3  | 2025/07/02 10:32:51 🔄 This node is Syncing...
2025-07-02 17:32:51 node3  | 2025/07/02 10:32:51 🌐 Syncing from peer: node1:50051
2025-07-02 17:32:51 node3  | 2025/07/02 10:32:51 🔍 Local height: 14 — starting sync from 15
2025-07-02 17:32:51 node3  | 2025/07/02 10:32:51 🌐 Peer has block height: 14
2025-07-02 17:32:51 node3  | 2025/07/02 10:32:51 🎉 Sync completed successfully.
2025-07-02 17:32:51 node3  | 2025/07/02 10:32:51 🎉 Sync completed successfully.
2025-07-02 17:32:51 node3  | 2025/07/02 10:32:51 🚀 gRPC server started on port 50051
2025-07-02 17:32:52 node1  | 2025/07/02 10:32:52 🗳️ No leader detected. Starting election...
2025-07-02 17:32:52 node1  | 2025/07/02 10:32:52 🗳️ No leader detected. Starting election...
2025-07-02 17:32:52 node1  | 2025/07/02 10:32:52 🎲 My priority is 528
2025-07-02 17:32:52 node1  | 2025/07/02 10:32:52 📡 Alive peers: [node1 node2 node3]
2025-07-02 17:32:52 node1  | 2025/07/02 10:32:52 🕒 Waiting for 3 priorities...
2025-07-02 17:32:52 node2  | 2025/07/02 10:32:52 🤝 Received priority 528 from node1
2025-07-02 17:32:52 node2  | 2025/07/02 10:32:52 📥 All priorities received so far: map[node1:528]
2025-07-02 17:32:52 node3  | 2025/07/02 10:32:52 🤝 Received priority 528 from node1
2025-07-02 17:32:52 node3  | 2025/07/02 10:32:52 📥 All priorities received so far: map[node1:528]
2025-07-02 17:32:53 node2  | 2025/07/02 10:32:53 🗳️ No leader detected. Starting election...
2025-07-02 17:32:53 node2  | 2025/07/02 10:32:53 🗳️ No leader detected. Starting election...
2025-07-02 17:32:53 node2  | 2025/07/02 10:32:53 🎲 My priority is 109
2025-07-02 17:32:53 node2  | 2025/07/02 10:32:53 📡 Alive peers: [node2 node1 node3]
2025-07-02 17:32:53 node2  | 2025/07/02 10:32:53 🕒 Waiting for 3 priorities...
2025-07-02 17:32:53 node3  | 2025/07/02 10:32:53 🤝 Received priority 109 from node2
2025-07-02 17:32:53 node3  | 2025/07/02 10:32:53 📥 All priorities received so far: map[node1:528 node2:109]
2025-07-02 17:32:53 node1  | 2025/07/02 10:32:53 🤝 Received priority 109 from node2
2025-07-02 17:32:53 node1  | 2025/07/02 10:32:53 📥 All priorities received so far: map[node1:528 node2:109]
2025-07-02 17:32:53 node3  | 2025/07/02 10:32:53 🗳️ No leader detected. Starting election...
2025-07-02 17:32:53 node3  | 2025/07/02 10:32:53 🗳️ No leader detected. Starting election...
2025-07-02 17:32:53 node3  | 2025/07/02 10:32:53 🎲 My priority is 34
2025-07-02 17:32:53 node3  | 2025/07/02 10:32:53 📡 Alive peers: [node3 node1 node2]
2025-07-02 17:32:53 node3  | 2025/07/02 10:32:53 🕒 Waiting for 3 priorities...
2025-07-02 17:32:53 node3  | 2025/07/02 10:32:53 🤖 I am a follower. Leader is node1
2025-07-02 17:32:53 node1  | 2025/07/02 10:32:53 🤝 Received priority 34 from node3
2025-07-02 17:32:53 node1  | 2025/07/02 10:32:53 📥 All priorities received so far: map[node1:528 node2:109 node3:34]
2025-07-02 17:32:53 node2  | 2025/07/02 10:32:53 🤝 Received priority 34 from node3
2025-07-02 17:32:53 node2  | 2025/07/02 10:32:53 📥 All priorities received so far: map[node1:528 node2:109 node3:34]
2025-07-02 17:32:53 node1  | 2025/07/02 10:32:53 👑 Elected as leader after full priority comparison
2025-07-02 17:32:53 node2  | 2025/07/02 10:32:53 🤖 I am a follower. Leader is node1
2025-07-02 17:32:58 node3  | 2025/07/02 10:32:58 ✅ Leader node1 still alive
2025-07-02 17:32:58 node1  | 2025/07/02 10:32:58 ⏳ Tick! Checking for pending transactions...
2025-07-02 17:32:58 node1  | 2025/07/02 10:32:58 🔍 No pending transactions. Skipping block creation.
2025-07-02 17:32:58 node2  | 2025/07/02 10:32:58 ✅ Leader node1 still alive
2025-07-02 17:33:03 node3  | 2025/07/02 10:33:03 ✅ Leader node1 still alive
2025-07-02 17:33:03 node1  | 2025/07/02 10:33:03 ⏳ Tick! Checking for pending transactions...
2025-07-02 17:33:03 node1  | 2025/07/02 10:33:03 🔍 No pending transactions. Skipping block creation.
```

### 2. Available CLI Tools
🧰 Create wallet:
```bash
$ docker exec -it node1 ./create_wallet --name Alice
$ docker exec -it node1 ./create_wallet --name Bob
```
```csharp
✅ The wallet has been created and saved at:  wallets/Alice_wallet.json
✅ The wallet has been created and saved at:  wallets/Bob_wallet.json
```
💸 Send transaction:
```bash
$ docker exec -it node1 ./send_tx --from Alice --to Bob --amount 10 --node localhost:50051
```
```csharp
2025-06-21 13:15:17 node1  | 2025/06/21 06:15:17 ⏳ Tick! Checking for pending transactions...
2025-06-21 13:15:17 node1  | 2025/06/21 06:15:17 🔍 No pending transactions. Skipping block creation.
2025-06-21 13:15:20 node1  | 2025/06/21 06:15:20 Received transaction from Alice to Bob (10.00 coins)
2025-06-21 13:15:20 node1  | 2025/06/21 06:15:20 📥 Transaction added to pending pool.
2025-06-21 13:15:22 node1  | 2025/06/21 06:15:22 ⏳ Tick! Checking for pending transactions...
2025-06-21 13:15:22 node1  | 2025/06/21 06:15:22 📨 Found 1 pending transaction(s). Creating new block...
2025-06-21 13:15:22 node1  | 2025/06/21 06:15:22 [Leader] Peer node2 voted true
2025-06-21 13:15:22 node1  | 2025/06/21 06:15:22 [Leader] Peer node3 voted true
2025-06-21 13:15:22 node3  | 2025/06/21 06:15:22 [Follower] Received proposed block: 960020616b25f25fbeb4055a2b1c48fcfbf89fbb23e334f05f73b828fdb56062
2025-06-21 13:15:22 node2  | 2025/06/21 06:15:22 [Follower] Received proposed block: 960020616b25f25fbeb4055a2b1c48fcfbf89fbb23e334f05f73b828fdb56062
2025-06-21 13:15:22 node2  | 2025/06/21 06:15:22 [Follower] Block committed: 960020616b25f25fbeb4055a2b1c48fcfbf89fbb23e334f05f73b828fdb56062
2025-06-21 13:15:22 node1  | 2025/06/21 06:15:22 Committed to node2:50051
2025-06-21 13:15:22 node1  | 2025/06/21 06:15:22 Committed to node3:50051
2025-06-21 13:15:22 node3  | 2025/06/21 06:15:22 [Follower] Block committed: 960020616b25f25fbeb4055a2b1c48fcfbf89fbb23e334f05f73b828fdb56062
2025-06-21 13:15:22 node1  | 2025/06/21 06:15:22 ✅ Committed block at height 1 with 1 txs
```
Checking for pending transactions to create a new block every 5 seconds

📊 View status block:
```bash
$ docker exec -it node1 ./status --node localhost:50051
```
```csharp
📦 The latest block:
👉 Height:        1
👉 Hash:          960020616b25f25fbeb4055a2b1c48fcfbf89fbb23e334f05f73b828fdb56062
👉 Prev Hash:     b50ad2d4bd47d6278d2b9387db537b221107d5f80f27954118a057d1b97af412
👉 Tx count:      1
```

📈 Check wallet balance:
```bash
$ docker exec -it node1 ./balance --name Alice
$ docker exec -it node1 ./balance --name Bob
```
```csharp
💰 Balance of Alice: -10.00
💰 Balance of Bob:   10.00
```

### 🔐 Transactions & Signing
Each transaction contains:
- Sender: Public Key (PEM encoded)
- Receiver: Wallet address (hex string)
- Amount, Timestamp, and Signature

Transactions are:
- Signed by sender's private key
- Verified by validator using public key before accepting into block

### 🔄 Leader Election & Fault Tolerance
- When no Leader is detected or the current Leader becomes unresponsive, the system automatically triggers a re-election.
- Each node generates a random priority and broadcasts it to currently alive peers only.
- Election only starts after receiving enough priorities from reachable nodes.
- The node with the highest priority (or lexicographically greater NODE_ID if tied) becomes the Leader.
- When a former Leader rejoins, it automatically becomes a Follower if a valid Leader already exists.

### 💾 Blockchain Storage
- Each node stores blockchain data locally using LevelDB in ./blockdata/<node-id>.
- The genesis block is only created if the database is empty.
- On startup, if the chain is outdated, the node auto-syncs from peers.

### ⚙️ Configuration
| ENV Variable | Description                                    |
| ------------ | ---------------------------------------------- |
| `NODE_ID`    | Unique identifier for the node (e.g., `node1`) |
| `PORT`       | gRPC listening port                            |
| `PEERS`      | Comma-separated list of peer addresses         |
| `DB_PATH`    | Directory for storing blockchain data          |

### 📌 Key Behavior
- Leader is dynamically elected — no need for IS_LEADER flag.
- Election only runs if no valid Leader exists.
- Only the Leader can accept new transactions.
- Re-election is triggered when the Leader goes down.
- Nodes recover and sync state automatically after downtime.

### 📚 Learnings & Key Concepts
- Implementing a basic Leader-Follower consensus system.
- Using ECDSA for digital signatures and wallet generation.
- Building and verifying Merkle roots for block data integrity.
- Handling inter-node communication via gRPC.
- Orchestrating multiple blockchain nodes with Docker Compose.
