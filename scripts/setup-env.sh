#!/bin/bash
# Setup script for local development
# Copies .env.example to .env and opens it for editing

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_EXAMPLE="$SCRIPT_DIR/.env.example"
ENV_FILE="$SCRIPT_DIR/.env"

if [ ! -f "$ENV_EXAMPLE" ]; then
    echo "❌ .env.example not found!"
    exit 1
fi

if [ -f "$ENV_FILE" ]; then
    echo "⚠️  .env file already exists"
    read -p "Overwrite? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Cancelled"
        exit 0
    fi
fi

cp "$ENV_EXAMPLE" "$ENV_FILE"
echo "✅ Created .env from .env.example"

# Set default local development values
sed -i '' 's|INTENT_API_URL=http://localhost:8080|INTENT_API_URL=http://localhost:8080|' "$ENV_FILE" 2>/dev/null || \
sed -i 's|INTENT_API_URL=http://localhost:8080|INTENT_API_URL=http://localhost:8080|' "$ENV_FILE"

echo ""
echo "📝 Edit .env file to configure:"
echo "   - INTENT_API_URL (default: http://localhost:8080)"
echo "   - INTENT_TOKEN (get from your local API)"
echo ""
echo "To edit:"
echo "   code .env    # VS Code"
echo "   vim .env     # Vim"
echo "   nano .env    # Nano"
echo ""
echo "After editing, reload your shell or run:"
echo "   source .env && export \$(cat .env | xargs)"

