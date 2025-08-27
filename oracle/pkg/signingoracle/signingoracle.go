package signingoracle

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"oracle/pkg/delegation"

	"github.com/ethereum/go-ethereum/crypto"
)

// SigningOracle holds the private key for signing
type SigningOracle struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	verifier   *delegation.Verifier
}

// NewSigningOracle creates a new signing oracle with a private key from environment
func NewSigningOracle() (*SigningOracle, error) {
	// Get private key from environment variable
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable is required")
	}

	// Remove "0x" prefix if present
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	// Decode the private key
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	// Create private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %v", err)
	}

	// Derive public key from private key
	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	// Get Polkadot RPC URL from environment
	rpcURL := os.Getenv("POLKADOT_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://rpc.polkadot.io" // Default to official Polkadot RPC
	}

	// Create delegation verifier
	verifier := delegation.NewVerifier(rpcURL)

	return &SigningOracle{
		privateKey: privateKey,
		publicKey:  publicKey,
		verifier:   verifier,
	}, nil
}

// GetPrivateKeyHex returns the private key as a hex string
func (so *SigningOracle) GetPrivateKeyHex() string {
	return hex.EncodeToString(crypto.FromECDSA(so.privateKey))
}

// GetPublicKeyHex returns the public key as a hex string
func (so *SigningOracle) GetPublicKeyHex() string {
	return hex.EncodeToString(crypto.FromECDSAPub(so.publicKey))
}

// GetAddress returns the Ethereum address derived from the public key
func (so *SigningOracle) GetAddress() string {
	return crypto.PubkeyToAddress(*so.publicKey).Hex()
}

// SignMessage signs the given message
func (so *SigningOracle) SignMessage(msg string) (string, error) {
	// Create the message hash
	msgHash := crypto.Keccak256Hash([]byte(msg))

	// Sign the hash
	signature, err := crypto.Sign(msgHash.Bytes(), so.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	// Return the signature as a hex string
	return hex.EncodeToString(signature), nil
}

// SignEthereumMessage signs the given message with Ethereum signed message format
func (so *SigningOracle) SignEthereumMessage(msg string) (string, error) {
	// Create the message hash
	msgHash := crypto.Keccak256Hash([]byte(msg))

	// Create Ethereum signed message hash
	// Ethereum signed message prefix: "\x19Ethereum Signed Message:\n32"
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	data := append(prefix, msgHash.Bytes()...)
	ethSignedMessageHash := crypto.Keccak256(data)

	// Sign the Ethereum signed message hash
	signature, err := crypto.Sign(ethSignedMessageHash, so.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign Ethereum message: %v", err)
	}

	// Return the signature as a hex string
	return hex.EncodeToString(signature), nil
}

// SignTriplet signs keccak256(abi.encodePacked(validator, nominator, msgText))
// with the EIP-191 "\x19Ethereum Signed Message:\n32" prefix.
func (so *SigningOracle) SignTriplet(validator, nominator, msgText string) (sig []byte, err error) {
	packed := append(append([]byte(validator), []byte(nominator)...), []byte(msgText)...)
	h := crypto.Keccak256(packed)

	// EIP-191 for bytes32
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	ethSigned := crypto.Keccak256(append(prefix, h...))

	return crypto.Sign(ethSigned, so.privateKey) // returns 65 bytes: r||s||v (v in {0,1})
}

// GetVerifier returns the delegation verifier
func (so *SigningOracle) GetVerifier() *delegation.Verifier {
	return so.verifier
}
