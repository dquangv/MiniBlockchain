package main

import (
	"context"
	"flag"
	"fmt"
	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/p2p/pb"
	"golang-chain/pkg/wallet"
	"log"
	"math/big"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	from := flag.String("from", "", "Tên ví người gửi")
	to := flag.String("to", "", "Tên người nhận (public key hoặc tên giả)")
	amount := flag.Float64("amount", 0, "Số lượng coin")
	nodeAddr := flag.String("node", "localhost:50051", "Địa chỉ node validator")
	flag.Parse()

	if !wallet.WalletExists(*from) {
		log.Fatalf("❌ Wallet %s does not exist.", *from)
	}
	if !wallet.WalletExists(*to) {
		log.Fatalf("❌ Wallet %s does not exist.", *to)
	}

	if *from == "" || *to == "" || *amount <= 0 {
		log.Fatalln("⚠️  Dùng đúng: --from Alice --to Bob --amount 10")
	}

	w, err := wallet.LoadWallet(*from)
	if err != nil {
		log.Fatalln("❌ Không load được ví:", err)
	}

	encodedSender, _ := wallet.EncodePublicKey(w.PublicKey)
	tx := blockchain.NewTransaction(encodedSender, []byte(*to), *amount)
	if err := tx.Sign(w.PrivateKey); err != nil {
		log.Fatalln("❌ Lỗi khi ký giao dịch:", err)
	}

	conn, err := grpc.Dial(*nodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("❌ Kết nối node thất bại:", err)
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)
	resp, err := client.SendTransaction(context.Background(), &pb.Transaction{
		Sender:    tx.Sender,
		Receiver:  tx.Receiver,
		Amount:    tx.Amount,
		Timestamp: tx.Timestamp,
		Signature: tx.Signature,
	})
	if err != nil {
		log.Fatalln("❌ Gửi transaction thất bại:", err)
	}

	balResp, err := client.GetBalance(context.Background(), &pb.BalanceRequest{Name: *from})
	if err != nil {
		log.Fatalf("❌ Failed to get balance: %v", err)
	}

	balance, _ := new(big.Float).SetString(balResp.Balance)
	amountFloat := big.NewFloat(*amount)

	if balance.Cmp(amountFloat) < 0 {
		log.Fatalf("❌ Insufficient balance. You have %s, trying to send %s", balance.Text('f', 2), amountFloat.Text('f', 2))
	}

	fmt.Println("📨", resp.Message)
}
