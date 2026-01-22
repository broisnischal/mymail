#!/bin/bash
# Setup script for Cloudflare Tunnel to expose SMTP server

set -e

echo "=== Cloudflare Tunnel Setup for MyMail ==="
echo ""

# Check if cloudflared is installed
if ! command -v cloudflared &> /dev/null; then
    echo "✗ cloudflared not found"
    echo ""
    echo "Install with:"
    echo "  sudo pacman -S cloudflared"
    echo "  # Or download from: https://github.com/cloudflare/cloudflared/releases"
    exit 1
fi

echo "✓ cloudflared found"
echo ""

# Check if logged in
if [ ! -f ~/.cloudflared/cert.pem ]; then
    echo "Please login to Cloudflare first:"
    echo "  cloudflared tunnel login"
    exit 1
fi

echo "✓ Cloudflare login detected"
echo ""

# Create tunnel
echo "Creating tunnel 'mymail-smtp'..."
TUNNEL_OUTPUT=$(cloudflared tunnel create mymail-smtp 2>&1)
TUNNEL_ID=$(echo "$TUNNEL_OUTPUT" | grep -oP 'Created tunnel \K[^ ]+' || echo "")

if [ -z "$TUNNEL_ID" ]; then
    echo "Tunnel might already exist. Checking..."
    TUNNEL_ID=$(cloudflared tunnel list | grep mymail-smtp | awk '{print $1}' || echo "")
fi

if [ -z "$TUNNEL_ID" ]; then
    echo "✗ Failed to get tunnel ID"
    exit 1
fi

echo "✓ Tunnel ID: $TUNNEL_ID"
echo ""

# Create config directory
mkdir -p ~/.cloudflared

# Create config file
CONFIG_FILE="$HOME/.cloudflared/config.yml"
cat > "$CONFIG_FILE" << EOF
tunnel: $TUNNEL_ID
credentials-file: $HOME/.cloudflared/$TUNNEL_ID.json

ingress:
  - hostname: smtp.jotko.site
    service: tcp://localhost:2525
  - hostname: api.jotko.site
    service: http://localhost:3000
  - service: http_status:404
EOF

echo "✓ Created config file: $CONFIG_FILE"
echo ""

echo "=== Next Steps ==="
echo ""
echo "1. Create DNS records in Cloudflare dashboard:"
echo "   - CNAME: smtp.jotko.site -> $TUNNEL_ID.cfargotunnel.com (Proxy: OFF)"
echo "   - MX: @ -> smtp.jotko.site (Priority: 10, Proxy: OFF)"
echo ""
echo "2. Add SPF record (TXT):"
echo "   - Name: @"
echo "   - Value: v=spf1 mx ~all"
echo ""
echo "3. Start the tunnel:"
echo "   cloudflared tunnel run mymail-smtp"
echo ""
echo "   Or run as service:"
echo "   sudo cloudflared service install"
echo "   sudo systemctl start cloudflared"
echo ""
echo "4. Set environment variable:"
echo "   export SMTP_DOMAIN=jotko.site"
echo ""
echo "5. Start your SMTP server:"
echo "   make dev-smtp"
echo ""
