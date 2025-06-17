package main

import (
	"fmt"
	"golang-chain/pkg/p2p"
	"os"
)

func main() {
	fmt.Println("Hello from validator node!")
	port := os.Args[1] // ví dụ: go run main.go 50052
	p2p.StartGRPCServer(port)
}
