package network

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	BlockTopic = "/modulax/blocks/1.0.0"
	TxTopic    = "/modulax/txs/1.0.0" // New topic for transactions
)

// PubSubService manages the pub/sub logic for the node.
type PubSubService struct {
	ps           *pubsub.PubSub
	blockTopic   *pubsub.Topic
	txTopic      *pubsub.Topic // New field for the transaction topic
	selfID       peer.ID
	blockHandler func(data []byte)
	txHandler    func(data []byte)
}

// NewPubSubService creates and initializes a new PubSubService.
func NewPubSubService(ctx context.Context, host host.Host) (*PubSubService, error) {
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub service: %w", err)
	}

	blockTopic, err := ps.Join(BlockTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to join block topic: %w", err)
	}

	txTopic, err := ps.Join(TxTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to join tx topic: %w", err)
	}

	return &PubSubService{
		ps:         ps,
		blockTopic: blockTopic,
		txTopic:    txTopic,
		selfID:     host.ID(),
	}, nil
}

// RegisterBlockHandler sets the callback function for new blocks.
func (pss *PubSubService) RegisterBlockHandler(handler func(data []byte)) {
	pss.blockHandler = handler
}

// RegisterTxHandler sets the callback function for new transactions.
func (pss *PubSubService) RegisterTxHandler(handler func(data []byte)) {
	pss.txHandler = handler
}

// Start begins the subscription loops for all topics.
func (pss *PubSubService) Start() {
	go pss.subscribeLoop(pss.blockTopic, pss.blockHandler)
	go pss.subscribeLoop(pss.txTopic, pss.txHandler)
}

// subscribeLoop handles reading messages from a given topic subscription.
func (pss *PubSubService) subscribeLoop(topic *pubsub.Topic, handler func(data []byte)) {
	if handler == nil {
		return // Don't subscribe if no handler is registered
	}

	sub, err := topic.Subscribe()
	if err != nil {
		fmt.Printf("Error subscribing to topic %s: %v\n", topic.String(), err)
		return
	}

	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			fmt.Println("Error reading from subscription:", err)
			continue
		}
		// Don't process messages we sent ourselves.
		if msg.ReceivedFrom == pss.selfID {
			continue
		}
		handler(msg.Data)
	}
}

// BroadcastBlock publishes a new block to all subscribed peers.
func (pss *PubSubService) BroadcastBlock(ctx context.Context, data []byte) error {
	return pss.blockTopic.Publish(ctx, data)
}

// BroadcastTransaction publishes a new transaction to all subscribed peers.
func (pss *PubSubService) BroadcastTransaction(ctx context.Context, data []byte) error {
	return pss.txTopic.Publish(ctx, data)
}

