package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// BlockHeader represents the header of a block.
type BlockHeader struct {
	ParentHash [32]byte
	Height     uint32
	Timestamp  int64
	// In a real blockchain, this would be a Merkle Root.
	TransactionsHash [32]byte
}

// Block represents a single block in the blockchain.
type Block struct {
	Header       *BlockHeader
	Transactions []*Transaction
	// Hash of the block header.
	Hash [32]byte
}

// NewBlock creates a new block.
func NewBlock(parentHash [32]byte, height uint32, timestamp int64, transactions []*Transaction) *Block {
	if timestamp == 0 {
		timestamp = time.Now().UnixNano()
	}

	header := &BlockHeader{
		ParentHash: parentHash,
		Height:     height,
		Timestamp:  timestamp,
		// For simplicity, we are not calculating a real Merkle Root yet.
		TransactionsHash: [32]byte{},
	}

	block := &Block{
		Header:       header,
		Transactions: transactions,
	}

	// Calculate and set the block's hash.
	hash := block.CalculateHash()
	block.Hash = hash

	return block
}

// CalculateHash calculates the SHA256 hash of the block's header.
func (b *Block) CalculateHash() [32]byte {
	// We use gob to serialize the header into a byte slice for hashing.
	// This is a deterministic way to represent the data.
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b.Header); err != nil {
		// In a real application, this should be handled more gracefully.
		panic(err)
	}
	return sha256.Sum256(buf.Bytes())
}

// Encode serializes the block into a byte slice using gob encoding.
func (b *Block) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeBlock deserializes a byte slice back into a Block pointer.
func DecodeBlock(data []byte) (*Block, error) {
	var block Block
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&block); err != nil {
		return nil, err
	}
	return &block, nil
}
