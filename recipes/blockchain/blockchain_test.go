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
		fmt.Printf("=================\n")
		newBlock := blockChain.MineBlock([]*Transaction{})


		fmt.Printf("%v", newBlock)
	}
}
