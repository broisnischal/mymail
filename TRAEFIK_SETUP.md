# Traefik Setup Guide for jotko.site

This guide explains how to set up Traefik with automatic SSL (Let's Encrypt) for your MyMail application.

## Prerequisites

1. **Domain DNS Configuration**
   - Point `jotko.site` and `www.jotko.site` to your server's IP address (A records)
   - Point `api.jotko.site` to your server's IP address (A record)
   - Point `traefik.jotko.site` to your server's IP address (A record) - optional, for dashboard

2. **Ports**
   - Port 80 (HTTP) - must be open for Let's Encrypt validation
   - Port 443 (HTTPS) - must be open for SSL traffic
   - Port 8080 (optional) - Traefik dashboard, can be closed in production

## Configuration

### 1. Update Email Address

Edit `compose.yml` and update the Let's Encrypt email:
```yaml
- --certificatesresolvers.letsencrypt.acme.email=admin@jotko.site
```
Change `admin@jotko.site` to your actual email address.

### 2. Set Traefik Dashboard Password (Optional)

Generate a password hash:
```bash
htpasswd -nb admin yourpassword
```

Set it in your `.env` file:
```bash
TRAEFIK_DASHBOARD_AUTH=admin:$$apr1$$yourhashedpassword
```

Or set it directly in `compose.yml` (line 32).

### 3. Start Services

```bash
docker compose up -d
```

### 4. Verify Setup

- **Main website**: https://jotko.site
- **API**: https://api.jotko.site
- **Traefik Dashboard**: https://traefik.jotko.site (if configured)

## DNS Records Required

Add these DNS records in your domain provider:

```
Type    Name    Value              TTL
A       @       YOUR_SERVER_IP     3600
A       www     YOUR_SERVER_IP     3600
A       api     YOUR_SERVER_IP     3600
A       traefik YOUR_SERVER_IP     3600  (optional)
```

## How It Works

1. **Traefik** acts as a reverse proxy and handles:
   - Automatic SSL certificate generation via Let's Encrypt
   - HTTP to HTTPS redirection
   - Routing based on domain names
   - Security headers

2. **Routing**:
   - `jotko.site` → UI (frontend)
   - `www.jotko.site` → Redirects to `jotko.site`
   - `api.jotko.site` → API backend
   - `traefik.jotko.site` → Traefik dashboard

3. **SSL Certificates**:
   - Automatically generated on first request
   - Stored in `traefik-letsencrypt` volume
   - Auto-renewed before expiration

## Troubleshooting

### Certificates Not Generating

1. Ensure ports 80 and 443 are open
2. Verify DNS records are pointing to your server
3. Check Traefik logs: `docker compose logs traefik`
4. Ensure domain is accessible: `curl -I http://jotko.site`

### Dashboard Not Accessible

1. Check if password is set correctly
2. Verify DNS record for `traefik.jotko.site`
3. Check Traefik logs for errors

### Services Not Routing

1. Verify services are on `traefik-network`: `docker network inspect traefik-network`
2. Check service labels in `compose.yml`
3. View Traefik dashboard to see registered routes

## Security Notes

- Remove port 8080 mapping in production if you don't need the dashboard
- Use strong passwords for Traefik dashboard
- Keep Traefik updated: `docker compose pull traefik`
- Monitor Let's Encrypt rate limits (50 certs/week per domain)

## Updating Configuration

After changing Traefik configuration:
```bash
docker compose up -d --force-recreate traefik
```

## Additional Domains

To add more domains, update the labels in `compose.yml`:

```yaml
- "traefik.http.routers.ui.rule=Host(`jotko.site`) || Host(`example.com`)"
```

Then update the Let's Encrypt email and restart.
