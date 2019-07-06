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
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(blockHeight int32, header *wire.BlockHeader, txList []*btcutil.Tx) {
			log.Printf("Block connected: %s, height: %d, tx count: %d", header.BlockHash().String(), blockHeight, len(txList))

			for _, tx := range txList {
				log.Printf("   Transaction: %s", tx.Hash().String())
			}
		},
		OnFilteredBlockDisconnected: func(int32, *wire.BlockHeader) {},
		OnRecvTx: func(tx *btcutil.Tx, block *btcjson.BlockDetails) {
			log.Printf("Transaction received : %v", tx.Hash, block.Height)
		},
		OnRedeemingTx: func(tx *btcutil.Tx, block *btcjson.BlockDetails) {
			log.Printf("Transaction redeemed : %v", tx.Hash, block.Height)
		},
		OnRelevantTxAccepted: func(tx []byte) {
			log.Printf("Relevant tx accepted")
			log.Printf("   %s", hex.EncodeToString(tx))
		},
		OnTxAccepted: func(hash *chainhash.Hash, amount btcutil.Amount) {
			log.Printf("Tx accepted: %s", hash.String())
		},
		OnTxAcceptedVerbose: func(tx *btcjson.TxRawResult) {
			log.Printf("Tx accepted (verbose): %s (Conf.: %d)", tx.Hash, tx.Confirmations)

			for _, vout := range tx.Vout {
				log.Printf("   Output amount: %f", vout.Value)
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
