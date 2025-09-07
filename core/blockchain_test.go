package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/Modulax-Protocol/go-modulax/crypto"
	"github.com/Modulax-Protocol/go-modulax/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockExecutor is a test implementation of the Executor interface.
// It allows us to test the Blockchain without a real EVM, breaking the import cycle.
type MockExecutor struct{}

func (e *MockExecutor) Execute(tx *Transaction) error {
	// For this mock, we assume all transactions are valid.
	// In more advanced tests, we could add logic here to simulate failures.
	return nil
}

// newTestBlockchain creates a new blockchain with a temporary LevelDB store for testing.
func newTestBlockchain(t *testing.T) *Blockchain {
	path := fmt.Sprintf("./test_db_%s", t.Name())
	os.RemoveAll(path)

	store, err := storage.NewLevelDBStore(path)
	require.NoError(t, err)

	// Create a new MockExecutor for the test.
	executor := &MockExecutor{}

	// Create the blockchain and pass in the mock executor.
	bc, err := NewBlockchain(store, executor)
	require.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(path)
	})

	return bc
}

// TestNewBlockchain tests the creation of a new blockchain and its genesis state.
func TestNewBlockchain(t *testing.T) {
	bc := newTestBlockchain(t)
	assert.NotNil(t, bc.state)
	latestBlock, err := bc.GetLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, uint32(0), latestBlock.Header.Height)

	// We must derive the genesis address from the known private key to verify the balance.
	pkBytes, _ := hex.DecodeString(GENESIS_PRIVATE_KEY)
	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = new(big.Int).SetBytes(pkBytes)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(pkBytes)
	genesisWallet := crypto.WalletFromPrivateKey(privateKey)
	genesisAddress := genesisWallet.Address()

	genesisBalance := bc.state.GetAccount(genesisAddress).Balance
	assert.Equal(t, uint64(1_000_000), genesisBalance)
}

// TestAddBlockWithTransfer tests adding a block.
// NOTE: With a MockExecutor, we are not testing the state change here, only that
// the block can be added to the chain without error. State change is tested
// in the EVM's own unit tests.
func TestAddBlockWithTransfer(t *testing.T) {
	bc := newTestBlockchain(t)

	senderWallet, _ := crypto.NewWallet()
	receiverWallet, _ := crypto.NewWallet()
	receiverAddress := receiverWallet.Address()

	println("Sender Balance:", bc.state.GetAccount(senderWallet.Address()).Balance)
	println("Receiver Balance:", bc.state.GetAccount(receiverAddress).Balance)

	tx := &Transaction{
		To:        receiverAddress,
		Value:     25,
		Nonce:     0,
		PublicKey: senderWallet.PublicKey(),
	}
	txHash, _ := tx.CalculateHash()
	signature, _ := senderWallet.Sign(txHash)
	tx.Signature = signature
	tx.Hash = txHash

	_, err := bc.AddBlock([]*Transaction{tx})
	require.NoError(t, err)

	latestBlock, _ := bc.GetLatestBlock()
	assert.Equal(t, uint32(1), latestBlock.Header.Height)

}
