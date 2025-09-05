package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Modulax-Protocol/go-modulax/core"
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
	rootCmd.AddCommand(runCmd)
}

var rootCmd = &cobra.Command{
	Use:   "modulax",
	Short: "Modulax is a quantum-resistant blockchain node",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the Modulax node",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := "./modulax_chain"
		listenAddr := "/ip4/0.0.0.0/tcp/4001"

		// Adjust paths and ports for the second node if connecting
		// This also ensures the API port is different if not specified.
		if connectNode != "" {
			dbPath = "./modulax_chain_2"
			listenAddr = "/ip4/0.0.0.0/tcp/4002"
			if apiPort == ":8080" { // Only override if it's the default
				apiPort = ":8081"
			}
		}

		// Initialize storage
		db, err := storage.NewLevelDBStore(dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create db store: %v\n", err)
			os.Exit(1)
		}

		// Initialize blockchain
		bc, err := core.NewBlockchain(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create blockchain: %v\n", err)
			os.Exit(1)
		}

		// Initialize P2P node
		p2pOpts := network.Options{ListenAddress: listenAddr}
		p2pNode, err := network.NewNode(context.Background(), p2pOpts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create P2P node: %v\n", err)
			os.Exit(1)
		}

		// Initialize server
		srv, err := server.NewServer(db, bc, p2pNode, apiPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
			os.Exit(1)
		}

		// Start the server
		if err := srv.Start(connectNode); err != nil {
			fmt.Fprintf(os.Stderr, "Server failed to start: %v\n", err)
			os.Exit(1)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing root command: %v\n", err)
		os.Exit(1)
	}
}

