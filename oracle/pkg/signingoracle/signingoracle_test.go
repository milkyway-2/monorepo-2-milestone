package signingoracle

import (
	"encoding/hex"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewSigningOracle(t *testing.T) {
	log.Printf("🧪 Starting TestNewSigningOracle")

	// Set a test private key
	testPrivateKey := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	os.Setenv("PRIVATE_KEY", testPrivateKey)
	defer os.Unsetenv("PRIVATE_KEY")

	// Set a test RPC URL
	testRPCURL := "https://rpc.polkadot.io"
	os.Setenv("POLKADOT_RPC_URL", testRPCURL)
	defer os.Unsetenv("POLKADOT_RPC_URL")

	log.Printf("🔧 Creating SigningOracle with test private key")
	oracle, err := NewSigningOracle()
	if err != nil {
		log.Printf("❌ Failed to create SigningOracle: %v", err)
		t.Fatalf("Expected no error, got: %v", err)
	}

	if oracle == nil {
		log.Printf("❌ SigningOracle is nil")
		t.Fatal("Expected SigningOracle to not be nil")
	}

	log.Printf("✅ SigningOracle created successfully")

	// Test public key
	publicKey := oracle.GetPublicKeyHex()
	log.Printf("📋 Public Key: %s", publicKey)
	if publicKey == "" {
		log.Printf("❌ Public key is empty")
		t.Error("Expected public key to not be empty")
	} else {
		log.Printf("✅ Public key retrieved successfully")
	}

	// Test address
	address := oracle.GetAddress()
	log.Printf("📋 Address: %s", address)
	if address == "" {
		log.Printf("❌ Address is empty")
		t.Error("Expected address to not be empty")
	} else {
		log.Printf("✅ Address retrieved successfully")
	}

	// Test private key (should be different from input due to processing)
	privateKey := oracle.GetPrivateKeyHex()
	log.Printf("📋 Private Key: %s", privateKey)
	if privateKey == "" {
		log.Printf("❌ Private key is empty")
		t.Error("Expected private key to not be empty")
	} else {
		log.Printf("✅ Private key retrieved successfully")
	}

	// Test message signing
	testMessage := "Hello, World!"
	log.Printf("🔍 Testing message signing with: %s", testMessage)
	signature, err := oracle.SignMessage(testMessage)
	if err != nil {
		log.Printf("❌ Failed to sign message: %v", err)
		t.Errorf("Expected no error, got: %v", err)
	} else {
		log.Printf("✅ Message signed successfully")
		log.Printf("📋 Signature: %s", signature)
	}

	// Test verifier
	verifier := oracle.GetVerifier()
	if verifier == nil {
		log.Printf("❌ Verifier is nil")
		t.Error("Expected verifier to not be nil")
	} else {
		log.Printf("✅ Verifier retrieved successfully")
	}

	log.Printf("🎉 TestNewSigningOracle completed successfully")
}

func TestNewSigningOracle_MissingPrivateKey(t *testing.T) {
	log.Printf("🧪 Starting TestNewSigningOracle_MissingPrivateKey")

	// Ensure PRIVATE_KEY is not set
	os.Unsetenv("PRIVATE_KEY")

	log.Printf("🔧 Attempting to create SigningOracle without private key")
	oracle, err := NewSigningOracle()
	if err == nil {
		log.Printf("❌ Expected error but got none")
		t.Error("Expected error when PRIVATE_KEY is not set")
	} else {
		log.Printf("✅ Expected error occurred: %v", err)
	}

	if oracle != nil {
		log.Printf("❌ Oracle should be nil when error occurs")
		t.Error("Expected oracle to be nil when error occurs")
	} else {
		log.Printf("✅ Oracle correctly nil when error occurs")
	}

	log.Printf("🎉 TestNewSigningOracle_MissingPrivateKey completed successfully")
}

// TestSignEthereumMessage tests the SignEthereumMessage method
func TestSignEthereumMessage(t *testing.T) {
	log.Printf("🧪 Starting TestSignEthereumMessage")

	// Set up environment for testing
	os.Setenv("PRIVATE_KEY", "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	os.Setenv("POLKADOT_RPC_URL", "https://rpc.polkadot.io")
	defer os.Unsetenv("PRIVATE_KEY")
	defer os.Unsetenv("POLKADOT_RPC_URL")

	// Create signing oracle
	signingOracle, err := NewSigningOracle()
	if err != nil {
		t.Fatalf("Failed to create signing oracle: %v", err)
	}

	// Test message
	testMessage := "Hello, Ethereum World!"
	log.Printf("🔍 Testing Ethereum message signing with: %s", testMessage)

	// Sign the message with Ethereum format
	signature, err := signingOracle.SignEthereumMessage(testMessage)
	if err != nil {
		t.Fatalf("Failed to sign Ethereum message: %v", err)
	}

	log.Printf("✅ Ethereum message signed successfully")
	log.Printf("📋 Signature: %s", signature)

	// Verify the signature can be recovered correctly
	// Create the message hash
	msgHash := crypto.Keccak256Hash([]byte(testMessage))

	// Create Ethereum signed message hash
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	data := append(prefix, msgHash.Bytes()...)
	ethSignedMessageHash := crypto.Keccak256(data)

	// Decode the signature
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		t.Fatalf("Failed to decode signature: %v", err)
	}

	// Recover the signer
	recoveredPubKey, err := crypto.Ecrecover(ethSignedMessageHash, signatureBytes)
	if err != nil {
		t.Fatalf("Failed to recover public key: %v", err)
	}

	// Convert to address
	pubKey, err := crypto.UnmarshalPubkey(recoveredPubKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal public key: %v", err)
	}
	recoveredAddress := crypto.PubkeyToAddress(*pubKey)

	expectedAddress := signingOracle.GetAddress()
	if recoveredAddress.Hex() != expectedAddress {
		t.Fatalf("Recovered address doesn't match: expected %s, got %s", expectedAddress, recoveredAddress.Hex())
	}

	log.Printf("✅ Signature verification successful")
	log.Printf("📋 Recovered Address: %s", recoveredAddress.Hex())
	log.Printf("📋 Expected Address: %s", expectedAddress)

	log.Printf("🎉 TestSignEthereumMessage completed successfully")
}
