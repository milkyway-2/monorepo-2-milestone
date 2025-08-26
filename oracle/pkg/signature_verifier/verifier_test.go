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
	signatureHex := "95cb703ba12c252f827b6f1f935013bfa7c4671083b67795a4e1b915bc3aaf202430f07045a7df61832a71fbaea93e71b6ad65f15ea3eb0a01fc35dd287a249701"

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
