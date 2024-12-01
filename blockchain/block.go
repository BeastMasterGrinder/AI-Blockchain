package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Hash []byte
	// Transactions []*Transaction
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{Hash: []byte{}, Data: []byte(data), PrevHash: prevHash}
	block.DeriveHash()
	return block
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	block := CreateBlock(data, prevBlock.Hash)

	chain.Blocks = append(chain.Blocks, block)
}

func GenesisBlock() *Block {
	genData := "Did chicken come first or the Egg?"
	return CreateBlock(genData, []byte{})
}

func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}
