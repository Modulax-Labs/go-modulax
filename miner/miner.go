package miner

import (
	"context"
	"fmt"
	"time"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/network"
)

// Miner is responsible for creating new blocks.
type Miner struct {
	bc     *core.Blockchain
	txPool *core.TxPool
	pubsub *network.PubSubService
	stopCh chan struct{}
}

// NewMiner creates a new Miner instance.
func NewMiner(bc *core.Blockchain, txPool *core.TxPool, pubsub *network.PubSubService) *Miner {
	return &Miner{
		bc:     bc,
		txPool: txPool,
		pubsub: pubsub,
		stopCh: make(chan struct{}),
	}
}

// Start begins the block proposing loop.
func (m *Miner) Start() {
	fmt.Println("⛏️ Starting Block Proposer...")

	// Create a ticker that fires every 5 seconds.
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-m.stopCh:
				ticker.Stop()
				return
			case <-ticker.C:
				m.createNewBlock()
			}
		}
	}()
}

// Stop halts the block proposing loop.
func (m *Miner) Stop() {
	close(m.stopCh)
}

// createNewBlock checks the txpool and creates a new block if there are pending transactions.
func (m *Miner) createNewBlock() {
	pendingTxs := m.txPool.Pending()
	if len(pendingTxs) == 0 {
		// No transactions, do nothing.
		return
	}

	fmt.Printf("💡 Found %d pending transactions. Proposing a new block...\n", len(pendingTxs))

	// Create and add the block to our own chain.
	newBlock, err := m.bc.AddBlock(pendingTxs)
	if err != nil {
		fmt.Printf("Error creating new block: %v\n", err)
		return
	}

	// Clear our local transaction pool.
	m.txPool.Clear()

	// Encode and broadcast the new block to the network.
	blockBytes, err := newBlock.Encode()
	if err != nil {
		fmt.Printf("Error encoding block for broadcast: %v\n", err)
		return
	}
	if err := m.pubsub.BroadcastBlock(context.Background(), blockBytes); err != nil {
		fmt.Printf("Error broadcasting block: %v\n", err)
	}

	fmt.Printf("✅ Successfully proposed and broadcasted new block. Hash: %x\n", newBlock.Hash)
}
