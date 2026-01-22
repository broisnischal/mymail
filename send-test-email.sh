#!/bin/bash
# Send test email using telnet (no additional tools required)

TO=${1:-test@mymail.com}
FROM=${2:-sender@example.com}
SMTP_HOST=${3:-localhost}
SMTP_PORT=${4:-2525}

echo "Sending test email to $TO via $SMTP_HOST:$SMTP_PORT"
echo ""

# Create email content
SUBJECT="Test Email - $(date '+%Y-%m-%d %H:%M:%S')"
BODY="This is a test email sent via telnet to test the MyMail SMTP server.
Time: $(date '+%Y-%m-%d %H:%M:%S')
This email was sent using the send-test-email.sh script."

# Use telnet to send email
{
    echo "EHLO test.com"
    echo "MAIL FROM: <$FROM>"
    echo "RCPT TO: <$TO>"
    echo "DATA"
    echo "From: $FROM"
    echo "To: $TO"
    echo "Subject: $SUBJECT"
    echo ""
    echo "$BODY"
    echo "."
    echo "QUIT"
} | telnet "$SMTP_HOST" "$SMTP_PORT" 2>/dev/null

echo ""
echo "Email sent! Check the database or API to verify it was received."
