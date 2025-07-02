package p2p

import (
	"context"
	"log"
	"time"

	"golang-chain/pkg/p2p/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MonitorLeader(server *NodeServer, peers []string) {
	go func() {
		for {
			time.Sleep(5 * time.Second)

			if *server.State != StateFollower || server.LeaderID == "" {
				continue
			}

			addr := peerAddressByID(server.LeaderID, peers)
			if addr == "" {
				continue
			}

			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Println("⚠️ Leader unreachable. Re-electing...")
				StartElection(server, peers)
				continue
			}
			defer conn.Close()

			client := pb.NewNodeServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			resp, err := client.Ping(ctx, &pb.Empty{})
			if err != nil || resp.Message != string(StateLeader) {
				log.Println("🚨 Leader unresponsive or demoted. Re-electing...")
				StartElection(server, peers)
				continue
			}

			log.Printf("✅ Leader %s still alive", server.LeaderID)
		}
	}()
}
