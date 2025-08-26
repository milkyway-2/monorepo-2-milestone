package signatureverifier

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Message represents the delegation message structure
type Message struct {
	ValidatorAddress string
	NominatorAddress string
	MsgText          string
}

// OracleVerifiedDelegation represents the verification logic from the smart contract
type OracleVerifiedDelegation struct {
	OracleAddress common.Address
}

// NewOracleVerifiedDelegation creates a new verifier instance
func NewOracleVerifiedDelegation(oracleAddressHex string) (*OracleVerifiedDelegation, error) {
	if !common.IsHexAddress(oracleAddressHex) {
		return nil, fmt.Errorf("invalid oracle address: %s", oracleAddressHex)
	}

	return &OracleVerifiedDelegation{
		OracleAddress: common.HexToAddress(oracleAddressHex),
	}, nil
}

// SubmitMessage verifies and processes a delegation message
// This mirrors the smart contract's submitMessage function
func (o *OracleVerifiedDelegation) SubmitMessage(
	validatorAddress string,
	nominatorAddress string,
	msgText string,
	signatureHex string,
) error {
	// Step 1: Decode the signature
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return fmt.Errorf("invalid signature hex: %w", err)
	}

	if len(signature) != 65 {
		return fmt.Errorf("invalid signature length: expected 65, got %d", len(signature))
	}

	// Step 2: Rebuild message hash (matches smart contract logic)
	messageHash := o.createMessageHash(validatorAddress, nominatorAddress, msgText)

	// Step 3: Create Ethereum signed message hash
	ethSignedMessageHash := o.toEthSignedMessageHash(messageHash)

	// Step 4: Recover signer from signature
	recoveredAddress, err := o.recoverSigner(ethSignedMessageHash, signature)
	if err != nil {
		return fmt.Errorf("failed to recover signer: %w", err)
	}

	// Step 5: Verify the recovered address matches the oracle address
	if recoveredAddress != o.OracleAddress {
		return fmt.Errorf("signature not from oracle: expected %s, got %s",
			o.OracleAddress.Hex(), recoveredAddress.Hex())
	}

	return nil
}

// createMessageHash creates the message hash from concatenated parameters
// This matches the smart contract's keccak256(abi.encodePacked(...)) logic
func (o *OracleVerifiedDelegation) createMessageHash(
	validatorAddress string,
	nominatorAddress string,
	msgText string,
) []byte {
	// Concatenate the parameters as they would be in abi.encodePacked
	message := validatorAddress + nominatorAddress + msgText

	// Create Keccak256 hash (Ethereum's standard hash function)
	hash := crypto.Keccak256([]byte(message))
	return hash
}

// toEthSignedMessageHash creates the Ethereum signed message hash
// This matches the smart contract's toEthSignedMessageHash function
func (o *OracleVerifiedDelegation) toEthSignedMessageHash(messageHash []byte) []byte {
	// Ethereum signed message prefix: "\x19Ethereum Signed Message:\n32"
	prefix := []byte("\x19Ethereum Signed Message:\n32")

	// Concatenate prefix with the message hash
	data := append(prefix, messageHash...)

	// Create hash of the concatenated data
	hash := crypto.Keccak256(data)
	return hash
}

// recoverSigner recovers the signer address from the signature
// This matches the smart contract's recoverSigner function
func (o *OracleVerifiedDelegation) recoverSigner(ethSignedMessageHash []byte, signature []byte) (common.Address, error) {
	if len(signature) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length: expected 65, got %d", len(signature))
	}

	// Use the signature directly with crypto.Ecrecover
	pubKey, err := crypto.Ecrecover(ethSignedMessageHash, signature)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover public key: %w", err)
	}

	// Extract the address from the public key
	var address common.Address
	copy(address[:], crypto.Keccak256(pubKey[1:])[12:])

	return address, nil
}

// VerifyMessage is a convenience function that combines all verification steps
func (o *OracleVerifiedDelegation) VerifyMessage(msg Message, signatureHex string) error {
	return o.SubmitMessage(msg.ValidatorAddress, msg.NominatorAddress, msg.MsgText, signatureHex)
}

// GetOracleAddress returns the oracle address
func (o *OracleVerifiedDelegation) GetOracleAddress() common.Address {
	return o.OracleAddress
}
