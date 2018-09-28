package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/yqsy/simplechain/recipes/merkletree"
	"math/rand"
	"time"
)

const (
	// max = 256,越大左移的越少,难度也就越大
	DifficultyDegreeBits = 20
)

func getTxsMerkleRoot(txs []*Transaction) []byte {
	txSigs := make([][]byte, 0)

	for i := 0; i < len(txs); i++ {
		txSigs = append(txSigs, txs[i].Id)
	}

	tree := merkletree.NewTree(txSigs)

	if tree == nil {
		return []byte{}
	}

	return tree.GetRoot().GetSig()
}

func getRandomBit(len int) []byte {
	token := make([]byte, len)
	rand.Read(token)
	return token
}

type Block struct {
	// 版本
	Version int32

	// 上一区块的hash
	PrevBlockHash []byte

	// 交易的merkleRoot
	MerkleRootHash []byte

	// 出块时的时间戳
	TimeStamp int64

	// 当前挖矿的难度,越小,难度越大
	DifficultyDegreeBits int64

	// 随机数
	Nonce int64

	// 本块计算出来的哈希 (比特币中没有这项)
	Hash []byte

	// 本块的高度 (比特币中没有这项)
	Height int32

	// 交易数据(全节点和矿工才会存储)
	Txs []*Transaction
}

func NewBlock(prevBlockHash []byte, height int32, txs []*Transaction) *Block {
	merkleRootHash := getTxsMerkleRoot(txs)


	block := &Block{
		Version:              0,
		PrevBlockHash:        prevBlockHash,
		MerkleRootHash:       merkleRootHash,
		TimeStamp:            time.Now().Unix(),    // TODO 是计算hash前就生成时间的吗?
		DifficultyDegreeBits: DifficultyDegreeBits, // TODO 固定难度. f(全网算力,难度) = 求解时间. 随着算力的增长,难度必须同时增长,才能保证出块时间稳定.
		Nonce:                0,
		Height:               height,
		Txs:                  txs,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(block); err != nil {
		panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	if err := decoder.Decode(&block); err != nil {
		panic(err)
	}
	return &block
}

func (block *Block) String() string {
	result := ""
	result += fmt.Sprintf("Version: %v\n", block.Version)
	result += fmt.Sprintf("PrevBlockHash: %x\n", block.PrevBlockHash)
	result += fmt.Sprintf("MerkleRootHash: %x\n", block.MerkleRootHash)
	result += fmt.Sprintf("TimeStamp: %v\n", block.TimeStamp)
	result += fmt.Sprintf("DifficultyDegreeBits: %v\n", block.DifficultyDegreeBits)
	result += fmt.Sprintf("Nonce: %v\n", block.Nonce)
	result += fmt.Sprintf("Height: %v\n", block.Height)
	result += fmt.Sprintf("txs: %v\n", len(block.Txs))
	return result
}
