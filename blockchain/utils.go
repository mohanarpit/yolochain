package blockchain

import (
	"github.com/mohanarpit/yolochain/models"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"math/rand"
	"log"
	"github.com/davecgh/go-spew/spew"
)

func CalculateStringHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CalculateHash(block models.Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Data) + block.PrevHash
	return CalculateStringHash(record)
}

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

func IsBlockValid(newBlock models.Block, oldBlock models.Block) bool {
	if (oldBlock.Index+1 != newBlock.Index) || (oldBlock.Hash != newBlock.PrevHash) || (CalculateHash(newBlock) != newBlock.Hash) {
		return false
	}

	return true
}

func ReplaceChains(newBlocks []models.Block) {
	if (len(newBlocks)) > len(models.Blockchain) {
		models.Blockchain = newBlocks
	}
}

func BootstrapBlockchain() models.Block {
	t := time.Now()
	genesisBlock := models.Block{0, t.String(), []byte(string(0)), "", "", ""}
	spew.Dump(genesisBlock)
	models.Blockchain = append(models.Blockchain, genesisBlock)
	return genesisBlock
}

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

			s := rand.NewSource(time.Now().Unix())
			r := rand.New(s)
			lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

			for _, block := range temp {
				if block.Validator == lotteryWinner {
					models.Mutex.Lock()
					models.Blockchain = append(models.Blockchain, block)
					models.Mutex.Unlock()
					for _ = range models.Validators {
						models.Announcements <- "\nWinning Validator: " + lotteryWinner + "\n"
					}
					break
				}
			}
		}
	}

	log.Println("Going to clean the tempBlocks array after picking the winner")
	models.Mutex.Lock()
	models.TempCandidateBlocks = []models.Block{}
	models.Mutex.Unlock()

}
