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
	"strings"

	"golang.org/x/crypto/ripemd160"
)

const (
	walletDir = "./wallets"
)

var (
	HexToCipherMap = map[rune]string{
		'0': "GG", '1': "GH", '2': "GJ", '3': "GK",
		'4': "HG", '5': "HH", '6': "HJ", '7': "HK",
		'8': "JG", '9': "JH",
		'a': "JJ", 'b': "JK", 'c': "KG", 'd': "KH",
		'e': "KJ", 'f': "KK",
	}
	CipherToHexMap = map[string]rune{
		"GG": '0', "GH": '1', "GJ": '2', "GK": '3',
		"HG": '4', "HH": '5', "HJ": '6', "HK": '7',
		"JG": '8', "JH": '9',
		"JJ": 'a', "JK": 'b', "KG": 'c', "KH": 'd',
		"KJ": 'e', "KK": 'f',
	}
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  []byte
}

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

func WalletFromPrivateKey(privateKey *ecdsa.PrivateKey) *Wallet {
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

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

func (w *Wallet) PublicKey() []byte {
	return w.publicKey
}

func (w *Wallet) Address() [20]byte {
	pubKeyHash := sha256.Sum256(w.publicKey)
	hasher := ripemd160.New()
	hasher.Write(pubKeyHash[:])
	ripeHash := hasher.Sum(nil)
	var address [20]byte
	copy(address[:], ripeHash)
	return address
}

func (w *Wallet) Sign(dataHash [32]byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, dataHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

func (w *Wallet) SaveToFile() (string, error) {
	if err := os.MkdirAll(walletDir, os.ModePerm); err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s/%x.wal", walletDir, w.Address())
	privateKeyBytes := w.privateKey.D.Bytes()
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	return fileName, os.WriteFile(fileName, []byte(privateKeyHex), 0644)
}

func AddressFromPublicKey(pubKey []byte) [20]byte {
	pubKeyHash := sha256.Sum256(pubKey)
	hasher := ripemd160.New()
	hasher.Write(pubKeyHash[:])
	ripeHash := hasher.Sum(nil)
	var address [20]byte
	copy(address[:], ripeHash)
	return address
}

func EncodeToCipher(hexString string) string {
	var builder strings.Builder
	for _, char := range strings.ToLower(hexString) {
		if cipherCode, ok := HexToCipherMap[char]; ok {
			builder.WriteString(cipherCode)
		} else {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func DecodeFromCipher(cipherString string) (string, error) {
	if len(cipherString)%2 != 0 {
		return "", fmt.Errorf("invalid cipher string length")
	}
	var builder strings.Builder
	for i := 0; i < len(cipherString); i += 2 {
		chunk := cipherString[i : i+2]
		if hexChar, ok := CipherToHexMap[chunk]; ok {
			builder.WriteRune(hexChar)
		} else {
			return "", fmt.Errorf("unknown cipher code: %s", chunk)
		}
	}
	return builder.String(), nil
}
