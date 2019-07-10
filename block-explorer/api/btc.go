package api

import (
	"encoding/json"
	"net/http"

	"github.com/farukterzioglu/btc-scraper/services"
	"github.com/gorilla/mux"
)

// BtcRoutes contains btc endpoints
type BtcDbRoutes struct {
	elasticService *services.ElasticService
}

// NewBtcRoutes create a new BtcRoutes instance
func NewBtcDbRoutes(e *services.ElasticService) *BtcDbRoutes {
	return &BtcDbRoutes{
		elasticService: e,
	}
}

// RegisterBtcRoutes registers routes for Btc
func (routes *BtcDbRoutes) RegisterRoutes(r *mux.Router, p string) {
	ur := r.PathPrefix(p).Subrouter()

	// swagger:route GET /btc/block Btc getBlocksReq
	// ---
	// Returns last 10 blocks.
	//
	// responses:
	//   200: blocksResp
	//   404: notFound
	//	 500: internal
	ur.HandleFunc("/block", routes.getBlocks).Methods("GET")

	// swagger:route GET /btc/block/{BlockID} Btc getBlockReq
	// ---
	// Returns a block by id.
	// If the block id is null, Error Bad Request will be returned.
	// responses:
	//   200: blockResp
	//   404: notFound
	//	 500: internal
	ur.HandleFunc("/block/{BlockID}", routes.getBlock).Methods("GET")

	// swagger:route GET /btc/tx/{TxHash} Btc getTransactionReq
	// ---
	// Returns a transaction by hash.
	// If the transaction hash is null, Error Bad Request will be returned.
	// responses:
	//   200: transactionResp
	//   404: notFound
	//	 500: internal
	ur.HandleFunc("/tx/{TxHash}", routes.getTransaction).Methods("GET")
}

func (routes *BtcDbRoutes) getBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	blockIDStr := params["BlockID"]

	block, err := routes.elasticService.GetBlock(blockIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(block); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (routes *BtcDbRoutes) getBlocks(w http.ResponseWriter, r *http.Request) {
	blockList, err := routes.elasticService.GetBlocks(10)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blockList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (routes *BtcDbRoutes) getTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	txHashStr := params["TxHash"]

	tx, err := routes.elasticService.GetTransaction(txHashStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
