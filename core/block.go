package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

type Header struct {
	ParentHash [32]byte
	Height     uint32
	Timestamp  int64
}
type Block struct {
	Header       *Header
	Transactions []*Transaction
	Hash         [32]byte
}

func NewBlock(parentHash [32]byte, height uint32, transactions []*Transaction) *Block {
	header := &Header{
		ParentHash: parentHash,
		Height:     height,
		Timestamp:  time.Now().UnixNano(),
	}
	block := &Block{
		Header:       header,
		Transactions: transactions,
	}
	block.Hash = block.CalculateHash()
	return block
}
func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}
func (b *Block) CalculateHash() [32]byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(b.Header); err != nil {
		panic(err)
	}
	return sha256.Sum256(buf.Bytes())
}
func (b *Block) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(b); err != nil {
		return nil, fmt.Errorf("failed to encode block: %w", err)
	}
	return buf.Bytes(), nil
}
func DecodeBlock(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		return nil, fmt.Errorf("failed to decode block: %w", err)
	}
	return &block, nil
}
