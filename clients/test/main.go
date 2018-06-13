package main

import (
	"google.golang.org/grpc"
	"log"
	"github.com/mohanarpit/yolochain/blockchainGrpc"
	"context"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:5050", grpc.WithInsecure())
	if err!= nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := blockchainGrpc.NewControlServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.AnnounceCandidates(ctx, &blockchainGrpc.AnnounceCandidateRequest{Message: "arpit"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp.Success)
}
