package models

import "sync"

type Block struct {
	Index     int
	Timestamp string
	Data      []byte
	Hash      string
	PrevHash  string
	Validator string
}

type Message struct {
	Data string
	Validator string
}

var Blockchain []Block

var BlockchainServer chan []Block
var InputMsgChan chan Message

// candidateBlocks handles incoming blocks for validation in POS blockchain
var CandidateBlocks = make(chan Block)
var TempCandidateBlocks []Block

// announcements broadcasts winning validator to all nodes in POS blockchain
var Announcements = make(chan string)

var Mutex = &sync.Mutex{}

// validators keeps track of open validators and balances in POS blockchain
var Validators = make(map[string]int)

