package main

import (
	"time"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"log"
	"github.com/mohanarpit/yolochain/handlers"
	"github.com/mohanarpit/yolochain/models"
	"github.com/joho/godotenv"
	"net"
	"github.com/mohanarpit/yolochain/blockchain"
	"flag"
)

func runHttpServer(httpHandler http.Handler) error {
	httpAddr := os.Getenv("HTTP_ADDR")
	log.Println("Listening on HTTP Port: ", httpAddr)
	s := &http.Server{
		Addr:           httpAddr,
		Handler:        httpHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func makeStandaloneHttpRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handlers.HandleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handlers.HandleWriteBlockchain).Methods("POST")
	return muxRouter
}

func makePOSHttpRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handlers.HandleGetBlockchain).Methods("GET")
	return muxRouter
}

// standaloneMain is the simplest blockchain that runs in standalone mode
// It exposes a HTTP server which can be used to query and write data to the Blockchain
func standaloneMain() {
	blockchain.BootstrapBlockchain()
	log.Fatal(runHttpServer(makeStandaloneHttpRouter()))
}

// networkMain supports the "network" mode in the blockchain. It allows clients to connect to it and create new blocks
// Currently it doesn't completely satisfy the POW as defined by the Blockchain paper
func networkMain() {

	models.BlockchainServer = make(chan []models.Block)
	models.InputMsgChan = make(chan models.Message)

	blockchain.BootstrapBlockchain()
	server, err := net.Listen("tcp", os.Getenv("TCP_ADDR"))
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handlers.HandleConn(conn)
	}

}

// posMain runs the blockchain server in Proof of Stake mode. It allows clients to connect to it, stake a certain number of
// tokens and then assigns the new block to the winner client node based on the number of tokens that are staked.
func posMain() {

	models.BlockchainServer = make(chan []models.Block)
	models.InputMsgChan = make(chan models.Message)

	// Genesis block to bootstrap the blockchain application
	blockchain.BootstrapBlockchain()

	// TCP Server to accept connections from clients
	tcpAddr := os.Getenv("TCP_ADDR")
	server, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	log.Println("Listening on TCP Port: ", tcpAddr)

	// Goroutine to handle the candidateBlocks from which the winner block will be chosen
	go blockchain.HandleCandidateBlocks()

	// Goroutine to pick the winners at regular intervals
	go func() {
		for {
			blockchain.PickPOSWinner()
		}
	}()

	// Goroutine to start the HTTP server for REST calls
	go func() {
		log.Fatal(runHttpServer(makePOSHttpRouter()))
	}()

	// Goroutine to handle the TCP connections for clients staking tokens
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handlers.HandlePOSConn(conn)
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var mode = flag.String("mode", "pos", "The mode in which you want to run the blockchain. Available modes are standalone/network/pos. Default is pos.")
	flag.Parse()
	log.Println("Going to run the Yolochain in mode: ", *mode)
	switch *mode {
	case "pos":
		posMain()
	case "standalone":
		standaloneMain()
	case "network":
		networkMain()
	default:
		break
	}
}
