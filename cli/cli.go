package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/farjad/AI-Blockchain/blockchain"
	"github.com/farjad/AI-Blockchain/ipfs"
	"github.com/farjad/AI-Blockchain/network"
	"github.com/farjad/AI-Blockchain/utils"
)

type CmdLine struct {
	Blockchain  *blockchain.Blockchain
	IPFSManager *ipfs.IPFSManager
}

func (cli *CmdLine) createTransaction(algorithmPath, inputPath string) {
	tx, err := blockchain.NewTransaction(algorithmPath, inputPath, cli.IPFSManager)
	if err != nil {
		fmt.Printf("Failed to create transaction: %v\n", err)
		runtime.Goexit()
	}

	// Add transaction to a new block
	block := blockchain.CreateBlock([]*blockchain.Transaction{tx}, cli.Blockchain.LastHash, cli.Blockchain.GetBestHeight()+1)
	cli.Blockchain.AddBlock(block)
	fmt.Println("Transaction added successfully!")
}

func (cli *CmdLine) printUse() {
	fmt.Println("Use:")
	fmt.Println("  create-tx -algorithm ALGORITHM_PATH -input INPUT_PATH - Create a new transaction")
	fmt.Println("  print - Print the blockchain")
}

func (cli *CmdLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUse()
		// To initiate a shitdown for the badgerdb to start garbage collecting and shutdown
		runtime.Goexit()
	}
}

func (cli *CmdLine) printBlockchain() {
	iter := cli.Blockchain.Parser()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data of the Block: %s\n", block.Transactions)

		fmt.Printf("Hash of the Block: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.IsValid()))
		fmt.Println()
	}
}
func (cli *CmdLine) startNode(nodeID string) {
	network.StartServer(nodeID, "gagaga")
	fmt.Printf("Node %s is up and running\n", nodeID)

}

func (cli *CmdLine) Run() {
	cli.validateArgs()

	// get node id as a flag in all
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("NODE_ID env is not set!")
		runtime.Goexit()
	}

	createTxCmd := flag.NewFlagSet("create-tx", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("start", flag.ExitOnError)

	// get node id as a flag to start the node

	// Create transaction flags
	algorithmPath := createTxCmd.String("algorithm", "", "Path to the algorithm file")
	inputPath := createTxCmd.String("input", "", "Path to the input data file")

	switch os.Args[1] {
	case "create-tx":
		err := createTxCmd.Parse(os.Args[2:])
		utils.ErrHandle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		utils.ErrHandle(err)
	case "start":
		err := startNodeCmd.Parse(os.Args[2:])
		utils.ErrHandle(err)
	default:
		cli.printUse()
		runtime.Goexit()
	}

	if createTxCmd.Parsed() {
		if *algorithmPath == "" || *inputPath == "" {
			createTxCmd.Usage()
			runtime.Goexit()
		}
		cli.createTransaction(*algorithmPath, *inputPath)
	}

	if printChainCmd.Parsed() {
		cli.printBlockchain()
	}

	if startNodeCmd.Parsed() {
		cli.startNode(nodeID)
	}
}
