package simplechain

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
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
}

func (b *Block) SetHash() {
	// 时间戳int64转换成字符串再hash? TODO
	timestamp := []byte(strconv.FormatInt(b.TimeStamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		TimeStamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}
	block.SetHash()
	return block
}

