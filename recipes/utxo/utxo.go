package main

import (
	"github.com/boltdb/bolt"
	"bytes"
	"encoding/gob"
	"os"
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

type TxIn struct {
	// 前一笔交易的: hash值
	PrevTxHashId []byte

	// 前一笔交易的: 输出index
	PrevOutIdx int
}

// 锁定 (在txOut中放置公钥哈希)
func (txOut *TxOut) Lock(publicKeyHash []byte) {
	txOut.PublicKeyHash = publicKeyHash
}

// 判断是否锁定 (判断txOut中的公钥哈希是否是指定公钥哈希)
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

// 1. 没有 -> 创建数据库 2. 使用
func NewUtxoDb(fileName string) *UtxoDb {
	existsOriginal := false
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		existsOriginal = true
	}

	if db, err := bolt.Open(fileName, 0600, nil); err != nil {
		panic(err)
	} else {
		if !existsOriginal {
			if err := db.Update(func(tx *bolt.Tx) error {
				_, err := tx.CreateBucket([]byte(utxoBucket))
				if err != nil {
					panic(err)
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}

		return &UtxoDb{db: db}
	}
}

func (utxoDb *UtxoDb) defenceTxOuts(txOuts []TxOut) {
	// 防御: 一笔tx的txOuts的数量范围在[1,2]
	if len(txOuts) < 1 || len(txOuts) > 2 {
		panic("err txOuts")
	}

	// 防御: 一笔tx的txOuts的输出地址只可能有一个
	if len(txOuts) == 2 {
		if bytes.Compare(txOuts[0].PublicKeyHash, txOuts[1].PublicKeyHash) == 0 {
			panic("err txOuts")
		}
	}
}

// 寻找满足转账金额(transferAmount)的可花费输出的 [txId]txOutIdx, 返回的remainAmount可能不满足transferAmount
func (utxoDb *UtxoDb) findSpendableTxOutIdx(publicKeyHash []byte, transferAmount int) (remainAmount int, spendableOuts map[string]int) {
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

			utxoDb.defenceTxOuts(txOuts)

			for txOutIdx, txOut := range txOuts {
				if txOut.IsLockedWithPublicKeyHash(publicKeyHash) {
					remainAmount += txOut.Amount
					spendableOuts[txId] = txOutIdx
					break
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return remainAmount, spendableOuts
}

// 获取新的txOuts的txInput (满足转账金额)
func (utxoDb *UtxoDb) getSpendableInputs(publicKeyHash []byte, transferAmount int) *[]TxIn {
	remainAmount, spendableOuts := utxoDb.findSpendableTxOutIdx(publicKeyHash, transferAmount)

	if remainAmount < transferAmount {
		return nil
	}

	result := make([]TxIn, 0)

	for k, v := range spendableOuts {
		txId := []byte(k)
		txOutIdx := v

		result = append(result, TxIn{txId, txOutIdx})
	}

	return &result
}

// 寻找所有的可花费输出, 返回 []TxOut
func (utxoDb *UtxoDb) findAllTxOut(publicKeyHash []byte) []TxOut {
	var txOutResult []TxOut

	if err := utxoDb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txOuts := DecodeTxOuts(v)

			utxoDb.defenceTxOuts(txOuts)

			for _, txOut := range txOuts {
				if txOut.IsLockedWithPublicKeyHash(publicKeyHash) {
					txOutResult = append(txOutResult, txOut)
					break
				}
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return txOutResult
}

// 获取余额
func (utxoDb *UtxoDb) getBalance(publicKeyHash []byte) int {
	// TODO: 其实可以优化的
	balance := 0
	allTxOuts := utxoDb.findAllTxOut(publicKeyHash)
	for _, txOut := range allTxOuts {
		balance += txOut.Amount
	}
	return balance
}

// 得到所有有效tx数量
func (utxoDb *UtxoDb) countTransactions() int {
	counter := 0

	if err := utxoDb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return counter
}

// 更新有效的utxo, TxIn: 用来清除之前的冗余txOut,  TxOut: 产生新的txOut // TODO: 搞不明白为什么Jeiwan/blockchain_go 要过滤掉Coinbase,Coinbase的txOut也是有效的!
func (utxoDb *UtxoDb) update(txIns []TxIn, txOuts []TxOut, txId []byte) {
	if err := utxoDb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		// 去除冗余fromTxOuts
		for _, txIn := range txIns {
			preTxOutsBytes := b.Get(txIn.PrevTxHashId)
			preTxOuts := DecodeTxOuts(preTxOutsBytes)

			utxoDb.defenceTxOuts(preTxOuts)

			updateTxOuts := make([]TxOut, 0)
			// 没有用到的txOut
			for txOutIdx, txOut := range preTxOuts {
				if txOutIdx != txIn.PrevOutIdx {
					updateTxOuts = append(updateTxOuts, txOut)
				}
			}

			// a. 0
			// b,c. 1
			if len(updateTxOuts) < 0 || len(updateTxOuts) > 1 {
				panic("err updateTxOuts")
			}

			// a. 直接删除冗余txOut
			if len(updateTxOuts) == 0 {
				if err := b.Delete(txIn.PrevTxHashId); err != nil {
					panic(err)
				}
			} else {
				// b,c. 去掉一个txOut
				if err := b.Put(txIn.PrevTxHashId, EncodeTxOuts(updateTxOuts)); err != nil {
					panic(err)
				}
			}
		}

		// 产生新的txOut
		if err := b.Put(txId, EncodeTxOuts(txOuts)); err != nil {
			panic(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}
