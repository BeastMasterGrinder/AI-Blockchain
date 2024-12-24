package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/farjad/AI-Blockchain/ipfs"
	"github.com/farjad/AI-Blockchain/utils"
)

type Transaction struct {
	ID            []byte
	AlgorithmCID  string
	AlgorithmHash string
	InputCID      string
	InputHash     string
	OutputHash    string
}

// NewTransaction creates a new transaction with algorithm and input data
func NewTransaction(algorithmPath, inputPath string, ipfsManager *ipfs.IPFSManager) (*Transaction, error) {

	// Upload algorithm to IPFS
	algCID, algHash, err := ipfsManager.UploadAlgorithm(algorithmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload algorithm: %v", err)
	}

	// Upload input data to IPFS
	inputCID, inputHash, err := ipfsManager.UploadData(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload input data: %v", err)
	}

	// Execute algorithm to get output hash
	outputHash, err := ipfsManager.ExecuteAlgorithm(algCID, inputCID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute algorithm: %v", err)
	}

	tx := &Transaction{
		AlgorithmCID:  algCID,
		AlgorithmHash: algHash,
		InputCID:      inputCID,
		InputHash:     inputHash,
		OutputHash:    outputHash,
	}

	// Generate transaction ID
	tx.ID = tx.Hash()

	return tx, nil
}

// Hash returns the hash of the transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}

	encoded := tx.Serialize()
	hash = sha256.Sum256(encoded)

	return hash[:]
}

// Serialize encodes the transaction into bytes
func (tx *Transaction) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(tx)
	utils.ErrHandle(err)
	return res.Bytes()
}

// Deserialize decodes the transaction from bytes
func DeserializeTransaction(data []byte) *Transaction {
	var transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	utils.ErrHandle(err)
	return &transaction
}

// Verify checks if the transaction is valid by executing the algorithm
func (tx *Transaction) Verify(ipfsManager *ipfs.IPFSManager) bool {
	calculatedHash, err := ipfsManager.ExecuteAlgorithm(tx.AlgorithmCID, tx.InputCID)
	if err != nil {
		fmt.Printf("Failed to verify transaction: %v\n", err)
		return false
	}
	return calculatedHash == tx.OutputHash
}
