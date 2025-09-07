package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"golang.org/x/crypto/ripemd160"
)

const (
	walletDir = "./wallets"
)

// Wallet holds a private/public key pair.
type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  []byte
}

// NewWallet creates a new Wallet instance with a newly generated key pair.
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// WalletFromPrivateKey creates a Wallet from an existing private key.
func WalletFromPrivateKey(privateKey *ecdsa.PrivateKey) *Wallet {
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// LoadWallet loads a wallet from a .wal file.
func LoadWallet(address string) (*Wallet, error) {
	fileName := fmt.Sprintf("%s/%s.wal", walletDir, address)
	privateKeyHex, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read wallet file: %w", err)
	}

	privateKeyBytes, err := hex.DecodeString(string(privateKeyHex))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKeyBytes)

	return WalletFromPrivateKey(privateKey), nil
}

// PublicKey returns the public key of the wallet.
func (w *Wallet) PublicKey() []byte {
	return w.publicKey
}

// Address derives the 20-byte address from the public key.
func (w *Wallet) Address() [20]byte {
	pubKeyHash := sha256.Sum256(w.publicKey)
	hasher := ripemd160.New()
	hasher.Write(pubKeyHash[:])
	ripeHash := hasher.Sum(nil)
	var address [20]byte
	copy(address[:], ripeHash)
	return address
}

// Sign signs a hash of data with the wallet's private key.
func (w *Wallet) Sign(dataHash [32]byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, dataHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

// SaveToFile saves the wallet's private key to a file in the wallets directory.
func (w *Wallet) SaveToFile() (string, error) {
	if err := os.MkdirAll(walletDir, os.ModePerm); err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s/%x.wal", walletDir, w.Address())
	privateKeyBytes := w.privateKey.D.Bytes()
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	return fileName, os.WriteFile(fileName, []byte(privateKeyHex), 0644)
}

// AddressFromPublicKey derives a 20-byte address from a given public key.
func AddressFromPublicKey(pubKey []byte) [20]byte {
	pubKeyHash := sha256.Sum256(pubKey)
	hasher := ripemd160.New()
	hasher.Write(pubKeyHash[:])
	ripeHash := hasher.Sum(nil)
	var address [20]byte
	copy(address[:], ripeHash)
	return address
}

