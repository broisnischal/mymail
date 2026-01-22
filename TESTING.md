# Testing Guide

## Sending Test Emails

### Using swaks (Recommended)

**Install swaks:**
```bash
# macOS
brew install swaks

# Debian/Ubuntu
sudo apt-get install swaks

# Arch Linux
sudo pacman -S swaks
```

**Send test email:**
```bash
swaks --to test@mymail.com \
      --from sender@example.com \
      --server localhost:2525 \
      --body "This is a test email body" \
      --subject "Test Subject"
```

**Or use the test script:**
```bash
./test-email.sh test@mymail.com
```

### Using telnet (Manual)

```bash
telnet localhost 2525
```

Then type:
```
EHLO test.com
MAIL FROM: <sender@example.com>
RCPT TO: <test@mymail.com>
DATA
Subject: Test Email

This is the email body.
.
QUIT
```

### Using curl (SMTP)

```bash
curl --url "smtp://localhost:2525" \
     --mail-from "sender@example.com" \
     --mail-rcpt "test@mymail.com" \
     --upload-file - <<EOF
From: sender@example.com
To: test@mymail.com
Subject: Test Email

This is a test email body.
EOF
```

### Using Python

```python
import smtplib
from email.mime.text import MIMEText

msg = MIMEText("This is a test email body")
msg['Subject'] = 'Test Email'
msg['From'] = 'sender@example.com'
msg['To'] = 'test@mymail.com'

server = smtplib.SMTP('localhost', 2525)
server.send_message(msg)
server.quit()
```

## Testing Temp Mail Feature

The SMTP server automatically creates mailboxes for any address if temp mail is enabled:

```bash
# This will auto-create abc@mymail.com if it doesn't exist
swaks --to abc@mymail.com \
      --from sender@example.com \
      --server localhost:2525 \
      --body "Test temp mail"
```

## Verifying Email Reception

### Check Database

```bash
# List all mailboxes
docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail -c "SELECT address, is_temp, created_at FROM mailboxes;"

# List emails for a mailbox
docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail -c "SELECT subject, \"from\", received_at FROM emails e JOIN mailboxes m ON e.mailbox_id = m.id WHERE m.address = 'test@mymail.com';"
```

### Check via API

```bash
# Login first
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@mymail.com","password":"yourpassword"}' | jq -r '.token')

# Get emails
curl http://localhost:3000/api/emails \
  -H "Authorization: Bearer $TOKEN"
```

### Check MinIO

Emails are stored in MinIO. Access the console at:
- URL: http://localhost:9001
- Username: `minioadmin`
- Password: `minioadmin`

## Testing Queue Processing

Check if worker is processing emails:

```bash
# Check queue jobs
docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail -c "SELECT type, status, attempts, created_at FROM queue_jobs ORDER BY created_at DESC LIMIT 10;"
```

## Common Test Scenarios

### 1. Basic Email
```bash
swaks --to test@mymail.com \
      --from sender@example.com \
      --server localhost:2525 \
      --body "Hello, this is a test!"
```

### 2. Email with HTML
```bash
swaks --to test@mymail.com \
      --from sender@example.com \
      --server localhost:2525 \
      --body "<h1>HTML Email</h1><p>This is HTML content</p>" \
      --header "Content-Type: text/html"
```

### 3. Email with Attachment
```bash
swaks --to test@mymail.com \
      --from sender@example.com \
      --server localhost:2525 \
      --body "Email with attachment" \
      --attach /path/to/file.txt
```

### 4. Multiple Recipients
```bash
swaks --to test@mymail.com,test2@mymail.com \
      --from sender@example.com \
      --server localhost:2525 \
      --body "Email to multiple recipients"
```

## Troubleshooting

### Email not received?

1. **Check SMTP server logs:**
   ```bash
   # If running with Air, check the terminal
   # Or check Docker logs
   docker compose -f compose.dev.yml logs smtp
   ```

2. **Check worker logs:**
   ```bash
   # If running with Air, check the terminal
   # Or check Docker logs
   docker compose -f compose.dev.yml logs worker
   ```

3. **Verify mailbox exists:**
   ```bash
   docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail -c "SELECT * FROM mailboxes WHERE address = 'test@mymail.com';"
   ```

4. **Check queue jobs:**
   ```bash
   docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail -c "SELECT * FROM queue_jobs ORDER BY created_at DESC LIMIT 5;"
   ```

### Connection refused?

- Ensure SMTP server is running: `make dev-smtp`
- Check port: Should be `2525` in development (not `25` which requires root)

### Rate limiting?

The server has rate limits. If you hit them:
- Wait a few minutes
- Or reset Redis: `docker compose -f compose.dev.yml restart redis`
