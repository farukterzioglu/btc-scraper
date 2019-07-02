package swagger

import (
	"github.com/farukterzioglu/btc-scraper/block-explorer/dtos"
)

// Request containing block id
// swagger:parameters getBlockReq
type swaggGetBlockReq struct {
	// in:path
	// description: id of the block
	// type: string
	// required: true
	BlockID string
}

// HTTP status code 200 and an array of block models in data
// swagger:response blocksResp
type swaggBlocksResp struct {
	// Array of block models
	// in:body
	Body []dtos.BlockDto
}

// HTTP status code 200 and a block model in data
// swagger:response blockResp
type swaggBlockResp struct {
	// A block models
	// in:body
	Body dtos.BlockDto
}
