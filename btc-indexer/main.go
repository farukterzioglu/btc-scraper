package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func main() {
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(blockHeight int32, header *wire.BlockHeader, txList []*btcutil.Tx) {
			fmt.Printf("New block. Height: %d, hash: %s\n", blockHeight, header.BlockHash())
		},
		OnRelevantTxAccepted: func(tx []byte) {
			fmt.Printf("New transaction for my address: %s\n", hex.EncodeToString(tx))
		},
		OnTxAcceptedVerbose: func(tx *btcjson.TxRawResult) {
			fmt.Printf("New transaction added to mempool: %s\n", tx.Hash)
		},
	}

	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, _ := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18556",
		Endpoint:     "ws",
		User:         "myuser",
		Pass:         "SomeDecentp4ssw0rd",
		Certificates: certs,
	}
	client, _ := rpcclient.New(connCfg, &ntfnHandlers)

	client.NotifyBlocks()
	client.NotifyNewTransactions(true)

	address1, _ := btcutil.DecodeAddress("SP6pS3u7RbtafiVKB9MAUNSPXaMc9U2YF5", &chaincfg.SimNetParams)
	client.LoadTxFilter(true, []btcutil.Address{address1}, nil)

	client.WaitForShutdown()
}
