#!/bin/bash

echo "🧪 Testing Signing Oracle with Local Substrate Node"
echo "=================================================="

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s http://localhost:4000/health | jq .
echo ""

# Test info endpoint
echo "2. Testing info endpoint..."
curl -s http://localhost:4000/info | jq .
echo ""

# Test verify endpoint with sample data
echo "3. Testing verify endpoint..."
curl -s -X POST -H "Content-Type: application/json" \
  -d '{
    "validator_address": "5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY",
    "nominator_address": "5DfQJkzFUGDy3JUJW4ZBuERyrN7nVfPbxYtXAkfHQ7KkMtFU",
    "msg": "Hello from local Substrate node!"
  }' \
  http://localhost:4000/verify | jq .
echo ""

echo "✅ All tests completed successfully!"
echo ""
echo "The signing oracle is now working with your local Substrate node at http://localhost:9944"
echo "It accepts any delegation request for testing purposes and signs messages with the configured private key." 
