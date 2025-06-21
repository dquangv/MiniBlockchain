package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang-chain/pkg/wallet"
	"os"
	"path/filepath"
)

func main() {
	name := flag.String("name", "", "Tên ví (VD: Alice, Bob)")
	flag.Parse()

	if *name == "" {
		fmt.Println("⚠️  Vui lòng nhập tên ví bằng flag --name")
		return
	}

	w, err := wallet.NewWallet()
	if err != nil {
		fmt.Println("❌ Lỗi khi tạo ví:", err)
		return
	}

	encodedPub, _ := wallet.EncodePublicKey(w.PublicKey)
	encodedPriv, _ := wallet.EncodePrivateKey(w.PrivateKey)

	// Tạo JSON và lưu vào file
	data := map[string]string{
		"publicKey":  string(encodedPub),
		"privateKey": string(encodedPriv),
	}

	walletsDir := "wallets"
	os.MkdirAll(walletsDir, os.ModePerm)

	filePath := filepath.Join(walletsDir, *name+"_wallet.json")
	file, _ := os.Create(filePath)
	defer file.Close()

	json.NewEncoder(file).Encode(data)

	fmt.Println("✅ The wallet has been created and saved at: ", filePath)
}
