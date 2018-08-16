package simplechain

import (
	"time"
	"bytes"
	"encoding/gob"
)

type Block struct {
	// 区块创建的时间戳
	TimeStamp int64

	// 区块存储的数据
	Data []byte

	// 先前块的hash
	PrevBlockHash []byte

	// 当前区块的hash
	Hash []byte

	// 随机数 (放在一起hash计算)
	Nonce int
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	return result.Bytes(), err
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	return &block, err
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		TimeStamp:     time.Now().Unix(),
		Data:          []byte (data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
