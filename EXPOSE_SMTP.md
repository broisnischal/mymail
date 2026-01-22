# Exposing SMTP Server to the Internet

This guide shows how to make your local SMTP server accessible from the internet and configure MX records for `jotko.site`.

## Option 1: Cloudflare Tunnel (Recommended - Free)

Cloudflare Tunnel is free and doesn't require port forwarding or public IP.

### Setup Steps:

1. **Install Cloudflare Tunnel (cloudflared)**:
   ```bash
   # Arch Linux
   sudo pacman -S cloudflared
   
   # Or download from: https://github.com/cloudflare/cloudflared/releases
   ```

2. **Login to Cloudflare**:
   ```bash
   cloudflared tunnel login
   ```

3. **Create a tunnel**:
   ```bash
   cloudflared tunnel create mymail-smtp
   ```

4. **Create config file** (`~/.cloudflared/config.yml`):
   ```yaml
   tunnel: <tunnel-id-from-step-3>
   credentials-file: /home/nees/.cloudflared/<tunnel-id>.json
   
   ingress:
     - hostname: smtp.jotko.site
       service: tcp://localhost:2525
     - hostname: api.jotko.site
       service: http://localhost:3000
     - service: http_status:404
   ```

5. **Create DNS records** (in Cloudflare dashboard):
   - Type: `CNAME`
   - Name: `smtp`
   - Target: `<tunnel-id>.cfargotunnel.com`
   - Proxy: Off (gray cloud) - Important for SMTP!

6. **Create MX record** (in Cloudflare dashboard):
   - Type: `MX`
   - Name: `@` (or `jotko.site`)
   - Priority: `10`
   - Target: `smtp.jotko.site`
   - Proxy: Off (gray cloud)

7. **Run tunnel**:
   ```bash
   cloudflared tunnel run mymail-smtp
   ```

8. **Or run as service** (systemd):
   ```bash
   sudo cloudflared service install
   sudo systemctl start cloudflared
   ```

## Option 2: ngrok (Quick Testing)

Good for testing, but has limitations on free tier.

### Setup:

1. **Install ngrok**:
   ```bash
   # Download from https://ngrok.com/download
   # Or use package manager
   ```

2. **Start tunnel**:
   ```bash
   ngrok tcp 2525
   ```

3. **Note the public address** (e.g., `0.tcp.ngrok.io:12345`)

4. **Update MX record**:
   - Point MX to the ngrok address
   - Note: ngrok free tier changes addresses on restart

## Option 3: Tailscale (VPN-based)

Good for personal use, creates a private network.

### Setup:

1. **Install Tailscale**:
   ```bash
   sudo pacman -S tailscale
   ```

2. **Start Tailscale**:
   ```bash
   sudo tailscale up
   ```

3. **Get your Tailscale IP**:
   ```bash
   tailscale ip -4
   ```

4. **Configure DNS** (Tailscale MagicDNS):
   - Use Tailscale IP or hostname
   - Set up MX records pointing to Tailscale address

## Option 4: VPS with Port Forwarding

If you have a VPS, you can forward ports.

### Setup:

1. **On VPS, install socat**:
   ```bash
   sudo apt-get install socat  # Debian/Ubuntu
   sudo pacman -S socat         # Arch
   ```

2. **Create SSH tunnel**:
   ```bash
   ssh -R 2525:localhost:2525 user@your-vps.com
   ```

3. **Or use autossh for persistent connection**:
   ```bash
   autossh -M 20000 -R 2525:localhost:2525 user@your-vps.com
   ```

## Configure Your SMTP Server

### 1. Update Domain Configuration

Update `smtp/src/config/config.go` or set environment variable:

```bash
export SMTP_DOMAIN=jotko.site
```

### 2. Update SMTP Handler

The handler already checks domain. Make sure it accepts `jotko.site`:

```go
// In smtp/src/handler/backend.go
allowedDomains := []string{s.backend.cfg.SMTP.Domain, "mail.localhost", "jotko.site"}
```

### 3. Port Considerations

- **Port 25**: Standard SMTP port, but often blocked by ISPs
- **Port 2525**: Alternative port (what you're using)
- **Port 587**: Submission port (for sending emails)
- **Port 465**: SMTPS port (SSL/TLS)

For receiving emails via MX records, you typically need port 25. However:
- Most ISPs block port 25
- Cloudflare Tunnel can help bypass this
- Consider using a VPS if port 25 is required

### 4. Update SMTP Server Port (if needed)

If you want to use port 25:

```bash
export SMTP_PORT=25
# Note: Requires root privileges or use authbind/setcap
```

## DNS Configuration

### MX Record Setup

In your DNS provider (Cloudflare, etc.):

```
Type: MX
Name: @ (or jotko.site)
Priority: 10
Target: smtp.jotko.site (or your tunnel address)
TTL: 3600
```

### SPF Record (Important!)

Add TXT record for SPF:

```
Type: TXT
Name: @
Value: v=spf1 mx ~all
TTL: 3600
```

### DKIM Record (Recommended)

1. Generate DKIM key (if not already done)
2. Add TXT record:

```
Type: TXT
Name: default._domainkey
Value: v=DKIM1; k=rsa; p=<your-public-key>
TTL: 3600
```

### DMARC Record (Recommended)

```
Type: TXT
Name: _dmarc
Value: v=DMARC1; p=none; rua=mailto:admin@jotko.site
TTL: 3600
```

## Testing

### 1. Test MX Record

```bash
dig MX jotko.site
# or
nslookup -type=MX jotko.site
```

### 2. Test SMTP Connection

From external server:
```bash
telnet smtp.jotko.site 2525
# or
swaks --to test@jotko.site --from sender@example.com --server smtp.jotko.site --port 2525
```

### 3. Send Test Email

```bash
# From external email service
# Send email to: test@jotko.site
```

## Security Considerations

1. **Rate Limiting**: Already implemented ✅
2. **TLS/SSL**: Configure TLS certificates
3. **Firewall**: Only expose necessary ports
4. **Authentication**: Consider adding auth for sending
5. **Monitoring**: Monitor for abuse

## Troubleshooting

### MX Record Not Working

1. Check DNS propagation: https://dnschecker.org
2. Verify MX record points to correct address
3. Check if port 25 is accessible (if using standard port)
4. Verify tunnel/proxy is running

### Emails Not Received

1. Check SMTP server logs
2. Check worker logs
3. Verify mailbox exists in database
4. Check queue jobs for errors

### Port 25 Blocked

- Use Cloudflare Tunnel (bypasses port blocking)
- Use alternative port (2525) and configure MX accordingly
- Use VPS with port 25 access

## Recommended Setup for jotko.site

1. **Use Cloudflare Tunnel** (free, reliable)
2. **Configure MX record** pointing to `smtp.jotko.site`
3. **Set up SPF, DKIM, DMARC** records
4. **Use port 2525** (avoid ISP port 25 blocking)
5. **Monitor logs** for incoming emails

## Quick Start Script

Create `start-tunnel.sh`:

```bash
#!/bin/bash
# Start Cloudflare Tunnel for SMTP

export SMTP_DOMAIN=jotko.site
export SMTP_PORT=2525

# Start tunnel
cloudflared tunnel run mymail-smtp &

# Start SMTP server
cd smtp && air
```

Make it executable:
```bash
chmod +x start-tunnel.sh
```
