package blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

type BlockChain struct {
	db *bolt.DB
}

const (
	DbFile      = "blockchain_%s.db"
	BlockBucket = "blocks"
)

// 在blot中创建创世块,并创建BlockChain对象
func CreateBlockChain(nodeId string) *BlockChain {
	dbFile := fmt.Sprintf(DbFile, nodeId)

	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		panic("blockchain file exist")
	}

	// 创世交易块,就不放奖励了
	coinBaseBlock := NewBlock([]byte{}, []*Transaction{NewTransaction()})

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BlockBucket))
		if err != nil {
			panic(err)
		}

		// 创世块存储
		blockData := coinBaseBlock.Serialize()
		if err = b.Put(coinBaseBlock.Hash, blockData); err != nil {
			panic(err)
		}

		// 最新的区块存储
		if err = b.Put([]byte("l"), coinBaseBlock.Hash); err != nil {
			panic(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	blockChain := &BlockChain{db: db}
	return blockChain
}

// 读取已有的区块信息,并创建BlockChain对象
func NewBlockChain(nodeId string) *BlockChain {
	dbFile := fmt.Sprintf(DbFile, nodeId)

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		panic("blockchain file not exist")
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))

		if b == nil {
			panic("err")
		}

		return nil
	}); err != nil {
		panic(err)
	}

	blockChain := &BlockChain{db: db}
	return blockChain
}

func (bc *BlockChain) AddBlock(block *Block) {
	// TODO
}

func (bc *BlockChain) FindTransaction(id []byte) {
	// TODO
}

// 返回 map[交易号]outs数组 (从全量的blocks数据中找出来未花费掉的out的copy)
func (bc *BlockChain) FindUTXO() {
	// TODO
}

// 迭代遍历所有的区块
func (bc *BlockChain) Iterator() {
	// TODO
}

// 获取最新的区块的高度
func (bc *BlockChain) GetBestHeight() {
	// TODO
}

func (bc *BlockChain) GetBlock(blockHash []byte) {
	// TODO
}

func (bc *BlockChain) GetBlockHashes() {
	// TODO
}

func (bc *BlockChain) MineBlock(txs []*Transaction) *Block {
	var lastHash []byte

	if err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		if b == nil {
			panic("err")
		}
		lastHash = b.Get([]byte("l"))
		return nil
	}); err != nil {
		panic(err)
	}

	// 挖矿!
	newBlock := NewBlock(lastHash, txs)

	if err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		if b == nil {
			panic("err")
		}
		if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			panic(err)
		}

		if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
			panic(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return newBlock
}

func (bc *BlockChain) SignTransaction() {
	// TODO
}

func (bc *BlockChain) VerifyTransaction() {
	// TODO
}
