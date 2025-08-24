package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"oracle/pkg/delegation"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Request represents the incoming request structure
type Request struct {
	ValidatorAddress string `json:"validator_address"`
	NominatorAddress string `json:"nominator_address"`
	Msg              string `json:"msg"`
}

// Response represents the response structure
type Response struct {
	ValidatorAddress string `json:"validator_address"`
	NominatorAddress string `json:"nominator_address"`
	Msg              string `json:"msg"`
	Signature        string `json:"signature"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

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

// VerifyHandler handles the /verify endpoint
func (so *SigningOracle) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST method
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ValidatorAddress == "" || req.NominatorAddress == "" || req.Msg == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Verify delegation
	isDelegated, err := so.verifier.VerifyDelegation(req.NominatorAddress, req.ValidatorAddress)
	if err != nil {
		log.Printf("Error verifying delegation: %v", err)
		errorResp := ErrorResponse{
			Error:   "verification_failed",
			Message: fmt.Sprintf("Failed to verify delegation: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	if !isDelegated {
		errorResp := ErrorResponse{
			Error:   "delegation_not_found",
			Message: "Nominator has not delegated to the specified validator",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// Sign the message
	signature, err := so.SignMessage(req.Msg)
	if err != nil {
		log.Printf("Error signing message: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the response
	response := Response{
		ValidatorAddress: req.ValidatorAddress,
		NominatorAddress: req.NominatorAddress,
		Msg:              req.Msg,
		Signature:        signature,
	}

	// Return the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// InfoHandler provides information about the oracle's keys
func (so *SigningOracle) InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	info := map[string]string{
		"public_key": so.GetPublicKeyHex(),
		"address":    so.GetAddress(),
		"status":     "ready",
	}

	json.NewEncoder(w).Encode(info)
}

// HealthHandler provides a simple health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Create a new signing oracle
	oracle, err := NewSigningOracle()
	if err != nil {
		log.Fatalf("Failed to create signing oracle: %v", err)
	}

	// Log oracle information
	log.Printf("Oracle initialized successfully")
	log.Printf("Public Key: %s", oracle.GetPublicKeyHex())
	log.Printf("Address: %s", oracle.GetAddress())

	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/verify", oracle.VerifyHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/info", oracle.InfoHandler).Methods("GET")
	r.HandleFunc("/health", HealthHandler).Methods("GET")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}

	// Start the server
	log.Printf("Starting signing oracle service on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  POST /verify - Sign a message (with delegation verification)")
	log.Printf("  GET  /info   - Get oracle information")
	log.Printf("  GET  /health - Health check")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
