package models

import (
	"sync"
	"github.com/mohanarpit/yolochain/blockchainGrpc"
)

type Block struct {
	Index     int64
	Timestamp string
	Data      []byte
	Hash      string
	PrevHash  string
	Validator string
}

type Message struct {
	Data      string
	Validator string
}

var Blockchain []Block

var BlockchainServer chan []Block
var InputMsgChan chan Message

// CandidateBlocks handles incoming blocks for validation in POS blockchain
var CandidateBlocks = make(chan Block)
var TempCandidateBlocks []Block

// Announcements broadcasts winning validator to all nodes in POS blockchain
var Announcements = make(chan blockchainGrpc.AnnounceCandidateRequest)

var Mutex = &sync.Mutex{}

// Validators keeps track of open validators and balances in POS blockchain
var Validators = make(map[string]int)

// Cluster is the array of servers that make up the blockchain system in the local network. Hardcoded for now
var Cluster []string