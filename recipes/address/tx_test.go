package main

import (
	"testing"
)

func TestTransaction(t *testing.T) {
	// 假装3笔输出 作为资金来源
	walletA := NewWallet()
	txPreOutA1Tx := NewTransaction(nil, []TxOut{TxOut{10, walletA.publicKeyHash}})
	txPreOutA2Tx := NewTransaction(nil, []TxOut{TxOut{20, walletA.publicKeyHash}})
	txPreOutA3Tx := NewTransaction(nil, []TxOut{TxOut{30, walletA.publicKeyHash}})

	// 55个币输出到钱包B, 5个币回退输出到钱包A

	walletB := NewWallet()

	// 构建当前交易的3个in和2个out
	txInA1 := TxIn{txPreOutA1Tx.Id, 0, nil, walletA.publicKeyBytes}
	txInA2 := TxIn{txPreOutA2Tx.Id, 0, nil, walletA.publicKeyBytes}
	txInA3 := TxIn{txPreOutA3Tx.Id, 0, nil, walletA.publicKeyBytes}
	txOutB1 := TxOut{55, walletB.publicKeyHash}
	txOutA1 := TxOut{5, walletA.publicKeyHash}

	// 当前的交易
	curTx := NewTransaction([]TxIn{txInA1, txInA2, txInA3}, []TxOut{txOutB1, txOutA1})

	// [txId]交易引用
	preTxMap := make(map[string]*Transaction)
	preTxMap[string(txPreOutA1Tx.Id)] = txPreOutA1Tx
	preTxMap[string(txPreOutA2Tx.Id)] = txPreOutA2Tx
	preTxMap[string(txPreOutA3Tx.Id)] = txPreOutA3Tx

	// 签名
	curTx.signTxs(preTxMap,walletA.privateKey)

	// 验证
	if !curTx.verifyTxs(preTxMap) {
		t.Fatal("err")
	}
}
