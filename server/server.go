package server

import (
	"context"
	"fmt"
	"time"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/network"
)

// Server is the main component of the node.
type Server struct {
	bc          *core.Blockchain
	apiServer   *APIServer
	p2pNode     *network.Node
	connectAddr string
}

// NewServer creates a new Server instance.
func NewServer(bc *core.Blockchain, apiPort string, connectAddr string) (*Server, error) {
	// Setup the P2P node.
	p2pOpts := network.Options{
		ListenAddress: "/ip4/0.0.0.0/tcp/4001", // Listen on TCP port 4001
	}
	// If we are connecting to another peer, listen on a different port.
	if connectAddr != "" {
		p2pOpts.ListenAddress = "/ip4/0.0.0.0/tcp/4002"
	}

	p2pNode, err := network.NewNode(context.Background(), p2pOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create p2p node: %w", err)
	}

	apiServer := NewAPIServer(apiPort, bc)
	return &Server{
		bc:          bc,
		apiServer:   apiServer,
		p2pNode:     p2pNode,
		connectAddr: connectAddr,
	}, nil
}

// Start begins the server, starting all its services.
func (s *Server) Start() {
	fmt.Println("Starting Modulax node...")
	s.p2pNode.Start() // Start the networking layer.

	// If a connect address is provided, try to connect to it.
	if s.connectAddr != "" {
		go func() {
			// Give the node a moment to start up before connecting.
			time.Sleep(1 * time.Second)
			if err := s.p2pNode.Connect(context.Background(), s.connectAddr); err != nil {
				fmt.Printf("Failed to connect to bootstrap node: %v\n", err)
			}
		}()
	}

	latestBlock, err := s.bc.GetLatestBlock()
	if err != nil {
		panic(fmt.Sprintf("Failed to get latest block: %v", err))
	}

	fmt.Printf("Current block height: %d\n", latestBlock.Header.Height)

	// Start the API server in a new goroutine so it doesn't block.
	go func() {
		if err := s.apiServer.Run(); err != nil {
			fmt.Printf("API Server failed to start: %v\n", err)
		}
	}()

	// The main loop will eventually handle P2P logic.
	select {} // Block forever.
}

