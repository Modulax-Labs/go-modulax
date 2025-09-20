package core

import (
	"fmt"
	"sync"

	"github.com/Modulax-Protocol/go-modulax/crypto"
)

type TxPool struct {
	mu         sync.RWMutex
	all        map[[32]byte]*Transaction
	state      StateReader
	maxPending int
}

func NewTxPool(state StateReader) *TxPool {
	return &TxPool{
		all:        make(map[[32]byte]*Transaction),
		state:      state,
		maxPending: 1024,
	}
}

func (p *TxPool) Add(tx *Transaction) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.all) >= p.maxPending {
		return fmt.Errorf("transaction pool is full")
	}
	if _, ok := p.all[tx.Hash]; ok {
		return fmt.Errorf("transaction %x already in pool", tx.Hash)
	}
	if valid, err := tx.Verify(); err != nil || !valid {
		return fmt.Errorf("invalid transaction signature")
	}
	senderAddr := crypto.AddressFromPublicKey(tx.PublicKey)
	senderAccount := p.state.GetAccount(senderAddr)
	if tx.Nonce != senderAccount.Nonce {
		return fmt.Errorf("invalid nonce for tx %x. want %d, got %d", tx.Hash, senderAccount.Nonce, tx.Nonce)
	}
	if senderAccount.Balance < tx.Value {
		return fmt.Errorf("insufficient funds for transfer")
	}
	p.all[tx.Hash] = tx
	return nil
}

func (p *TxPool) Pending() []*Transaction {
	p.mu.RLock()
	defer p.mu.RUnlock()
	txs := make([]*Transaction, 0, len(p.all))
	for _, tx := range p.all {
		txs = append(txs, tx)
	}
	return txs
}

func (p *TxPool) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.all = make(map[[32]byte]*Transaction)
}
