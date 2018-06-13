package handlers

import (
	"context"
	"github.com/mohanarpit/yolochain/blockchainGrpc"
	"log"
	"github.com/mohanarpit/yolochain/models"
)

// server is used to implement helloworld.GreeterServer.
type Server struct{}

func (s *Server) AnnounceCandidates(ctx context.Context, in *blockchainGrpc.AnnounceCandidateRequest) (*blockchainGrpc.AnnounceCandidateReply, error) {
	log.Println("In the AnnounceCandidates server with msg: ", in.Message)
	// Append this new block to the existing Blockchain on the local server
	models.Mutex.Lock()
	defer models.Mutex.Unlock()

	localBlock := models.Block{
		Index: in.Block.Index,
		Timestamp: in.Block.Timestamp,
		PrevHash: in.Block.PrevHash,
		Hash: in.Block.Hash,
		Validator: in.Block.Validator,
		Data: in.Block.Data,
	}
	models.Blockchain = append(models.Blockchain, localBlock)
	log.Println("Appended the block to the local blockchain")
	return &blockchainGrpc.AnnounceCandidateReply{Success: true}, nil
}

func (s *Server) MulticastPing(ctx context.Context, in *blockchainGrpc.PingRequest) (*blockchainGrpc.PingReply, error) {
	return &blockchainGrpc.PingReply{}, nil
}

