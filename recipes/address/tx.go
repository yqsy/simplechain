package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/rand"
	"math/big"
	"bytes"
	"encoding/gob"
)

type Transaction struct {
	Id     []byte
	TxIns  []TxIn
	TxOuts []TxOut
}

// 不包含in的 1. signature 2. prevPublicKeyHash
// 生成交易时使用
func (tx *Transaction) hash() []byte {
	txTrimCopy := tx.trimCopy()
	txTrimCopy.Id = nil
	return txTrimCopy.serialize()
}

func (tx *Transaction) serialize() []byte {
	var encode bytes.Buffer
	enc := gob.NewEncoder(&encode)
	if err := enc.Encode(tx); err != nil {
		panic(err)
	}
	return encode.Bytes()
}

// 去除in的 1. signature 2. prevPublicKeyHash
func (tx *Transaction) trimCopy() *Transaction {
	txInsCopy := make([]TxIn, len(tx.TxIns))
	TxOutsCopy := make([]TxOut, len(tx.TxOuts))
	copy(txInsCopy, tx.TxIns)
	copy(TxOutsCopy, tx.TxOuts)

	for idx := range txInsCopy {
		txInsCopy[idx].Signature = nil
		txInsCopy[idx].PrevPublicKeyHash = nil
	}

	idCopy := make([]byte, len(tx.Id))
	copy(idCopy, tx.Id)

	return &Transaction{idCopy, txInsCopy, TxOutsCopy}
}

func (tx *Transaction) signTxs(preTxMap map[string]*Transaction, privateKey *ecdsa.PrivateKey) {
	curTxTrimCopy := tx.trimCopy()
	for idx, in := range curTxTrimCopy.TxIns {
		// 从来源处获得公钥哈希值
		prevPublicKeyHash := preTxMap[string(in.PrevTxHashId)].TxOuts[in.PrevOutIdx].PublicKeyHash

		// 填充来区分交易
		curTxTrimCopy.TxIns[idx].PrevPublicKeyHash = prevPublicKeyHash

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

		// 回退copy的prevPublicKeyHash
		curTxTrimCopy.TxIns[idx].PrevPublicKeyHash = nil
	}
}

func (tx *Transaction) verifyTxs(preTxMap map[string]*Transaction, publicKey *ecdsa.PublicKey) bool {
	curTxTrimCopy := tx.trimCopy()

	for idx, in := range curTxTrimCopy.TxIns {
		// 从来源处获得公钥哈希值
		prevPublicKeyHash := preTxMap[string(in.PrevTxHashId)].TxOuts[in.PrevOutIdx].PublicKeyHash

		// 填充来区分交易
		curTxTrimCopy.TxIns[idx].PrevPublicKeyHash = prevPublicKeyHash

		// 公钥解密(签名) == hash(证书)

		// 证书
		certificate := curTxTrimCopy.serialize()

		// hash
		sha256Sum := sha256.Sum256([]byte(certificate))

		// 签名
		r := big.Int{}
		s := big.Int{}
		sigLen := len(tx.TxIns[idx].Signature)
		r.SetBytes(tx.TxIns[idx].Signature[:(sigLen / 2)])
		s.SetBytes(tx.TxIns[idx].Signature[(sigLen / 2):])

		// 公钥验证
		if !ecdsa.Verify(publicKey, sha256Sum[:], &r, &s) {
			return false
		}

		// 回退copy的prevPublicKeyHash
		curTxTrimCopy.TxIns[idx].PrevPublicKeyHash = nil
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

	// 前一笔交易的: 输出到的公钥哈希
	PrevPublicKeyHash []byte
}

type TxOut struct {
	// 本次交易的: 转账金额
	Amount int

	// 本次交易的: 输出到的公钥哈希
	PublicKeyHash []byte
}
