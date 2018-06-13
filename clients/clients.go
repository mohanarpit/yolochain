package clients

import (
	"log"
	"github.com/mohanarpit/yolochain/models"
	"google.golang.org/grpc"
	"github.com/mohanarpit/yolochain/blockchainGrpc"
	"time"
	"context"
)

func AnnounceCandidates(req blockchainGrpc.AnnounceCandidateRequest) {
	var conn *grpc.ClientConn
	var err error

	for _, address := range models.Cluster {
		log.Println("Going to Announce candidates to ", address)
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err!= nil {
			log.Fatalln(err)
		}
		client := blockchainGrpc.NewControlServiceClient(conn)
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.AnnounceCandidates(ctx, &req)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(resp.Success)
		conn.Close()
	}
}

