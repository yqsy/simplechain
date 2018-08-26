package main

import (
	"testing"
	"math/rand"
	"os"
	"github.com/boltdb/bolt"
	"reflect"
)

func getRandomBit(len int) []byte {
	token := make([]byte, len)
	rand.Read(token)
	return token
}

// a. WalletAddr-B 全部转出,那么这笔交易在存储池中就被删除
func TestSimpleA(t *testing.T) {
	os.Remove("utxoDb.db")

	utxoDb := NewUtxoDb("utxoDb.db")
	defer func() { os.Remove("utxoDb.db") }()

	// 交易1: 创建公钥哈希A
	publicKeyHashA, txIdHash1 := createTxOnlyOut(utxoDb, 10, t)

	// 交易2: 创建公钥哈希B, A => B
	publicKeyHashB, txIdHash2 := createTxWithIn(utxoDb, TxIn{txIdHash1, 0}, 10, nil)

	// 1. A可花费输出为空 2. B可花费输出为10个
	remainAmount, spendableOuts := utxoDb.findSpendableTxOutIdx(publicKeyHashA, 10)
	if remainAmount != 0 {
		t.Fatal("err")
	}
	spendableOutsEqual := make(map[string]int)
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	remainAmount, spendableOuts = utxoDb.findSpendableTxOutIdx(publicKeyHashB, 10)
	if remainAmount != 10 {
		t.Fatal("err")
	}
	spendableOutsEqual = make(map[string]int)
	spendableOutsEqual[string(txIdHash2)] = 0
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	// 1. A余额为空 2. B余额为10
	if utxoDb.getBalance(publicKeyHashA) != 0  || utxoDb.getBalance(publicKeyHashB) != 10{
		t.Fatal("err")
	}

	// 交易数量为 1
	if utxoDb.countTransactions() != 1 {
		t.Fatal("err")
	}

	// 交易3: 创建公钥哈希C, B => C 全部转
	publicKeyHashC, txIdHash3 := createTxWithIn(utxoDb, TxIn{txIdHash2, 0}, 10, nil)

	// 1. B可花费输出为空 2. C可花费输出为10个
	remainAmount, spendableOuts = utxoDb.findSpendableTxOutIdx(publicKeyHashB, 10)
	if remainAmount != 0 {
		t.Fatal("err")
	}
	spendableOutsEqual = make(map[string]int)
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	remainAmount, spendableOuts = utxoDb.findSpendableTxOutIdx(publicKeyHashC, 10)
	if remainAmount != 10 {
		t.Fatal("err")
	}
	spendableOutsEqual = make(map[string]int)
	spendableOutsEqual[string(txIdHash3)] = 0
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	// 1. B余额为空 2. C余额为10
	if utxoDb.getBalance(publicKeyHashB) != 0  || utxoDb.getBalance(publicKeyHashC) != 10{
		t.Fatal("err")
	}

	// 交易数量为 1
	if utxoDb.countTransactions() != 1 {
		t.Fatal("err")
	}
}

// 创造一个输出 (依赖in)
func createTxWithIn(utxoDb *UtxoDb, txIn TxIn, transferAmount int, backOut *TxOut /*返还的*/) (publicKeyHash []byte, txId []byte) {
	txId = getRandomBit(32)
	publicKeyHash = getRandomBit(20)
	txOuts := []TxOut{TxOut{transferAmount, publicKeyHash}}

	if backOut != nil {
		txOuts = append(txOuts, *backOut)
	}

	utxoDb.update([]TxIn{txIn}, txOuts, txId)
	return publicKeyHash, txId
}

// 创造一个输出 (仅有输出)
func createTxOnlyOut(utxoDb *UtxoDb, amount int, t *testing.T) (publicKeyHash []byte, txId []byte) {
	txId = getRandomBit(32)
	publicKeyHash = getRandomBit(20)

	txOuts := []TxOut{TxOut{amount, publicKeyHash}}

	if err := utxoDb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		if err := b.Put(txId, EncodeTxOuts(txOuts)); err != nil {
			t.Fatal(err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return publicKeyHash, txId
}
