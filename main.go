package main

import (
	"fmt"

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
	}

}
