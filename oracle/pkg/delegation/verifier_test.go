package delegation

import (
	"log"
	"testing"
)

func TestVerifyV2_RealPolkadotAddresses(t *testing.T) {
	log.Printf("🧪 Starting TestVerifyV2_RealPolkadotAddresses")

	rpcURL := "https://rpc.polkadot.io"
	log.Printf("📡 Creating verifier with RPC URL: %s", rpcURL)
	verifier := NewVerifier(rpcURL)
	log.Printf("✅ Verifier created successfully")

	// Real Polkadot addresses for testing
	// These are actual addresses from the Polkadot network
	nominator := "0x73479ae11533f4e717e3f7b45a8f54d95021785395df62abbe68ff9af32e40cc"
	validator := "12GTt3pfM3SjTU6UL6dQ3SMgMSvdw94PnRoF6osU6hPvxbUZ"

	log.Printf("🔍 Testing VerifyV2 with real Polkadot addresses:")
	log.Printf("   Nominator: %s", nominator)
	log.Printf("   Validator: %s", validator)

	log.Printf("🚀 Calling VerifyV2...")
	result, err := verifier.VerifyV2(nominator, validator)
	log.Printf("📋 VerifyV2 returned - result: %+v, error: %v", result, err)

	if err != nil {
		log.Printf("❌ Unexpected error occurred: %v", err)
		t.Errorf("Expected no error, got: %v", err)
	} else {
		log.Printf("✅ No errors occurred during verification")
	}

	if result == nil {
		log.Printf("❌ Result is nil")
		t.Fatal("Expected result to not be nil")
	}

	// Log detailed result information
	log.Printf("📊 Verification Result Details:")
	log.Printf("   Is Valid: %t", result.IsValid)
	log.Printf("   Address Validation: %t", result.AddressValidation)
	log.Printf("   Extrinsic Validation: %t", result.ExtrinsicValidation)
	log.Printf("   Storage Validation: %t", result.StorageValidation)
	log.Printf("   Active Era Validation: %t", result.ActiveEraValidation)
	log.Printf("   Error: %s", result.Error)
	log.Printf("   Additional Info: %s", result.AdditionalInfo)
	log.Printf("   Timestamp: %s", result.Timestamp.Format("2006-01-02 15:04:05"))

	// Verify that address validation passed
	if !result.AddressValidation {
		log.Printf("❌ Address validation failed")
		t.Error("Expected address validation to pass")
	} else {
		log.Printf("✅ Address validation passed")
	}

	// Verify that storage validation passed
	if !result.StorageValidation {
		log.Printf("❌ Storage validation failed")
		t.Error("Expected storage validation to pass")
	} else {
		log.Printf("✅ Storage validation passed")
	}

	// Verify that active era validation passed
	if !result.ActiveEraValidation {
		log.Printf("❌ Active era validation failed")
		t.Error("Expected active era validation to pass")
	} else {
		log.Printf("✅ Active era validation passed")
	}

	// Verify that extrinsic validation is false (as expected in V2)
	if result.ExtrinsicValidation {
		log.Printf("❌ Extrinsic validation should be false in V2")
		t.Error("Expected extrinsic validation to be false in V2")
	} else {
		log.Printf("✅ Extrinsic validation correctly false in V2")
	}

	// Verify overall validity
	if !result.IsValid {
		log.Printf("❌ Overall verification failed")
		t.Error("Expected overall verification to pass")
	} else {
		log.Printf("✅ Overall verification successful")
	}

	log.Printf("🎉 TestVerifyV2_RealPolkadotAddresses completed successfully")
}
