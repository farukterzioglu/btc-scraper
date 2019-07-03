package dtos

// BlockDto is a DTO for api requests
type BlockDto struct {
	Hash   string `json:"hash"`
	Height int64  `json:"height"`
}
