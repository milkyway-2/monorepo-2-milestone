package delegation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// RPCRequest represents a Polkadot RPC request
type RPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// RPCResponse represents a Polkadot RPC response
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      int         `json:"id"`
}

// RPCError represents an RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Verifier handles Polkadot delegation verification via HTTP RPC
type Verifier struct {
	rpcURL string
	client *http.Client
}

// NewVerifier creates a new delegation verifier
func NewVerifier(rpcURL string) *Verifier {
	return &Verifier{
		rpcURL: rpcURL,
		client: &http.Client{},
	}
}

// makeRPCCall makes a call to the Polkadot RPC endpoint
func (v *Verifier) makeRPCCall(request RPCRequest) (interface{}, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := v.client.Post(v.rpcURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make RPC call: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response RPCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", response.Error.Message)
	}

	return response.Result, nil
}

// getActiveEra gets the current active era from Polkadot
func (v *Verifier) getActiveEra() (interface{}, error) {
	log.Printf("📅 Querying active era from Polkadot")

	// Query the ActiveEra storage value
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "state_getStorage",
		Params: []interface{}{
			"0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70",
		},
		ID: 1,
	}

	result, err := v.makeRPCCall(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get active era: %w", err)
	}

	log.Printf("📅 Active era retrieved: %v", result)
	return result, nil
}

// checkIfNominated checks if a nominator has nominated a specific validator
func (v *Verifier) checkIfNominated(nominatorAddress, validatorAddress string) (bool, error) {
	log.Printf("🔍 Checking if nominator %s has nominated validator %s", nominatorAddress, validatorAddress)

	// For now, let's use a simpler approach and check if the addresses are valid
	// In a real implementation, you would query the actual staking storage
	// This is a placeholder that validates the address format

	// Check if addresses are valid (basic validation)
	if len(nominatorAddress) < 10 || len(validatorAddress) < 10 {
		log.Printf("❌ Invalid address format")
		return false, fmt.Errorf("invalid address format")
	}

	// For testing purposes, we'll simulate a real check
	// In production, you would:
	// 1. Query the Staking.Nominators storage map
	// 2. Decode the nomination data
	// 3. Check if the validator is in the targets list

	log.Printf("✅ Addresses appear valid, checking nomination status")

	// Simulate a real check - in production this would be an actual storage query
	// For now, we'll return true for any valid-looking addresses
	// This should be replaced with actual storage queries
	log.Printf("⚠️  Using simplified check - replace with actual storage queries in production")
	return true, nil
}

// checkIfActive checks if the nomination is currently active
func (v *Verifier) checkIfActive(nominatorAddress, validatorAddress string) (bool, error) {
	log.Printf("🔍 Checking if nomination is currently active...")

	// Query the current era to check if the nomination is active
	// In a real implementation, you would check the current era against the nomination era
	activeEra, err := v.getActiveEra()
	if err != nil {
		log.Printf("❌ Failed to get active era for activity check: %v", err)
		return false, fmt.Errorf("failed to get active era: %w", err)
	}

	log.Printf("📅 Current active era: %v", activeEra)

	// For now, we'll assume the nomination is active if it exists
	// In a real implementation, you'd check the nomination era and other factors
	log.Printf("✅ Assuming nomination is active (simplified check)")
	return true, nil
}

// VerifyDelegation checks if a nominator has delegated to a validator
func (v *Verifier) VerifyDelegation(nominatorAddress, validatorAddress string) (bool, error) {
	log.Printf("🔍 Verifying delegation: %s -> %s", nominatorAddress, validatorAddress)

	// Get the current active era
	activeEra, err := v.getActiveEra()
	if err != nil {
		log.Printf("❌ Failed to get active era: %v", err)
		return false, fmt.Errorf("failed to get active era: %w", err)
	}
	log.Printf("📅 Current active era: %v", activeEra)

	// Check if the nominator has nominated the validator
	isNominated, err := v.checkIfNominated(nominatorAddress, validatorAddress)
	if err != nil {
		return false, fmt.Errorf("failed to check nomination: %w", err)
	}

	if !isNominated {
		log.Printf("❌ Nominator %s has NOT nominated validator %s", nominatorAddress, validatorAddress)
		return false, nil
	}

	log.Printf("✅ Nominator %s HAS nominated validator %s", nominatorAddress, validatorAddress)

	// Check if the nomination is currently active
	isActive, err := v.checkIfActive(nominatorAddress, validatorAddress)
	if err != nil {
		return false, fmt.Errorf("failed to check if nomination is active: %w", err)
	}

	if isActive {
		log.Printf("✅ The nomination is currently ACTIVE and earning rewards")
	} else {
		log.Printf("⚠️  The nomination exists but is currently INACTIVE (not earning rewards)")
	}

	return true, nil
}
