package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

// Header represents the header of a block.
type Header struct {
	ParentHash [32]byte
	Height     uint32
	Timestamp  int64
}

// Block represents a block in the blockchain.
type Block struct {
	Header       *Header
	Transactions []*Transaction
	Hash         [32]byte
}

// NewBlock creates a new block.
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

// AddTransaction adds a transaction to the block.
func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

// CalculateHash calculates the SHA256 hash of the block's header.
func (b *Block) CalculateHash() [32]byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	// We only hash the header for the block's hash.
	if err := encoder.Encode(b.Header); err != nil {
		// This should not happen with our simple struct.
		panic(err)
	}
	return sha256.Sum256(buf.Bytes())
}

// Encode serializes the block into a byte slice using gob.
func (b *Block) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(b); err != nil {
		return nil, fmt.Errorf("failed to encode block: %w", err)
	}
	return buf.Bytes(), nil
}

// DecodeBlock deserializes a byte slice into a Block.
func DecodeBlock(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		return nil, fmt.Errorf("failed to decode block: %w", err)
	}
	return &block, nil
}

