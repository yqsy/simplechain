package main

import (
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"github.com/yqsy/simplechain/base58"
	"fmt"
	"crypto/x509"
	"encoding/pem"
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

func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}

func main() {
	privateKey, publicKey := generateKeyPair()
	privateKeyEncoded, publicKeyEncoded := encode(privateKey, publicKey)
	fmt.Printf("%v \n\n%v\n\n", privateKeyEncoded, publicKeyEncoded)

	fmt.Printf("version: %v\n", Version)

	publicKeyHash := generatePublicKeyHash(publicKey)
	fmt.Printf("publickeyHash: %x\n", publicKeyHash)

	publicKeyHashCheckSum := generatePublicKeyHashCheckSum(Version, publicKeyHash)
	publicKeyHashCheckSumCut := publicKeyHashCheckSum[:CheckSumLen]
	fmt.Printf("publicKeyHashCheckSumCut: %x\n", publicKeyHashCheckSumCut)
	walletAddress := generateWalletAddress(Version, publicKeyHash, publicKeyHashCheckSumCut)
	fmt.Printf("walletAddress: %x\n", walletAddress)
}
