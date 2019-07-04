package api

import (
	"encoding/json"
	"net/http"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/farukterzioglu/btc-scraper/block-explorer/dtos"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
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
	// Returns last 10 blocks.
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

	// swagger:route GET /btc/tx/{TxHash} BtcAPI getTransactionReq
	// ---
	// Returns a transaction by hash.
	// If the transaction hash is null, Error Bad Request will be returned.
	// responses:
	//   200: transactionResp
	//   404: notFound
	//	 500: internal
	ur.HandleFunc("/tx/{TxHash}", routes.getTransaction).Methods("GET")
}

func (routes *BtcRoutes) getBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	blockIDStr := params["BlockID"]

	hash, err := chainhash.NewHashFromStr(blockIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	block, err := routes.client.GetBlockVerbose(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	blockDto := dtos.BlockDto{
		Hash:   block.Hash,
		Height: block.Height,
		Time: block.Time,
		Confirmations: block.Confirmations,
		Tx: block.Tx,
		NextHash: block.NextHash,
		PreviousHash: block.PreviousHash,
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blockDto); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (routes *BtcRoutes) getBlocks(w http.ResponseWriter, r *http.Request) {
	blockCount, err := routes.client.GetBlockCount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	blockListLenght := 10
	if blockCount < 10 {
		blockListLenght = int(blockCount)
	}

	blockList := make([]dtos.BlockDto, blockListLenght, 10)

	for i := 0; i < blockListLenght; i++ {
		height := blockCount - int64(i)
		blockHash, err := routes.client.GetBlockHash(height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		blockList[i] = dtos.BlockDto{
			Hash:   blockHash.String(),
			Height: height,
		}
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blockList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (routes *BtcRoutes) getTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	txHashStr := params["TxHash"]

	txHash, err := chainhash.NewHashFromStr(txHashStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tx, err := routes.client.GetRawTransactionVerbose(txHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// TODO : map tx to TransactionDto

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}