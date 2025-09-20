package evm

import (
	"fmt"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/crypto"
)

// EVM is the Ethereum Virtual Machine that processes transactions.
// It implements the core.Executor interface.
type EVM struct {
	state *core.State
}

// NewEVM creates a new EVM instance.
func NewEVM(state *core.State) *EVM {
	return &EVM{
		state: state,
	}
}

// Execute processes a standard, signed transaction and applies its changes to the state.
// This method signature now matches the core.Executor interface.
func (e *EVM) Execute(tx *core.Transaction) error {
	fmt.Printf("Executing transaction %x in the EVM...\n", tx.Hash)

	// Genesis transaction is a special case handled by the blockchain's bootstrap logic.
	// The EVM only executes standard, user-signed transactions.
	if tx.PublicKey == nil {
		return fmt.Errorf("EVM cannot execute transactions without a public key")
	}

	valid, err := tx.Verify()
	if err != nil || !valid {
		return fmt.Errorf("invalid transaction signature")
	}

	senderAddr := crypto.AddressFromPublicKey(tx.PublicKey)
	senderAccount := e.state.GetAccount(senderAddr)

	if tx.Nonce != senderAccount.Nonce {
		return fmt.Errorf("invalid nonce. want %d, got %d", senderAccount.Nonce, tx.Nonce)
	}

	if err := e.state.Transfer(senderAddr, tx.To, tx.Value); err != nil {
		return err
	}

	senderAccount.Nonce++
	return nil
}
