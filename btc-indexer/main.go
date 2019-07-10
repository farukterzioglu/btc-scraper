package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"encoding/hex"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/farukterzioglu/btc-scraper/services"
	"github.com/olivere/elastic/v7"
)

func main() {
	blockChannel := make(chan BlockNotification)
	txChannel := make(chan TxNotification)
	txHexChannel := make(chan string)

	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(blockHeight int32, header *wire.BlockHeader, txList []*btcutil.Tx) {
			notification := BlockNotification{
				BlockHeight: blockHeight,
				Header:      header,
				TxList:      txList,
			}

			select {
			case blockChannel <- notification:
				break
			default:
				log.Println("couldn't sent the block. no consumer.")
			}
		},
		OnFilteredBlockDisconnected: func(int32, *wire.BlockHeader) {},
		OnRelevantTxAccepted: func(tx []byte) {
			select {
			case txHexChannel <- hex.EncodeToString(tx):
				break
			default:
				log.Println("couldn't sent the transaction hex. no consumer.")
			}
		},
		OnTxAccepted: func(hash *chainhash.Hash, amount btcutil.Amount) {},
		OnTxAcceptedVerbose: func(tx *btcjson.TxRawResult) {
			select {
			case txChannel <- TxNotification{Tx: tx}:
				break
			default:
				log.Println("couldn't send the transaction. no consumer.")
			}

		},
	}

	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18556",
		Endpoint:     "ws",
		User:         "myuser",
		Pass:         "SomeDecentp4ssw0rd",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	address1, err := btcutil.DecodeAddress("SRMzZD8Ar1ipDjgkxspmiG3uDhdESPjVvL", &chaincfg.SimNetParams)
	address2, err := btcutil.DecodeAddress("Sa5SQHMZ9on5ZjG13AceP4dzFmycTTzkmq", &chaincfg.SimNetParams)
	if err != nil {
		client.Shutdown()
		log.Fatal(err)
	}

	err = client.NotifyBlocks()
	err = client.NotifyNewTransactions(true)
	err = client.LoadTxFilter(true, []btcutil.Address{address1, address2}, nil)

	if err != nil {
		client.Shutdown()
		log.Fatal(err)
	}

	// Elastic client
	esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	ctx := context.Background()
	info, code, err := esClient.Ping("http://localhost:9200").Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	elasticService := services.NewElasticService("btc", esClient)

	// Notification handler
	handler := NewNotificationHandler(client, elasticService)
	go handler.ConsumeBlocks(blockChannel)
	go handler.ConsumeTx(txChannel)
	go handler.ConsumeRelevantTxHex(txHexChannel)

	// time.AfterFunc(time.Second*10, func() {
	// 	log.Println("Client shutting down...")
	// 	client.Shutdown()
	// 	log.Println("Client shutdown complete.")
	// })

	client.WaitForShutdown()
}
