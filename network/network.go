package network

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Options holds the configuration for the network node.
type Options struct {
	ListenAddress string
}

// Node is a wrapper around a libp2p host.
type Node struct {
	host host.Host
	id   peer.ID

}

// NewNode creates a new network node.
func NewNode(ctx context.Context, opts Options) (*Node, error) {
	// Create a new libp2p host.
	// 0.0.0.0 listens on all available interfaces.
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(opts.ListenAddress),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}

	// Add a diagnostic print statement to immediately confirm the Peer ID.
	fmt.Printf("✅ New node created with Peer ID: %s\n", host.ID())

	return &Node{
		host: host,
		id:   host.ID(),
	}, nil
}

// Start begins the node's network services.
func (n *Node) Start() {
	// Get the host's listening addresses.
	addrs := n.host.Addrs()
	fmt.Println("Node is listening on:")
	for _, addr := range addrs {
		// The address includes the node's unique Peer ID.
		fullAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("%s/p2p/%s", addr, n.host.ID()))
		if err != nil {
			fmt.Printf("Error creating full multiaddr: %v\n", err)
			continue
		}
		fmt.Println(fullAddr)
	}
}

// Connect connects the node to a given peer.
func (n *Node) Connect(ctx context.Context, peerAddr string) error {
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return fmt.Errorf("failed to parse multiaddr: %w", err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return fmt.Errorf("failed to get peer info from multiaddr: %w", err)
	}

	if err := n.host.Connect(ctx, *peerInfo); err != nil {
		return fmt.Errorf("failed to connect to peer: %w", err)
	}

	fmt.Printf("Successfully connected to peer: %s\n", peerInfo.ID)
	return nil
}


