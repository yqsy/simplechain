package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/rand"
	"bytes"
	"encoding/gob"
)

type Transaction struct {
	Id     []byte
	TxIns  []TxIn
	TxOuts []TxOut
}

// 生成交易Id时使用
func (tx *Transaction) hash() []byte {
	// 去除 1. signature 2. PublicKey
	txTrimCopy := tx.trimCopy()

	// 去除 3. Id
	txTrimCopy.Id = nil
	return txTrimCopy.serialize()
}

// a. 生成交易Id时使用: 不包括1. signature 2. PublicKey 3. Id
// b. 生成交易证书时使用: 不包括 signature 包括 publicKey(此时publicKey是prePublicKeyHash)
func (tx *Transaction) serialize() []byte {
	var encode bytes.Buffer
	enc := gob.NewEncoder(&encode)
	if err := enc.Encode(tx); err != nil {
		panic(err)
	}
	return encode.Bytes()
}

// 去除in的 1. signature 2. PublicKey
func (tx *Transaction) trimCopy() *Transaction {
	txInsCopy := make([]TxIn, len(tx.TxIns))
	TxOutsCopy := make([]TxOut, len(tx.TxOuts))
	copy(txInsCopy, tx.TxIns)
	copy(TxOutsCopy, tx.TxOuts)

	for idx := range txInsCopy {
		txInsCopy[idx].Signature = nil
		txInsCopy[idx].PublicKey = nil
	}

	idCopy := make([]byte, len(tx.Id))
	copy(idCopy, tx.Id)

	return &Transaction{idCopy, txInsCopy, TxOutsCopy}
}

func (tx *Transaction) signTxs(preTxMap map[string]*Transaction, privateKey *ecdsa.PrivateKey) {
	curTxTrimCopy := tx.trimCopy()
	for idx, in := range curTxTrimCopy.TxIns {
		// PublicKey填为上一笔输出的prevPublicKeyHash,用来做特殊签名
		prevPublicKeyHash := preTxMap[string(in.PrevTxHashId)].TxOuts[in.PrevOutIdx].PublicKeyHash
		curTxTrimCopy.TxIns[idx].PublicKey = prevPublicKeyHash

		// 私钥签名(hash(证书)) -> 签名

		// 证书
		certificate := curTxTrimCopy.serialize()

		// hash
		sha256Sum := sha256.Sum256([]byte(certificate))

		// 私钥签名
		r, s, _ := ecdsa.Sign(rand.Reader, privateKey, sha256Sum[:])
		signature := append(r.Bytes(), s.Bytes()...)

		// 填充 // 输出
		tx.TxIns[idx].Signature = signature

		// 回退为nil,不影响其他in的签名认证
		curTxTrimCopy.TxIns[idx].PublicKey = nil
	}
}

func (tx *Transaction) verifyTxs(preTxMap map[string]*Transaction) bool {
	curTxTrimCopy := tx.trimCopy()
	for idx, in := range curTxTrimCopy.TxIns {
		// PublicKey填为上一笔输出的prevPublicKeyHash,用来做特殊签名
		prevPublicKeyHash := preTxMap[string(in.PrevTxHashId)].TxOuts[in.PrevOutIdx].PublicKeyHash
		curTxTrimCopy.TxIns[idx].PublicKey = prevPublicKeyHash

		// 公钥解密(签名) == hash(证书)

		// 证书
		certificate := curTxTrimCopy.serialize()

		// hash
		sha256Sum := sha256.Sum256([]byte(certificate))

		// 签名
		r, s := convertSignatureTors(tx.TxIns[idx].Signature)

		// 公钥
		publicKey := convertBytesToPublicKey(tx.TxIns[idx].PublicKey)

		// 公钥验证
		if !ecdsa.Verify(publicKey, sha256Sum[:], r, s) {
			return false
		}

		// 回退为nil,不影响其他in的签名认证
		curTxTrimCopy.TxIns[idx].PublicKey = nil
	}

	return true
}

func NewTransaction(txIns []TxIn, txOuts []TxOut) *Transaction {
	tx := &Transaction{}
	tx.TxIns = txIns
	tx.TxOuts = txOuts
	tx.Id = tx.hash()
	return tx
}

type TxIn struct {
	// 前一笔交易的: hash值
	PrevTxHashId []byte

	// 前一笔交易的: 输出index
	PrevOutIdx int

	// 本次交易的: 签名
	Signature []byte

	// 1. 本次交易的: 公钥 2: 前一笔交易的: 输出到的公钥哈希(sign,verify时用)
	PublicKey []byte
}

type TxOut struct {
	// 本次交易的: 转账金额
	Amount int

	// 本次交易的: 输出到的公钥哈希
	PublicKeyHash []byte
}
