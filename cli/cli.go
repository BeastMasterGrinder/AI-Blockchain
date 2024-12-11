package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/farjad/AI-Blockchain/blockchain"
	"github.com/farjad/AI-Blockchain/utils"
)

type CmdLine struct {
	Blockchain *blockchain.Blockchain
}

func (cli *CmdLine) printUse() {
	fmt.Println("Use:")
	fmt.Println("Add a block")
	fmt.Println("printing")
}

func (cli *CmdLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUse()
		// To initiate a shitdown for the badgerdb to start garbage collecting and shutdown
		runtime.Goexit()
	}
}

func (cli *CmdLine) AddBlock(data string) {
	cli.Blockchain.AddBlock(data)
	fmt.Println("Succesfully! Added a block to the blockchain")
}

func (cli *CmdLine) printBlockchain() {
	iter := cli.Blockchain.Parser()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data of the Block: %s\n", block.Data)

		fmt.Printf("Hash of the Block: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.IsValid()))
		fmt.Println()
	}
}

func (cli *CmdLine) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printBlockchain := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block Data")

	switch os.Args[1] {
	case "add":
		errors := addBlockCmd.Parse(os.Args[2:])
		utils.ErrHandle(errors)
	case "print":
		errors := printBlockchain.Parse(os.Args[2:])
		utils.ErrHandle(errors)
	default:
		cli.printUse()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlock(*addBlockData)
	}
	if printBlockchain.Parsed() {
		cli.printBlockchain()
	}
}
