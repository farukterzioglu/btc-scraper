package main

import (
	"log"

	"github.com/ahmetb/go-linq"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/farukterzioglu/btc-scraper/models"
	"github.com/farukterzioglu/btc-scraper/services"
)

type NotificationHandler struct {
	client  *rpcclient.Client
	service *services.ElasticService
}

func NewNotificationHandler(c *rpcclient.Client, s *services.ElasticService) *NotificationHandler {
	return &NotificationHandler{
		client:  c,
		service: s,
	}
}

func (handler *NotificationHandler) ConsumeBlocks(chn <-chan BlockNotification) {
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

		// Map tx to dto
		var transactionList []models.TransactionDto
		linq.From(block.RawTx).WhereT(func(tx btcjson.TxRawResult) bool {
			// Don't index coinbase tx
			return tx.Vin[0].Coinbase == ""
		}).SelectT(func(tx btcjson.TxRawResult) models.TransactionDto {
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
			handler.service.InsertTx(transaction)
			// fmt.Printf("Indexed transaciton %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		}

		// Insert block (w/ tx hashes)
		var txHashes []string
		linq.From(transactionList).SelectT(func(tx models.TransactionDto) string {
			return tx.Hash
		}).ToSlice(&txHashes)

		blockDto := models.BlockDto{
			Hash:          block.Hash,
			Confirmations: block.Confirmations,
			Height:        block.Height,
			Time:          block.Time,
			PreviousHash:  block.PreviousHash,
			NextHash:      block.NextHash,
			Tx:            txHashes,
		}

		err = handler.service.InsertBlock(blockDto)
		if err != nil {
			log.Print(err)
		}
	}
}

func (handler *NotificationHandler) ConsumeTx(chn <-chan TxNotification) {
	log.Println("Started to consume transactions...")

	for {
		txNtf := <-chn
		tx := txNtf.Tx

		log.Printf("Tx accepted (verbose): %s (Conf.: %d)", tx.Hash, tx.Confirmations)

		for _, vout := range tx.Vout {
			log.Printf("   Output amount: %f", vout.Value)
		}
	}
}

func (handler *NotificationHandler) ConsumeRelevantTxHex(chn <-chan string) {
	log.Println("Started to consume relevant transactions...")

	for {
		hex := <-chn

		log.Printf("Relevant tx accepted")
		log.Printf("   %s", hex)
	}
}
