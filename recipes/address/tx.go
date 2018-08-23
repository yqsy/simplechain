package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/rand"
	"math/big"
	"fmt"
)

type Transaction struct {
	id     []byte
	txIns  []TxIn
	txOuts []TxOut
}

// 去除in的 1. signature 2. prevPublicKeyHash
func (tx *Transaction) TrimCopy() *Transaction {
	txInsCopy := make([]TxIn, len(tx.txIns))
	TxOutsCopy := make([]TxOut, len(tx.txOuts))
	copy(txInsCopy, tx.txIns)
	copy(TxOutsCopy, tx.txOuts)

	for idx := range txInsCopy {
		txInsCopy[idx].signature = nil
		txInsCopy[idx].prevPublicKeyHash = nil
	}

	idCopy := make([]byte, len(tx.id))
	copy(idCopy, tx.id)

	return &Transaction{idCopy, txInsCopy, TxOutsCopy}
}

func (tx *Transaction) signTxs(preTxMap map[string]*Transaction, privateKey *ecdsa.PrivateKey) {
	curTxTrimCopy := tx.TrimCopy()
	for idx, in := range curTxTrimCopy.txIns {
		// 从来源处获得公钥哈希值
		prevPublicKeyHash := preTxMap[string(in.prevTxHashId)].txOuts[in.prevOutIdx].publicKeyHash

		// 填充来区分交易
		curTxTrimCopy.txIns[idx].prevPublicKeyHash = prevPublicKeyHash

		// 私钥签名(hash(证书)) -> 签名

		// 证书
		certificate := fmt.Sprintf("%x\n", curTxTrimCopy)

		// hash
		sha256Sum := sha256.Sum256([]byte(certificate))

		// 私钥签名
		r, s, _ := ecdsa.Sign(rand.Reader, privateKey, sha256Sum[:])
		signature := append(r.Bytes(), s.Bytes()...)

		// 填充 // 输出
		tx.txIns[idx].signature = signature

		// 回退copy的prevPublicKeyHash
		curTxTrimCopy.txIns[idx].prevPublicKeyHash = nil
	}
}

func (tx *Transaction) verifyTxs(preTxMap map[string]*Transaction, publicKey *ecdsa.PublicKey) bool {
	curTxTrimCopy := tx.TrimCopy()

	for idx, in := range curTxTrimCopy.txIns {
		// 从来源处获得公钥哈希值
		prevPublicKeyHash := preTxMap[string(in.prevTxHashId)].txOuts[in.prevOutIdx].publicKeyHash

		// 填充来区分交易
		curTxTrimCopy.txIns[idx].prevPublicKeyHash = prevPublicKeyHash

		// 公钥解密(签名) == hash(证书)

		// 证书
		certificate := fmt.Sprintf("%x\n", curTxTrimCopy)

		// hash
		sha256Sum := sha256.Sum256([]byte(certificate))

		// 签名
		r := big.Int{}
		s := big.Int{}
		sigLen := len(tx.txIns[idx].signature)
		r.SetBytes(tx.txIns[idx].signature[:(sigLen / 2)])
		s.SetBytes(tx.txIns[idx].signature[(sigLen / 2):])

		// 公钥验证
		if !ecdsa.Verify(publicKey, sha256Sum[:], &r, &s) {
			return false
		}

		// 回退copy的prevPublicKeyHash
		curTxTrimCopy.txIns[idx].prevPublicKeyHash = nil
	}

	return true
}

func NewTransaction(txIns []TxIn, txOuts []TxOut) *Transaction {
	tx := &Transaction{}
	tx.id = getRandomBit(160)
	tx.txIns = txIns
	tx.txOuts = txOuts
	return tx
}

type TxIn struct {
	// 前一笔交易的: hash值
	prevTxHashId []byte

	// 前一笔交易的: 输出index
	prevOutIdx int

	// 本次交易的: 签名
	signature []byte

	// 前一笔交易的: 输出到的公钥哈希
	prevPublicKeyHash []byte
}

type TxOut struct {
	// 本次交易的: 转账金额
	amount int

	// 本次交易的: 输出到的公钥哈希
	publicKeyHash []byte
}
