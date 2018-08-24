package main

import (
	"github.com/boltdb/bolt"
	"bytes"
	"encoding/gob"
)

const (
	utxoBucket = "utxoBucket"
)

// [txId] -> 有效的outs
type UtxoDb struct {
	db *bolt.DB
}

type TxOut struct {
	// 本次交易的: 转账金额
	Amount int

	// 本次交易的: 输出到的公钥哈希
	PublicKeyHash []byte
}

func EncodeTxOuts(txOuts []TxOut) []byte {
	var encode bytes.Buffer
	enc := gob.NewEncoder(&encode)
	if err := enc.Encode(txOuts); err != nil {
		panic(err)
	}
	return encode.Bytes()
}

func DecodeTxOuts(txOutsBytes []byte) []TxOut {
	var txOuts []TxOut
	dec := gob.NewDecoder(bytes.NewReader(txOutsBytes))
	if err := dec.Decode(&txOuts); err != nil {
		panic(err)
	}
	return txOuts
}

// 1. 没有 -> 创建数据库
// 2. 使用
func NewUtxoDb(fileName string) *UtxoDb {
	if db, err := bolt.Open(fileName, 0600, nil); err != nil {
		panic(err)
	} else {
		return &UtxoDb{db: db}
	}
}

// 寻找满足转账金额(transferAmount)的可花费输出
func (utxoDb *UtxoDb) findSpendableOuts(publicKeyHash []byte, transferAmount int) (remainAmount int, spendableOuts map[string]TxOut) {
	if err := utxoDb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOuts := DecodeTxOuts(v)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}
