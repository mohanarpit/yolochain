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
	"github.com/davecgh/go-spew/spew"
	"net"
	"github.com/mohanarpit/yolochain/blockchain"
	"flag"
)

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handlers.HandleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handlers.HandleWriteBlockchain).Methods("POST")
	return muxRouter
}

func standaloneMain() {
	go func() {
		t := time.Now()
		genesisBlock := models.Block{0, t.String(), []byte(string(0)), "", "", ""}
		spew.Dump(genesisBlock)
		models.Blockchain = append(models.Blockchain, genesisBlock)
	}()

	log.Fatal(run())
}

func networkMain() {

	models.BlockchainServer = make(chan []models.Block)
	models.InputMsgChan = make(chan models.Message)

	t := time.Now()
	genesisBlock := models.Block{0, t.String(), []byte(string(0)), "", "", ""}
	spew.Dump(genesisBlock)
	models.Blockchain = append(models.Blockchain, genesisBlock)
	server, err := net.Listen("tcp", os.Getenv("ADDR"))
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

func posMain() {

	models.BlockchainServer = make(chan []models.Block)
	models.InputMsgChan = make(chan models.Message)

	t := time.Now()
	genesisBlock := models.Block{0, t.String(), []byte(string(0)), "", "", ""}
	spew.Dump(genesisBlock)
	models.Blockchain = append(models.Blockchain, genesisBlock)

	server, err := net.Listen("tcp", os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	go blockchain.HandleCandidateBlocks()

	go func() {
		for {
			blockchain.PickPOSWinner()
		}
	}()

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
