package core

import (
	"os"
	"testing"

	"github.com/Modulax-Protocol/go-modulax/crypto"
	"github.com/Modulax-Protocol/go-modulax/storage"
	"github.com/stretchr/testify/assert"
)

// newTestBlockchain creates a new blockchain with a temporary LevelDB store for testing.
func newTestBlockchain(t *testing.T) *Blockchain {
	// Create a temporary directory for the test database.
	path := "./test_db"
	// Clean up any previous test database.
	os.RemoveAll(path)

	store, err := storage.NewLevelDBStore(path)
	assert.NoError(t, err)

	bc, err := NewBlockchain(store)
	assert.NoError(t, err)

	// Add a cleanup function to remove the database after the test.
	t.Cleanup(func() {
		os.RemoveAll(path)
	})

	return bc
}

// TestNewBlockchain tests the creation of a new blockchain and its genesis state.
func TestNewBlockchain(t *testing.T) {
	bc := newTestBlockchain(t)
	assert.NotNil(t, bc.state)
	assert.Equal(t, uint32(0), len(bc.latestBlockHash))

	// Check that the genesis address was correctly funded.
	genesisBalance := bc.state.GetBalance(GENESIS_ADDRESS)
	assert.Equal(t, uint64(1_000_000), genesisBalance)
}

// TestAddBlockWithTransfer tests adding a block with a valid value transfer transaction.
func TestAddBlockWithTransfer(t *testing.T) {
	bc := newTestBlockchain(t)

	// Create two wallets for the test.
	senderWallet, _ := crypto.NewWallet()
	receiverWallet, _ := crypto.NewWallet()
	senderAddress := senderWallet.Address()
	receiverAddress := receiverWallet.Address()

	// Manually fund the sender's account directly in the state for this test.
	bc.state.AddBalance(senderAddress, 100)
	assert.Equal(t, uint64(100), bc.state.GetBalance(senderAddress))

	// Create a transaction from the sender to the receiver.
	tx := &Transaction{
		To:        receiverAddress,
		Value:     25,
		Nonce:     0, // First transaction from this account
		PublicKey: senderWallet.PublicKey(),
	}
	txHash, _ := tx.CalculateHash()
	signature, _ := senderWallet.Sign(txHash)
	tx.Signature = signature
	tx.Hash = txHash

	// Add a new block containing this transaction.
	_, err := bc.AddBlock([]*Transaction{tx})
	assert.NoError(t, err)

	// Check the final balances.
	senderBalance := bc.state.GetBalance(senderAddress)
	receiverBalance := bc.state.GetBalance(receiverAddress)

	assert.Equal(t, uint64(75), senderBalance)
	assert.Equal(t, uint64(25), receiverBalance)

	// Check that the sender's nonce was incremented.
	senderAccount := bc.state.GetAccount(senderAddress)
	assert.Equal(t, uint64(1), senderAccount.Nonce)
}

