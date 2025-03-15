package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/farjad/AI-Blockchain/utils"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func CreateBlock(transactions []*Transaction, prevHash []byte) *Block {
	block := &Block{
		Hash:         []byte{},
		Transactions: transactions,
		PrevHash:     prevHash,
	}
	pow := NewProof(block)
	nonce, hash := pow.Proof()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func GenesisBlock() *Block {
	return CreateBlock([]*Transaction{}, []byte{})
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode((b))

	utils.ErrHandle(err)

	return result.Bytes()
}

func (b *Block) Deserialize(data []byte) *Block {
	var block Block
	res := gob.NewDecoder(bytes.NewReader(data))

	err := res.Decode(&block)

	utils.ErrHandle(err)

	return &block
}
