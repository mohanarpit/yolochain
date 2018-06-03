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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := models.Block{0, t.String(), []byte(string(0)), "", "", ""}
		spew.Dump(genesisBlock)
		models.Blockchain = append(models.Blockchain, genesisBlock)
	}()

	log.Fatal(run())
}

func networkMain() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	models.BlockchainServer = make(chan []models.Block)
	models.InputMsgChan = make(chan string)

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

func main() {
	//standaloneMain()
	networkMain()
}
