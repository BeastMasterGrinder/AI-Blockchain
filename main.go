package main

import (
	"fmt"
	"strconv"

	"github.com/farjad/AI-Blockchain/blockchain"
)

func main() {

	chain := blockchain.InitBlockChain()

	chain.AddBlock("1 block")
	chain.AddBlock("2 block")
	chain.AddBlock("3 block")

	for _, block := range chain.Blocks {
		fmt.Println(block.PrevHash)
		fmt.Println(block.Data)

		fmt.Println(block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.IsValid()))
		fmt.Println()
	}

}
