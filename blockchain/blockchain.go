package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/farjad/AI-Blockchain/utils"
)

const DBPath = "./db/blocks"

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

// Helps iterate through the blockchain
type BlockchainParser struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Parsing through the blockchain
func (chain *Blockchain) Parser() *BlockchainParser {
	i := &BlockchainParser{chain.LastHash, chain.Database}

	return i
}

func (i *BlockchainParser) Next() *Block {
	var block *Block

	err := i.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.CurrentHash)
		utils.ErrHandle(err)

		encodedBlock, err := item.Value()
		utils.ErrHandle(err)
		block.Deserialize(encodedBlock)

		utils.ErrHandle(err)

		return err
	})

	utils.ErrHandle(err)

	// GOing backwards by getting the previous block
	i.CurrentHash = block.PrevHash

	return block
}

func InitBlockChain(address string) *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = DBPath
	opts.ValueDir = DBPath

	db, err := badger.Open(opts)
	utils.ErrHandle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No blockchain found")
			genesis := GenesisBlock()
			fmt.Println("Genesis Proved!!!")

			err = txn.Set([]byte("lh"), genesis.Hash)
			utils.ErrHandle(err)
			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			lastHash, err = item.Value()
			utils.ErrHandle(err)
			return err
		}
	})

	utils.ErrHandle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func (chain *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.ErrHandle(err)
		lastHash, err = item.Value()
		return err
	})
	utils.ErrHandle(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.ErrHandle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	utils.ErrHandle(err)
}
