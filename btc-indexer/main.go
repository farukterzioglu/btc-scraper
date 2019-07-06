package main

import (
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
)

func main() {
	type BlockNotification struct {
		BlockHeight int32
		Header      *wire.BlockHeader
		TxList      []*btcutil.Tx
	}

	blockChannel := make(chan BlockNotification)
	go func(chn chan BlockNotification) {
		for {
			block := <-chn

			log.Printf("Block connected: %s, height: %d, tx count: %d", block.Header.BlockHash().String(), block.BlockHeight, len(block.TxList))

			for _, tx := range block.TxList {
				log.Printf("   Transaction: %s", tx.Hash().String())
			}
		}
	}(blockChannel)

	type TxNotification struct {
		Tx *btcjson.TxRawResult
	}

	txChannel := make(chan TxNotification)
	go func(chn chan TxNotification) {
		for {
			txNtf := <-chn
			tx := txNtf.Tx

			log.Printf("Tx accepted (verbose): %s (Conf.: %d)", tx.Hash, tx.Confirmations)

			for _, vout := range tx.Vout {
				log.Printf("   Output amount: %f", vout.Value)
			}
		}
	}(txChannel)

	txHexChannel := make(chan string)
	go func(chn chan string) {
		for {
			hex := <-chn

			log.Printf("Relevant tx accepted")
			log.Printf("   %s", hex)
		}
	}(txHexChannel)

	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(blockHeight int32, header *wire.BlockHeader, txList []*btcutil.Tx) {
			blockChannel <- BlockNotification{
				BlockHeight: blockHeight,
				Header:      header,
				TxList:      txList,
			}
		},
		OnFilteredBlockDisconnected: func(int32, *wire.BlockHeader) {},
		OnRelevantTxAccepted: func(tx []byte) {
			txHexChannel <- hex.EncodeToString(tx)
		},
		OnTxAccepted: func(hash *chainhash.Hash, amount btcutil.Amount) {},
		OnTxAcceptedVerbose: func(tx *btcjson.TxRawResult) {
			txChannel <- TxNotification{
				Tx: tx,
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

	// time.AfterFunc(time.Second*10, func() {
	// 	log.Println("Client shutting down...")
	// 	client.Shutdown()
	// 	log.Println("Client shutdown complete.")
	// })

	client.WaitForShutdown()
}
