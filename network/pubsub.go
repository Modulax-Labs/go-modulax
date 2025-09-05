package network

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

const BlockTopic = "/modulax/blocks/1.0.0"

// PubSubService manages the pub/sub logic for the node.
type PubSubService struct {
	ps    *pubsub.PubSub
	topic *pubsub.Topic
}

// NewPubSubService creates and initializes a new PubSubService.
func NewPubSubService(ctx context.Context, host host.Host) (*PubSubService, error) {
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub service: %w", err)
	}

	topic, err := ps.Join(BlockTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to join block topic: %w", err)
	}

	return &PubSubService{
		ps:    ps,
		topic: topic,
	}, nil
}

// Subscribe allows the node to start listening for new blocks from peers.
func (pss *PubSubService) Subscribe(handler func(data []byte)) (*pubsub.Subscription, error) {
	sub, err := pss.topic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to block topic: %w", err)
	}

	// Start a background goroutine to read messages from the subscription.
	go func() {
		for {
			msg, err := sub.Next(context.Background())
			if err != nil {
				fmt.Println("Error reading from subscription:", err)
				continue
			}
			// Don't process messages we sent ourselves.
			if msg.ReceivedFrom == pss.ps.ID() {
				continue
			}
			handler(msg.Data)
		}
	}()

	return sub, nil
}

// BroadcastBlock publishes a new block to all subscribed peers.
func (pss *PubSubService) BroadcastBlock(ctx context.Context, data []byte) error {
	return pss.topic.Publish(ctx, data)
}
