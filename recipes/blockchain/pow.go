package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
)

const (
	MaxNonce = math.MaxInt64
)

// 转换成大端法
func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

type ProofOfWork struct {
	block *Block

	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.DifficultyDegreeBits))

	pow := &ProofOfWork{b, target}
	return pow
}

// 准备需要hash计算的数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.MerkleRootHash,
		intToHex(pow.block.TimeStamp),
		intToHex(pow.block.DifficultyDegreeBits),
		intToHex(nonce),
	}, []byte{}, )

	return data
}

func (pow *ProofOfWork) Run() (int64, []byte) {
	var nonce int64
	var hash [32]byte

	for nonce < MaxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		var hashInt big.Int
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	// TODO 会无解吗?
	if nonce >= MaxNonce {
		panic("unsolvable")
	}

	return nonce, hash[:]
}
