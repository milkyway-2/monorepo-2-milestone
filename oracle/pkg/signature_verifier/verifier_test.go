package signatureverifier

import (
	"encoding/hex"
	"log"
	"os"
	"testing"

	"oracle/pkg/signingoracle"
)

// TestSignAndVerifySuccess demonstrates a complete successful flow of signing and verifying
func TestSignAndVerifySuccess(t *testing.T) {
	log.Printf("🧪 Starting Sign and Verify Success Test")

	// Step 1: Set up the signing oracle with a test private key
	privateKeyHex := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	os.Setenv("PRIVATE_KEY", privateKeyHex)
	defer os.Unsetenv("PRIVATE_KEY")

	os.Setenv("POLKADOT_RPC_URL", "https://rpc.polkadot.io")
	defer os.Unsetenv("POLKADOT_RPC_URL")

	signingOracle, err := signingoracle.NewSigningOracle()
	if err != nil {
		t.Fatalf("Failed to create signing oracle: %v", err)
	}

	oracleAddress := signingOracle.GetAddress()
	log.Printf("📋 Oracle Address: %s", oracleAddress)

	// Step 2: Create the verifier
	verifier, err := NewOracleVerifiedDelegation(oracleAddress)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// Step 3: Define the message to sign
	validatorAddress := "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY"
	nominatorAddress := "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty"
	msgText := "I want to delegate 100 DOT to this validator"

	log.Printf("📋 Validator: %s", validatorAddress)
	log.Printf("📋 Nominator: %s", nominatorAddress)
	log.Printf("📋 Message: %s", msgText)

	// Step 4: Sign the message using the signing oracle
	fullMessage := validatorAddress + nominatorAddress + msgText
	signatureHex, err := signingOracle.SignEthereumMessage(fullMessage)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	log.Printf("📋 Signature: %s", signatureHex)

	// Step 5: Verify the signature using the verifier
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		t.Fatalf("Signature verification failed: %v", err)
	}

	log.Printf("✅ SUCCESS: Message signed and verified successfully!")
	log.Printf("🎉 The complete flow works correctly!")
}

// TestEnvironmentVariableHandling tests the signing oracle with environment variables
func TestEnvironmentVariableHandling(t *testing.T) {
	log.Printf("🧪 Testing Environment Variable Handling")

	// Test with missing PRIVATE_KEY
	os.Unsetenv("PRIVATE_KEY")
	_, err := signingoracle.NewSigningOracle()
	if err == nil {
		t.Fatalf("Expected error when PRIVATE_KEY is missing")
	}
	log.Printf("✅ Correctly handled missing PRIVATE_KEY: %v", err)

	// Test with invalid PRIVATE_KEY
	os.Setenv("PRIVATE_KEY", "invalid_hex")
	_, err = signingoracle.NewSigningOracle()
	if err == nil {
		t.Fatalf("Expected error when PRIVATE_KEY is invalid")
	}
	log.Printf("✅ Correctly handled invalid PRIVATE_KEY: %v", err)

	// Test with valid PRIVATE_KEY
	os.Setenv("PRIVATE_KEY", "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	os.Setenv("POLKADOT_RPC_URL", "https://rpc.polkadot.io")

	signingOracle, err := signingoracle.NewSigningOracle()
	if err != nil {
		t.Fatalf("Failed to create signing oracle with valid environment: %v", err)
	}

	log.Printf("✅ Successfully created signing oracle with valid environment")
	log.Printf("📋 Oracle Address: %s", signingOracle.GetAddress())
	log.Printf("📋 Public Key: %s", signingOracle.GetPublicKeyHex())

	// Clean up
	os.Unsetenv("PRIVATE_KEY")
	os.Unsetenv("POLKADOT_RPC_URL")
}

// TestSigningOracleIntegration tests the integration between signing oracle and verifier
func TestSigningOracleIntegration(t *testing.T) {
	log.Printf("🧪 Testing Signing Oracle Integration")

	// Set up environment for testing
	privateKeyHex := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	os.Setenv("PRIVATE_KEY", privateKeyHex)
	defer os.Unsetenv("PRIVATE_KEY")

	os.Setenv("POLKADOT_RPC_URL", "https://rpc.polkadot.io")
	defer os.Unsetenv("POLKADOT_RPC_URL")

	// Create signing oracle
	signingOracle, err := signingoracle.NewSigningOracle()
	if err != nil {
		t.Fatalf("Failed to create signing oracle: %v", err)
	}

	// Create verifier
	verifier, err := NewOracleVerifiedDelegation(signingOracle.GetAddress())
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// Test data
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "Test delegation message"

	log.Printf("📋 Oracle Address: %s", signingOracle.GetAddress())
	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)

	// Create message hash
	messageHash := verifier.createMessageHash(validatorAddress, nominatorAddress, msgText)
	log.Printf("📋 Message Hash: %s", hex.EncodeToString(messageHash))

	// Create Ethereum signed message hash
	ethSignedMessageHash := verifier.toEthSignedMessageHash(messageHash)
	log.Printf("📋 Ethereum Signed Message Hash: %s", hex.EncodeToString(ethSignedMessageHash))

	// Sign using signing oracle
	fullMessage := validatorAddress + nominatorAddress + msgText
	signatureHex, err := signingOracle.SignEthereumMessage(fullMessage)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	log.Printf("📋 Signature: %s", signatureHex)

	// Verify using verifier
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}

	log.Printf("✅ Signing oracle integration test passed!")
}
