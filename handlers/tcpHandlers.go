package handlers

import (
	"net"
	"bufio"
	"io"
	"github.com/mohanarpit/yolochain/models"
	"log"
	"time"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/mohanarpit/yolochain/blockchain"
	"strconv"
)

func handleAccouncements(conn net.Conn) {

	for {
		msg := <-models.Announcements
		io.WriteString(conn, msg)
	}

}

func HandlePOSConn(conn net.Conn) {
	defer conn.Close()

	// Write all announcements to the all the connections
	go handleAccouncements(conn)

	go HandlePOSInputData()
	InputData(conn)

	// Simulating the receiving broadcast
	for {
		time.Sleep(time.Minute)
		models.Mutex.Lock()
		output, err := json.MarshalIndent(models.Blockchain, "", "    ")
		models.Mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, "Broadcast: " + string(output)+"\n")
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()

	go HandlePOWInputData()
	InputData(conn)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			output, err := json.Marshal(models.Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range models.BlockchainServer {
		log.Println("In the range loop")
		spew.Dump(models.Blockchain)
	}
	log.Println("HandleConn end")
}

func InputData(conn net.Conn) {
	var address string

	io.WriteString(conn, "Input the token balance for this account: ")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v\n", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = blockchain.CalculateStringHash(t.String())
		models.Validators[address] = balance
		log.Printf("Got the balance as %v for address: %v\n", balance, address)
		break
	}

	io.WriteString(conn, "Input the new data: ")
	scanner := bufio.NewScanner(conn)

	go func() {

		for scanner.Scan() {
			data := scanner.Text()
			log.Println("Got data from scanner")
			msg := models.Message{
				Data:      data,
				Validator: address,
			}
			models.InputMsgChan <- msg
			log.Println("Pushed the msg to inputMsgChan")
			io.WriteString(conn, "Input MOAAR data: ")
		}
	}()
}

func HandlePOSInputData() {
	log.Println("In the HandlePOSInputData")
	for {
		select {
		case msg := <-models.InputMsgChan:
			data := msg.Data
			address := msg.Validator
			log.Printf("Got msg in the InputMsgChan: %s\n", data)
			models.Mutex.Lock()
			oldBlock := models.Blockchain[len(models.Blockchain)-1 ]
			models.Mutex.Unlock()

			newBlock, err := blockchain.GenerateBlock(oldBlock, []byte(data), address)
			if err != nil {
				log.Println(err)
				continue
			}
			if blockchain.IsBlockValid(newBlock, oldBlock) {
				models.CandidateBlocks <- newBlock
			}
			log.Println("Completed generating the candidate block")
		}
	}
}

func HandlePOWInputData() {
	log.Println("In the HandlePOWInputData")
	for {
		select {
		case msg := <-models.InputMsgChan:
			data := msg.Data
			log.Println("Got msg in InputMsgChan: " + data)
			oldBlock := models.Blockchain[len(models.Blockchain)-1]
			newBlock, err := blockchain.GenerateBlock(oldBlock, []byte(data), "")
			if err != nil {
				log.Println(err)
			}
			if blockchain.IsBlockValid(newBlock, oldBlock) {
				newBlockchain := append(models.Blockchain, newBlock)
				blockchain.ReplaceChains(newBlockchain)
			}

			// TODO: When we get a message on the BlockchainServer channel, we should also broadcast it to all the ports
			// in the localhost network. This ensures that we can run the Blockchain server on multiple machines in the cluster
			models.BlockchainServer <- models.Blockchain
		}
	}

}
