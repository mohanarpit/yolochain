package blockchain

import (
	"github.com/mohanarpit/yolochain/models"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func calculateStringHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CalculateHash(block models.Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Data) + block.PrevHash
	return calculateStringHash(record)
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
