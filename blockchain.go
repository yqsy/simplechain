package main

import (
	"github.com/boltdb/bolt"
	"encoding/hex"
)

const (
	// TODO
	DbFile              = "simplechain.db"
	BlocksBucket        = "blocks"
	GenesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

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

type BlockChain struct {
	// 最新的区块的hash值
	tip []byte

	db *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
	return bci
}

func (bc *BlockChain) AddBlock(transactions []*Transaction) {
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

	newBlock := NewBlock(transactions, lastHash)

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

// 获取未消费的交易输出
func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {

	var unspentTXs []Transaction

	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		if block == nil {
			break
		}

		for _, tx := range block.Transactions {

			txID := hex.EncodeToString(tx.Id)

		Outputs:

			for outIdx, out := range tx.VOut {

				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if outIdx == spentOut {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx) // TODO 放的应该不是这个数据吧 = = 应该是,因为有address
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.VIn {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.TxId)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Prevout)
					}
				}
			}
		}

	}

	return unspentTXs
}

func (bc *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput

	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.VOut {

			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}

		}
	}

	return UTXOs
}

func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	// [交易号] outIdx
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.Id)

		for outIdx, out := range tx.VOut {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func NewBlockChain(address string) *BlockChain {
	var tip []byte

	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		if b == nil {
			// 创建创世区块
			cbtx := NewCoinbaseTx(address, GenesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)

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
