package main

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type BlockNotification struct {
	BlockHeight int32
	Header      *wire.BlockHeader
	TxList      []*btcutil.Tx
}

type TxNotification struct {
	Tx *btcjson.TxRawResult
}
