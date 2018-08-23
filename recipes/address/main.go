package main

import (
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"github.com/yqsy/simplechain/base58"
	"fmt"
)

const (
	VersionLen  = 1
	CheckSumLen = 4
)

var (
	Version = []byte{0}
)

// 生成一对公私钥
// Elliptic Curve Digital Signature Algorithm
// https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm
func generateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	p256 := elliptic.P256()
	if privateKey, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
		panic(err)
	} else {
		return privateKey, &privateKey.PublicKey
	}
}

// 获取公钥的bytes
func getBytes(publicKey *ecdsa.PublicKey) []byte {
	return append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
}

// 生成公钥哈希
func generatePublicKeyHash(publicKey *ecdsa.PublicKey) []byte {
	publicKeyBytes := getBytes(publicKey)
	sha256Sum := sha256.Sum256(publicKeyBytes)
	md := ripemd160.New()
	if _, err := md.Write(sha256Sum[:]); err != nil {
		panic(err)
	}
	ripemd160Sum := md.Sum(nil)
	return ripemd160Sum
}

// 生成公钥哈希的checksum
// 1.版本 2.公钥哈希
func generatePublicKeyHashCheckSum(version []byte, publicKeyHash []byte) []byte {
	payload := append(version, publicKeyHash...)
	firstSha256Sum := sha256.Sum256(payload)
	secondSha256Sum := sha256.Sum256(firstSha256Sum[:])
	return secondSha256Sum[:]
}

// 生成钱包地址
// 1.版本 2.公钥哈希 3.checksum截断值
func generateWalletAddress(version []byte, publicKeyHash []byte, checksumCut []byte) []byte {
	payload := append(append(version, publicKeyHash...), checksumCut...)
	return base58.Base58Encode(payload)
}

// 钱包地址提取公钥哈希
func getPublicKeyFromWalletAddress(walletAddress []byte) []byte {
	walletAddressDecoded := base58.Base58Decode(walletAddress)
	if len(walletAddressDecoded) < VersionLen+CheckSumLen {
		panic("walletAddressDecoded too short")
	}
	return walletAddress[VersionLen : len(walletAddressDecoded)-CheckSumLen]
}

// 这里为了把所有的成员放置在一起
// 实际上wallet只需要 1. privateKey 2. publicKey
type Wallet struct {
	privateKey               *ecdsa.PrivateKey // 私钥
	publicKey                *ecdsa.PublicKey  // 公钥
	privateKeyEncoded        string            // 私钥x509 encode
	publicKeyEncoded         string            // 公钥x509 encode
	publicKeyHash            []byte            // 公钥hash
	publicKeyHashCheckSum    []byte            // [版本+公钥哈希]校验值
	publicKeyHashCheckSumCut []byte            // [版本+公钥哈希]校验值前4字节
	walletAddress            []byte            // 钱包地址
}

func NewWallet() *Wallet {
	wallet := &Wallet{}
	wallet.privateKey, wallet.publicKey = generateKeyPair()
	wallet.privateKeyEncoded, wallet.publicKeyEncoded = encode(wallet.privateKey, wallet.publicKey)
	wallet.publicKeyHash = generatePublicKeyHash(wallet.publicKey)
	wallet.publicKeyHashCheckSum = generatePublicKeyHashCheckSum(Version, wallet.publicKeyHash)
	wallet.publicKeyHashCheckSumCut = wallet.publicKeyHashCheckSum[:CheckSumLen]
	wallet.walletAddress = generateWalletAddress(Version, wallet.publicKeyHash, wallet.publicKeyHashCheckSumCut)
	return wallet
}

func (wallet *Wallet) String() string {
	result := ""
	result += fmt.Sprintf("%v\n\n", wallet.privateKeyEncoded)
	result += fmt.Sprintf("%v\n\n", wallet.publicKeyEncoded)
	result += fmt.Sprintf("version: %v\n", Version)
	result += fmt.Sprintf("publickeyHash: %x\n", wallet.publicKeyHash)
	result += fmt.Sprintf("publicKeyHashCheckSumCut: %x\n", wallet.publicKeyHashCheckSumCut)
	result += fmt.Sprintf("walletAddress: %s\n", wallet.walletAddress)
	return result
}

func main() {
	// 假装3笔输出 作为资金来源
	walletA := NewWallet()
	txPreOutA1Tx := NewTransaction(nil, []TxOut{TxOut{10, walletA.publicKeyHash}})
	txPreOutA2Tx := NewTransaction(nil, []TxOut{TxOut{20, walletA.publicKeyHash}})
	txPreOutA3Tx := NewTransaction(nil, []TxOut{TxOut{30, walletA.publicKeyHash}})

	// 55个币输出到钱包B, 5个币回退输出到钱包A

	walletB := NewWallet()

	// 构建当前交易的3个in和2个out
	txInA1 := TxIn{txPreOutA1Tx.Id, 0, nil, nil}
	txInA2 := TxIn{txPreOutA2Tx.Id, 0, nil, nil}
	txInA3 := TxIn{txPreOutA3Tx.Id, 0, nil, nil}
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
	curTx.signTxs(preTxMap, walletA.privateKey)

	// 验证
	if !curTx.verifyTxs(preTxMap, walletA.publicKey) {
		panic("verify error")
	}
}
