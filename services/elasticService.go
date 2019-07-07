package services

import (
	"github.com/elastic/go-elasticsearch"
	"github.com/farukterzioglu/btc-scraper/block-explorer/dtos"
)

type ElasticService struct {
	elastic *elasticsearch.Client
}

func NewElasticService(e *elasticsearch.Client) *ElasticService {
	return &ElasticService{
		elastic: e,
	}
}

// TODO : Impleöemt elastic gets

// GetBlocks returns last 'blockCount' blocks
func (e *ElasticService) GetBlocks(cryptoCode string, blockCount int64) ([]dtos.BlockDto, error) {
	return []dtos.BlockDto{}, nil
}

func (e *ElasticService) GetBlock(cryptoCode string, blockHash string) (dtos.BlockDto, error) {
	return dtos.BlockDto{}, nil
}

func (e *ElasticService) GetTransaction(cryptoCode string, txHash string) (dtos.TransactionDto, error) {
	return dtos.TransactionDto{}, nil
}
