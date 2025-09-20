package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/Modulax-Protocol/go-modulax/client"
	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/crypto"
	"github.com/Modulax-Protocol/go-modulax/evm"
	"github.com/Modulax-Protocol/go-modulax/network"
	"github.com/Modulax-Protocol/go-modulax/server"
	"github.com/Modulax-Protocol/go-modulax/storage"
	"github.com/spf13/cobra"
)

var (
	connectNode string
	apiPort     string
)

func init() {
	runCmd.Flags().StringVar(&connectNode, "connect", "", "Address of a peer to connect to")
	runCmd.Flags().StringVar(&apiPort, "apiport", ":8080", "Port for the JSON-RPC API server")
	walletCmd.AddCommand(newWalletCmd)
	walletCmd.AddCommand(balanceCmd)
	walletCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(walletCmd)
}

var rootCmd = &cobra.Command{
	Use:   "modulax",
	Short: "Modulax is a quantum-resistant blockchain node and wallet CLI",
}
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Manage your Modulax wallet",
}
var newWalletCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates and saves a new wallet key pair",
	Run: func(cmd *cobra.Command, args []string) {
		wallet, err := crypto.NewWallet()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create wallet: %v\n", err)
			os.Exit(1)
		}
		fileName, err := wallet.SaveToFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save wallet: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("🎉 New Modulax Wallet Created!")
		fmt.Printf("Address: %s\n", wallet.CipherAddress())
		fmt.Printf("Wallet saved to: %s\n", fileName)
	},
}
var balanceCmd = &cobra.Command{
	Use:   "balance [address]",
	Short: "Gets the balance of a given address in Modulax-Cipher format",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addressStr := args[0]
		hexAddress, err := crypto.DecodeFromCipher(addressStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid Modulax address format: %v\n", err)
			os.Exit(1)
		}
		client := client.New("http://localhost:8080/rpc")
		account, err := client.GetAccount(hexAddress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Balance for %s: %d\n", addressStr, account.Balance)
		fmt.Printf("Nonce for   %s: %d\n", addressStr, account.Nonce)
	},
}
var sendCmd = &cobra.Command{
	Use:   "send [from_address] [to_address] [amount]",
	Short: "Send tokens from one address to another (addresses in Cipher format)",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fromAddrCipher := args[0]
		toAddrCipher := args[1]
		amountStr := args[2]
		amount, err := strconv.ParseUint(amountStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid amount: %v\n", err)
			os.Exit(1)
		}
		senderWallet, err := crypto.LoadWallet(fromAddrCipher)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not load sender wallet: %v\n", err)
			os.Exit(1)
		}
		toHexAddress, err := crypto.DecodeFromCipher(toAddrCipher)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid recipient address format.\n")
			os.Exit(1)
		}
		toAddrBytes, _ := hex.DecodeString(toHexAddress)
		var toAddr [20]byte
		copy(toAddr[:], toAddrBytes)
		client := client.New("http://localhost:8080/rpc")
		fromAddrBytes := senderWallet.Address()
		fromHexAddress := hex.EncodeToString(fromAddrBytes[:])
		senderAccount, err := client.GetAccount(fromHexAddress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not get sender account details: %v\n", err)
			os.Exit(1)
		}
		tx := &core.Transaction{
			To:        toAddr,
			Value:     amount,
			Nonce:     senderAccount.Nonce,
			PublicKey: senderWallet.PublicKey(),
		}
		txHash, _ := tx.CalculateHash()
		signature, _ := senderWallet.Sign(txHash)
		tx.Signature = signature
		tx.Hash = txHash
		txHashStr, err := client.SendTransaction(tx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending transaction: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Transaction sent successfully!\nHash (hex): %s\n", txHashStr)
	},
}
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the Modulax node",
	Run: func(cmd *cobra.Command, args []string) {
		setupGenesisWallet()
		dbPath := "./modulax_chain"
		listenAddr := "/ip4/0.0.0.0/tcp/4001"

		if connectNode != "" {
			dbPath = "./modulax_chain_2"
			listenAddr = "/ip4/0.0.0.0/tcp/4002"
			if apiPort == ":8080" {
				apiPort = ":8081"
			}
		}

		db, err := storage.NewLevelDBStore(dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create db store: %v\n", err)
			os.Exit(1)
		}

		// Correct Initialization Order:
		state, err := core.NewState(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create state: %v\n", err)
			os.Exit(1)
		}
		executor := evm.NewEVM(state)
		bc, err := core.NewBlockchain(db, executor)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create blockchain: %v\n", err)
			os.Exit(1)
		}

		p2pOpts := network.Options{ListenAddress: listenAddr}
		p2pNode, err := network.NewNode(context.Background(), p2pOpts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create P2P node: %v\n", err)
			os.Exit(1)
		}

		srv, err := server.NewServer(db, bc, p2pNode, apiPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
			os.Exit(1)
		}

		if err := srv.Start(connectNode); err != nil {
			fmt.Fprintf(os.Stderr, "Server failed to start: %v\n", err)
			os.Exit(1)
		}
	},
}

func setupGenesisWallet() {
	pkBytes, err := hex.DecodeString(core.GENESIS_PRIVATE_KEY)
	if err != nil {
		panic(err)
	}
	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = new(big.Int).SetBytes(pkBytes)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(pkBytes)
	genesisWallet := crypto.WalletFromPrivateKey(privateKey)
	fileName := fmt.Sprintf("./wallets/%s.wal", genesisWallet.CipherAddress())
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		savedFile, err := genesisWallet.SaveToFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not save genesis wallet: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🔑 Genesis wallet file created at: %s\n", savedFile)
		fmt.Printf("   - Genesis Address: %s\n", genesisWallet.CipherAddress())
	}
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing root command: %v\n", err)
		os.Exit(1)
	}
}
