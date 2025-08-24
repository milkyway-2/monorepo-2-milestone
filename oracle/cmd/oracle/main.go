package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"oracle/pkg/signingoracle"

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

// VerifyHandler handles the /verify endpoint
func VerifyHandler(so *signingoracle.SigningOracle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		verifier := so.GetVerifier()
		isDelegated, err := verifier.VerifyDelegation(req.NominatorAddress, req.ValidatorAddress)
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
			Signature:        "0x" + signature,
		}

		// Return the response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// InfoHandler provides information about the oracle's keys
func InfoHandler(so *signingoracle.SigningOracle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		info := map[string]string{
			"public_key": so.GetPublicKeyHex(),
			"address":    so.GetAddress(),
			"status":     "ready",
		}

		json.NewEncoder(w).Encode(info)
	}
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
	oracle, err := signingoracle.NewSigningOracle()
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
	r.HandleFunc("/verify", VerifyHandler(oracle)).Methods("POST", "OPTIONS")
	r.HandleFunc("/info", InfoHandler(oracle)).Methods("GET")
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
