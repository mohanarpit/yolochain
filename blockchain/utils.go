package blockchain

import (
	"github.com/mohanarpit/yolochain/models"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"math/rand"
	"log"
	"github.com/davecgh/go-spew/spew"
	"github.com/mohanarpit/yolochain/blockchainGrpc"
)

func CalculateStringHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// CalculateHash simply calculates the hash for a block and stores it in that block. Used to check if the data is valid
// and the block hasn't been tampered with.
func CalculateHash(block models.Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Data) + block.PrevHash
	return CalculateStringHash(record)
}

// GenerateBlock creates a new block based on the arbitrary data and the address of the client node that proposed it
func GenerateBlock(oldBlock models.Block, data []byte, address string) (models.Block, error) {

	var newBlock models.Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)
	newBlock.Validator = address

	return newBlock, nil
}

// IsBlockValid checks if the new block is valid as compared to the last old block
func IsBlockValid(newBlock models.Block, oldBlock models.Block) bool {
	if (oldBlock.Index+1 != newBlock.Index) || (oldBlock.Hash != newBlock.PrevHash) || (CalculateHash(newBlock) != newBlock.Hash) {
		return false
	}

	return true
}

// ReplaceChains is used to overwrite the local copy of the blockchain if a longer chain is found in the network.
func ReplaceChains(newBlocks []models.Block) {
	if (len(newBlocks)) > len(models.Blockchain) {
		models.Blockchain = newBlocks
	}
}

// BootstrapBlockchain bootstraps a blockchain with the genesis block. It's typically empty values.
// Returns the genesis block
func BootstrapBlockchain() models.Block {
	t := time.Now()
	genesisBlock := models.Block{0, t.String(), []byte(string(0)), "", "", ""}
	spew.Dump(genesisBlock)
	models.Blockchain = append(models.Blockchain, genesisBlock)
	return genesisBlock
}

// HandleCandidateBlocks simply receives the candidate blocks via a Go channel and adds them to an internal temp array
func HandleCandidateBlocks() {
	// Add all the candidate blocks into a temp array
	for {
		select {
		case candidateBlock := <-models.CandidateBlocks:
			models.Mutex.Lock()
			log.Println("Going to append the candidateBlock to the list of tempBlocks")
			models.TempCandidateBlocks = append(models.TempCandidateBlocks, candidateBlock)
			models.Mutex.Unlock()
		}
	}
}

// PickPOSWinner picks the winner node for the new block based on the number of tokens that have been staked by
// individual clients. It does this every 30 seconds.
func PickPOSWinner() {
	time.Sleep(30 * time.Second)
	log.Println("Going to the pick the winner")

	models.Mutex.Lock()
	temp := models.TempCandidateBlocks
	models.Mutex.Unlock()

	lotteryPool := []string{}
	if len(temp) > 0 {
	OUTER:
		for _, block := range temp {

			// If the node is already in the lottery pool, skip it
			for _, node := range lotteryPool {
				if block.Validator == node {
					log.Printf("Got the node in lotteryPool as %v. Skipping", node)
					continue OUTER
				}
			}

			models.Mutex.Lock()
			setValidators := models.Validators
			models.Mutex.Unlock()

			// Based on the number of tokens staked, add those many items of the Validator address node to the list
			// This will ensure that when we randomly pick nodes, the probability of picking the node changes based
			// on the number of tokens that have been staked
			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

		for _, block := range temp {
			if block.Validator == lotteryWinner {

				// Appending to the local blockchain is done by the AnnounceCandidates GRPC handler
				for _ = range models.Validators {
					// Transform the local block to the grpcChain Block so that we can push it over the wire
					// TODO: Change this to use the GRPC blockchain only so that we don't have to keep transforming the values
					grpcBlock := blockchainGrpc.Block{
						Data: block.Data,
						Validator: block.Validator,
						Hash: block.Hash,
						PrevHash: block.PrevHash,
						Timestamp: block.Timestamp,
						Index: block.Index,
					}
					req := blockchainGrpc.AnnounceCandidateRequest{
						Message: "Winning Validator " + lotteryWinner,
						Block: &grpcBlock,
					}
					models.Announcements <- req
				}
				break
			}
		}
	}

	log.Println("Going to clean the tempBlocks array after picking the winner")
	models.Mutex.Lock()
	models.TempCandidateBlocks = []models.Block{}
	models.Mutex.Unlock()

}
