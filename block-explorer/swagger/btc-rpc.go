package swagger

import (
	"github.com/farukterzioglu/btc-scraper/block-explorer/dtos"
)

// Request containing block id
// swagger:parameters getBlockRpcReq
type swaggGetBlockRpcReq struct {
	// in:path
	// description: id of the block
	// type: string
	// required: true
	BlockID string
}

// HTTP status code 200 and an array of block models in data
// swagger:response blocksRpcResp
type swaggBlocksRpcResp struct {
	// Array of block models
	// in:body
	Body []dtos.BlockDto
}

// HTTP status code 200 and a block model in data
// swagger:response blockRpcResp
type swaggBlockRpcResp struct {
	// A block model
	// in:body
	Body dtos.BlockDto
}

// Request containing transaction hash
// swagger:parameters getTransactionRpcReq
type swaggGetTxRpcReq struct {
	// in:path
	// description: hash of the tx
	// type: string
	// required: true
	TxHash string
}

// HTTP status code 200 and a transaction model in data
// swagger:response transactionRpcResp
type swaggTransactionRpcResp struct {
	// A transaction model
	// in:body
	Body dtos.TransactionDto
}
