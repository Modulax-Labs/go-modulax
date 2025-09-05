package core

import (
	"fmt"
	"sync"
)

// TxPool holds all pending transactions that are waiting to be included in a block.
// It provides a simple, thread-safe way to manage transactions.
type TxPool struct {
	mu           sync.RWMutex
	transactions map[[32]byte]*Transaction
}

// NewTxPool creates a new transaction pool.
func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[[32]byte]*Transaction),
	}
}

// Add adds a transaction to the pool.
func (p *TxPool) Add(tx *Transaction) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if the transaction already exists in the pool.
	if _, ok := p.transactions[tx.Hash]; ok {
		return fmt.Errorf("transaction %x already in the pool", tx.Hash)
	}

	p.transactions[tx.Hash] = tx
	return nil
}

// Has returns true if the transaction exists in the pool.
func (p *TxPool) Has(hash [32]byte) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, ok := p.transactions[hash]
	return ok
}

// Pending returns all transactions currently in the pool.
func (p *TxPool) Pending() []*Transaction {
	p.mu.RLock()
	defer p.mu.RUnlock()

	txs := make([]*Transaction, 0, len(p.transactions))
	for _, tx := range p.transactions {
		txs = append(txs, tx)
	}
	return txs
}

// Clear clears all transactions from the pool.
func (p *TxPool) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.transactions = make(map[[32]byte]*Transaction)
}

// Count returns the number of pending transactions.
func (p *TxPool) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.transactions)
}

