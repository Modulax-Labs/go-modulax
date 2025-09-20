package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/Modulax-Protocol/go-modulax/crypto"
	"github.com/Modulax-Protocol/go-modulax/storage"
)

// Executor defines the interface for a state transition machine.
// The EVM will implement this interface.
type Executor interface {
	Execute(tx *Transaction) error
}

const GENESIS_PRIVATE_KEY = "c1850f2b53d1e1f7cf655513970b13c847796a4b1054b1509a2a7a42140a33a5"

const lastBlockHashKey = "l"

// Blockchain now holds a generic Executor.
type Blockchain struct {
	store           storage.Storer
	state           *State
	executor        Executor
	latestBlockHash [32]byte
}

// NewBlockchain now accepts the Executor as a parameter.
func NewBlockchain(store storage.Storer, executor Executor) (*Blockchain, error) {
	state, err := NewState(store)
	if err != nil {
		return nil, fmt.Errorf("failed to create state: %w", err)
	}
	bc := &Blockchain{
		store:    store,
		state:    state,
		executor: executor,
	}

	latestHashBytes, err := store.Get([]byte(lastBlockHashKey))
	if err != nil {
		genesis := createGenesisBlock()
		// Manually process the genesis block to initialize the state.
		if err := state.AddBalance(genesis.Transactions[0].To, genesis.Transactions[0].Value); err != nil {
			return nil, fmt.Errorf("failed to apply genesis state: %w", err)
		}
		if err := state.Persist(); err != nil {
			return nil, err
		}

		genesisBytes, _ := genesis.Encode()
		store.Put(genesis.Hash[:], genesisBytes)
		store.Put([]byte(lastBlockHashKey), genesis.Hash[:])
		bc.latestBlockHash = genesis.Hash
	} else {
		copy(bc.latestBlockHash[:], latestHashBytes)
	}
	return bc, nil
}

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
	newBlockBytes, _ := newBlock.Encode()
	bc.store.Put(newBlock.Hash[:], newBlockBytes)
	bc.store.Put([]byte(lastBlockHashKey), newBlock.Hash[:])
	bc.latestBlockHash = newBlock.Hash
	return newBlock, nil
}

func (bc *Blockchain) AddExistingBlock(block *Block) error {
	if block.Header.ParentHash != bc.latestBlockHash {
		return fmt.Errorf("received block has invalid parent hash")
	}
	if err := bc.processBlock(block); err != nil {
		return err
	}
	blockBytes, _ := block.Encode()
	bc.store.Put(block.Hash[:], blockBytes)
	bc.store.Put([]byte(lastBlockHashKey), block.Hash[:])
	bc.latestBlockHash = block.Hash
	return nil
}

func (bc *Blockchain) State() *State {
	return bc.state
}

func (bc *Blockchain) GetBlockByHash(hash [32]byte) (*Block, error) {
	blockBytes, err := bc.store.Get(hash[:])
	if err != nil {
		return nil, fmt.Errorf("could not get block by hash %s: %w", hex.EncodeToString(hash[:]), err)
	}
	return DecodeBlock(blockBytes)
}

// processBlock now only handles non-genesis blocks.
func (bc *Blockchain) processBlock(block *Block) error {
	if block.Header.Height == 0 {
		return nil
	}
	for _, tx := range block.Transactions {
		if err := bc.executor.Execute(tx); err != nil {
			return fmt.Errorf("Execution failed for tx %x: %w", tx.Hash, err)
		}
	}
	fmt.Printf("✅ Processed %d transactions in Block %d\n", len(block.Transactions), block.Header.Height)
	return bc.state.Persist()
}

func (bc *Blockchain) GetLatestBlock() (*Block, error) {
	if bc.latestBlockHash == [32]byte{} {
		return nil, fmt.Errorf("blockchain is empty")
	}
	return bc.GetBlockByHash(bc.latestBlockHash)
}

func createGenesisBlock() *Block {
	pkBytes, _ := hex.DecodeString(GENESIS_PRIVATE_KEY)
	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = new(big.Int).SetBytes(pkBytes)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(pkBytes)
	genesisWallet := crypto.WalletFromPrivateKey(privateKey)

	genesisTx := &Transaction{
		To:        genesisWallet.Address(),
		Value:     1_000_000,
		Nonce:     0,
		PublicKey: genesisWallet.PublicKey(),
	}
	txHash, _ := genesisTx.CalculateHash()
	signature, _ := genesisWallet.Sign(txHash)
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
