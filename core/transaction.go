package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
)

type Transaction struct {
	Hash      [32]byte
	To        [20]byte
	Value     uint64
	Nonce     uint64
	PublicKey []byte
	Signature []byte
}

func (tx *Transaction) Verify() (bool, error) {
	if tx.PublicKey == nil || tx.Signature == nil {
		return false, fmt.Errorf("transaction has no signature or public key")
	}
	r := &big.Int{}
	s := &big.Int{}
	sigLen := len(tx.Signature)
	r.SetBytes(tx.Signature[:(sigLen / 2)])
	s.SetBytes(tx.Signature[(sigLen / 2):])
	x := &big.Int{}
	y := &big.Int{}
	keyLen := len(tx.PublicKey)
	x.SetBytes(tx.PublicKey[:(keyLen / 2)])
	y.SetBytes(tx.PublicKey[(keyLen / 2):])
	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	txHash, err := tx.CalculateHash()
	if err != nil {
		return false, err
	}
	return ecdsa.Verify(pubKey, txHash[:], r, s), nil
}
func (tx *Transaction) CalculateHash() ([32]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(struct {
		To    [20]byte
		Value uint64
		Nonce uint64
	}{
		To:    tx.To,
		Value: tx.Value,
		Nonce: tx.Nonce,
	})
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to encode transaction for hashing: %w", err)
	}
	return sha256.Sum256(buf.Bytes()), nil
}
func (tx *Transaction) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(tx); err != nil {
		return nil, fmt.Errorf("failed to encode transaction: %w", err)
	}
	return buf.Bytes(), nil
}
func DecodeTransaction(data []byte) (*Transaction, error) {
	var tx Transaction
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&tx); err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}
	return &tx, nil
}
