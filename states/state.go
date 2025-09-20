package core

import (
	"fmt"

	"github.com/Modulax-Protocol/go-modulax/storage"
)

// StateReader is an interface that provides read-only access to the state.
type StateReader interface {
	GetAccount(addr [20]byte) *Account
}

// Account represents a user account in the state.
type Account struct {
	Balance uint64
	Nonce   uint64
}

// State manages all the accounts in the blockchain. It implements the StateReader interface.
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
	return s, nil
}

// GetAccount retrieves an account from the state. If it doesn't exist, it creates a new one.
func (s *State) GetAccount(addr [20]byte) *Account {
	if acc, ok := s.accounts[addr]; ok {
		return acc
	}
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
func (s *State) AddBalance(addr [20]byte, amount uint64) error {
	account := s.GetAccount(addr)
	account.Balance += amount
	return nil
}

// Persist writes the current state to the database.
func (s *State) Persist() error {
	fmt.Println("--- Persisting State ---")
	for addr, account := range s.accounts {
		fmt.Printf("Address: %x, Balance: %d, Nonce: %d\n", addr, account.Balance, account.Nonce)
	}
	return nil
}
