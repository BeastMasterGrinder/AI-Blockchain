package blockchain

// Transaction holds the Hash of the Algorithm, Data, output, sender
type Transaction struct {
	ID  []byte
	Vin []TXInput
	// Vout []TXOutput
}

type TXInput struct {
	TxID      []byte
	Out       int
	Signature []byte
}
