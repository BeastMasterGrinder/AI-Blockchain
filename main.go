package main

import (
	"os"

	"github.com/farjad/AI-Blockchain/blockchain"
	"github.com/farjad/AI-Blockchain/cli"
)

func main() {
	// make sure that app is closed properly to makesure that the db closes properly
	defer os.Exit(0)
	chain := blockchain.InitBlockChain("Sigma boys Sigma Boys")

	defer chain.Database.Close()

	cli := cli.CmdLine{chain}
	cli.Run()
}
