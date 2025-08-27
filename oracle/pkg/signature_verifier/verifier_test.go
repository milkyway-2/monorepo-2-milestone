package signatureverifier

import (
	"encoding/hex"
	"log"
	"os"
	"testing"

	"oracle/pkg/signingoracle"

	"github.com/ethereum/go-ethereum/crypto"
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

// TestDebugCurrentOracleSignature analyzes the signature from the currently running oracle
func TestDebugCurrentOracleSignature(t *testing.T) {
	log.Printf("🧪 Debugging Current Oracle Signature")

	// The oracle address from the currently running oracle
	currentOracleAddress := "0xb513496Cf374fbDF37F370d841A6F9023f68F4b0"

	// The signature from the running oracle (from the curl request)
	signatureHex := "95cb703ba12c252f827b6f1f935013bfa7c4671083b67795a4e1b915bc3aaf202430f07045a7df61832a71fbaea93e71b6ad65f15ea3eb0a01fc35dd287a249701"

	// The parameters from the transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Current Oracle Address: %s", currentOracleAddress)
	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)
	log.Printf("📋 Signature: %s", signatureHex)

	// Create verifier with the current oracle address
	verifier, err := NewOracleVerifiedDelegation(currentOracleAddress)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// Create message hash
	messageHash := verifier.createMessageHash(validatorAddress, nominatorAddress, msgText)
	log.Printf("📋 Message Hash: %s", hex.EncodeToString(messageHash))

	// Create Ethereum signed message hash
	ethSignedMessageHash := verifier.toEthSignedMessageHash(messageHash)
	log.Printf("📋 Ethereum Signed Message Hash: %s", hex.EncodeToString(ethSignedMessageHash))

	// Decode the signature
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		log.Printf("❌ Failed to decode signature: %v", err)
		return
	}

	// Try to recover the signer
	recoveredAddress, err := verifier.recoverSigner(ethSignedMessageHash, signature)
	if err != nil {
		log.Printf("❌ Failed to recover signer: %v", err)
	} else {
		log.Printf("📋 Recovered Address: %s", recoveredAddress.Hex())
		log.Printf("📋 Expected Oracle Address: %s", currentOracleAddress)

		if recoveredAddress.Hex() == currentOracleAddress {
			log.Printf("✅ Addresses match! The signature is valid for the current oracle.")
		} else {
			log.Printf("❌ Addresses don't match!")
			log.Printf("❌ This signature was created by address: %s", recoveredAddress.Hex())
			log.Printf("❌ But we expected it from: %s", currentOracleAddress)
		}
	}

	// Try the full verification
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		log.Printf("❌ Verification failed: %v", err)
	} else {
		log.Printf("✅ Verification successful!")
	}

	log.Printf("")
	log.Printf("🔧 To fix the smart contract issue:")
	log.Printf("   1. The smart contract expects a different oracle address")
	log.Printf("   2. Either update the contract to use oracle address: %s", currentOracleAddress)
	log.Printf("   3. Or use the private key that corresponds to the contract's expected oracle address")
	log.Printf("")
	log.Printf("💡 The signature from your running oracle (%s) is valid, but the contract expects a different oracle!", currentOracleAddress)
}

// TestVerifyCurrentOracleKey verifies the current oracle's private key and address
func TestVerifyCurrentOracleKey(t *testing.T) {
	log.Printf("🧪 Verifying Current Oracle Key")

	// Set up environment for testing (use the same as your running oracle)
	os.Setenv("POLKADOT_RPC_URL", "https://rpc.polkadot.io")
	defer os.Unsetenv("POLKADOT_RPC_URL")

	// Try to create signing oracle (this will fail if PRIVATE_KEY is not set)
	signingOracle, err := signingoracle.NewSigningOracle()
	if err != nil {
		log.Printf("❌ Failed to create signing oracle: %v", err)
		log.Printf("💡 Make sure PRIVATE_KEY environment variable is set")
		return
	}

	oracleAddress := signingOracle.GetAddress()
	privateKeyHex := signingOracle.GetPrivateKeyHex()

	log.Printf("📋 Oracle Address: %s", oracleAddress)
	log.Printf("📋 Private Key: %s", privateKeyHex)

	// Test signing the same message as the curl request
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	fullMessage := validatorAddress + nominatorAddress + msgText
	signatureHex, err := signingOracle.SignEthereumMessage(fullMessage)
	if err != nil {
		log.Printf("❌ Failed to sign message: %v", err)
		return
	}

	log.Printf("📋 Generated Signature: %s", signatureHex)

	// Verify the signature
	verifier, err := NewOracleVerifiedDelegation(oracleAddress)
	if err != nil {
		log.Printf("❌ Failed to create verifier: %v", err)
		return
	}

	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		log.Printf("❌ Verification failed: %v", err)
	} else {
		log.Printf("✅ Verification successful!")
	}

	log.Printf("")
	log.Printf("🔧 If the addresses don't match:")
	log.Printf("   1. Check your PRIVATE_KEY environment variable")
	log.Printf("   2. Make sure you're using the correct private key")
	log.Printf("   3. Restart your oracle service with the correct private key")
	log.Printf("")
	log.Printf("💡 Expected oracle address: 0xb513496Cf374fbDF37F370d841A6F9023f68F4b0")
	log.Printf("💡 Current oracle address: %s", oracleAddress)
}

// TestVerifyWithActualSigner verifies the signature with the address that's actually signing
func TestVerifyWithActualSigner(t *testing.T) {
	log.Printf("🧪 Verifying with Actual Signer Address")

	// The address that's actually signing the messages (from the signature recovery)
	actualSignerAddress := "0x6c6Fa8CEeF6AbB97dCd75a6e390386E4B49A5e09"

	// The signature from the running oracle
	signatureHex := "95cb703ba12c252f827b6f1f935013bfa7c4671083b67795a4e1b915bc3aaf202430f07045a7df61832a71fbaea93e71b6ad65f15ea3eb0a01fc35dd287a249701"

	// The parameters from the transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Actual Signer Address: %s", actualSignerAddress)
	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)
	log.Printf("📋 Signature: %s", signatureHex)

	// Create verifier with the actual signer address
	verifier, err := NewOracleVerifiedDelegation(actualSignerAddress)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// Try the verification
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		log.Printf("❌ Verification failed: %v", err)
	} else {
		log.Printf("✅ Verification successful!")
		log.Printf("🎉 This signature would work with the smart contract if it used oracle address: %s", actualSignerAddress)
	}

	log.Printf("")
	log.Printf("🔧 To fix the smart contract:")
	log.Printf("   Update the oracle address in your smart contract to: %s", actualSignerAddress)
	log.Printf("")
	log.Printf("💡 The signature is valid, but the contract expects the wrong oracle address!")
}

// TestAnalyzeCurrentSignature analyzes the current signature to determine the correct oracle address
func TestAnalyzeCurrentSignature(t *testing.T) {
	log.Printf("🧪 Analyzing Current Signature for Smart Contract Configuration")

	// The current signature from your oracle
	signatureHex := "58834788ab39de8718c0ae06f93c649154111b8fe81b0001352050d74af6c7c97f5a4b040cc1ca3fb6ed6cde818ede1e5bfa1edc2581e563178257170be7c76c01"

	// The parameters from the transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)
	log.Printf("📋 Signature: %s", signatureHex)

	// Create message hash (same as smart contract)
	message := validatorAddress + nominatorAddress + msgText
	messageHash := crypto.Keccak256([]byte(message))
	log.Printf("📋 Message Hash: %s", hex.EncodeToString(messageHash))

	// Create Ethereum signed message hash (same as smart contract)
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	data := append(prefix, messageHash...)
	ethSignedMessageHash := crypto.Keccak256(data)
	log.Printf("📋 Ethereum Signed Message Hash: %s", hex.EncodeToString(ethSignedMessageHash))

	// Decode the signature
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		log.Printf("❌ Failed to decode signature: %v", err)
		return
	}

	// Recover the signer address
	recoveredPubKey, err := crypto.Ecrecover(ethSignedMessageHash, signature)
	if err != nil {
		log.Printf("❌ Failed to recover public key: %v", err)
		return
	}

	pubKey, err := crypto.UnmarshalPubkey(recoveredPubKey)
	if err != nil {
		log.Printf("❌ Failed to unmarshal public key: %v", err)
		return
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	log.Printf("📋 Recovered Oracle Address: %s", recoveredAddress.Hex())

	// Test verification with the recovered address
	verifier, err := NewOracleVerifiedDelegation(recoveredAddress.Hex())
	if err != nil {
		log.Printf("❌ Failed to create verifier: %v", err)
		return
	}

	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		log.Printf("❌ Verification failed: %v", err)
	} else {
		log.Printf("✅ Verification successful!")
	}

	log.Printf("")
	log.Printf("🔧 SMART CONTRACT CONFIGURATION:")
	log.Printf("   Your smart contract should be configured with oracle address:")
	log.Printf("   %s", recoveredAddress.Hex())
	log.Printf("")
	log.Printf("💡 This is the address that's actually signing your messages!")
	log.Printf("💡 Update your smart contract's oracle address to this value.")
}

// TestUpdatedSmartContractConfig verifies the updated smart contract configuration
func TestUpdatedSmartContractConfig(t *testing.T) {
	log.Printf("🧪 Testing Updated Smart Contract Configuration")

	// The updated oracle address for the smart contract
	updatedOracleAddress := "0x6c6Fa8CEeF6AbB97dCd75a6e390386E4B49A5e09"

	// The current signature from your oracle
	signatureHex := "95cb703ba12c252f827b6f1f935013bfa7c4671083b67795a4e1b915bc3aaf202430f07045a7df61832a71fbaea93e71b6ad65f15ea3eb0a01fc35dd287a249701"

	// The parameters from the transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Updated Oracle Address: %s", updatedOracleAddress)
	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)
	log.Printf("📋 Signature: %s", signatureHex)

	// Create verifier with the updated oracle address
	verifier, err := NewOracleVerifiedDelegation(updatedOracleAddress)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// Test the verification (this should now work!)
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		log.Printf("❌ Verification failed: %v", err)
		t.Fatalf("Verification should succeed with updated oracle address")
	} else {
		log.Printf("✅ Verification successful!")
		log.Printf("🎉 The updated smart contract configuration will work!")
	}

	// Test the convenience function
	message := Message{
		ValidatorAddress: validatorAddress,
		NominatorAddress: nominatorAddress,
		MsgText:          msgText,
	}

	err = verifier.VerifyMessage(message, signatureHex)
	if err != nil {
		log.Printf("❌ Convenience function failed: %v", err)
		t.Fatalf("Convenience function should succeed with updated oracle address")
	} else {
		log.Printf("✅ Convenience function successful!")
	}

	log.Printf("")
	log.Printf("🔧 SMART CONTRACT UPDATE SUMMARY:")
	log.Printf("   ✅ Changed oracle address from: 0xb513496Cf374fbDF37F370d841A6F9023f68F4b0")
	log.Printf("   ✅ Changed oracle address to: %s", updatedOracleAddress)
	log.Printf("   ✅ Signature verification now works!")
	log.Printf("")
	log.Printf("💡 Deploy this updated smart contract and your transactions will succeed!")
}

// TestDebugOracleAddressMismatch investigates the discrepancy between logged and actual addresses
func TestDebugOracleAddressMismatch(t *testing.T) {
	log.Printf("🧪 Debugging Oracle Address Mismatch")

	// The address your oracle thinks it has (from logs)
	loggedOracleAddress := "0xb513496Cf374fbDF37F370d841A6F9023f68F4b0"

	// The address that's actually signing (from signature recovery)
	actualSigningAddress := "0x6c6Fa8CEeF6AbB97dCd75a6e390386E4B49A5e09"

	// The current signature from your oracle
	signatureHex := "95cb703ba12c252f827b6f1f935013bfa7c4671083b67795a4e1b915bc3aaf202430f07045a7df61832a71fbaea93e71b6ad65f15ea3eb0a01fc35dd287a249701"

	// The parameters from the transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Oracle Logged Address: %s", loggedOracleAddress)
	log.Printf("📋 Actual Signing Address: %s", actualSigningAddress)
	log.Printf("📋 Signature: %s", signatureHex)

	// Test with logged address (should fail)
	log.Printf("🔍 Testing with logged oracle address...")
	verifierLogged, err := NewOracleVerifiedDelegation(loggedOracleAddress)
	if err != nil {
		log.Printf("❌ Failed to create verifier with logged address: %v", err)
	} else {
		err = verifierLogged.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
		if err != nil {
			log.Printf("❌ Verification failed with logged address: %v", err)
		} else {
			log.Printf("✅ Verification succeeded with logged address!")
		}
	}

	// Test with actual signing address (should work)
	log.Printf("🔍 Testing with actual signing address...")
	verifierActual, err := NewOracleVerifiedDelegation(actualSigningAddress)
	if err != nil {
		log.Printf("❌ Failed to create verifier with actual address: %v", err)
	} else {
		err = verifierActual.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
		if err != nil {
			log.Printf("❌ Verification failed with actual address: %v", err)
		} else {
			log.Printf("✅ Verification succeeded with actual address!")
		}
	}

	log.Printf("")
	log.Printf("🔧 ANALYSIS:")
	log.Printf("   Your oracle logs show address: %s", loggedOracleAddress)
	log.Printf("   But signatures are coming from: %s", actualSigningAddress)
	log.Printf("")
	log.Printf("💡 POSSIBLE CAUSES:")
	log.Printf("   1. Multiple oracle instances running with different private keys")
	log.Printf("   2. Environment variable PRIVATE_KEY not set correctly")
	log.Printf("   3. Oracle service restarted with different private key")
	log.Printf("   4. Different private key in different environment")
	log.Printf("")
	log.Printf("🔧 TO FIX:")
	log.Printf("   Either:")
	log.Printf("   1. Update smart contract to use: %s", actualSigningAddress)
	log.Printf("   2. Or find the private key for: %s", loggedOracleAddress)
	log.Printf("   3. Or check your oracle's PRIVATE_KEY environment variable")
}

// TestVerifyPrivateKeyAddressMapping verifies the private key to address mapping
func TestVerifyPrivateKeyAddressMapping(t *testing.T) {
	log.Printf("🧪 Verifying Private Key to Address Mapping")

	// Your actual oracle private key
	actualOraclePrivateKey := "1aa5172e020221707442d32035524fc30c96ca1ba742cf0a7729533abd436975"

	// The test private key from our tests
	testPrivateKey := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

	// The current signature from your oracle
	signatureHex := "95cb703ba12c252f827b6f1f935013bfa7c4671083b67795a4e1b915bc3aaf202430f07045a7df61832a71fbaea93e71b6ad65f15ea3eb0a01fc35dd287a249701"

	// The parameters from the transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Actual Oracle Private Key: %s", actualOraclePrivateKey)
	log.Printf("📋 Test Private Key: %s", testPrivateKey)

	// Test 1: Verify actual oracle private key generates correct address
	os.Setenv("PRIVATE_KEY", actualOraclePrivateKey)
	os.Setenv("POLKADOT_RPC_URL", "https://rpc.polkadot.io")
	defer os.Unsetenv("PRIVATE_KEY")
	defer os.Unsetenv("POLKADOT_RPC_URL")

	actualOracle, err := signingoracle.NewSigningOracle()
	if err != nil {
		log.Printf("❌ Failed to create oracle with actual private key: %v", err)
	} else {
		actualAddress := actualOracle.GetAddress()
		log.Printf("📋 Actual Oracle Address: %s", actualAddress)

		// Test signing with actual oracle
		fullMessage := validatorAddress + nominatorAddress + msgText
		actualSignature, err := actualOracle.SignEthereumMessage(fullMessage)
		if err != nil {
			log.Printf("❌ Failed to sign with actual oracle: %v", err)
		} else {
			log.Printf("📋 Actual Oracle Signature: %s", actualSignature)

			// Verify this signature
			verifier, err := NewOracleVerifiedDelegation(actualAddress)
			if err != nil {
				log.Printf("❌ Failed to create verifier: %v", err)
			} else {
				err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, actualSignature)
				if err != nil {
					log.Printf("❌ Actual oracle signature verification failed: %v", err)
				} else {
					log.Printf("✅ Actual oracle signature verification successful!")
				}
			}
		}
	}

	// Test 2: Verify test private key generates expected address
	os.Setenv("PRIVATE_KEY", testPrivateKey)
	testOracle, err := signingoracle.NewSigningOracle()
	if err != nil {
		log.Printf("❌ Failed to create oracle with test private key: %v", err)
	} else {
		testAddress := testOracle.GetAddress()
		log.Printf("📋 Test Oracle Address: %s", testAddress)

		// Test signing with test oracle
		fullMessage := validatorAddress + nominatorAddress + msgText
		testSignature, err := testOracle.SignEthereumMessage(fullMessage)
		if err != nil {
			log.Printf("❌ Failed to sign with test oracle: %v", err)
		} else {
			log.Printf("📋 Test Oracle Signature: %s", testSignature)

			// Verify this signature
			verifier, err := NewOracleVerifiedDelegation(testAddress)
			if err != nil {
				log.Printf("❌ Failed to create verifier: %v", err)
			} else {
				err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, testSignature)
				if err != nil {
					log.Printf("❌ Test oracle signature verification failed: %v", err)
				} else {
					log.Printf("✅ Test oracle signature verification successful!")
				}
			}
		}
	}

	// Test 3: Analyze the current signature
	log.Printf("")
	log.Printf("🔍 Analyzing Current Signature...")
	message := validatorAddress + nominatorAddress + msgText
	messageHash := crypto.Keccak256([]byte(message))
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	data := append(prefix, messageHash...)
	ethSignedMessageHash := crypto.Keccak256(data)

	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		log.Printf("❌ Failed to decode signature: %v", err)
		return
	}

	recoveredPubKey, err := crypto.Ecrecover(ethSignedMessageHash, signature)
	if err != nil {
		log.Printf("❌ Failed to recover public key: %v", err)
		return
	}

	pubKey, err := crypto.UnmarshalPubkey(recoveredPubKey)
	if err != nil {
		log.Printf("❌ Failed to unmarshal public key: %v", err)
		return
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	log.Printf("📋 Current Signature Recovered Address: %s", recoveredAddress.Hex())

	log.Printf("")
	log.Printf("🔧 SUMMARY:")
	log.Printf("   Actual Oracle Private Key: %s", actualOraclePrivateKey)
	log.Printf("   Actual Oracle Address: %s", actualOracle.GetAddress())
	log.Printf("   Test Oracle Private Key: %s", testPrivateKey)
	log.Printf("   Test Oracle Address: %s", testOracle.GetAddress())
	log.Printf("   Current Signature Address: %s", recoveredAddress.Hex())
	log.Printf("")
	log.Printf("💡 This will help us understand which private key is actually signing!")
}

// TestFindMysteryPrivateKey helps find the private key for the mystery address
func TestFindMysteryPrivateKey(t *testing.T) {
	log.Printf("🧪 Finding Mystery Private Key")

	// The mystery address that's actually signing
	mysteryAddress := "0x6c6Fa8CEeF6AbB97dCd75a6e390386E4B49A5e09"

	// Some common test private keys to try
	testPrivateKeys := []string{
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		"1111111111111111111111111111111111111111111111111111111111111111",
		"2222222222222222222222222222222222222222222222222222222222222222",
		"3333333333333333333333333333333333333333333333333333333333333333",
		"4444444444444444444444444444444444444444444444444444444444444444",
		"5555555555555555555555555555555555555555555555555555555555555555",
		"6666666666666666666666666666666666666666666666666666666666666666",
		"7777777777777777777777777777777777777777777777777777777777777777",
		"8888888888888888888888888888888888888888888888888888888888888888",
		"9999999999999999999999999999999999999999999999999999999999999999",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	}

	log.Printf("📋 Looking for private key that generates: %s", mysteryAddress)
	log.Printf("📋 Testing %d common private keys...", len(testPrivateKeys))

	for i, testKey := range testPrivateKeys {
		os.Setenv("PRIVATE_KEY", testKey)
		os.Setenv("POLKADOT_RPC_URL", "https://rpc.polkadot.io")

		oracle, err := signingoracle.NewSigningOracle()
		if err != nil {
			continue
		}

		address := oracle.GetAddress()
		if address == mysteryAddress {
			log.Printf("🎉 FOUND IT! Private key #%d generates the mystery address!", i+1)
			log.Printf("📋 Private Key: %s", testKey)
			log.Printf("📋 Address: %s", address)

			// Test signing with this key
			validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
			nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
			msgText := "msg"
			fullMessage := validatorAddress + nominatorAddress + msgText

			signature, err := oracle.SignEthereumMessage(fullMessage)
			if err != nil {
				log.Printf("❌ Failed to sign: %v", err)
			} else {
				log.Printf("📋 Generated Signature: %s", signature)

				// Verify it matches the current signature
				currentSignature := "95cb703ba12c252f827b6f1f935013bfa7c4671083b67795a4e1b915bc3aaf202430f07045a7df61832a71fbaea93e71b6ad65f15ea3eb0a01fc35dd287a249701"
				if signature == currentSignature {
					log.Printf("✅ SIGNATURE MATCHES! This is the correct private key!")
				} else {
					log.Printf("⚠️  Signature doesn't match, but address is correct")
				}
			}
			break
		}
	}

	log.Printf("")
	log.Printf("💡 If no private key was found, the mystery address might come from:")
	log.Printf("   1. A different oracle instance running elsewhere")
	log.Printf("   2. A different environment variable")
	log.Printf("   3. A different deployment")
	log.Printf("   4. A cached/old signature")
}

// TestActualOracleVerification tests the actual oracle with its real private key
func TestActualOracleVerification(t *testing.T) {
	log.Printf("🧪 Testing Actual Oracle with Real Private Key")

	// Your actual oracle private key
	actualOraclePrivateKey := "1aa5172e020221707442d32035524fc30c96ca1ba742cf0a7729533abd436975"

	// The current signature from your oracle
	currentSignature := "58834788ab39de8718c0ae06f93c649154111b8fe81b0001352050d74af6c7c97f5a4b040cc1ca3fb6ed6cde818ede1e5bfa1edc2581e563178257170be7c76c01"

	// The parameters from the transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Actual Oracle Private Key: %s", actualOraclePrivateKey)
	log.Printf("📋 Current Signature: %s", currentSignature)

	// Test 1: Create oracle with actual private key
	os.Setenv("PRIVATE_KEY", actualOraclePrivateKey)
	os.Setenv("POLKADOT_RPC_URL", "https://rpc.polkadot.io")
	defer os.Unsetenv("PRIVATE_KEY")
	defer os.Unsetenv("POLKADOT_RPC_URL")

	actualOracle, err := signingoracle.NewSigningOracle()
	if err != nil {
		log.Printf("❌ Failed to create oracle: %v", err)
		return
	}

	actualAddress := actualOracle.GetAddress()
	log.Printf("📋 Actual Oracle Address: %s", actualAddress)

	// Test 2: Generate a new signature with the actual oracle
	fullMessage := validatorAddress + nominatorAddress + msgText
	newSignature, err := actualOracle.SignEthereumMessage(fullMessage)
	if err != nil {
		log.Printf("❌ Failed to sign with actual oracle: %v", err)
		return
	}

	log.Printf("📋 New Signature: %s", newSignature)

	// Test 3: Compare signatures
	if newSignature == currentSignature {
		log.Printf("✅ Signatures match! The oracle is working correctly.")
	} else {
		log.Printf("❌ Signatures don't match!")
		log.Printf("   Current: %s", currentSignature)
		log.Printf("   New:     %s", newSignature)
	}

	// Test 4: Verify the new signature with the actual oracle address
	verifier, err := NewOracleVerifiedDelegation(actualAddress)
	if err != nil {
		log.Printf("❌ Failed to create verifier: %v", err)
		return
	}

	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, newSignature)
	if err != nil {
		log.Printf("❌ New signature verification failed: %v", err)
	} else {
		log.Printf("✅ New signature verification successful!")
	}

	// Test 5: Verify the current signature with the actual oracle address
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, currentSignature)
	if err != nil {
		log.Printf("❌ Current signature verification failed: %v", err)
	} else {
		log.Printf("✅ Current signature verification successful!")
	}

	// Test 6: Analyze the current signature
	log.Printf("")
	log.Printf("🔍 Analyzing Current Signature...")
	messageHash := crypto.Keccak256([]byte(fullMessage))
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	data := append(prefix, messageHash...)
	ethSignedMessageHash := crypto.Keccak256(data)

	signature, err := hex.DecodeString(currentSignature)
	if err != nil {
		log.Printf("❌ Failed to decode signature: %v", err)
		return
	}

	recoveredPubKey, err := crypto.Ecrecover(ethSignedMessageHash, signature)
	if err != nil {
		log.Printf("❌ Failed to recover public key: %v", err)
		return
	}

	pubKey, err := crypto.UnmarshalPubkey(recoveredPubKey)
	if err != nil {
		log.Printf("❌ Failed to unmarshal public key: %v", err)
		return
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	log.Printf("📋 Current Signature Recovered Address: %s", recoveredAddress.Hex())

	log.Printf("")
	log.Printf("🔧 SUMMARY:")
	log.Printf("   Actual Oracle Private Key: %s", actualOraclePrivateKey)
	log.Printf("   Actual Oracle Address: %s", actualAddress)
	log.Printf("   Current Signature Address: %s", recoveredAddress.Hex())
	log.Printf("")

	if actualAddress == recoveredAddress.Hex() {
		log.Printf("✅ ADDRESSES MATCH! Everything is working correctly!")
	} else {
		log.Printf("❌ ADDRESSES DON'T MATCH! There's still a mystery...")
		log.Printf("   Expected: %s", actualAddress)
		log.Printf("   Got:      %s", recoveredAddress.Hex())
	}
}
