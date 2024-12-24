package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/farjad/AI-Blockchain/utils"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Height       int
}

func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), []byte{}, txs, prevHash, 0, height}
	pow := NewProof(block)
	nonce, hash := pow.Proof()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func GenesisBlock() *Block {
	return CreateBlock([]*Transaction{}, []byte{}, 0)
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode((b))

	utils.ErrHandle(err)

	return result.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	res := gob.NewDecoder(bytes.NewReader(data))

	err := res.Decode(&block)

	utils.ErrHandle(err)

	return &block
}
