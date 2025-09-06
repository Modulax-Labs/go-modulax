package server

import (
	"context"
	"encoding/hex"
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
	Params  []interface{} `json:"params,omitempty"`
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
	case "getAccount":
		if len(req.Params) < 1 {
			resp.Error = &RPCError{Code: -32602, Message: "Invalid params: requires address"}
			break
		}
		addrStr, _ := req.Params[0].(string)
		addrBytes, err := hex.DecodeString(addrStr)
		if err != nil || len(addrBytes) != 20 {
			resp.Error = &RPCError{Code: -32602, Message: "Invalid address format"}
			break
		}
		var address [20]byte
		copy(address[:], addrBytes)

		account := s.bc.State().GetAccount(address)
		resp.Result = account

	case "sendTransaction":
		if len(req.Params) < 1 {
			resp.Error = &RPCError{Code: -32602, Message: "Invalid params: requires raw tx hex"}
			break
		}
		txHex, _ := req.Params[0].(string)
		txBytes, err := hex.DecodeString(txHex)
		if err != nil {
			resp.Error = &RPCError{Code: -32602, Message: "Invalid transaction hex"}
			break
		}

		tx, err := core.DecodeTransaction(txBytes)
		if err != nil {
			resp.Error = &RPCError{Code: -32000, Message: "Failed to decode transaction"}
			break
		}

		valid, err := tx.Verify()
		if err != nil || !valid {
			resp.Error = &RPCError{Code: -32000, Message: "Invalid transaction signature"}
			break
		}

		if err := s.txPool.Add(tx); err != nil {
			resp.Error = &RPCError{Code: -32004, Message: "Failed to add transaction to pool"}
		} else {
			s.pubsub.BroadcastTransaction(context.Background(), txBytes)
			resp.Result = fmt.Sprintf("Transaction accepted: %x", tx.Hash)
		}

	default:
		resp.Error = &RPCError{Code: -32601, Message: "Method not found"}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

