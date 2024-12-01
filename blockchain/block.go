package blockchain

type Block struct {
	Hash []byte
	// Transactions []*Transaction
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{Hash: []byte{}, Data: []byte(data), PrevHash: prevHash}
	pow := NewProof(block)

	nonce, hash := pow.Proof()

	block.Hash = hash[:]
	block.Nonce = nonce

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
