package main

import "log"

type NotificationHandler struct {
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (handler *NotificationHandler) ConsumeBlocks(chn <-chan BlockNotification) *NotificationHandler {
	log.Println("Started to consume blocks...")

	for {
		block := <-chn

		log.Printf("Block connected: %s, height: %d, tx count: %d", block.Header.BlockHash().String(), block.BlockHeight, len(block.TxList))

		for _, tx := range block.TxList {
			log.Printf("   Transaction: %s", tx.Hash().String())
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
