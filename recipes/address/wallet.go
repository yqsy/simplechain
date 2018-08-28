package main

import (
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"github.com/yqsy/simplechain/base58"
	"fmt"
	"math/big"
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

// 生成公钥哈希
func generatePublicKeyHash(publicKey *ecdsa.PublicKey) []byte {
	publicKeyBytes := getPublicKeyBytes(publicKey)
	sha256Sum := sha256.Sum256(publicKeyBytes)
	md := ripemd160.New()
	if _, err := md.Write(sha256Sum[:]); err != nil {
		panic(err)
	}
	ripemd160Sum := md.Sum(nil)
	return ripemd160Sum
}

// 生成checksum // 1.版本 2.payload
func double256Sum(version []byte, payload []byte) []byte {
	versionPayload := append(version, payload...)
	firstSha256Sum := sha256.Sum256(versionPayload)
	secondSha256Sum := sha256.Sum256(firstSha256Sum[:])
	return secondSha256Sum[:]
}

// 生成地址 // 1.版本 2.payload 3.checksum截断值
func base58PayloadWithVersionAndChecksumCut(version []byte, payload []byte, checksumCut []byte) []byte {
	versionPayloadChecksumCut := append(append(version, payload...), checksumCut...)
	return base58.Base58Encode(versionPayloadChecksumCut)
}

// 1. 生成checksum 2. base58 (version,payload,checkSumCut)
func convertVersionAndPayloadToAddress(version []byte, payload []byte) []byte{
	checkSum := double256Sum(version, payload)
	checkSumCut := checkSum[:CheckSumLen]
	return base58PayloadWithVersionAndChecksumCut(version, payload, checkSumCut)
}

// 钱包地址提取公钥哈希
func getPublicKeyFromWalletAddress(walletAddress []byte) []byte {
	walletAddressDecoded := base58.Base58Decode(walletAddress)
	if len(walletAddressDecoded) < VersionLen+CheckSumLen {
		panic("walletAddressDecoded too short")
	}
	return walletAddress[VersionLen : len(walletAddressDecoded)-CheckSumLen]
}

// ecdsa.PublicKey -> publicKeyBytes
func convertPublicKeyToBytes(publicKey *ecdsa.PublicKey) []byte {
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	return publicKeyBytes
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

// 这里为了把所有的成员放置在一起 // 实际上wallet只需要 1. privateKey 2. publicKey
type Wallet struct {
	// --- 我自己实现的部分
	privateKey        *ecdsa.PrivateKey // 私钥
	publicKey         *ecdsa.PublicKey  // 公钥
	privateKeyBytes   []byte            // 私钥二进制数组
	publicKeyBytes    []byte            // 公钥二进制数组
	privateKeyEncoded string            // 私钥x509 encode
	publicKeyEncoded  string            // 公钥x509 encode
	publicKeyHash     []byte            // 公钥hash
	walletAddress     []byte            // 钱包地址

	// --- 比特币
	btcPrivateKeyBytes    []byte // 私钥
	btcPrivateKeyWif      []byte // 私钥的Wallet import format  (5 前缀)
	btcPublicKeyBytes     []byte // 公钥 (04 前缀)
	btcWalletAddressP2PKH []byte // P2PKH钱包地址 (1 前缀)
	btcWalletAddressP2SH  []byte // P2SH钱包地址 (3 前缀)

	// --- 压缩
	btcPrivateKeyBytesCompressed    []byte // 私钥 (01 末尾)
	btcPrivateKeyWifCompressed      []byte // 私钥的Wallet import format (K,L前缀)
	btcPublicKeyBytesCompressed     []byte // 公钥 (02 03 前缀)
	btcWalletAddressP2PKHCompressed []byte // P2PKH钱包地址 (1 前缀)
	btcWalletAddressP2SHCompressed  []byte // P2SH钱包地址 (3 前缀)
}

func NewWallet() *Wallet {
	wallet := &Wallet{}

	// --- 我自己实现的部分
	wallet.privateKey, wallet.publicKey = generateKeyPair()
	wallet.privateKeyEncoded, wallet.publicKeyEncoded = encode(wallet.privateKey, wallet.publicKey)
	wallet.privateKeyBytes = getPrivateKeyBytes(wallet.privateKey)
	wallet.publicKeyBytes = getPublicKeyBytes(wallet.publicKey)
	wallet.publicKeyHash = generatePublicKeyHash(wallet.publicKey)
	wallet.walletAddress = convertVersionAndPayloadToAddress(P2PKHVersion, wallet.publicKeyHash)

	// --- 比特币
	wallet.btcPrivateKeyBytes = wallet.privateKeyBytes
	wallet.btcPrivateKeyWif = convertVersionAndPayloadToAddress(WIFVersion, wallet.btcPrivateKeyBytes)
	wallet.btcPublicKeyBytes = append([]byte{0, 4}, wallet.publicKeyBytes...)
	wallet.btcWalletAddressP2PKH = wallet.walletAddress


	return wallet
}

func (wallet *Wallet) String() string {
	result := ""
	result += fmt.Sprintf("in my Realization:\n")
	result += fmt.Sprintf("%v\n\n", wallet.privateKeyEncoded)
	result += fmt.Sprintf("%v\n\n", wallet.publicKeyEncoded)
	result += fmt.Sprintf("privateKeyBytes: %x\n", wallet.privateKeyBytes)
	result += fmt.Sprintf("publicKeyBytes: %x\n", wallet.publicKeyBytes)
	result += fmt.Sprintf("publickeyHash: %x\n", wallet.publicKeyHash)
	result += fmt.Sprintf("walletAddress: %s\n", wallet.walletAddress)

	result += fmt.Sprintf("in <bitcoin>:\n")

	return result
}
