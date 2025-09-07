package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Modulax-Protocol/go-modulax/core"
)

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	ID      int           `json:"id"`
	Params  []interface{} `json:"params,omitempty"`
}
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type Client struct {
	endpoint string
}

func New(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}
func (c *Client) GetAccount(address string) (*core.Account, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		Method:  "getAccount",
		ID:      1,
		Params:  []interface{}{address},
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpResp, err := http.Post(c.endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	var resp RPCResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("rpc error: %s", resp.Error.Message)
	}
	accountData, ok := resp.Result.(map[string]interface{})
	if !ok {
		if resp.Result == nil {
			return &core.Account{Balance: 0, Nonce: 0}, nil
		}
		return nil, fmt.Errorf("unexpected type for account result: %T", resp.Result)
	}
	acc := &core.Account{}
	if balance, ok := accountData["Balance"].(float64); ok {
		acc.Balance = uint64(balance)
	}
	if nonce, ok := accountData["Nonce"].(float64); ok {
		acc.Nonce = uint64(nonce)
	}
	return acc, nil
}
func (c *Client) SendTransaction(tx *core.Transaction) (string, error) {
	txBytes, err := tx.Encode()
	if err != nil {
		return "", err
	}
	txHex := hex.EncodeToString(txBytes)
	req := RPCRequest{
		JSONRPC: "2.0",
		Method:  "sendTransaction",
		ID:      2,
		Params:  []interface{}{txHex},
	}
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	httpResp, err := http.Post(c.endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()
	var resp RPCResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("rpc error: %s", resp.Error.Message)
	}
	return resp.Result.(string), nil
}
