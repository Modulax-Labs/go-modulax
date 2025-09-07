package core

import (
	"fmt"

	"github.com/Modulax-Protocol/go-modulax/storage"
)

type StateReader interface {
	GetAccount(addr [20]byte) *Account
}
type Account struct {
	Balance uint64
	Nonce   uint64
}
type State struct {
	store    storage.Storer
	accounts map[[20]byte]*Account
}

func NewState(store storage.Storer) (*State, error) {
	s := &State{
		store:    store,
		accounts: make(map[[20]byte]*Account),
	}
	return s, nil
}
func (s *State) GetAccount(addr [20]byte) *Account {
	if acc, ok := s.accounts[addr]; ok {
		return acc
	}
	s.accounts[addr] = &Account{}
	return s.accounts[addr]
}
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
func (s *State) AddBalance(addr [20]byte, amount uint64) error {
	account := s.GetAccount(addr)
	account.Balance += amount
	return nil
}
func (s *State) Persist() error {
	fmt.Println("--- Persisting State ---")
	for addr, account := range s.accounts {
		fmt.Printf("Address: %x, Balance: %d, Nonce: %d\n", addr, account.Balance, account.Nonce)
	}
	return nil
}
