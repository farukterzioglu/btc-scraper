package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func main() {
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(blockHeight int32, header *wire.BlockHeader, txList []*btcutil.Tx) {
			log.Printf("Block connected: %v (%d)", header.BlockHash(), blockHeight)
			log.Printf("Transaction count: %d", len(txList))

			for _, tx := range txList {
				log.Printf("Transaction: %s", tx.Hash().String())
			}
		},
		OnFilteredBlockDisconnected: func(int32, *wire.BlockHeader) {
		},
		OnRecvTx: func(tx *btcutil.Tx, block *btcjson.BlockDetails) {
			log.Printf("Transaction received : %v", tx.Hash, block.Height)
		},
		OnRedeemingTx: func(tx *btcutil.Tx, block *btcjson.BlockDetails) {
			log.Printf("Transaction redeemed : %v", tx.Hash, block.Height)
		},
		OnRelevantTxAccepted: func(tx []byte) {

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

	if err := client.NotifyBlocks(); err != nil {
		client.Shutdown()
		log.Fatal(err)
	}

	address, err := btcutil.DecodeAddress("SRMzZD8Ar1ipDjgkxspmiG3uDhdESPjVvL", &chaincfg.SimNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	client.LoadTxFilter(true, []btcutil.Address{address}, nil)

	// time.AfterFunc(time.Second*10, func() {
	// 	log.Println("Client shutting down...")
	// 	client.Shutdown()
	// 	log.Println("Client shutdown complete.")
	// })

	client.WaitForShutdown()
}
