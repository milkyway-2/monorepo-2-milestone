# Delegation Verification Module

This module provides functionality to verify Polkadot delegation relationships between nominators and validators.

## Features

- Verify if a nominator has delegated to a specific validator
- Check if the delegation is currently active
- RPC communication with Polkadot nodes
- Simulation mode for testing purposes

## Usage

```go
package main

import (
    "log"
    "signing-oracle/pkg/delegation"
)

func main() {
    // Create a new verifier with RPC URL
    verifier := delegation.NewVerifier("https://rpc.polkadot.io")
    
    // Verify delegation
    isDelegated, err := verifier.VerifyDelegation(
        "12ztGE9cY2p7kPJFpfvMrL6NsCUeqoiaBY3jciMqYFuFNJ2o", // nominator
        "12GTt3pfM3SjTU6UL6dQ3SMgMSvdw94PnRoF6osU6hPvxbUZ", // validator
    )
    
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    
    if isDelegated {
        log.Println("Delegation verified successfully")
    } else {
        log.Println("Delegation not found")
    }
}
```

## API Reference

### `NewVerifier(rpcURL string) *Verifier`

Creates a new delegation verifier instance.

**Parameters:**
- `rpcURL`: The URL of the Polkadot RPC endpoint

**Returns:**
- `*Verifier`: A new verifier instance

### `VerifyDelegation(nominatorAddress, validatorAddress string) (bool, error)`

Verifies if a nominator has delegated to a validator.

**Parameters:**
- `nominatorAddress`: The address of the nominator
- `validatorAddress`: The address of the validator

**Returns:**
- `bool`: `true` if delegation exists and is active, `false` otherwise
- `error`: Any error that occurred during verification

## Testing

Run the tests with:

```bash
go test ./pkg/delegation -v
```

## Notes

- The module includes simulation mode for testing purposes
- In production, you should implement actual Polkadot staking storage queries
- The module is designed to be easily extensible for different verification strategies
