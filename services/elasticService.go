package services

import (
	"context"
	"encoding/json"

	"github.com/ahmetb/go-linq"
	"github.com/farukterzioglu/btc-scraper/models"
	"github.com/olivere/elastic/v7"
)

type ElasticService struct {
	cryptoCode string
	client     *elastic.Client
}

func NewElasticService(cryptoCode string, e *elastic.Client) *ElasticService {
	return &ElasticService{
		cryptoCode: cryptoCode,
		client:     e,
	}
}

func (e *ElasticService) InsertBlock(block models.BlockDto) error {
	ctx := context.Background()
	_, err := e.client.Index().
		Index(e.cryptoCode + "-block").
		Type("block").
		Id(block.Hash).
		BodyJson(block).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e *ElasticService) InsertTx(tx models.TransactionDto) error {
	ctx := context.Background()
	_, err := e.client.Index().
		Index(e.cryptoCode + "-transaction").
		Type("transaction").
		Id(tx.Hash).
		BodyJson(tx).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetBlocks returns last 'blockCount' blocks
func (e *ElasticService) GetBlocks(blockCount int) ([]models.BlockDto, error) {
	searchResult, err := e.client.Search().
		Index(e.cryptoCode+"-block").
		Sort("height", false).
		From(0).Size(blockCount).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	if searchResult.Hits.TotalHits.Value > 0 {
		var blockList []models.BlockDto
		linq.From(searchResult.Hits.Hits).WhereT(func(hit *elastic.SearchHit) bool {
			return true
		}).SelectT(func(hit *elastic.SearchHit) models.BlockDto {
			var b models.BlockDto
			err := json.Unmarshal(hit.Source, &b)
			if err != nil {
				// TODO : Deserialization failed
			}

			return b
		}).ToSlice(&blockList)

		return blockList, nil
	} else {
		return nil, nil
	}
}

func (e *ElasticService) GetBlock(blockHash string) (models.BlockDto, error) {
	termQuery := elastic.NewTermQuery("hash", blockHash)
	searchResult, err := e.client.Search().
		Index(e.cryptoCode + "-block").
		Query(termQuery).
		From(0).Size(1).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return models.BlockDto{}, err
	}

	if searchResult.Hits.TotalHits.Value == 0 {
		return models.BlockDto{}, nil
	}

	var block models.BlockDto
	err = json.Unmarshal(searchResult.Hits.Hits[0].Source, &block)
	if err != nil {
		return models.BlockDto{}, err
	}

	return block, nil
}

func (e *ElasticService) GetTransaction(txHash string) (models.TransactionDto, error) {
	termQuery := elastic.NewTermQuery("hash", txHash)
	searchResult, err := e.client.Search().
		Index(e.cryptoCode + "-transaction").
		Query(termQuery).
		From(0).Size(1).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return models.TransactionDto{}, err
	}

	if searchResult.Hits.TotalHits.Value == 0 {
		return models.TransactionDto{}, nil
	}

	var tx models.TransactionDto
	err = json.Unmarshal(searchResult.Hits.Hits[0].Source, &tx)
	if err != nil {
		return models.TransactionDto{}, err
	}

	return tx, nil
}
