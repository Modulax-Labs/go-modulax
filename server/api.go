package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Modulax-Protocol/go-modulax/core"
	"github.com/gorilla/mux"
)

// APIServer is responsible for handling JSON-RPC requests.
type APIServer struct {
	listenAddr string
	bc         *core.Blockchain
}

// NewAPIServer creates a new instance of the APIServer.
func NewAPIServer(listenAddr string, bc *core.Blockchain) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		bc:         bc,
	}
}

// Run starts the HTTP server.
func (s *APIServer) Run() error {
	router := mux.NewRouter()

	// The main endpoint for all RPC requests.
	router.HandleFunc("/rpc", s.handleRPC)

	fmt.Printf("JSON-RPC server listening on %s\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

// handleRPC is the central handler for all incoming JSON-RPC requests.
func (s *APIServer) handleRPC(w http.ResponseWriter, r *http.Request) {
	var req RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	// Route the request to the correct method handler.
	switch req.Method {
	case "getBlockHeight":
		s.handleGetBlockHeight(w, req)
	case "addBlock":
		s.handleAddBlock(w, req)
	default:
		writeError(w, http.StatusBadRequest, fmt.Sprintf("method %s not found", req.Method))
	}
}

// handleGetBlockHeight handles the getBlockHeight RPC method.
func (s *APIServer) handleGetBlockHeight(w http.ResponseWriter, req RPCRequest) {
	latestBlock, err := s.bc.GetLatestBlock()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  latestBlock.Header.Height,
	}

	writeResponse(w, res)
}

// handleAddBlock handles the addBlock RPC method.
func (s *APIServer) handleAddBlock(w http.ResponseWriter, req RPCRequest) {
	// For now, we create a dummy transaction to include in the new block.
	// In the future, this would come from a transaction pool (mempool).
	dummyTx := &core.Transaction{
		From:      []byte("sender"),
		To:        []byte("receiver"),
		Value:     100,
		Timestamp: time.Now().UnixNano(),
		Nonce:     0,
	}
	// Sign and hash the transaction
	dummyTx.Sign()
	hash, _ := dummyTx.CalculateHash()
	dummyTx.Hash = hash

	newBlock, err := s.bc.AddBlock([]*core.Transaction{dummyTx})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the hash of the newly created block as confirmation.
	res := RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  fmt.Sprintf("%x", newBlock.Hash), // Return hash as a hex string
	}

	writeResponse(w, res)
}

// writeResponse is a helper function to marshal and write a JSON-RPC response.
func writeResponse(w http.ResponseWriter, res RPCResponse) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
	}
}

// writeError is a helper function to write a JSON-RPC error response.
func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	res := RPCResponse{
		JSONRPC: "2.0",
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
	writeResponse(w, res)
}

// RPCRequest defines the structure for an incoming JSON-RPC request.
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// RPCResponse defines the structure for a JSON-RPC response.
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError defines the structure for a JSON-RPC error object.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

