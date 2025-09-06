package core

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Modulax-Protocol/go-modulax/crypto"
	"github.com/Modulax-Protocol/go-modulax/storage"
)

// GENESIS_ADDRESS is a placeholder for a well-known address to receive initial funds.
var GENESIS_ADDRESS = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

const (
	// lastBlockHashKey will be the key for storing the hash of the latest block in the DB.
	lastBlockHashKey = "l"
)

// Blockchain represents the chain of blocks, now backed by a state machine.
type Blockchain struct {
	store           storage.Storer
	state           *State
	latestBlockHash [32]byte
}

// NewBlockchain creates a new blockchain instance, loading from storage if it exists.
func NewBlockchain(store storage.Storer) (*Blockchain, error) {
	state, err := NewState(store)
	if err != nil {
		return nil, fmt.Errorf("failed to create state: %w", err)
	}

	bc := &Blockchain{
		store: store,
		state: state,
	}

	latestHashBytes, err := store.Get([]byte(lastBlockHashKey))
	if err != nil {
		// Database is new, create and process the genesis block.
		genesis := createGenesisBlock()
		if err := bc.processBlock(genesis); err != nil {
			return nil, err
		}

		genesisBytes, err := genesis.Encode()
		if err != nil {
			return nil, err
		}
		if err := store.Put(genesis.Hash[:], genesisBytes); err != nil {
			return nil, err
		}
		if err := store.Put([]byte(lastBlockHashKey), genesis.Hash[:]); err != nil {
			return nil, err
		}

		bc.latestBlockHash = genesis.Hash
	} else {
		copy(bc.latestBlockHash[:], latestHashBytes)
	}

	return bc, nil
}

// AddBlock adds a new block to the blockchain, processing its transactions.
func (bc *Blockchain) AddBlock(transactions []*Transaction) (*Block, error) {
	prevBlockHash := bc.latestBlockHash
	prevBlock, err := bc.GetBlockByHash(prevBlockHash)
	if err != nil {
		return nil, err
	}

	newBlock := NewBlock(prevBlockHash, prevBlock.Header.Height+1, transactions)

	if err := bc.processBlock(newBlock); err != nil {
		return nil, err
	}

	newBlockBytes, err := newBlock.Encode()
	if err != nil {
		return nil, err
	}
	if err := bc.store.Put(newBlock.Hash[:], newBlockBytes); err != nil {
		return nil, err
	}
	if err := bc.store.Put([]byte(lastBlockHashKey), newBlock.Hash[:]); err != nil {
		return nil, err
	}
	bc.latestBlockHash = newBlock.Hash

	return newBlock, nil
}

// AddExistingBlock adds a block from the network, processing its transactions.
func (bc *Blockchain) AddExistingBlock(block *Block) error {
	if block.Header.ParentHash != bc.latestBlockHash {
		return fmt.Errorf("received block has invalid parent hash")
	}

	if err := bc.processBlock(block); err != nil {
		return err
	}

	blockBytes, err := block.Encode()
	if err != nil {
		return err
	}
	if err := bc.store.Put(block.Hash[:], blockBytes); err != nil {
		return err
	}
	if err := bc.store.Put([]byte(lastBlockHashKey), block.Hash[:]); err != nil {
		return err
	}
	bc.latestBlockHash = block.Hash

	return nil
}

func (bc *Blockchain) State() *State {
	return bc.state
}




// GetBlockByHash retrieves a block from storage by its hash.
func (bc *Blockchain) GetBlockByHash(hash [32]byte) (*Block, error) {
	blockBytes, err := bc.store.Get(hash[:])
	if err != nil {
		return nil, fmt.Errorf("could not get block by hash %s: %w", hex.EncodeToString(hash[:]), err)
	}
	return DecodeBlock(blockBytes)
}

// processBlock is the state transition function.
func (bc *Blockchain) processBlock(block *Block) error {
	for _, tx := range block.Transactions {
		valid, err := tx.Verify()
		if err != nil {
			return fmt.Errorf("failed to verify tx %x: %w", tx.Hash, err)
		}
		if !valid {
			return fmt.Errorf("invalid signature on tx %x", tx.Hash)
		}

		// Derive the sender's address from their public key.
		senderAddr := crypto.AddressFromPublicKey(tx.PublicKey)

		// Handle the special case for the genesis transaction.
		if block.Header.Height == 0 {
			bc.state.AddBalance(GENESIS_ADDRESS, 1_000_000)
			continue // Skip the rest of the logic for the genesis tx
		}

		// For regular transactions, process the transfer.
		senderAccount := bc.state.GetAccount(senderAddr)
		if tx.Nonce != senderAccount.Nonce {
			return fmt.Errorf("invalid nonce for tx %x. want %d, got %d", tx.Hash, senderAccount.Nonce, tx.Nonce)
		}

		if err := bc.state.Transfer(senderAddr, tx.To, tx.Value); err != nil {
			return fmt.Errorf("failed to transfer for tx %x: %w", tx.Hash, err)
		}

		// Increment the sender's nonce to prevent replay attacks.
		senderAccount.Nonce++
	}

	fmt.Printf("✅ Processed %d transactions in Block %d\n", len(block.Transactions), block.Header.Height)

	return bc.state.Persist() // Save the updated state.
}

// GetLatestBlock returns the most recent block on the chain from storage.
func (bc *Blockchain) GetLatestBlock() (*Block, error) {
	if bc.latestBlockHash == [32]byte{} {
		return nil, fmt.Errorf("blockchain is empty")
	}
	return bc.GetBlockByHash(bc.latestBlockHash)
}

// createGenesisBlock creates the very first, deterministic block in the chain.
func createGenesisBlock() *Block {
	// A placeholder transaction is still needed to create a valid block structure.
	// The actual minting happens in processBlock.
	genesisTx := &Transaction{
		To:    [20]byte{}, // No recipient
		Value: 0,
		Nonce: 0,
	}
	// Sign with a dummy key, as there's no real sender.
	dummyWallet, _ := crypto.NewWallet()
	genesisTx.PublicKey = dummyWallet.PublicKey()

	// Correctly sign the transaction using the new method
	txHash, _ := genesisTx.CalculateHash()
	signature, _ := dummyWallet.Sign(txHash)
	genesisTx.Signature = signature
	genesisTx.Hash = txHash

	parentHash := [32]byte{}
	genesisHeader := &Header{
		ParentHash: parentHash,
		Height:     0,
		Timestamp:  time.Unix(0, 0).UnixNano(),
	}
	genesisBlock := &Block{Header: genesisHeader}
	genesisBlock.AddTransaction(genesisTx)
	genesisBlock.Hash = genesisBlock.CalculateHash()

	return genesisBlock
}

