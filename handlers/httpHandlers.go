package handlers

import (
	"net/http"
	"encoding/json"
	"io"
	"github.com/mohanarpit/yolochain/models"
	"github.com/mohanarpit/yolochain/blockchain"
	"github.com/davecgh/go-spew/spew"
)

func HandleGetBlockchain(w http.ResponseWriter, r *http.Request) {

	bytes, err := json.MarshalIndent(models.Blockchain, "", "	")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(bytes))
}

func HandleWriteBlockchain(w http.ResponseWriter, r *http.Request) {
	var m models.Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJson(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()

	oldBlock := models.Blockchain[len(models.Blockchain)-1]
	newBlock, err := blockchain.GenerateBlock(oldBlock, []byte(m.Data), "")
	if err != nil {
		respondWithJson(w, r, http.StatusInternalServerError, r.Body)
		return
	}

	if blockchain.IsBlockValid(newBlock, oldBlock) {
		newBlockchain := append(models.Blockchain, newBlock)
		blockchain.ReplaceChains(newBlockchain)
		spew.Dump(models.Blockchain)
	}

	respondWithJson(w, r, http.StatusCreated, newBlock)
}

func respondWithJson(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "		")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write([]byte(response))
	return
}
