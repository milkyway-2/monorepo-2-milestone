package delegation

import (
	"log"
	"testing"
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
