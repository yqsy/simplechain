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
	publicKeyHashA := getRandomBit(20)
	txIdHash1 := getRandomBit(32)

	txOutsA := []TxOut{TxOut{10, publicKeyHashA}}
	utxoDb := NewUtxoDb("utxoDb.db")
	defer func() { os.Remove("utxoDb.db") }()

	// 交易1 手动放
	if err := utxoDb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		if err := b.Put(txIdHash1, EncodeTxOuts(txOutsA)); err != nil {
			t.Fatal(err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	// 交易2: A => B
	publicKeyHashB := getRandomBit(20)
	txIdHash2 := getRandomBit(32)

	txInsB := []TxIn{TxIn{txIdHash1, 0}}
	txOutsB := []TxOut{TxOut{10, publicKeyHashB}}

	utxoDb.update(txInsB, txOutsB, txIdHash2)

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
	if utxoDb.getBalance(publicKeyHashA) != 0 {
		t.Fatal("err")
	}

	if utxoDb.getBalance(publicKeyHashB) != 10 {
		t.Fatal("err")
	}

	// 交易数量为 1
	if utxoDb.countTransactions() != 1 {
		t.Fatal("err")
	}
}
