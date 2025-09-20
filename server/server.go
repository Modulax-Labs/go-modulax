package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/miner"
	"github.com/Modulax-Protocol/go-modulax/network"
	"github.com/Modulax-Protocol/go-modulax/storage"
	"github.com/gorilla/mux"
)

// Server holds the state of the node.
type Server struct {
	db      storage.Storer
	bc      *core.Blockchain
	txPool  *core.TxPool
	miner   *miner.Miner
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

	txPool := core.NewTxPool(bc.State())
	miner := miner.NewMiner(bc, txPool, pubsub)

	return &Server{
		db:      db,
		bc:      bc,
		txPool:  txPool,
		miner:   miner,
		p2pNode: p2pNode,
		pubsub:  pubsub,
		apiPort: apiPort,
	}, nil
}

// Start runs the server.
func (s *Server) Start(bootstrapNode string) error {
	fmt.Println("Starting Modulax node...")
	s.p2pNode.Start()

	// Register the handlers for incoming network messages.
	s.pubsub.RegisterBlockHandler(s.handleNewBlock)
	s.pubsub.RegisterTxHandler(s.handleNewTransaction)

	// Start the pub/sub service to begin listening on topics.
	s.pubsub.Start()

	// Start the automated block proposer.
	s.miner.Start()

	// If a bootstrap node is provided, connect to it.
	if bootstrapNode != "" {
		if err := s.p2pNode.Connect(context.Background(), bootstrapNode); err != nil {
			fmt.Printf("Failed to connect to bootstrap node: %v\n", err)
		}
	}

	// Start the JSON-RPC API server.
	router := mux.NewRouter()
	api := NewAPIServer(s.bc, s.pubsub, s.txPool)
	router.HandleFunc("/rpc", api.handleRPC).Methods("POST")

	latestBlock, _ := s.bc.GetLatestBlock()
	fmt.Printf("Current block height: %d\n", latestBlock.Header.Height)
	fmt.Printf("JSON-RPC server listening on %s\n", s.apiPort)

	return http.ListenAndServe(s.apiPort, router)
}

// handleNewBlock is the callback function for when a new block is received.
func (s *Server) handleNewBlock(data []byte) {
	fmt.Println("\n--- 📣 Received New Block Message from Network! ---")
	block, err := core.DecodeBlock(data)
	if err != nil {
		fmt.Printf("[DEBUG] Error decoding block: %v\n", err)
		return
	}
	if err := s.bc.AddExistingBlock(block); err != nil {
		fmt.Printf("[DEBUG] Error adding block: %v\n", err)
		return
	}
	// When we sync a new block, we should clear our transaction pool.
	s.txPool.Clear()
	newLatestBlock, _ := s.bc.GetLatestBlock()
	fmt.Printf("--- ✅ Successfully Synced Block! New Height: %d ---\n\n", newLatestBlock.Header.Height)
}

// handleNewTransaction is the callback for when a new tx is received.
func (s *Server) handleNewTransaction(data []byte) {
	tx, err := core.DecodeTransaction(data)
	if err != nil {
		fmt.Printf("Error decoding transaction: %v\n", err)
		return
	}
	if err := s.txPool.Add(tx); err != nil {
		// This is the primary validation point. An error here means the tx was invalid.
		return
	}
	fmt.Printf("--- 📥 Received & Validated New Transaction! Hash: %x ---\n", tx.Hash)
}
