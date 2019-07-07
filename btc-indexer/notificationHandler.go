package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/farukterzioglu/btc-scraper/models"
)

type NotificationHandler struct {
	client  *rpcclient.Client
	elastic *elasticsearch.Client
}

func NewNotificationHandler(c *rpcclient.Client, e *elasticsearch.Client) *NotificationHandler {
	return &NotificationHandler{
		client:  c,
		elastic: e,
	}
}

func (handler *NotificationHandler) ConsumeBlocks(chn <-chan BlockNotification) *NotificationHandler {
	log.Println("Started to consume blocks...")

	for {
		notification := <-chn

		log.Printf("Block connected: %s, height: %d, tx count: %d", notification.Header.BlockHash().String(), notification.BlockHeight, len(notification.TxList))
		for _, tx := range notification.TxList {
			log.Printf("   New transaction for tracked address: %s", tx.Hash().String())
		}

		hash := notification.Header.BlockHash()

		// Get block details with rpc
		var block *btcjson.GetBlockVerboseResult
		block, err := handler.client.GetBlockVerboseTx(&hash)
		if err != nil {
			log.Printf("couldn't get block details")
			continue
		}

		// log.Printf("Block details: %+v", block)

		blockDto := models.BlockDto{
			Hash:          block.Hash,
			Confirmations: block.Confirmations,
			Height:        block.Height,
			Time:          block.Time,
			PreviousHash:  block.PreviousHash,
			NextHash:      block.NextHash,
		}

		data, _ := json.Marshal(blockDto)
		req := esapi.IndexRequest{
			Index:      "block",
			DocumentID: blockDto.Hash,
			Body:       strings.NewReader(string(data)),
			Refresh:    "true",
		}
		res, err := req.Do(context.Background(), handler.elastic)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document ID=%d", res.Status(), blockDto.Hash)
		} else {
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				log.Printf("[%s] block %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}

		// Map tx to dto
		var transactionList []models.TransactionDto
		linq.From(block.RawTx).SelectT(func(tx btcjson.TxRawResult) models.TransactionDto {
			var vinList []models.Vin

			linq.From(tx.Vin).SelectT(func(vin btcjson.Vin) models.Vin {
				return models.Vin{
					Txid:     vin.Txid,
					Vout:     vin.Vout,
					Sequence: vin.Sequence,
				}
			}).ToSlice(&vinList)

			var voutList []models.Vout
			linq.From(tx.Vout).SelectT(func(vout btcjson.Vout) models.Vout {
				return models.Vout{
					Value: vout.Value,
					N:     vout.N,
					ScriptPubKey: models.ScriptPubKeyResult{
						Addresses: vout.ScriptPubKey.Addresses,
						Asm:       vout.ScriptPubKey.Asm,
						Hex:       vout.ScriptPubKey.Hex,
						ReqSigs:   vout.ScriptPubKey.ReqSigs,
						Type:      vout.ScriptPubKey.Type,
					},
				}
			}).ToSlice(&voutList)

			txDto := models.TransactionDto{
				Hash:          tx.Hash,
				BlockHash:     tx.BlockHash,
				Blocktime:     tx.Blocktime,
				Confirmations: tx.Confirmations,
				LockTime:      tx.LockTime,
				Time:          tx.Time,
				Txid:          tx.Txid,
				Version:       tx.Version,
				Vsize:         tx.Vsize,
				Vin:           vinList,
				Vout:          voutList,
			}

			return txDto
		}).ToSlice(&transactionList)

		// TODO : Insert bulk
		for _, transaction := range transactionList {
			data, _ := json.Marshal(transaction)
			req := esapi.IndexRequest{
				Index:      "transaction",
				DocumentID: transaction.Hash,
				Body:       strings.NewReader(string(data)),
				Refresh:    "true",
			}
			res, err := req.Do(context.Background(), handler.elastic)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), transaction.Hash)
			} else {
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					log.Printf("[%s] transaction %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}
	}

	return handler
}

func (handler *NotificationHandler) ConsumeTx(chn <-chan TxNotification) *NotificationHandler {
	log.Println("Started to consume transactions...")

	for {
		txNtf := <-chn
		tx := txNtf.Tx

		log.Printf("Tx accepted (verbose): %s (Conf.: %d)", tx.Hash, tx.Confirmations)

		for _, vout := range tx.Vout {
			log.Printf("   Output amount: %f", vout.Value)
		}
	}

	return handler
}

func (handler *NotificationHandler) ConsumeRelevantTxHex(chn <-chan string) *NotificationHandler {
	log.Println("Started to consume relevant transactions...")

	for {
		hex := <-chn

		log.Printf("Relevant tx accepted")
		log.Printf("   %s", hex)
	}

	return handler
}
