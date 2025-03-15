package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/farjad/AI-Blockchain/utils"
)

//
/*
Take the data from the block
Create a counter (nonce) which starts with 0s
Create Hash of the data plus the counter
Check if the Hash meets the requirement
Requirement
first 4 must be consecutive 0s
*/

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

const numZeros = 7

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	target.Lsh(target, uint(256-numZeros))

	newProofOW := &ProofOfWork{b, target}

	return newProofOW

}

func (pow *ProofOfWork) MergeData(Nonce int) []byte {
	var txBytes []byte
	for _, tx := range pow.Block.Transactions {
		txBytes = append(txBytes, tx.Serialize()...)
	}

	joinedData := bytes.Join(
		[][]byte{
			txBytes,
			pow.Block.PrevHash,
			utils.ToHex(int64(Nonce)),
			utils.ToHex(int64(numZeros)),
		},
		[]byte{})

	return joinedData
}

func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int

	data := pow.MergeData(pow.Block.Nonce)

	hash := sha256.Sum256(data)

	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.Target) == -1
}

func (pow *ProofOfWork) Proof() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.MergeData(nonce)
		hash = sha256.Sum256(data)

		// @Todo: remove this after testing
		fmt.Println(hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}

	}

	return nonce, hash[:]

}
