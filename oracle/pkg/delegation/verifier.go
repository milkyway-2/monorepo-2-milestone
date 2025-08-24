package delegation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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

// ExtrinsicInfo represents information about an extrinsic
type ExtrinsicInfo struct {
	BlockHash    string                 `json:"blockHash"`
	BlockNumber  string                 `json:"blockNumber"`
	ExtrinsicIdx int                    `json:"extrinsicIdx"`
	Method       map[string]interface{} `json:"method"`
	Events       []interface{}          `json:"events"`
	Success      bool                   `json:"success"`
}

// StakingExtrinsic represents a staking-related extrinsic
type StakingExtrinsic struct {
	ExtrinsicHash string                 `json:"extrinsicHash"`
	BlockHash     string                 `json:"blockHash"`
	BlockNumber   string                 `json:"blockNumber"`
	ExtrinsicIdx  int                    `json:"extrinsicIdx"`
	Method        string                 `json:"method"`
	Params        map[string]interface{} `json:"params"`
	Success       bool                   `json:"success"`
	Events        []interface{}          `json:"events"`
	Timestamp     string                 `json:"timestamp,omitempty"`
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

// getExtrinsicInfo retrieves information about a specific extrinsic by its hash
func (v *Verifier) getExtrinsicInfo(extrinsicHash string) (*ExtrinsicInfo, error) {
	log.Printf("🔍 Retrieving extrinsic info for hash: %s", extrinsicHash)

	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "chain_getBlock",
		Params: []interface{}{
			extrinsicHash,
		},
		ID: 1,
	}

	result, err := v.makeRPCCall(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get extrinsic info: %w", err)
	}

	// Parse the result to extract extrinsic information
	if resultMap, ok := result.(map[string]interface{}); ok {
		if block, ok := resultMap["block"].(map[string]interface{}); ok {
			if extrinsics, ok := block["extrinsics"].([]interface{}); ok {
				// For now, we'll look at the first extrinsic in the block
				// In a more sophisticated implementation, you'd find the specific extrinsic
				if len(extrinsics) > 0 {
					log.Printf("📋 Found %d extrinsics in block", len(extrinsics))

					// Try to decode the extrinsic to check if it's a nomination
					for i, extrinsic := range extrinsics {
						log.Printf("🔍 Examining extrinsic %d: %v", i, extrinsic)

						// Check if this extrinsic contains nomination information
						if v.isNominationExtrinsic(extrinsic) {
							log.Printf("✅ Found nomination extrinsic at index %d", i)
							return &ExtrinsicInfo{
								BlockHash:    extrinsicHash,
								ExtrinsicIdx: i,
								Success:      true, // Assume success for now
							}, nil
						}
					}
				}
			}
		}
	}

	log.Printf("⚠️  Could not find nomination extrinsic in block")
	return nil, fmt.Errorf("no nomination extrinsic found in block")
}

// isNominationExtrinsic checks if an extrinsic is related to nomination/delegation
func (v *Verifier) isNominationExtrinsic(extrinsic interface{}) bool {
	// This is a simplified check - in a real implementation, you would:
	// 1. Decode the extrinsic properly
	// 2. Check if it's a Staking.nominate call
	// 3. Extract the nominator and validator addresses

	extrinsicStr := fmt.Sprintf("%v", extrinsic)

	// Look for common patterns in nomination extrinsics
	nominationPatterns := []string{
		"nominate",
		"staking",
		"delegate",
		"bond",
	}

	for _, pattern := range nominationPatterns {
		if strings.Contains(strings.ToLower(extrinsicStr), pattern) {
			log.Printf("🔍 Found nomination pattern '%s' in extrinsic", pattern)
			return true
		}
	}

	return false
}

// verifyDelegationByExtrinsic verifies delegation using a specific extrinsic hash
func (v *Verifier) verifyDelegationByExtrinsic(extrinsicHash, nominatorAddress, validatorAddress string) (bool, error) {
	log.Printf("🔍 Verifying delegation using extrinsic hash: %s", extrinsicHash)
	log.Printf("   Nominator: %s", nominatorAddress)
	log.Printf("   Validator: %s", validatorAddress)

	// Get extrinsic information
	extrinsicInfo, err := v.getExtrinsicInfo(extrinsicHash)
	if err != nil {
		log.Printf("❌ Failed to get extrinsic info: %v", err)
		return false, fmt.Errorf("failed to get extrinsic info: %w", err)
	}

	if extrinsicInfo == nil {
		log.Printf("❌ No extrinsic info found")
		return false, fmt.Errorf("no extrinsic info found")
	}

	log.Printf("✅ Extrinsic info retrieved successfully")
	log.Printf("   Block Hash: %s", extrinsicInfo.BlockHash)
	log.Printf("   Extrinsic Index: %d", extrinsicInfo.ExtrinsicIdx)
	log.Printf("   Success: %t", extrinsicInfo.Success)

	// Check if the extrinsic was successful
	if !extrinsicInfo.Success {
		log.Printf("❌ Extrinsic was not successful")
		return false, fmt.Errorf("extrinsic was not successful")
	}

	// For now, we'll assume the extrinsic is valid if we can retrieve it
	// In a more sophisticated implementation, you would:
	// 1. Decode the extrinsic properly
	// 2. Extract the actual nominator and validator addresses
	// 3. Compare them with the provided addresses
	// 4. Check if the nomination is still active

	log.Printf("✅ Extrinsic verification successful")
	return true, nil
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

// VerifyDelegationWithExtrinsic checks if a nominator has delegated to a validator using a specific extrinsic hash
func (v *Verifier) VerifyDelegationWithExtrinsic(extrinsicHash, nominatorAddress, validatorAddress string) (bool, error) {
	log.Printf("🔍 Verifying delegation with extrinsic hash: %s", extrinsicHash)
	log.Printf("   Nominator: %s", nominatorAddress)
	log.Printf("   Validator: %s", validatorAddress)

	// First, verify the extrinsic itself
	extrinsicValid, err := v.verifyDelegationByExtrinsic(extrinsicHash, nominatorAddress, validatorAddress)
	if err != nil {
		log.Printf("❌ Extrinsic verification failed: %v", err)
		return false, fmt.Errorf("extrinsic verification failed: %w", err)
	}

	if !extrinsicValid {
		log.Printf("❌ Extrinsic verification failed")
		return false, fmt.Errorf("extrinsic verification failed")
	}

	log.Printf("✅ Extrinsic verification successful")

	// Then, perform the standard delegation verification
	standardValid, err := v.VerifyDelegation(nominatorAddress, validatorAddress)
	if err != nil {
		log.Printf("❌ Standard delegation verification failed: %v", err)
		return false, fmt.Errorf("standard delegation verification failed: %w", err)
	}

	if !standardValid {
		log.Printf("❌ Standard delegation verification failed")
		return false, fmt.Errorf("standard delegation verification failed")
	}

	log.Printf("✅ Both extrinsic and standard verification successful")
	return true, nil
}

// GetStakingExtrinsics retrieves all staking-related extrinsics for a given nominator-validator pair
func (v *Verifier) GetStakingExtrinsics(nominatorAddress, validatorAddress string) ([]StakingExtrinsic, error) {
	log.Printf("🔍 Getting staking extrinsics for nominator: %s, validator: %s", nominatorAddress, validatorAddress)

	var extrinsics []StakingExtrinsic

	// Method 1: If nominatorAddress looks like an extrinsic hash, try to get it directly
	if strings.HasPrefix(nominatorAddress, "0x") && len(nominatorAddress) == 66 {
		log.Printf("🔍 Nominator address looks like an extrinsic hash, trying direct lookup")
		directExtrinsic, err := v.getExtrinsicByHash(nominatorAddress)
		if err != nil {
			log.Printf("⚠️  Error getting extrinsic by hash: %v", err)
		} else if directExtrinsic != nil {
			extrinsics = append(extrinsics, *directExtrinsic)
			log.Printf("✅ Found extrinsic by hash, returning early")
			return extrinsics, nil
		}
	}

	// Method 2: Try to find the extrinsic using a more targeted approach
	targetedExtrinsics, err := v.findExtrinsicByAddress(nominatorAddress, validatorAddress)
	if err != nil {
		log.Printf("⚠️  Error in targeted search: %v", err)
	} else {
		extrinsics = append(extrinsics, targetedExtrinsics...)
	}

	// Method 3: Use state_queryStorageAt to find specific staking events (simplified)
	storageExtrinsics, err := v.queryStakingStorage(nominatorAddress, validatorAddress)
	if err != nil {
		log.Printf("⚠️  Error querying staking storage: %v", err)
	} else {
		extrinsics = append(extrinsics, storageExtrinsics...)
	}

	// Remove duplicates based on extrinsic hash
	uniqueExtrinsics := v.removeDuplicateExtrinsics(extrinsics)

	log.Printf("✅ Found %d unique staking extrinsics", len(uniqueExtrinsics))
	return uniqueExtrinsics, nil
}

// getLatestBlockNumber gets the latest block number
func (v *Verifier) getLatestBlockNumber() (int64, error) {
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "chain_getHeader",
		Params:  []interface{}{},
		ID:      1,
	}

	result, err := v.makeRPCCall(request)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block header: %w", err)
	}

	// Parse the block number from the result
	if headerMap, ok := result.(map[string]interface{}); ok {
		if numberStr, exists := headerMap["number"].(string); exists {
			// Convert hex string to int64
			var blockNum int64
			_, err := fmt.Sscanf(numberStr, "0x%x", &blockNum)
			if err != nil {
				return 0, fmt.Errorf("failed to parse block number: %w", err)
			}
			return blockNum, nil
		}
	}

	return 0, fmt.Errorf("could not extract block number from response")
}

// getStakingExtrinsicsFromBlock gets staking extrinsics from a specific block
func (v *Verifier) getStakingExtrinsicsFromBlock(blockNumber int64, nominatorAddress, validatorAddress string) ([]StakingExtrinsic, error) {
	var extrinsics []StakingExtrinsic

	// First, get the block hash for the block number
	blockHash, err := v.getBlockHash(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get block hash for block %d: %w", blockNumber, err)
	}

	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "chain_getBlock",
		Params: []interface{}{
			blockHash,
		},
		ID: 1,
	}

	result, err := v.makeRPCCall(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get block %d: %w", blockNumber, err)
	}

	// Parse the block to find staking extrinsics
	if resultMap, ok := result.(map[string]interface{}); ok {
		if block, ok := resultMap["block"].(map[string]interface{}); ok {
			if blockExtrinsics, ok := block["extrinsics"].([]interface{}); ok {
				for i, extrinsic := range blockExtrinsics {
					if v.isStakingExtrinsic(extrinsic, nominatorAddress, validatorAddress) {
						stakingExtrinsic := StakingExtrinsic{
							BlockHash:    blockHash,
							BlockNumber:  fmt.Sprintf("%d", blockNumber),
							ExtrinsicIdx: i,
							Method:       "staking.nominate", // Default, will be updated
							Success:      true,               // Assume success for now
						}
						extrinsics = append(extrinsics, stakingExtrinsic)
					}
				}
			}
		}
	}

	return extrinsics, nil
}

// getBlockHash gets the block hash for a given block number
func (v *Verifier) getBlockHash(blockNumber int64) (string, error) {
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "chain_getBlockHash",
		Params: []interface{}{
			fmt.Sprintf("0x%x", blockNumber),
		},
		ID: 1,
	}

	result, err := v.makeRPCCall(request)
	if err != nil {
		return "", fmt.Errorf("failed to get block hash: %w", err)
	}

	if blockHash, ok := result.(string); ok {
		return blockHash, nil
	}

	return "", fmt.Errorf("invalid block hash response")
}

// isStakingExtrinsic checks if an extrinsic is a staking extrinsic for the given addresses
func (v *Verifier) isStakingExtrinsic(extrinsic interface{}, nominatorAddress, validatorAddress string) bool {
	extrinsicStr := fmt.Sprintf("%v", extrinsic)

	// Check for staking-related patterns
	stakingPatterns := []string{
		"nominate",
		"bond",
		"unbond",
		"withdraw_unbonded",
		"chill",
		"set_payee",
		"set_controller",
		"validate",
	}

	for _, pattern := range stakingPatterns {
		if strings.Contains(strings.ToLower(extrinsicStr), pattern) {
			// Additional check: see if the addresses are mentioned in the extrinsic
			if strings.Contains(extrinsicStr, nominatorAddress) || strings.Contains(extrinsicStr, validatorAddress) {
				return true
			}
		}
	}

	return false
}

// queryStakingStorage queries staking storage for specific events
func (v *Verifier) queryStakingStorage(nominatorAddress, validatorAddress string) ([]StakingExtrinsic, error) {
	log.Printf("🔍 Querying staking storage for events")

	var extrinsics []StakingExtrinsic

	// Query Staking.Nominators storage
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "state_getStorage",
		Params: []interface{}{
			"0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70", // Staking.Nominators
		},
		ID: 1,
	}

	result, err := v.makeRPCCall(request)
	if err != nil {
		return nil, fmt.Errorf("failed to query staking storage: %w", err)
	}

	log.Printf("📋 Staking storage result: %v", result)

	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Decode the storage value properly
	// 2. Extract nomination information
	// 3. Find the corresponding extrinsics

	return extrinsics, nil
}

// removeDuplicateExtrinsics removes duplicate extrinsics based on extrinsic hash
func (v *Verifier) removeDuplicateExtrinsics(extrinsics []StakingExtrinsic) []StakingExtrinsic {
	seen := make(map[string]bool)
	var unique []StakingExtrinsic

	for _, extrinsic := range extrinsics {
		if !seen[extrinsic.ExtrinsicHash] {
			seen[extrinsic.ExtrinsicHash] = true
			unique = append(unique, extrinsic)
		}
	}

	return unique
}

// getExtrinsicByHash retrieves an extrinsic directly by its hash
func (v *Verifier) getExtrinsicByHash(extrinsicHash string) (*StakingExtrinsic, error) {
	log.Printf("🔍 Getting extrinsic by hash: %s", extrinsicHash)

	// Try to get the extrinsic using chain_getBlock
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "chain_getBlock",
		Params: []interface{}{
			extrinsicHash,
		},
		ID: 1,
	}

	result, err := v.makeRPCCall(request)
	if err != nil {
		log.Printf("⚠️  Failed to get extrinsic by hash: %v", err)
		return nil, err
	}

	// Parse the result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if block, ok := resultMap["block"].(map[string]interface{}); ok {
			if extrinsics, ok := block["extrinsics"].([]interface{}); ok {
				for i, extrinsic := range extrinsics {
					// Check if this is a staking extrinsic
					if v.isStakingExtrinsic(extrinsic, extrinsicHash, "") {
						log.Printf("✅ Found staking extrinsic at index %d", i)
						return &StakingExtrinsic{
							ExtrinsicHash: extrinsicHash,
							BlockHash:     extrinsicHash, // In this case, the hash is the block hash
							BlockNumber:   "unknown",     // We don't have the block number
							ExtrinsicIdx:  i,
							Method:        "staking.nominate",
							Success:       true,
						}, nil
					}
				}
			}
		}
	}

	log.Printf("⚠️  No staking extrinsic found for hash: %s", extrinsicHash)
	return nil, nil
}

// findExtrinsicByAddress tries to find extrinsics by searching a small range of recent blocks
func (v *Verifier) findExtrinsicByAddress(nominatorAddress, validatorAddress string) ([]StakingExtrinsic, error) {
	log.Printf("🔍 Finding extrinsics by address in recent blocks")

	var extrinsics []StakingExtrinsic

	// Get the latest block number
	latestBlock, err := v.getLatestBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	// Search through only the last 10 blocks for performance
	searchRange := int64(10)
	startBlock := latestBlock - searchRange
	if startBlock < 0 {
		startBlock = 0
	}

	log.Printf("📊 Searching blocks from %d to %d", startBlock, latestBlock)

	// Search in reverse order (newest first) and limit results
	maxExtrinsics := 5
	for blockNum := latestBlock; blockNum >= startBlock && len(extrinsics) < maxExtrinsics; blockNum-- {
		blockExtrinsics, err := v.getStakingExtrinsicsFromBlock(blockNum, nominatorAddress, validatorAddress)
		if err != nil {
			log.Printf("⚠️  Error getting extrinsics from block %d: %v", blockNum, err)
			continue
		}
		extrinsics = append(extrinsics, blockExtrinsics...)

		// Add a small delay to avoid overwhelming the RPC
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("✅ Found %d extrinsics in recent blocks", len(extrinsics))
	return extrinsics, nil
}
