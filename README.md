# 🔗 Golang Blockchain Mini System

A lightweight blockchain simulation built with Golang. This system demonstrates coin transfers between two users (**Alice** and **Bob**) with a leader-follower consensus mechanism across 3 validator nodes—all containerized with Docker.

---

## 🧠 Architecture & Technologies

### ✅ Tech Stack:
- **Golang**: Main programming language.
- **ECDSA** (`crypto/ecdsa`): For digital signatures.
- **LevelDB** (`github.com/syndtr/goleveldb/leveldb`): Embedded key-value store for blocks.
- **gRPC**: Communication between validator nodes.
- **Docker + Docker Compose**: Containerized environment for the full network.

### 🏗️ System Architecture:
- **3 validator nodes**:
  - `node1` = **Leader**
  - `node2`, `node3` = **Followers**
- **Leader-Follower consensus**:
  - Leader receives transactions, creates blocks, and proposes them.
  - Followers validate the block and respond with a vote.
  - A block is committed only if **at least 2 votes** are received.
- **Recovery capability**:
  - Any node that restarts will automatically **synchronize missing blocks** from peers to catch up with the network.

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
2025-06-21 13:13:07 node3  | Hello from validator node!
2025-06-21 13:13:07 node3  | 2025/06/21 06:13:07 🔄 This node is Syncing...
2025-06-21 13:13:07 node3  | 2025/06/21 06:13:07 🌐 Syncing from peer: node1:50051
2025-06-21 13:13:07 node3  | 2025/06/21 06:13:07 📭 No local block found — full sync from height 0
2025-06-21 13:13:07 node3  | 2025/06/21 06:13:07 🌐 Peer has block height: 0
2025-06-21 13:13:07 node3  | 2025/06/21 06:13:07 ✅ Synced block at height 0 (hash: b50ad2d4bd47d6278d2b9387db537b221107d5f80f27954118a057d1b97af412)
2025-06-21 13:13:07 node3  | 2025/06/21 06:13:07 🎉 Sync completed successfully.
2025-06-21 13:13:07 node3  | 2025/06/21 06:13:07 🔁 Sync complete. Now acting as Follower.
2025-06-21 13:13:07 node3  | gRPC server listening on port 50051
2025-06-21 13:13:07 node3  | 🔎 Current state: Follower
2025-06-21 13:13:07 node1  | Hello from validator node!
2025-06-21 13:13:07 node1  | 2025/06/21 06:13:07 📦 No blocks found. Creating genesis block...
2025-06-21 13:13:07 node1  | 2025/06/21 06:13:07 ✅ Genesis block created.
2025-06-21 13:13:07 node1  | 2025/06/21 06:13:07 🧠 This node is the Leader.
2025-06-21 13:13:07 node1  | gRPC server listening on port 50051
2025-06-21 13:13:07 node1  | 🔎 Current state: Leader
2025-06-21 13:13:07 node2  | Hello from validator node!
2025-06-21 13:13:07 node2  | 2025/06/21 06:13:07 🔄 This node is Syncing...
2025-06-21 13:13:07 node2  | 2025/06/21 06:13:07 🌐 Syncing from peer: node1:50051
2025-06-21 13:13:07 node2  | 2025/06/21 06:13:07 📭 No local block found — full sync from height 0
2025-06-21 13:13:07 node2  | 2025/06/21 06:13:07 🌐 Peer has block height: 0
2025-06-21 13:13:07 node2  | 2025/06/21 06:13:07 ✅ Synced block at height 0 (hash: b50ad2d4bd47d6278d2b9387db537b221107d5f80f27954118a057d1b97af412)
2025-06-21 13:13:07 node2  | 2025/06/21 06:13:07 🎉 Sync completed successfully.
2025-06-21 13:13:07 node2  | 2025/06/21 06:13:07 🔁 Sync complete. Now acting as Follower.
2025-06-21 13:13:07 node2  | gRPC server listening on port 50051
2025-06-21 13:13:07 node2  | 🔎 Current state: Follower
2025-06-21 13:13:12 node1  | 2025/06/21 06:13:12 ⏳ Tick! Checking for pending transactions...
2025-06-21 13:13:12 node1  | 2025/06/21 06:13:12 🔍 No pending transactions. Skipping block creation.
2025-06-21 13:13:17 node1  | 2025/06/21 06:13:17 ⏳ Tick! Checking for pending transactions...
2025-06-21 13:13:17 node1  | 2025/06/21 06:13:17 🔍 No pending transactions. Skipping block creation.
```

### 2. Available CLI Tools
✅ Create wallet:
```bash
$ docker exec -it node1 ./create_wallet --name Alice
$ docker exec -it node1 ./create_wallet --name Bob
```
```csharp
✅ The wallet has been created and saved at:  wallets/Alice_wallet.json
✅ The wallet has been created and saved at:  wallets/Bob_wallet.json
```
✅ Send transaction:
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

✅ View status block:
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
