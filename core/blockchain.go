package core

import (
	"fmt"

	"github.com/Modulax-Protocol/go-modulax/storage"
)

const (
	// lastBlockHashKey will be the key for storing the hash of the latest block in the DB.
	lastBlockHashKey = "l"
)

// Blockchain represents the chain of blocks, now backed by persistent storage.
type Blockchain struct {
	store           storage.Storer
	latestBlockHash [32]byte
}

// NewBlockchain creates a new blockchain instance, loading from storage if it exists.
func NewBlockchain(store storage.Storer) (*Blockchain, error) {
	bc := &Blockchain{store: store}

	// Check if a latest block hash is already stored.
	latestHashBytes, err := store.Get([]byte(lastBlockHashKey))
	// A "not found" error is expected for a new database.
	if err != nil {
		// If the key is not found, it means the database is new.
		// We need to create and store the genesis block.
		genesis := createGenesisBlock()
		genesisBytes, err := genesis.Encode()
		if err != nil {
			return nil, err
		}

		// Store the genesis block by its hash.
		if err := store.Put(genesis.Hash[:], genesisBytes); err != nil {
			return nil, err
		}

		// Store the genesis block's hash as the "latest block hash" pointer.
		if err := store.Put([]byte(lastBlockHashKey), genesis.Hash[:]); err != nil {
			return nil, err
		}

		bc.latestBlockHash = genesis.Hash
	} else {
		// If the key exists, we load the latest block hash from the database.
		copy(bc.latestBlockHash[:], latestHashBytes)
	}

	return bc, nil
}

// AddBlock adds a new block to the blockchain and persists it to storage.
func (bc *Blockchain) AddBlock(transactions []*Transaction) (*Block, error) {
	// Get the latest block hash to use as the parent hash for the new block.
	prevBlockHash := bc.latestBlockHash

	// Retrieve the full previous block from the database to get its height.
	prevBlockBytes, err := bc.store.Get(prevBlockHash[:])
	if err != nil {
		return nil, fmt.Errorf("could not get previous block: %w", err)
	}
	prevBlock, err := DecodeBlock(prevBlockBytes)
	if err != nil {
		return nil, fmt.Errorf("could not decode previous block: %w", err)
	}

	// Create the new block.
	newBlock := NewBlock(prevBlockHash, prevBlock.Header.Height+1, 0, transactions)

	// Encode the new block for storage.
	newBlockBytes, err := newBlock.Encode()
	if err != nil {
		return nil, fmt.Errorf("could not encode new block: %w", err)
	}

	// Store the new block in the database, keyed by its hash.
	if err := bc.store.Put(newBlock.Hash[:], newBlockBytes); err != nil {
		return nil, err
	}

	// Update the "latest block hash" pointer in the database.
	// NOTE: In a real implementation, this would be part of a DB transaction
	// to ensure atomicity.
	if err := bc.store.Put([]byte(lastBlockHashKey), newBlock.Hash[:]); err != nil {
		return nil, err
	}

	// Update the latest block hash in memory.
	bc.latestBlockHash = newBlock.Hash

	return newBlock, nil
}

// GetLatestBlock returns the most recent block on the chain from storage.
func (bc *Blockchain) GetLatestBlock() (*Block, error) {
	latestHash := bc.latestBlockHash
	if latestHash == [32]byte{} {
		return nil, fmt.Errorf("blockchain is empty")
	}

	blockBytes, err := bc.store.Get(latestHash[:])
	if err != nil {
		return nil, fmt.Errorf("could not get latest block from store: %w", err)
	}

	return DecodeBlock(blockBytes)
}

// createGenesisBlock creates the very first block in the chain.
func createGenesisBlock() *Block {
	// The genesis block has no parent, so its parent hash is all zeros.
	parentHash := [32]byte{}
	// It contains a single, special transaction.
	genesisTx := &Transaction{
		From:      []byte("genesis"),
		To:        []byte("genesis"),
		Value:     1000000, // Initial supply
		Timestamp: 0,
		Nonce:     0,
	}
	// Sign and hash the genesis transaction (using placeholder functions for now)
	genesisTx.Sign()
	hash, _ := genesisTx.CalculateHash()
	genesisTx.Hash = hash

	return NewBlock(parentHash, 0, 0, []*Transaction{genesisTx})
}

