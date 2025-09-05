package core

import (
	"os"
	"testing"

	"github.com/Modulax-Protocol/go-modulax/storage"
)

// newTestBlockchain is a helper function to create a new blockchain with a temporary
// LevelDB store for testing purposes. It also returns a cleanup function.
func newTestBlockchain(t *testing.T) (*Blockchain, func()) {
	// Create a temporary directory for the database.
	dir, err := os.MkdirTemp("", "modulax_test_db")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create a cleanup function to remove the directory after the test.
	cleanup := func() {
		os.RemoveAll(dir)
	}

	// Create a new LevelDB store in the temporary directory.
	store, err := storage.NewLevelDBStore(dir)
	if err != nil {
		t.Fatalf("Failed to create new LevelDB store: %v", err)
	}

	// Create a new blockchain with the test store.
	bc, err := NewBlockchain(store)
	if err != nil {
		t.Fatalf("Failed to create new blockchain: %v", err)
	}

	return bc, cleanup
}

// TestNewBlockchain tests that a new blockchain is created correctly with a genesis block.
func TestNewBlockchain(t *testing.T) {
	bc, cleanup := newTestBlockchain(t)
	defer cleanup() // Ensure cleanup is called at the end of the test.

	latestBlock, err := bc.GetLatestBlock()
	if err != nil {
		t.Fatalf("Failed to get latest block: %v", err)
	}

	if latestBlock.Header.Height != 0 {
		t.Errorf("Expected genesis block height to be 0, but got %d", latestBlock.Header.Height)
	}
}

// TestAddBlock tests adding a block and retrieving it.
func TestAddBlock(t *testing.T) {
	bc, cleanup := newTestBlockchain(t)
	defer cleanup()

	genesisBlock, _ := bc.GetLatestBlock()

	// Create a sample transaction.
	tx := &Transaction{From: []byte("a"), To: []byte("b"), Value: 1}
	tx.Sign()
	hash, _ := tx.CalculateHash()
	tx.Hash = hash

	// Add a new block.
	newBlock, err := bc.AddBlock([]*Transaction{tx})
	if err != nil {
		t.Fatalf("Failed to add block: %v", err)
	}

	if newBlock.Header.Height != 1 {
		t.Errorf("Expected new block height to be 1, but got %d", newBlock.Header.Height)
	}

	if newBlock.Header.ParentHash != genesisBlock.Hash {
		t.Errorf("Expected new block's parent hash to match genesis block's hash")
	}
}

// TestBlockchainPersistence tests that the blockchain state is saved and reloaded correctly.
func TestBlockchainPersistence(t *testing.T) {
	// Create a temporary directory for the database.
	dir, err := os.MkdirTemp("", "modulax_persistence_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Cleanup at the end.

	// --- First Instance ---
	store1, err := storage.NewLevelDBStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	bc1, err := NewBlockchain(store1)
	if err != nil {
		t.Fatal(err)
	}

	// Add a block to the first instance.
	_, err = bc1.AddBlock([]*Transaction{})
	if err != nil {
		t.Fatal(err)
	}

	// Get the hash of the latest block before closing.
	latestHashBeforeClose := bc1.latestBlockHash

	// Close the first database connection.
	if err := store1.Close(); err != nil {
		t.Fatalf("Failed to close first store: %v", err)
	}

	// --- Second Instance ---
	// Create a new store and blockchain instance using the same database path.
	store2, err := storage.NewLevelDBStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	bc2, err := NewBlockchain(store2)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the reloaded blockchain has the same latest block hash.
	if bc2.latestBlockHash != latestHashBeforeClose {
		t.Errorf("Expected reloaded blockchain to have latest hash %x, but got %x", latestHashBeforeClose, bc2.latestBlockHash)
	}

	latestBlock, _ := bc2.GetLatestBlock()
	if latestBlock.Header.Height != 1 {
		t.Errorf("Expected reloaded blockchain height to be 1, but got %d", latestBlock.Header.Height)
	}
}

