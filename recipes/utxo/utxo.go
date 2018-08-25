package main

import (
	"github.com/boltdb/bolt"
	"bytes"
	"encoding/gob"
)

const (
	utxoBucket = "utxoBucket"
)

// [txId] -> 有效的txOuts(serialized)
type UtxoDb struct {
	db *bolt.DB
}

type TxOut struct {
	// 本次交易的: 转账金额
	Amount int

	// 本次交易的: 输出到的公钥哈希
	PublicKeyHash []byte
}

func (txOut *TxOut) Lock(publicKeyHash []byte) {
	txOut.PublicKeyHash = publicKeyHash
}

func (txOut *TxOut) IsLockedWithPublicKeyHash(publicKeyHash []byte) bool {
	return bytes.Compare(txOut.PublicKeyHash, publicKeyHash) == 0
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

// 寻找满足转账金额(transferAmount)的可花费输出的 [txId]txOutIdx, 返回的remainAmount可能不满足transferAmount
func (utxoDb *UtxoDb) findSpendableOuts(publicKeyHash []byte, transferAmount int) (remainAmount int, spendableOuts map[string]int) {
	if err := utxoDb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		// 累积金额
		remainAmount = 0

		// 可花费输出[txId]txOutIdx
		spendableOuts = make(map[string]int)

		for k, v := c.First(); k != nil; k, v = c.Next() {

			// 满足转账金额了
			if remainAmount >= transferAmount {
				break
			}

			txId := string(k) // 数据库存储是[]byte,和map交互时用string
			txOuts := DecodeTxOuts(v)

			// 防御: 一笔tx的txOuts的范围在[1,2]
			if len(txOuts) < 1 || len(txOuts) > 2 {
				panic("err txOuts")
			}

			for txOutIdx, txOut := range txOuts {
				if txOut.IsLockedWithPublicKeyHash(publicKeyHash) {
					remainAmount += txOut.Amount
					spendableOuts[txId] = txOutIdx
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return remainAmount, spendableOuts
}
