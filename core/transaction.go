package core

import (
	"crypto/sha256"
	"encoding/json"
)

// Transaction represents a state change on the Modulax blockchain.
// For now, it's a simple value transfer.
type Transaction struct {
	// In a real implementation, From/To would be public keys or addresses.
	From      []byte `json:"from"`
	To        []byte `json:"to"`
	Value     uint64 `json:"value"`
	Nonce     uint64 `json:"nonce"` // Transaction nonce to prevent replay attacks
	Timestamp int64  `json:"timestamp"`

	// In the PQ-EVM, this signature would use a quantum-resistant algorithm.
	Signature []byte `json:"signature"`
	Hash      [32]byte `json:"hash"`
}

// CalculateHash calculates the SHA256 hash of the transaction data (excluding signature and hash).
func (tx *Transaction) CalculateHash() ([32]byte, error) {
	// Create a temporary struct to marshal without the signature and hash fields
	txData := struct {
		From      []byte `json:"from"`
		To        []byte `json:"to"`
		Value     uint64 `json:"value"`
		Nonce     uint64 `json:"nonce"`
		Timestamp int64  `json:"timestamp"`
	}{
		From:      tx.From,
		To:        tx.To,
		Value:     tx.Value,
		Nonce:     tx.Nonce,
		Timestamp: tx.Timestamp,
	}

	txBytes, err := json.Marshal(txData)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(txBytes), nil
}

// Sign signs the transaction with a private key (placeholder function).
// In a real implementation, this would take a private key and produce a PQ-secure signature.
func (tx *Transaction) Sign() error {
	// Placeholder for signing logic.
	tx.Signature = []byte("placeholder_signature")
	return nil
}
