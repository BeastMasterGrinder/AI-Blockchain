package blockchain

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dgraph-io/badger"
	"github.com/farjad/AI-Blockchain/utils"
)

const DBPath = "./db/blocks_%s"

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

func DBexists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

func (chain *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.ErrHandle(err)
		lastHash, _ := item.Value()

		item, err = txn.Get(lastHash)
		utils.ErrHandle(err)
		lastBlockData, _ := item.Value()

		lastBlock = *Deserialize(lastBlockData)

		return nil
	})
	utils.ErrHandle(err)

	return lastBlock.Height
}

func (chain *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte

	iter := chain.Parser()

	for {
		block := iter.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks
}

func (chain *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := chain.Database.View(func(txn *badger.Txn) error {
		if item, err := txn.Get(blockHash); err != nil {
			return errors.New("Block is not found")
		} else {
			blockData, _ := item.Value()

			block = *Deserialize(blockData)
		}
		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
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

		block = Deserialize(encodedBlock)

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

func (chain *Blockchain) AddBlock(block *Block) {
	err := chain.Database.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(block.Hash); err == nil {
			return nil
		}

		blockData := block.Serialize()
		err := txn.Set(block.Hash, blockData)
		utils.ErrHandle(err)

		item, err := txn.Get([]byte("lh"))
		utils.ErrHandle(err)
		lastHash, _ := item.Value()

		item, err = txn.Get(lastHash)
		utils.ErrHandle(err)
		lastBlockData, _ := item.Value()

		lastBlock := Deserialize(lastBlockData)

		if block.Height > lastBlock.Height {
			err = txn.Set([]byte("lh"), block.Hash)
			utils.ErrHandle(err)
			chain.LastHash = block.Hash
		}

		return nil
	})
	utils.ErrHandle(err)
}
func (chain *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if chain.VerifyTransaction(tx) != true {
			log.Panic("Invalid Transaction")
		}
	}

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.ErrHandle(err)
		lastHash, err = item.Value()

		item, err = txn.Get(lastHash)
		utils.ErrHandle(err)
		lastBlockData, _ := item.Value()

		lastBlock := Deserialize(lastBlockData)

		lastHeight = lastBlock.Height

		return err
	})
	utils.ErrHandle(err)

	newBlock := CreateBlock(transactions, lastHash, lastHeight+1)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.ErrHandle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	utils.ErrHandle(err)

	return newBlock
}

func (chain *Blockchain) VerifyTransaction(tx *Transaction) bool {
	// prevTxs := make(map[string]Transaction)

	// for _, in := range tx.Inputs {
	// 	prevTx, err := chain.FindTransaction(in.ID)
	// 	utils.ErrHandle(err)
	// 	prevTxs[prevTx.ID] = prevTx
	// }

	// return tx.Verify(prevTxs)
	return true
}

func ContinueBlockChain(nodeId string) *Blockchain {
	path := fmt.Sprintf(DBPath, nodeId)
	if DBexists(path) == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path

	db, err := openDB(path, opts)
	utils.ErrHandle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.ErrHandle(err)
		lastHash, err = item.Value()

		return err
	})
	utils.ErrHandle(err)

	chain := Blockchain{lastHash, db}

	return &chain
}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	return db, err
}

func openDB(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	} else {
		return db, nil
	}
}
