package main

import (
	"log"
	"os"

	"github.com/farjad/AI-Blockchain/blockchain"
	"github.com/farjad/AI-Blockchain/cli"
	"github.com/farjad/AI-Blockchain/ipfs"
)

func main() {
	defer os.Exit(0)
	
	// Initialize IPFS manager
	ipfsManager, err := ipfs.NewIPFSManager("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize blockchain
	chain := blockchain.InitBlockChain("Sigma boys Sigma Boys")
	defer chain.Database.Close()

	// Initialize CLI with both blockchain and IPFS manager
	cli := cli.CmdLine{
		Blockchain:  chain,
		IPFSManager: ipfsManager,
	}
	cli.Run()
}
