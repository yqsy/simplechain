package blockchain

import (
	"fmt"
	"os"
	"testing"
)

func TestMine(t *testing.T) {
	blockChain := CreateBlockChain("01")

	defer func() {
		os.Remove("blockchain_01.db")
	}()

	for i := 0; i < 10; i++ {
		newBlock := blockChain.MineBlock([]*Transaction{})

		fmt.Printf("nonce: %v hash: %x\n", newBlock.Nonce, newBlock.Hash)
	}
}
