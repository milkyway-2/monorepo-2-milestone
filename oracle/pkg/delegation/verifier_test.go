package delegation

import (
	"testing"
)

func TestNewVerifier(t *testing.T) {
	rpcURL := "https://rpc.polkadot.io"
	verifier := NewVerifier(rpcURL)

	if verifier == nil {
		t.Fatal("NewVerifier returned nil")
	}

	if verifier.rpcURL != rpcURL {
		t.Errorf("Expected RPC URL %s, got %s", rpcURL, verifier.rpcURL)
	}

	if verifier.client == nil {
		t.Error("HTTP client should not be nil")
	}
}

func TestVerifyDelegation_Simulated(t *testing.T) {
	verifier := NewVerifier("https://rpc.polkadot.io")

	// Test with known addresses that should work in simulation mode
	nominator := "12ztGE9cY2p7kPJFpfvMrL6NsCUeqoiaBY3jciMqYFuFNJ2o"
	validator := "12GTt3pfM3SjTU6UL6dQ3SMgMSvdw94PnRoF6osU6hPvxbUZ"

	isDelegated, err := verifier.VerifyDelegation(nominator, validator)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !isDelegated {
		t.Error("Expected delegation to be verified for known addresses")
	}
}
