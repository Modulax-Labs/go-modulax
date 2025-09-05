package core

import (
	"os"
	"testing"

	"github.com/Modulax-Protocol/go-modulax/storage"
	"github.com/stretchr/testify/assert"
)

// newTestBlockchain creates a new blockchain with a temporary LevelDB store.
func newTestBlockchain(t *testing.T) (*Blockchain, func()) {
	// Create a temporary directory for the test database.
	tempDir, err := os.MkdirTemp("", "modulax_test_db")
	assert.NoError(t, err)

	// Create a new LevelDB store in the temporary directory.
	store, err := storage.NewLevelDBStore(tempDir)
	assert.NoError(t, err)

	// Create a new blockchain instance.
	bc, err := NewBlockchain(store)
	assert.NoError(t, err)

	// Teardown function to clean up the database after the test.
	teardown := func() {
		store.Close()
		os.RemoveAll(tempDir)
	}

	return bc, teardown
}

func TestNewBlockchain(t *testing.T) {
	bc, teardown := newTestBlockchain(t)
	defer teardown()

	// Test that the blockchain is not nil.
	assert.NotNil(t, bc)

	// Test that a latest block hash exists (it should be the genesis block).
	assert.NotEqual(t, [32]byte{}, bc.latestBlockHash)

	// Test that the genesis block has height 0.
	latestBlock, err := bc.GetLatestBlock()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), latestBlock.Header.Height)
}

func TestAddBlock(t *testing.T) {
	bc, teardown := newTestBlockchain(t)
	defer teardown()

	// Get the initial latest block (genesis).
	genesisBlock, err := bc.GetLatestBlock()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), genesisBlock.Header.Height)

	// Create a sample transaction for the new block.
	sampleTx := &Transaction{
		Data: []byte("hello world"),
	}
	sampleTx.Sign()
	hash, _ := sampleTx.CalculateHash()
	sampleTx.Hash = hash

	// Add a new block.
	newBlock, err := bc.AddBlock([]*Transaction{sampleTx})
	assert.NoError(t, err)
	assert.NotNil(t, newBlock)

	// Test that the new block has the correct height.
	assert.Equal(t, uint32(1), newBlock.Header.Height)

	// Test that the new block's parent hash is the genesis block's hash.
	assert.Equal(t, genesisBlock.Hash, newBlock.Header.ParentHash)

	// Test that the blockchain's latest block hash has been updated.
	assert.Equal(t, newBlock.Hash, bc.latestBlockHash)

	// Test that GetLatestBlock now returns our new block.
	latestBlock, err := bc.GetLatestBlock()
	assert.NoError(t, err)
	assert.Equal(t, newBlock.Hash, latestBlock.Hash)
}

