package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/Modulax-Protocol/go-modulax/network"
)

// RPCRequest represents a JSON-RPC request.
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	ID      int           `json:"id"`
	Params  []interface{} `json:"params"`
}

// RPCResponse represents a JSON-RPC response.
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error object.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// APIServer handles JSON-RPC requests.
type APIServer struct {
	bc     *core.Blockchain
	pubsub *network.PubSubService
	txPool *core.TxPool
}

// NewAPIServer creates a new APIServer instance.
func NewAPIServer(bc *core.Blockchain, pubsub *network.PubSubService, txPool *core.TxPool) *APIServer {
	return &APIServer{
		bc:     bc,
		pubsub: pubsub,
		txPool: txPool,
	}
}

// handleRPC processes incoming JSON-RPC requests.
func (s *APIServer) handleRPC(w http.ResponseWriter, r *http.Request) {
	var req RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
		return
	}

	var resp RPCResponse
	resp.JSONRPC = "2.0"
	resp.ID = req.ID

	switch req.Method {
	case "getBlockHeight":
		latestBlock, err := s.bc.GetLatestBlock()
		if err != nil {
			resp.Error = &RPCError{Code: -32001, Message: "Could not get latest block"}
		} else {
			resp.Result = latestBlock.Header.Height
		}

	case "sendTransaction":
		if len(req.Params) == 0 {
			resp.Error = &RPCError{Code: -32602, Message: "Invalid params"}
			break
		}
		// Assume the param is a simple string for the transaction data.
		txData, ok := req.Params[0].(string)
		if !ok {
			resp.Error = &RPCError{Code: -32602, Message: "Invalid params: expected a string"}
			break
		}

		tx := &core.Transaction{Data: []byte(txData)}
		tx.Sign()
		hash, _ := tx.CalculateHash()
		tx.Hash = hash

		if err := s.txPool.Add(tx); err != nil {
			resp.Error = &RPCError{Code: -32004, Message: "Failed to add transaction to pool"}
		} else {
			resp.Result = fmt.Sprintf("Transaction accepted: %x", tx.Hash)
		}

	case "addBlock":
		// Get pending transactions from the pool.
		pendingTxs := s.txPool.Pending()
		if len(pendingTxs) == 0 {
			resp.Error = &RPCError{Code: -32005, Message: "No pending transactions to add"}
			break
		}

		// Add a new block containing the pending transactions.
		newBlock, err := s.bc.AddBlock(pendingTxs)
		if err != nil {
			resp.Error = &RPCError{Code: -32002, Message: "Could not add block"}
		} else {
			// Clear the transaction pool after including the transactions in a block.
			s.txPool.Clear()

			blockBytes, err := newBlock.Encode()
			if err != nil {
				resp.Error = &RPCError{Code: -32003, Message: "Could not encode block for broadcast"}
			} else {
				if err := s.pubsub.BroadcastBlock(context.Background(), blockBytes); err != nil {
					fmt.Printf("Error broadcasting block: %v\n", err)
				}
				resp.Result = fmt.Sprintf("%x", newBlock.Hash)
			}
		}

	default:
		resp.Error = &RPCError{Code: -32601, Message: "Method not found"}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

