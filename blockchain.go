package main

import (
	"github.com/boltdb/bolt"
)

const (
	// TODO
	DbFile       = "simplechain.db"
	BlocksBucket = "blocks"
)

type BlockChain struct {
	// 最新的区块的hash值
	tip []byte

	db *bolt.DB
}

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// 获取当前的block,并向前迭代(重置currentHash)
func (bci *BlockChainIterator) Next() *Block {
	var block *Block

	if len(bci.currentHash) == 0 {
		return nil
	}

	err := bci.db.View(func(tx *bolt.Tx) error {
		// 取出currentHash的结构体
		b := tx.Bucket([]byte(BlocksBucket))

		if b == nil {
			panic("tx.Bucket return nil")
		}

		encodedBlock := b.Get(bci.currentHash)
		blockTemp, err := DeserializeBlock(encodedBlock)
		if err != nil {
			panic(err)
		}
		block = blockTemp
		return nil
	})

	if err != nil {
		panic(err)
	}
	// hash向前迭代
	bci.currentHash = block.PrevBlockHash
	return block
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
	return bci
}

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		// 从数据库更新一下最新的区块的hash值
		b := tx.Bucket([]byte(BlocksBucket))

		if b == nil {
			panic("tx.Bucket return nil")
		}

		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		// 更新最新区块
		b := tx.Bucket([]byte(BlocksBucket))

		if b == nil {
			panic("tx.Bucket return nil")
		}

		serialized, err := newBlock.Serialize()
		if err != nil {
			panic(err)
		}

		if err := b.Put(newBlock.Hash, serialized); err != nil {
			panic(err)
		}

		if err = b.Put([]byte("l"), newBlock.Hash); err != nil {
			panic(err)
		}

		bc.tip = newBlock.Hash
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func NewBlockChain() *BlockChain {
	var tip []byte

	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		if b == nil {
			// 创建创世区块
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(BlocksBucket))
			if err != nil {
				panic(err)
			}
			serialized, err := genesis.Serialize()
			if err != nil {
				panic(err)
			}
			if err = b.Put(genesis.Hash, serialized); err != nil {
				panic(err)
			}
			if err = b.Put([]byte("l"), genesis.Hash); err != nil {
				panic(err)
			}

			tip = genesis.Hash
		} else {
			// 得到最后一个已有的区块
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	bc := BlockChain{
		tip: tip,
		db:  db}
	return &bc
}
