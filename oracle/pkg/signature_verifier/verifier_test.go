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
