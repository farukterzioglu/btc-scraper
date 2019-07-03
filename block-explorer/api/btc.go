package api

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"log"
	_ "log"
	"net/http"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/farukterzioglu/btc-scraper/block-explorer/dtos"
	_ "github.com/farukterzioglu/btc-scraper/block-explorer/dtos"
	"github.com/gorilla/mux"
)

// BtcRoutes contains btc endpoints
type BtcRoutes struct {
	client *rpcclient.Client
}

// NewBtcRoutes create a new BtcRoutes instance
func NewBtcRoutes(client *rpcclient.Client) *BtcRoutes {
	return &BtcRoutes{
		client: client,
	}
}

// RegisterBtcRoutes registers routes for Btc
func (routes *BtcRoutes) RegisterBtcRoutes(r *mux.Router, p string) {
	ur := r.PathPrefix(p).Subrouter()

	// swagger:route GET /btc/block BtcAPI blockList
	// ---
	// Returns all blocks.
	//
	// responses:
	//   200: blocksResp
	//   404: notFound
	//	 500: internal
	ur.HandleFunc("/block", routes.getBlocks).Methods("GET")

	// swagger:route GET /btc/block/{BlockID} BtcAPI getBlockReq
	// ---
	// Returns a block by id.
	// If the block id is null, Error Bad Request will be returned.
	// responses:
	//   200: blockResp
	//   404: notFound
	//	 500: internal
	ur.HandleFunc("/block/{BlockID}", routes.getBlock).Methods("GET")
}

// TODO : Implement this
func (routes *BtcRoutes) getBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	blockIDStr := params["BlockID"]
	fmt.Printf("Getting block %s...\n", blockIDStr)

	block := dtos.BlockDto{
		Hash:   blockIDStr,
		Height: 0,
	}
	if err := json.NewEncoder(w).Encode(block); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

// TODO : Implement this
func (routes *BtcRoutes) getBlocks(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting blocks...\n")

	// Get the current block count.
	blockCount, err := routes.client.GetBlockCount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("Block count: %d", blockCount)

	if blockCount > 10 {
		blockCount = 10
	}

	blockList := make([]dtos.BlockDto, blockCount, 10)

	for i := blockCount; i > 0; i-- {
		blockHash, err := routes.client.GetBlockHash(i)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("Block %d, hash: %s", i, blockHash)

		blockList = append(blockList, dtos.BlockDto{
			Hash:   blockHash,
			Height: i,
		})
	}

	if err := json.NewEncoder(w).Encode(blockList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
