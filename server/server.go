package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/network"
	"github.com/Modulax-Protocol/go-modulax/storage"
	"github.com/gorilla/mux"
)

// Server holds the state of the node.
type Server struct {
	db      storage.Storer
	bc      *core.Blockchain
	txPool  *core.TxPool // Add the transaction pool
	p2pNode *network.Node
	pubsub  *network.PubSubService
	apiPort string
}

// NewServer creates a new server instance.
func NewServer(db storage.Storer, bc *core.Blockchain, p2pNode *network.Node, apiPort string) (*Server, error) {
	ctx := context.Background()
	pubsub, err := network.NewPubSubService(ctx, p2pNode.Host())
	if err != nil {
		return nil, err
	}

	return &Server{
		db:      db,
		bc:      bc,
		txPool:  core.NewTxPool(), // Initialize the new transaction pool
		p2pNode: p2pNode,
		pubsub:  pubsub,
		apiPort: apiPort,
	}, nil
}

// Start runs the server.
func (s *Server) Start(bootstrapNode string) error {
	fmt.Println("Starting Modulax node...")
	s.p2pNode.Start()

	// If a bootstrap node is provided, connect to it.
	if bootstrapNode != "" {
		if err := s.p2pNode.Connect(context.Background(), bootstrapNode); err != nil {
			fmt.Printf("Failed to connect to bootstrap node: %v\n", err)
		}
	}

	// Subscribe to block announcements from the network.
	_, err := s.pubsub.Subscribe(s.handleNewBlock)
	if err != nil {
		return fmt.Errorf("failed to subscribe to block topic: %w", err)
	}

	// Start the JSON-RPC API server.
	router := mux.NewRouter()
	api := NewAPIServer(s.bc, s.pubsub, s.txPool) // Pass the txPool to the API server
	router.HandleFunc("/rpc", api.handleRPC).Methods("POST")

	latestBlock, _ := s.bc.GetLatestBlock()
	fmt.Printf("Current block height: %d\n", latestBlock.Header.Height)
	fmt.Printf("JSON-RPC server listening on %s\n", s.apiPort)

	return http.ListenAndServe(s.apiPort, router)
}

// handleNewBlock is the callback function for when a new block is received from the network.
func (s *Server) handleNewBlock(data []byte) {
	fmt.Println("\n--- 📣 Received New Block Message from Network! ---")

	block, err := core.DecodeBlock(data)
	if err != nil {
		fmt.Printf("[DEBUG] Error decoding block from network: %v\n", err)
		return
	}

	err = s.bc.AddExistingBlock(block)
	if err != nil {
		fmt.Printf("[DEBUG] Error adding block from network to our chain: %v\n", err)
		return
	}

	newLatestBlock, _ := s.bc.GetLatestBlock()
	fmt.Printf("--- ✅ Successfully Synced Block! New Height: %d ---\n\n", newLatestBlock.Header.Height)
}

