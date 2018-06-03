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
)

func HandleConn(conn net.Conn) {
	defer conn.Close()

	go HandleInputData()
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
	io.WriteString(conn, "Input the new data")
	scanner := bufio.NewScanner(conn)

	go func() {

		for scanner.Scan() {
			data := scanner.Text()
			log.Println("Got data from scanner")
			models.InputMsgChan <- data
			log.Println("Pushed the msg to inputMsgChan")
			io.WriteString(conn, "Input MOAAR data")
		}
	}()

}

func HandleInputData() {
	log.Println("In the HandleInputData")
	for {
		select {
		case data := <-models.InputMsgChan:
			log.Println("Got msg in InputMsgChan: " + data)
			oldBlock := models.Blockchain[len(models.Blockchain) -1]
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