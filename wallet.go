package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

const (
	AddressChecksumLen = 4
	Version            = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func (w *Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{Version}, pubKeyHash...)

	checksum := GetChecksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)

	address := Base58Encode(fullPayload)
	return address
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func GetChecksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:AddressChecksumLen]
}
