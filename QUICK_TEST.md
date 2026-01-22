# Quick Test Commands

## Send Test Email

### Using Python (No Installation Required)
```bash
python3 send-test-email-python.py test@mail.localhost
```

### Using swaks (After Installation)
```bash
# Install first:
sudo pacman -S swaks

# Then use:
swaks --to test@mail.localhost \
      --from sender@example.com \
      --server localhost:2525 \
      --body "Test email body" \
      --subject "Test Subject"
```

## Check Results

### View Emails in Database
```bash
docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail -c "SELECT e.id, e.subject, e.\"from\", e.received_at FROM emails e JOIN mailboxes m ON e.mailbox_id = m.id WHERE m.address = 'test@mail.localhost' ORDER BY e.received_at DESC LIMIT 5;"
```

### Check Queue Jobs
```bash
docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail -c "SELECT id, type, status, attempts, created_at FROM queue_jobs ORDER BY created_at DESC LIMIT 5;"
```

### Remove Failed Jobs
```bash
docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail -c "DELETE FROM queue_jobs WHERE status = 'failed' AND attempts >= 3;"
```

## Common Issues

### "swaks: command not found"
- Install: `sudo pacman -S swaks`
- Or use Python script instead

### "swacks: command not found"
- Typo! It's `swaks` (not `swacks`)

### Email not received
- Check if worker is running: `make dev-worker`
- Check if SMTP server is running: `make dev-smtp`
- Check queue jobs for errors


swaks --to test@jotko.site --from sender@example.com --server smtp.jotko.site --port 2525