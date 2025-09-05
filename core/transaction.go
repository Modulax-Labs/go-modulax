package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// Transaction represents a single transaction in the blockchain.
// For now, it's simplified to just hold arbitrary data.
type Transaction struct {
	Hash      [32]byte
	Data      []byte
	Signature []byte
}

// Sign simulates signing the transaction.
// In a real implementation, this would involve a private key.
func (tx *Transaction) Sign() {
	// For now, we'll just use a placeholder signature.
	tx.Signature = []byte("placeholder_signature")
}

// CalculateHash calculates the SHA256 hash of the transaction data.
func (tx *Transaction) CalculateHash() ([32]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// Encode the parts of the transaction that should be hashed.
	// We don't hash the Signature or the Hash itself.
	err := encoder.Encode(struct {
		Data []byte
	}{
		Data: tx.Data,
	})

	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to encode transaction for hashing: %w", err)
	}

	return sha256.Sum256(buf.Bytes()), nil
}

