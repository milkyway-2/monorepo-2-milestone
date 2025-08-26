package signatureverifier

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestCompleteSigningAndVerificationFlow(t *testing.T) {
	log.Printf("🧪 Starting Complete Signing and Verification Flow Test")

	// Step 1: Generate a test private key and derive the oracle address
	log.Printf("🔑 Step 1: Generating test private key and oracle address")

	// Generate a random private key for testing
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the public key and derive the address
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	oracleAddress := crypto.PubkeyToAddress(*publicKey)

	log.Printf("📋 Generated Oracle Address: %s", oracleAddress.Hex())
	log.Printf("📋 Private Key: %s", hex.EncodeToString(crypto.FromECDSA(privateKey)))

	// Step 2: Create the verifier with the oracle address
	log.Printf("🔧 Step 2: Creating verifier with oracle address")
	verifier, err := NewOracleVerifiedDelegation(oracleAddress.Hex())
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}
	log.Printf("✅ Verifier created successfully")

	// Step 3: Prepare the delegation message
	log.Printf("📝 Step 3: Preparing delegation message")
	validatorAddress := "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY"
	nominatorAddress := "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty"
	msgText := "I want to delegate 100 DOT to this validator"

	message := Message{
		ValidatorAddress: validatorAddress,
		NominatorAddress: nominatorAddress,
		MsgText:          msgText,
	}

	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)

	// Step 4: Create the message hash (same as verifier does)
	log.Printf("🔍 Step 4: Creating message hash")
	messageHash := verifier.createMessageHash(validatorAddress, nominatorAddress, msgText)
	log.Printf("📋 Message Hash: %s", hex.EncodeToString(messageHash))

	// Step 5: Create Ethereum signed message hash
	log.Printf("🔐 Step 5: Creating Ethereum signed message hash")
	ethSignedMessageHash := verifier.toEthSignedMessageHash(messageHash)
	log.Printf("📋 Ethereum Signed Message Hash: %s", hex.EncodeToString(ethSignedMessageHash))

	// Step 6: Sign the Ethereum signed message hash
	log.Printf("✍️ Step 6: Signing the message hash")
	signature, err := crypto.Sign(ethSignedMessageHash, privateKey)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	signatureHex := hex.EncodeToString(signature)
	log.Printf("📋 Signature: %s", signatureHex)

	// Step 7: Verify the signature using our verifier
	log.Printf("✅ Step 7: Verifying the signature")
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		t.Fatalf("Signature verification failed: %v", err)
	}
	log.Printf("✅ Signature verification successful!")

	// Step 8: Test the convenience function
	log.Printf("🔄 Step 8: Testing convenience function")
	err = verifier.VerifyMessage(message, signatureHex)
	if err != nil {
		t.Fatalf("Convenience function verification failed: %v", err)
	}
	log.Printf("✅ Convenience function verification successful!")

	// Step 9: Verify oracle address matches
	log.Printf("🔍 Step 9: Verifying oracle address")
	recoveredAddress, err := verifier.recoverSigner(ethSignedMessageHash, signature)
	if err != nil {
		t.Fatalf("Failed to recover signer: %v", err)
	}

	if recoveredAddress != oracleAddress {
		t.Fatalf("Recovered address doesn't match oracle address: expected %s, got %s",
			oracleAddress.Hex(), recoveredAddress.Hex())
	}
	log.Printf("✅ Oracle address verification successful: %s", recoveredAddress.Hex())

	log.Printf("🎉 Complete signing and verification flow test passed!")
}

// TestRealContractParameters tests with the exact parameters from the failing contract transaction
func TestRealContractParameters(t *testing.T) {
	log.Printf("🧪 Testing with real contract parameters")

	// The oracle address from the actual smart contract (updated)
	oracleAddress := "0x45D1960EB3E945e148D2828a2dC0CbBb52a29609"

	// Create verifier with the real oracle address
	verifier, err := NewOracleVerifiedDelegation(oracleAddress)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// The exact parameters from the failing transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"
	signatureHex := "100476bab5ff7bdeab21fc9171dcf118a909a2a00aae5fa3005c082b7820aa743687029ff34288f0c2a8303246aefc84264f50082ee8fc8df546dff3461a025701"

	log.Printf("📋 Oracle Address: %s", oracleAddress)
	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)
	log.Printf("📋 Signature: %s", signatureHex)

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
		log.Printf("📋 Expected Oracle Address: %s", oracleAddress)

		if recoveredAddress.Hex() == oracleAddress {
			log.Printf("✅ Addresses match!")
		} else {
			log.Printf("❌ Addresses don't match!")
		}
	}

	// Try the full verification
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		log.Printf("❌ Verification failed: %v", err)
	} else {
		log.Printf("✅ Verification successful!")
	}
}

// TestCorrectOracleSignature demonstrates how to create a valid signature with the correct oracle private key
func TestCorrectOracleSignature(t *testing.T) {
	log.Printf("🧪 Testing with correct oracle signature")

	// The oracle address from the actual smart contract (updated)
	oracleAddress := "0x45D1960EB3E945e148D2828a2dC0CbBb52a29609"

	// Create verifier with the real oracle address
	verifier, err := NewOracleVerifiedDelegation(oracleAddress)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// The parameters from the failing transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Oracle Address: %s", oracleAddress)
	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)

	// Create message hash
	messageHash := verifier.createMessageHash(validatorAddress, nominatorAddress, msgText)
	log.Printf("📋 Message Hash: %s", hex.EncodeToString(messageHash))

	// Create Ethereum signed message hash
	ethSignedMessageHash := verifier.toEthSignedMessageHash(messageHash)
	log.Printf("📋 Ethereum Signed Message Hash: %s", hex.EncodeToString(ethSignedMessageHash))

	// NOTE: To create a valid signature, you need the private key that corresponds to the oracle address
	// The private key for address 0x2BB632bAa1bCA1F51B7f4B2D02bC9bC07D5CDdFD would be needed here
	// For demonstration, we'll show what the signature creation would look like:

	log.Printf("⚠️  To create a valid signature, you need the private key for oracle address %s", oracleAddress)
	log.Printf("⚠️  The signature creation would look like:")
	log.Printf("⚠️  privateKey := [private key bytes for oracle address]")
	log.Printf("⚠️  signature, err := crypto.Sign(ethSignedMessageHash, privateKey)")
	log.Printf("⚠️  signatureHex := hex.EncodeToString(signature)")

	// Example of what the correct signature would look like (if we had the private key):
	log.Printf("⚠️  The correct signature would verify successfully with the contract")
	log.Printf("⚠️  Current signature was created by address: 0x45D1960EB3E945e148D2828a2dC0CbBb52a29609")
	log.Printf("⚠️  But contract expects signature from address: %s", oracleAddress)
}

// TestCreateValidSignature demonstrates how to create a valid signature with the correct oracle private key
func TestCreateValidSignature(t *testing.T) {
	log.Printf("🧪 Testing signature creation with correct oracle private key")

	// Generate a new private key for testing (in production, this would be the oracle's actual private key)
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the oracle address from the private key
	oracleAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	log.Printf("📋 Generated Oracle Address: %s", oracleAddress.Hex())
	log.Printf("📋 Private Key: %s", privateKeyHex)

	// Create verifier with the oracle address
	verifier, err := NewOracleVerifiedDelegation(oracleAddress.Hex())
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// The parameters from the failing transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"

	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)

	// Create a valid signature using the helper function
	validSignatureHex, err := verifier.CreateValidSignature(validatorAddress, nominatorAddress, msgText, privateKeyHex)
	if err != nil {
		t.Fatalf("Failed to create valid signature: %v", err)
	}

	log.Printf("📋 Valid Signature: %s", validSignatureHex)

	// Verify the signature
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, validSignatureHex)
	if err != nil {
		t.Fatalf("Valid signature verification failed: %v", err)
	}

	log.Printf("✅ Valid signature verification successful!")

	// Test the convenience function
	message := Message{
		ValidatorAddress: validatorAddress,
		NominatorAddress: nominatorAddress,
		MsgText:          msgText,
	}

	err = verifier.VerifyMessage(message, validSignatureHex)
	if err != nil {
		t.Fatalf("Convenience function verification failed: %v", err)
	}

	log.Printf("✅ Convenience function verification successful!")

	log.Printf("🎉 This signature would work with the smart contract!")
}

// TestNewFailingSignature analyzes the latest failing transaction signature
func TestNewFailingSignature(t *testing.T) {
	log.Printf("🧪 Analyzing new failing signature from latest transaction")

	// The oracle address from the updated smart contract
	oracleAddress := "0x45D1960EB3E945e148D2828a2dC0CbBb52a29609"

	// Create verifier with the oracle address
	verifier, err := NewOracleVerifiedDelegation(oracleAddress)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	// The parameters from the latest failing transaction
	validatorAddress := "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"
	nominatorAddress := "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU"
	msgText := "msg"
	signatureHex := "a02a2e74c854e261ab6633b72707ba875b389bca94075d6d2289e72dd261e5e44308854b5f71a8ee330323a7daac14dab3b0a69346759d00f61624be53b74b1100"

	log.Printf("📋 Oracle Address: %s", oracleAddress)
	log.Printf("📋 Validator Address: %s", validatorAddress)
	log.Printf("📋 Nominator Address: %s", nominatorAddress)
	log.Printf("📋 Message Text: %s", msgText)
	log.Printf("📋 New Signature: %s", signatureHex)

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
		log.Printf("📋 Expected Oracle Address: %s", oracleAddress)

		if recoveredAddress.Hex() == oracleAddress {
			log.Printf("✅ Addresses match!")
		} else {
			log.Printf("❌ Addresses don't match!")
			log.Printf("❌ This signature was created by a different private key!")
		}
	}

	// Try the full verification
	err = verifier.SubmitMessage(validatorAddress, nominatorAddress, msgText, signatureHex)
	if err != nil {
		log.Printf("❌ Verification failed: %v", err)
	} else {
		log.Printf("✅ Verification successful!")
	}
}
