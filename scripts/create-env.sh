#!/bin/bash
# Create .env file from .env.example with generated JWT secret

set -e

if [ -f .env ]; then
    echo "⚠️  .env already exists!"
    read -p "Overwrite? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
fi

echo "Creating .env from .env.example..."

# Copy example file
cp .env.example .env

# Generate JWT secret if openssl is available
if command -v openssl &> /dev/null; then
    JWT_SECRET=$(openssl rand -hex 32)
    # Replace JWT_SECRET line (works on both Linux and macOS)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s/^JWT_SECRET=.*/JWT_SECRET=$JWT_SECRET/" .env
    else
        sed -i "s/^JWT_SECRET=.*/JWT_SECRET=$JWT_SECRET/" .env
    fi
    echo "✓ Generated JWT_SECRET"
else
    echo "⚠️  openssl not found - please update JWT_SECRET manually"
fi

echo ""
echo "✓ Created .env file"
echo ""
echo "Next steps:"
echo "1. Review .env and update any values if needed"
echo "2. Load variables: direnv allow (or export manually)"
echo ""
