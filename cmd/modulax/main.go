package main

import (
	"fmt"
	"os"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/server"
	"github.com/Modulax-Protocol/go-modulax/storage"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "modulax",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the Modulax node",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the flag values.
		connectAddr, _ := cmd.Flags().GetString("connect")
		apiPort, _ := cmd.Flags().GetString("apiport")

		// Define the path for the database.
		dbPath := "./modulax_chain"
		// If connecting to another node, use a different DB path to simulate a separate node.
		if connectAddr != "" {
			dbPath = "./modulax_chain_2"
		}

		// Create a new persistent storage instance.
		store, err := storage.NewLevelDBStore(dbPath)
		if err != nil {
			return fmt.Errorf("failed to create leveldb store: %w", err)
		}

		// Create a new blockchain with the persistent store.
		bc, err := core.NewBlockchain(store)
		if err != nil {
			return fmt.Errorf("failed to create blockchain: %w", err)
		}

		// Create and start the server.
		srv, err := server.NewServer(bc, apiPort, connectAddr)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		srv.Start()

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add flags to the "run" command.
	runCmd.Flags().String("connect", "", "Address of a peer to connect to")
	runCmd.Flags().String("apiport", ":8080", "Port for the JSON-RPC API server")
	rootCmd.AddCommand(runCmd)
}

func main() {
	Execute()
}

