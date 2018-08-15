package simplechain

import (
	"math/big"
	"bytes"
	"fmt"
	"math"
	"crypto/sha256"
)

var (
	maxNonce = math.MaxInt64
)

// 值越大 -> target 越小 -> 难度越大
const targetBits = 24

type ProofOfWork struct {
	block *Block

	// 每次算的hash都要比这个值小,数字越小难度越大
	target *big.Int
}

func NewProofWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.TimeStamp), // 时间戳int64转换成大端法再hash? TODO
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		// 转换成bit.Int
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			// 找到区块了
			break
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")
	return nonce, hash[:]
}
