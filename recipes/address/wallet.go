package main

import (
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"fmt"

	"golang.org/x/crypto/ripemd160"
	"github.com/yqsy/simplechain/base58"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	VersionLen  = 1
	CheckSumLen = 4
)

var (
	P2PKHVersion = []byte{0}
	WIFVersion   = []byte{0x80}
)

// 生成一对公私钥 // Elliptic Curve Digital Signature Algorithm // https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm
func generateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	p256 := elliptic.P256()
	if privateKey, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
		panic(err)
	} else {
		return privateKey, &privateKey.PublicKey
	}
}

// 获取公钥的bytes
func getPublicKeyBytes(publicKey *ecdsa.PublicKey) []byte {
	return append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
}

// 获取私钥的bytes
func getPrivateKeyBytes(privateKey *ecdsa.PrivateKey) []byte {
	return privateKey.D.Bytes()
}

// 生成公钥bytes的哈希
func generatePublicKeyHash(publicKeyBytes []byte) []byte {
	sha256Sum := sha256.Sum256(publicKeyBytes)
	md := ripemd160.New()
	if _, err := md.Write(sha256Sum[:]); err != nil {
		panic(err)
	}
	ripemd160Sum := md.Sum(nil)
	return ripemd160Sum
}

// version + payload 生成checksum
func double256Sum(version []byte, payload []byte) []byte {
	versionPayload := append(version, payload...)
	firstSha256Sum := sha256.Sum256(versionPayload)
	secondSha256Sum := sha256.Sum256(firstSha256Sum[:])
	return secondSha256Sum[:]
}

// version + payload + checksum截断值 生成地址
func base58PayloadWithVersionAndChecksumCut(version []byte, payload []byte, checksumCut []byte) []byte {
	versionPayloadChecksumCut := append(append(version, payload...), checksumCut...)
	return base58.Base58Encode(versionPayloadChecksumCut)
}

// version + payload 生成地址
func convertVersionAndPayloadToAddress(version []byte, payload []byte) []byte {
	checkSum := double256Sum(version, payload)
	checkSumCut := checkSum[:CheckSumLen]
	return base58PayloadWithVersionAndChecksumCut(version, payload, checkSumCut)
}

// 地址提取payload
func getPayloadFromAddress(address []byte) []byte {
	walletAddressDecoded := base58.Base58Decode(address)
	if len(walletAddressDecoded) < VersionLen+CheckSumLen {
		panic("walletAddressDecoded too short")
	}
	return address[VersionLen : len(walletAddressDecoded)-CheckSumLen]
}

// privateKeyBytes -> ecdsa.PrivateKey
func convertBytesToPrivateKey(d []byte) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = crypto.S256()

	if 8*len(d) != priv.Params().BitSize {
		panic("invalid length")
	}

	priv.D = new(big.Int).SetBytes(d)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		panic("invalid private key")
	}
	return priv
}

// publicKeyBytes -> ecdsa.PublicKey
func convertBytesToPublicKey(publicKeyBytes []byte) *ecdsa.PublicKey {
	x := big.Int{}
	y := big.Int{}
	keyLen := len(publicKeyBytes)
	x.SetBytes(publicKeyBytes[:(keyLen / 2)])
	y.SetBytes(publicKeyBytes[(keyLen / 2):])
	curve := elliptic.P256()
	return &ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
}

// signature -> r,s
func convertSignatureTors(signature []byte) (*big.Int, *big.Int) {
	r := &big.Int{}
	s := &big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])
	return r, s
}

type Wallet struct {
	privateKey      *ecdsa.PrivateKey // 私钥
	publicKey       *ecdsa.PublicKey  // 公钥
	privateKeyBytes []byte            // 私钥二进制
	publicKeyBytes  []byte            // 公钥二进制
	publicKeyHash   []byte            // 公钥hash
	walletAddress   []byte            // 钱包地址
}

func NewWallet() *Wallet {
	wallet := &Wallet{}
	wallet.privateKey, wallet.publicKey = generateKeyPair()
	wallet.privateKeyBytes = getPrivateKeyBytes(wallet.privateKey)
	wallet.publicKeyBytes = getPublicKeyBytes(wallet.publicKey)
	wallet.publicKeyHash = generatePublicKeyHash(wallet.publicKeyBytes)
	wallet.walletAddress = convertVersionAndPayloadToAddress(P2PKHVersion, wallet.publicKeyHash)
	return wallet
}

func (wallet *Wallet) String() string {
	result := ""
	result += fmt.Sprintf("privateKeyBytes: %x\n", wallet.privateKeyBytes)
	result += fmt.Sprintf("publicKeyBytes: %x\n", wallet.publicKeyBytes)
	result += fmt.Sprintf("publickeyHash: %x\n", wallet.publicKeyHash)
	result += fmt.Sprintf("walletAddress: %s\n", wallet.walletAddress)
	return result
}
