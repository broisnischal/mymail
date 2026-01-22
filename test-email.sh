#!/bin/bash
# Test email script for MyMail

TO=${1:-test@mymail.com}
FROM=${2:-sender@example.com}
SMTP_SERVER=${3:-localhost:2525}

echo "Sending test email..."
echo "  To: $TO"
echo "  From: $FROM"
echo "  Server: $SMTP_SERVER"
echo ""

if command -v swaks &> /dev/null; then
    swaks --to "$TO" \
          --from "$FROM" \
          --server "$SMTP_SERVER" \
          --body "This is a test email sent via swaks to test the MyMail SMTP server." \
          --subject "Test Email - $(date '+%Y-%m-%d %H:%M:%S')"
else
    echo "swaks not found. Install with:"
    echo "  brew install swaks        # macOS"
    echo "  sudo apt-get install swaks # Debian/Ubuntu"
    echo ""
    echo "Or use telnet method (see below)"
fi
