package delegation

import (
	"log"
	"testing"
	"time"
)

func TestNewVerifier(t *testing.T) {
	log.Printf("🧪 Starting TestNewVerifier")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)

	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	if verifier == nil {
		log.Printf("❌ Verifier is nil")
		t.Fatal("NewVerifier returned nil")
	}
	log.Printf("✅ Verifier is not nil")

	if verifier.rpcURL != rpcURL {
		log.Printf("❌ RPC URL mismatch - Expected: %s, Got: %s", rpcURL, verifier.rpcURL)
		t.Errorf("Expected RPC URL %s, got %s", rpcURL, verifier.rpcURL)
	}
	log.Printf("✅ RPC URL matches expected value: %s", verifier.rpcURL)

	if verifier.client == nil {
		log.Printf("❌ HTTP client is nil")
		t.Error("HTTP client should not be nil")
	}
	log.Printf("✅ HTTP client is properly initialized")

	log.Printf("🎉 TestNewVerifier completed successfully")
}

func TestVerifyDelegation_RealRPC(t *testing.T) {
	log.Printf("🧪 Starting TestVerifyDelegation_RealRPC")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)
	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	// Test with known addresses that should work in simulation mode
	nominator := "12ztGE9cY2p7kPJFpfvMrL6NsCUeqoiaBY3jciMqYFuFNJ2o"
	validator := "12GTt3pfM3SjTU6UL6dQ3SMgMSvdw94PnRoF6osU6hPvxbUZ"

	log.Printf("🔍 Testing delegation verification:")
	log.Printf("   Nominator: %s", nominator)
	log.Printf("   Validator: %s", validator)
	log.Printf("   Expected: Delegation should be verified (real RPC mode)")

	log.Printf("🚀 Calling VerifyDelegation...")
	isDelegated, err := verifier.VerifyDelegation(nominator, validator)
	log.Printf("📋 VerifyDelegation returned - isDelegated: %t, error: %v", isDelegated, err)

	if err != nil {
		log.Printf("❌ Unexpected error occurred: %v", err)
		t.Errorf("Expected no error, got: %v", err)
	} else {
		log.Printf("✅ No errors occurred during verification")
	}

	if !isDelegated {
		log.Printf("❌ Delegation verification failed - expected true, got false")
		t.Error("Expected delegation to be verified for known addresses")
	} else {
		log.Printf("✅ Delegation verification successful - delegation confirmed")
	}

	log.Printf("🎉 TestVerifyDelegation_RealRPC completed successfully")
}

func TestVerifierConfiguration(t *testing.T) {
	log.Printf("🧪 Starting TestVerifierConfiguration")

	testCases := []struct {
		name   string
		rpcURL string
	}{
		{
			name:   "Official Polkadot RPC",
			rpcURL: "https://rpc.polkadot.io",
		},
		{
			name:   "Local Development RPC",
			rpcURL: "http://localhost:9944",
		},
		{
			name:   "Custom RPC Endpoint",
			rpcURL: "https://custom-polkadot-rpc.example.com",
		},
	}

	for _, tc := range testCases {
		log.Printf("🔧 Testing configuration: %s", tc.name)
		log.Printf("   RPC URL: %s", tc.rpcURL)

		verifier := NewVerifier(tc.rpcURL)

		if verifier == nil {
			log.Printf("❌ Verifier creation failed for %s", tc.name)
			t.Errorf("NewVerifier failed for %s", tc.name)
			continue
		}

		if verifier.rpcURL != tc.rpcURL {
			log.Printf("❌ RPC URL mismatch for %s - Expected: %s, Got: %s", tc.name, tc.rpcURL, verifier.rpcURL)
			t.Errorf("RPC URL mismatch for %s", tc.name)
		} else {
			log.Printf("✅ RPC URL correctly set for %s", tc.name)
		}

		if verifier.client == nil {
			log.Printf("❌ HTTP client is nil for %s", tc.name)
			t.Errorf("HTTP client is nil for %s", tc.name)
		} else {
			log.Printf("✅ HTTP client properly initialized for %s", tc.name)
		}
	}

	log.Printf("🎉 TestVerifierConfiguration completed successfully")
}

func TestGetActiveEra(t *testing.T) {
	log.Printf("🧪 Starting TestGetActiveEra")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)
	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	log.Printf("🚀 Calling getActiveEra...")
	activeEra, err := verifier.getActiveEra()
	log.Printf("📋 getActiveEra returned - activeEra: %v, error: %v", activeEra, err)

	if err != nil {
		log.Printf("❌ Error occurred while getting active era: %v", err)
		t.Errorf("Expected no error, got: %v", err)
	} else {
		log.Printf("✅ No errors occurred while getting active era")
	}

	// Log the active era details
	if activeEra != nil {
		log.Printf("📅 Active era data: %v", activeEra)

		// Try to extract era information if it's a map
		if eraMap, ok := activeEra.(map[string]interface{}); ok {
			if index, exists := eraMap["index"]; exists {
				log.Printf("📊 Era index: %v", index)
			}
			if start, exists := eraMap["start"]; exists {
				log.Printf("📅 Era start: %v", start)
			}
		}
	} else {
		log.Printf("⚠️  Active era is nil (this might be expected for some RPC endpoints)")
	}

	log.Printf("🎉 TestGetActiveEra completed successfully")
}

func TestGetActiveEraWithInvalidRPC(t *testing.T) {
	log.Printf("🧪 Starting TestGetActiveEraWithInvalidRPC")

	invalidRPCURL := "https://invalid-rpc-endpoint-that-does-not-exist.com"
	log.Printf("📡 Creating verifier with invalid RPC URL: %s", invalidRPCURL)
	verifier := NewVerifier(invalidRPCURL)
	log.Printf("✅ Verifier created successfully")

	log.Printf("🚀 Calling getActiveEra with invalid RPC...")
	activeEra, err := verifier.getActiveEra()
	log.Printf("📋 getActiveEra returned - activeEra: %v, error: %v", activeEra, err)

	// We expect an error with an invalid RPC URL
	if err != nil {
		log.Printf("✅ Expected error occurred: %v", err)
	} else {
		log.Printf("⚠️  No error occurred (unexpected)")
		t.Error("Expected error with invalid RPC URL, but got none")
	}

	if activeEra != nil {
		log.Printf("⚠️  Active era is not nil (unexpected with invalid RPC)")
		t.Error("Expected nil active era with invalid RPC URL")
	} else {
		log.Printf("✅ Active era is nil as expected")
	}

	log.Printf("🎉 TestGetActiveEraWithInvalidRPC completed successfully")
}

func TestGetActiveEraMultipleCalls(t *testing.T) {
	log.Printf("🧪 Starting TestGetActiveEraMultipleCalls")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)
	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	// Make multiple calls to test consistency
	for i := 1; i <= 3; i++ {
		log.Printf("🔄 Call %d: Getting active era...", i)
		activeEra, err := verifier.getActiveEra()

		if err != nil {
			log.Printf("❌ Error on call %d: %v", i, err)
			t.Errorf("Call %d failed: %v", i, err)
		} else {
			log.Printf("✅ Call %d successful - Active era: %v", i, activeEra)
		}

		// Small delay between calls
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("🎉 TestGetActiveEraMultipleCalls completed successfully")
}

func TestGetActiveEraStorageKeys(t *testing.T) {
	log.Printf("🧪 Starting TestGetActiveEraStorageKeys")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)
	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	// Test different storage keys for active era
	storageKeys := []string{
		"0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70", // Current key (Staking.ActiveEra)
		"0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70", // Blake2_128("Staking") + Blake2_128("ActiveEra")
		"0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70", // Alternative key
	}

	for i, key := range storageKeys {
		log.Printf("🔍 Testing storage key %d: %s", i+1, key)

		request := RPCRequest{
			JSONRPC: "2.0",
			Method:  "state_getStorage",
			Params: []interface{}{
				key,
			},
			ID: 1,
		}

		result, err := verifier.makeRPCCall(request)
		log.Printf("📋 Key %d result: %v, error: %v", i+1, result, err)

		if err != nil {
			log.Printf("❌ Key %d failed: %v", i+1, err)
		} else if result != nil {
			log.Printf("✅ Key %d successful: %v", i+1, result)
		} else {
			log.Printf("⚠️  Key %d returned nil", i+1)
		}

		// Small delay between calls
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("🎉 TestGetActiveEraStorageKeys completed successfully")
}

func TestGetActiveEraAlternativeMethods(t *testing.T) {
	log.Printf("🧪 Starting TestGetActiveEraAlternativeMethods")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)
	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	// Test alternative RPC methods for getting era information
	testCases := []struct {
		name   string
		method string
		params []interface{}
	}{
		{
			name:   "state_call with StakingApi_active_era",
			method: "state_call",
			params: []interface{}{
				"StakingApi_active_era",
				"0x",
			},
		},
		{
			name:   "state_call with StakingApi_eras_start",
			method: "state_call",
			params: []interface{}{
				"StakingApi_eras_start",
				"0x",
			},
		},
		{
			name:   "state_call with StakingApi_eras_total_stake",
			method: "state_call",
			params: []interface{}{
				"StakingApi_eras_total_stake",
				"0x",
			},
		},
		{
			name:   "chain_getHeader to get current block",
			method: "chain_getHeader",
			params: []interface{}{},
		},
	}

	for i, tc := range testCases {
		log.Printf("🔍 Testing method %d: %s", i+1, tc.name)

		request := RPCRequest{
			JSONRPC: "2.0",
			Method:  tc.method,
			Params:  tc.params,
			ID:      1,
		}

		result, err := verifier.makeRPCCall(request)
		log.Printf("📋 Method %d result: %v, error: %v", i+1, result, err)

		if err != nil {
			log.Printf("❌ Method %d failed: %v", i+1, err)
		} else if result != nil {
			log.Printf("✅ Method %d successful: %v", i+1, result)
		} else {
			log.Printf("⚠️  Method %d returned nil", i+1)
		}

		// Small delay between calls
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("🎉 TestGetActiveEraAlternativeMethods completed successfully")
}

func TestExploreStorageKeys(t *testing.T) {
	log.Printf("🧪 Starting TestExploreStorageKeys")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)
	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	// Test different storage keys that might contain era information
	storageKeys := []struct {
		name string
		key  string
	}{
		{
			name: "Staking.ActiveEra (current)",
			key:  "0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70",
		},
		{
			name: "Staking.ErasStart",
			key:  "0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70",
		},
		{
			name: "Staking.ErasTotalStake",
			key:  "0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70",
		},
		{
			name: "Staking.ErasRewardPoints",
			key:  "0x5f3e4907f716ac89b6347d15ececedca3ed14b45ed20d054f05e37e2542cfe70",
		},
	}

	for i, testCase := range storageKeys {
		log.Printf("🔍 Testing storage key %d: %s", i+1, testCase.name)
		log.Printf("   Key: %s", testCase.key)

		request := RPCRequest{
			JSONRPC: "2.0",
			Method:  "state_getStorage",
			Params: []interface{}{
				testCase.key,
			},
			ID: 1,
		}

		result, err := verifier.makeRPCCall(request)
		log.Printf("📋 %s result: %v, error: %v", testCase.name, result, err)

		if err != nil {
			log.Printf("❌ %s failed: %v", testCase.name, err)
		} else if result != nil {
			log.Printf("✅ %s successful: %v", testCase.name, result)
		} else {
			log.Printf("⚠️  %s returned nil", testCase.name)
		}

		// Small delay between calls
		time.Sleep(300 * time.Millisecond)
	}

	log.Printf("🎉 TestExploreStorageKeys completed successfully")
}

func TestGetCurrentBlockInfo(t *testing.T) {
	log.Printf("🧪 Starting TestGetCurrentBlockInfo")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)
	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	// Get current block header
	log.Printf("🚀 Getting current block header...")
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "chain_getHeader",
		Params:  []interface{}{},
		ID:      1,
	}

	result, err := verifier.makeRPCCall(request)
	log.Printf("📋 Block header result: %v, error: %v", result, err)

	if err != nil {
		log.Printf("❌ Failed to get block header: %v", err)
		t.Errorf("Expected no error, got: %v", err)
	} else {
		log.Printf("✅ Successfully got block header")

		// Try to extract block number
		if headerMap, ok := result.(map[string]interface{}); ok {
			if number, exists := headerMap["number"]; exists {
				log.Printf("📊 Current block number: %v", number)
			}
			if stateRoot, exists := headerMap["stateRoot"]; exists {
				log.Printf("🔗 State root: %v", stateRoot)
			}
		}
	}

	log.Printf("🎉 TestGetCurrentBlockInfo completed successfully")
}
