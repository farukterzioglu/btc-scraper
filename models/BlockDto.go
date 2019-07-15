package models

// BlockDto
type BlockDto struct {
	Hash          string   `json:"hash"`
	Confirmations int64    `json:"confirmations"`
	Height        int64    `json:"height"`
	Tx            []string `json:"tx,omitempty"`
	Time          int64    `json:"time"`
	PreviousHash  string   `json:"previousblockhash"`
	NextHash      string   `json:"nextblockhash,omitempty"`
}
