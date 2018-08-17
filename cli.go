package main

import (
	"os"
	"fmt"
	"flag"
	"strconv"
)

type Cli struct {
	bc *BlockChain
}

var usage = `Usage:
addblock - Add block to the block chain
printchain - Print block chain
`

func (cli *Cli) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsageAndExit()
	}
}

func (cli *Cli) printUsageAndExit() {
	fmt.Println(usage)
	os.Exit(1)
}

func (cli *Cli) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")
	createBlockChainAddr := createBlockChainCmd.String("address", "value", "Genesis Block Chain coin address")

	_ = createBlockChainAddr

	switch os.Args[1] {
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
	case "createblockchain":
		if err := createBlockChainCmd.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
	default:
		cli.printUsageAndExit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	} else if printChainCmd.Parsed() {
		cli.printChain()
	} else if createBlockChainCmd.Parsed() {
		cli.createBlockChain()
	}

}

func (cli *Cli) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("success")
}

func (cli *Cli) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()
		if block == nil {
			break
		}

		fmt.Print(block.String())
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
	}
}

func (cli *Cli) createBlockChain() {
	// TODO
}

func (cli *Cli) getBalance(address string) {
	bc := NewBlockChain(address)
	defer bc.db.Close()

	balance := 0

	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
