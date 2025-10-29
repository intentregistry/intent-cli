#!/bin/bash
# Generate ed25519 key pair for intent package signing

echo "🔑 Generating ed25519 key pair for intent package signing..."
echo ""

# Generate private key using OpenSSL
openssl genpkey -algorithm ed25519 -out private_key.pem

# Extract private key in hex format (64 bytes = 128 hex chars)
# This is what intent CLI expects
PRIV_HEX=$(openssl pkey -in private_key.pem -noout -text | \
  grep -A 32 "priv:" | \
  tail -n +2 | \
  tr -d ' :\n' | \
  head -c 128)

echo "$PRIV_HEX" > private_key.hex

# Extract public key in hex format (32 bytes = 64 hex chars)
PUB_HEX=$(openssl pkey -in private_key.pem -pubout -outform DER 2>/dev/null | \
  tail -c 32 | \
  xxd -p -c 32 | \
  tr -d '\n')

echo "$PUB_HEX" > public_key.hex

echo "✅ Key pair generated!"
echo ""
echo "📁 Files created:"
echo "  - private_key.pem (PEM format - keep SECRET!)"
echo "  - private_key.hex (hex format - use with --sign-key)"
echo "  - public_key.hex (hex format - share for verification)"
echo ""
echo "🔐 Usage examples:"
echo ""
echo "  # Option 1: Use environment variable"
echo "  export INTENT_SIGN_KEY=\$(pwd)/private_key.hex"
echo "  intent package ."
echo ""
echo "  # Option 2: Use --sign-key flag"
echo "  intent package . --sign-key private_key.hex"
echo ""
echo "⚠️  Keep private_key.pem and private_key.hex SECRET!"
echo "✅ You can share public_key.hex for verification"

