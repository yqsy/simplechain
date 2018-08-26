package main

import (
	"testing"
	"math/rand"
	"os"
	"github.com/boltdb/bolt"
	"reflect"
)

const (
	DbFileName = "utxoDb.db"
)

func getRandomBit(len int) []byte {
	token := make([]byte, len)
	rand.Read(token)
	return token
}

// a. WalletAddr-B 全部转出,那么这笔交易在存储池中就被删除
func TestSimpleA(t *testing.T) {
	os.Remove(DbFileName)

	utxoDb := NewUtxoDb(DbFileName)
	defer func() { os.Remove(DbFileName) }()

	// 交易1: 创建公钥哈希A
	publicKeyHashA, txIdHash1 := createTxOnlyOut(utxoDb, 10, t)

	// 交易2: 创建公钥哈希B, A => B
	publicKeyHashB, txIdHash2 := createTxWithIn(utxoDb, TxIn{txIdHash1, 0}, 10, nil)

	// 1. A可花费输出为0 2. B可花费输出为10
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

	// 1. A余额为0 2. B余额为10
	if utxoDb.getBalance(publicKeyHashA) != 0 || utxoDb.getBalance(publicKeyHashB) != 10 {
		t.Fatal("err")
	}

	// 交易数量为 1
	if utxoDb.countTransactions() != 1 {
		t.Fatal("err")
	}

	// 交易3: 创建公钥哈希C, B => C
	publicKeyHashC, txIdHash3 := createTxWithIn(utxoDb, TxIn{txIdHash2, 0}, 10, nil)

	// 1. B可花费输出为0 2. C可花费输出为10
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

	// 1. B余额为0 2. C余额为10
	if utxoDb.getBalance(publicKeyHashB) != 0 || utxoDb.getBalance(publicKeyHashC) != 10 {
		t.Fatal("err")
	}

	// 交易数量为 1
	if utxoDb.countTransactions() != 1 {
		t.Fatal("err")
	}
}

// b. WalletAddr-B 部分转出, 那么这笔交易会保留部分
func TestSimpleB(t *testing.T) {
	os.Remove(DbFileName)

	utxoDb := NewUtxoDb(DbFileName)
	defer func() { os.Remove(DbFileName) }()

	// 交易1: 创建公钥哈希A
	publicKeyHashA, txIdHash1 := createTxOnlyOut(utxoDb, 10, t)

	// 交易2: 创建公钥哈希B, A => B (转6余4)
	publicKeyHashB, txIdHash2 := createTxWithIn(utxoDb, TxIn{txIdHash1, 0}, 6,
		&TxOut{4, publicKeyHashA})

	// 1. A 可花费输出为4 2. 可花费输出为6
	remainAmount, spendableOuts := utxoDb.findSpendableTxOutIdx(publicKeyHashA, 4)
	if remainAmount != 4 {
		t.Fatal("err")
	}
	spendableOutsEqual := make(map[string]int)
	spendableOutsEqual[string(txIdHash2)] = 1
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	remainAmount, spendableOuts = utxoDb.findSpendableTxOutIdx(publicKeyHashB, 6)
	if remainAmount != 6 {
		t.Fatal("err")
	}
	spendableOutsEqual = make(map[string]int)
	spendableOutsEqual[string(txIdHash2)] = 0
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	// 1. A余额为 4 2. B余额为 6
	if utxoDb.getBalance(publicKeyHashA) != 4 || utxoDb.getBalance(publicKeyHashB) != 6 {
		t.Fatal("err")
	}

	// 交易数量为 1
	if utxoDb.countTransactions() != 1 {
		t.Fatal("err")
	}

	// 交易3: 创建公钥哈希C, B => C
	publicKeyHashC, txIdHash3 := createTxWithIn(utxoDb, TxIn{txIdHash2, 0}, 6, nil)

	// 1. A可花费输出为4 2.B可花费输出为0 3.C可花费输出为6
	remainAmount, spendableOuts = utxoDb.findSpendableTxOutIdx(publicKeyHashA, 4)
	if remainAmount != 4 {
		t.Fatal("err")
	}
	spendableOutsEqual = make(map[string]int)
	spendableOutsEqual[string(txIdHash2)] = 0
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	remainAmount, spendableOuts = utxoDb.findSpendableTxOutIdx(publicKeyHashB, 6)
	if remainAmount != 0 {
		t.Fatal("err")
	}
	spendableOutsEqual = make(map[string]int)
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	remainAmount, spendableOuts = utxoDb.findSpendableTxOutIdx(publicKeyHashC, 6)
	if remainAmount != 6 {
		t.Fatal("err")
	}
	spendableOutsEqual = make(map[string]int)
	spendableOutsEqual[string(txIdHash3)] = 0
	if !reflect.DeepEqual(spendableOuts, spendableOutsEqual) {
		t.Fatal("err")
	}

	// 1. A余额为4 2. B余额为0 3. C余额为6
	if utxoDb.getBalance(publicKeyHashA) != 4 || utxoDb.getBalance(publicKeyHashB) != 0 || utxoDb.getBalance(publicKeyHashC) != 6 {
		t.Fatal("err")
	}

	// 交易数量为 2
	if utxoDb.countTransactions() != 2 {
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
