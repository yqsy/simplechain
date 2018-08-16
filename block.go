package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"fmt"
	"crypto/sha256"
)

type Block struct {
	// 区块创建的时间戳
	TimeStamp int64

	// 交易数据
	Transactions []*Transaction

	// 先前块的hash
	PrevBlockHash []byte

	// 当前区块的hash
	Hash []byte

	// 随机数 (放在一起hash计算)
	Nonce int
}

func (b *Block) String() string {
	var result string
	result += fmt.Sprintf("Prev.Hash: %x\n", b.PrevBlockHash)
	// result += fmt.Sprintf("Data: %s\n", b.Data) TODO
	result += fmt.Sprintf("TimeStamp: %s\n", time.Unix(b.TimeStamp, 0).Format(time.RFC822))
	result += fmt.Sprintf("Hash: %x\n", b.Hash)
	result += fmt.Sprintf("Nonce: %d\n", b.Nonce)
	return result
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	return result.Bytes(), err
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Id)
	}

	joined := bytes.Join(txHashes, []byte{})
	txHash = sha256.Sum256(joined)
	return txHash[:]
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	return &block, err
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		TimeStamp:     time.Now().Unix(),
		Transactions:  transactions,
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

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
