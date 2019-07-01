package api

import (
	_ "encoding/json"
	_ "fmt"
	_ "log"
	"net/http"

	_ "github.com/farukterzioglu/btc-scraper/block-explorer/dtos"
	"github.com/gorilla/mux"
)

// BtcRoutes contains btc endpoints
type BtcRoutes struct {
	// TODO: rpc client
}

// NewBtcRoutes create a new BtcRoutes instance
func NewBtcRoutes() *BtcRoutes {
	return &BtcRoutes{}
}

// RegisterBtcRoutes registers routes for Btc
func (routes *BtcRoutes) RegisterBtcRoutes(r *mux.Router, p string) {
	ur := r.PathPrefix(p).Subrouter()

	// swagger:route GET /block BtcAPI blockList
	// ---
	// Returns all blocks.
	//
	// responses:
	//   200: blockssResp
	//   404: notFound
	//	 500: internal
	ur.HandleFunc("", routes.getBlocks).Methods("GET")

	// swagger:route GET /block/{BlockID} BtcAPI getBlockReq
	// ---
	// Returns a block by id.
	// If the block id is null, Error Bad Request will be returned.
	// responses:
	//   200: blockResp
	//   404: notFound
	//	 500: internal
	ur.HandleFunc("/{BlockID}", controller.getBlock).Methods("GET")
}

func (route *BtcRoutes) getBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reviewIDStr := params["BlockID"]

	// TODO : Implement this

	w.WriteHeader(http.StatusOK)
}

func (routes *BtcRoutes) getBlocks(w http.ResponseWriter, r *http.Request) {
	// TODO : Implement this

	w.WriteHeader(http.StatusOK)
}
