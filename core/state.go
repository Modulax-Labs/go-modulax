package core

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/Modulax-Protocol/go-modulax/storage"
)

// Account represents a user account in the state.
type Account struct {
	Balance uint64
	Nonce   uint64
}

// State manages all the accounts in the blockchain.
type State struct {
	store    storage.Storer
	accounts map[[20]byte]*Account
}

// NewState creates a new State instance from storage.
func NewState(store storage.Storer) (*State, error) {
	s := &State{
		store:    store,
		accounts: make(map[[20]byte]*Account),
	}
	// In a real implementation, we would load the state from the database here.
	return s, nil
}

// GetAccount retrieves an account from the state. If it doesn't exist, it creates a new one.
func (s *State) GetAccount(addr [20]byte) *Account {
	if acc, ok := s.accounts[addr]; ok {
		return acc
	}
	// If account does not exist, create it.
	s.accounts[addr] = &Account{}
	return s.accounts[addr]
}

// Transfer moves funds from one account to another.
func (s *State) Transfer(from [20]byte, to [20]byte, amount uint64) error {
	fromAccount := s.GetAccount(from)
	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}
	fromAccount.Balance -= amount

	toAccount := s.GetAccount(to)
	toAccount.Balance += amount

	return nil
}

// AddBalance adds a balance to an account (used for minting).
func (s *State) AddBalance(addr [20]byte, amount uint64) {
	account := s.GetAccount(addr)
	account.Balance += amount
}

// GetBalance returns the balance of an account.
func (s *State) GetBalance(addr [20]byte) uint64 {
	if account, ok := s.accounts[addr]; ok {
		return account.Balance
	}
	return 0
}

// Persist writes the current state to the database.
// NOTE: This is a simple, inefficient implementation for demonstration.
func (s *State) Persist() error {
	// A real blockchain would use a Merkle Patricia Trie for efficient state storage.
	// For now, we will just print the state to the console for visibility.
	fmt.Println("--- Persisting State ---")
	for addr, account := range s.accounts {
		fmt.Printf("Address: %x, Balance: %d, Nonce: %d\n", addr, account.Balance, account.Nonce)
	}
	return nil
}

// Encode encodes the Account to a writer.
func (a *Account) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(a)
}

// Decode decodes an Account from a reader.
func (a *Account) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(a)
}

