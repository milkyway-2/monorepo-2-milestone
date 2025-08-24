#!/bin/bash

echo "=== Generating ECDSA Key Pair ==="
echo ""

# Generate private key in PEM format
openssl ecparam -genkey -name secp256k1 -out private_key.pem

# Extract private key in hex format using a more reliable method
echo "=== Private Key (Hex Format) ==="
PRIVATE_KEY_HEX=$(openssl ec -in private_key.pem -text -noout | sed -n '/priv:/,/pub:/p' | grep -v "priv:" | grep -v "pub:" | tr -d ':' | tr -d ' ' | tr -d '\n')
echo $PRIVATE_KEY_HEX
echo ""

# Extract public key in hex format
echo "=== Public Key (Hex Format) ==="
PUBLIC_KEY_HEX=$(openssl ec -in private_key.pem -text -noout | sed -n '/pub:/,/ASN1/p' | grep -v "pub:" | grep -v "ASN1" | tr -d ':' | tr -d ' ' | tr -d '\n')
echo $PUBLIC_KEY_HEX
echo ""

# Compute Ethereum address from public key
echo "=== Ethereum Address ==="
# Remove the '04' prefix from the public key (uncompressed format)
PUBLIC_KEY_WITHOUT_PREFIX=$(echo $PUBLIC_KEY_HEX | sed 's/^04//')
# Take the last 40 characters (20 bytes) of the Keccak-256 hash
ETHEREUM_ADDRESS=$(echo -n $PUBLIC_KEY_WITHOUT_PREFIX | xxd -r -p | openssl dgst -sha3-256 -binary | tail -c 20 | xxd -p | tr -d '\n')
# Add '0x' prefix and convert to checksum address
ETHEREUM_ADDRESS_FULL="0x$ETHEREUM_ADDRESS"
echo $ETHEREUM_ADDRESS_FULL
echo ""

# Create .env file
echo "=== Creating .env file ==="
cat > .env << EOF
# Signing Oracle Environment Variables
PRIVATE_KEY=$PRIVATE_KEY_HEX
PUBLIC_KEY=$PUBLIC_KEY_HEX
ETHEREUM_ADDRESS=$ETHEREUM_ADDRESS_FULL
POLKADOT_RPC_URL=https://rpc.polkadot.io
PORT=4000
EOF

echo "Created .env file with the following contents:"
echo ""
cat .env
echo ""

echo "=== Usage Instructions ==="
echo "1. The .env file has been created with your private key"
echo ""
echo "2. Run the signing oracle:"
echo "   go run main.go"
echo ""
echo "3. The oracle will be available at:"
echo "   http://localhost:4000"
echo ""
echo "4. Get oracle info:"
echo "   curl http://localhost:4000/info"
echo ""
echo "=== Files Created ==="
echo "- private_key.pem (PEM format)"
echo "- public_key.pem (PEM format)"
echo "- .env (environment variables)"
echo ""
echo "=== Security Note ==="
echo "Keep your private key secure and never share it!"
echo "The private key is stored in private_key.pem and .env"
echo "Make sure to add .env to your .gitignore file!" s
