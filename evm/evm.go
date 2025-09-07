package evm

import (
	"fmt"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/crypto"
)

// EVM is the Ethereum Virtual Machine that processes transactions.
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
func (e *EVM) Execute(tx *core.Transaction) error {
	fmt.Printf("Executing transaction %x in the EVM...\n", tx.Hash)

	// 1. Verify the signature.
	valid, err := tx.Verify()
	if err != nil || !valid {
		return fmt.Errorf("invalid transaction signature")
	}

	// 2. Derive the sender's address.
	senderAddr := crypto.AddressFromPublicKey(tx.PublicKey)
	senderAccount := e.state.GetAccount(senderAddr)

	// 3. Check the nonce.
	if tx.Nonce != senderAccount.Nonce {
		return fmt.Errorf("invalid nonce. want %d, got %d", senderAccount.Nonce, tx.Nonce)
	}

	// 4. Perform the transfer.
	if err := e.state.Transfer(senderAddr, tx.To, tx.Value); err != nil {
		return err
	}

	// 5. Increment the sender's nonce.
	senderAccount.Nonce++
	return nil
}
