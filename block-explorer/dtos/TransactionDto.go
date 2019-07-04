package dtos

type Vin struct {
	Coinbase  string     `json:"coinbase"`
	Txid      string     `json:"txid"`
	Vout      uint32     `json:"vout"`
	Sequence  uint32     `json:"sequence"`
	Witness   []string   `json:"txinwitness"`
}

type Vout struct {
	Value        float64            `json:"value"`
	N            uint32             `json:"n"`
}

type TransactionDto struct {
	Txid          string `json:"txid"`
	Hash          string `json:"hash,omitempty"`
	Vsize         int32  `json:"vsize,omitempty"`
	Version       int32  `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}
