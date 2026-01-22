# Environment Configuration Guide

## Quick Start

1. **Copy the example file**:
   ```bash
   cp .env.example .env
   ```

2. **Update values** for your environment (especially `JWT_SECRET`)

3. **Load environment variables**:
   ```bash
   # Using direnv (recommended)
   direnv allow
   
   # Or manually source
   export $(cat .env | xargs)
   ```

## Configuration Sections

### API Server
- `API_PORT`: Port for the API server (default: 3000)
- `API_HOST`: Host to bind to (default: 0.0.0.0)
- `JWT_SECRET`: Secret key for JWT tokens (⚠️ CHANGE IN PRODUCTION)
- `JWT_EXPIRY`: JWT token expiration (default: 7d)

### Database (PostgreSQL)
- `DATABASE_URL`: PostgreSQL connection string
  - Development: `postgresql://postgres:postgres@localhost:5432/mymail?sslmode=disable`
  - Production: Use SSL and strong credentials
- `DB_MAX_CONNECTIONS`: Maximum database connections (default: 10)

### Redis
- `REDIS_URL`: Redis connection URL
  - Development: `redis://localhost:6379`
- `REDIS_PREFIX`: Key prefix for Redis (default: `mymail:`)

### MinIO (Object Storage)
- `MINIO_ENDPOINT`: MinIO server endpoint
  - Development: `localhost:9000`
  - Production: Your S3-compatible storage endpoint
- `MINIO_ACCESS_KEY`: Access key
- `MINIO_SECRET_KEY`: Secret key
- `MINIO_BUCKET`: Bucket name for emails (default: `mails`)
- `MINIO_USE_SSL`: Use SSL/TLS (default: false)

### SMTP Server
- `SMTP_HOST`: Host to bind SMTP server (default: 0.0.0.0)
- `SMTP_PORT`: SMTP port (default: 2525)
  - Port 25 requires root privileges
  - Use 2525 for development
- `SMTP_DOMAIN`: Your email domain (default: `jotko.site`)
- `SMTP_MAX_SIZE`: Maximum email size in bytes (default: 10MB)

### Worker
- `WORKER_CONCURRENCY`: Number of concurrent workers (default: 10)
- `WORKER_BATCH_SIZE`: Batch size for processing jobs (default: 100)

### Rate Limiting
- `RATE_LIMIT_EMAILS_PER_USER`: Max emails per user (default: 1000)
- `RATE_LIMIT_EMAILS_PER_HOUR`: Max emails per hour (default: 100)
- `RATE_LIMIT_CONNECTIONS_PER_IP`: Max connections per IP (default: 10)

### DKIM (Email Authentication)
- `DKIM_ENABLED`: Enable DKIM signing (default: false)
- `DKIM_SELECTOR`: DKIM selector (default: `default`)
- `DKIM_DOMAIN`: Domain for DKIM (default: `jotko.site`)
- `DKIM_PRIVATE_KEY`: DKIM private key (generate with OpenSSL)

**Generate DKIM key**:
```bash
openssl genrsa -out dkim_private.pem 2048
openssl rsa -in dkim_private.pem -pubout -out dkim_public.pem
# Set DKIM_PRIVATE_KEY to the contents of dkim_private.pem
```

### TLS/SSL
- `TLS_ENABLED`: Enable TLS (default: false)
- `TLS_CERT_FILE`: Path to SSL certificate
- `TLS_KEY_FILE`: Path to SSL private key

**For production**, use Let's Encrypt:
```bash
certbot certonly --standalone -d smtp.jotko.site
# Set TLS_CERT_FILE=/etc/letsencrypt/live/smtp.jotko.site/fullchain.pem
# Set TLS_KEY_FILE=/etc/letsencrypt/live/smtp.jotko.site/privkey.pem
```

### Temp Mail
- `TEMP_MAIL_ENABLED`: Enable temporary email addresses (default: true)
- `TEMP_MAIL_TTL`: Time-to-live in seconds (default: 86400 = 24 hours)

## Production Checklist

Before deploying to production:

- [ ] Change `JWT_SECRET` to a strong random string
- [ ] Update `DATABASE_URL` with production credentials
- [ ] Enable SSL for `MINIO_USE_SSL` if using external storage
- [ ] Set `SMTP_DOMAIN` to your actual domain
- [ ] Configure DKIM keys and enable `DKIM_ENABLED`
- [ ] Set up TLS certificates and enable `TLS_ENABLED`
- [ ] Review and adjust rate limiting values
- [ ] Set `NODE_ENV=production`

## Environment-Specific Files

You can create environment-specific files:
- `.env.development` - Development settings
- `.env.production` - Production settings
- `.env.test` - Test settings

Load them with:
```bash
export $(cat .env.production | xargs)
```

## Using direnv (Recommended)

If you have `direnv` installed (see `.envrc`), environment variables are automatically loaded when you enter the directory:

```bash
# Install direnv
sudo pacman -S direnv

# Add to your shell config (~/.zshrc or ~/.bashrc)
eval "$(direnv hook zsh)"

# Allow direnv in this directory
direnv allow
```

## Security Notes

1. **Never commit `.env` to git** - It's already in `.gitignore`
2. **Use strong secrets** - Generate with `openssl rand -hex 32`
3. **Use SSL in production** - Enable TLS for SMTP and SSL for MinIO
4. **Restrict database access** - Use strong passwords and SSL
5. **Enable DKIM** - Improves email deliverability

## Troubleshooting

### Variables not loading
- Check if `.env` file exists
- Verify syntax (no spaces around `=`)
- Use `direnv allow` if using direnv
- Manually export: `export $(cat .env | xargs)`

### Port conflicts
- Change `API_PORT` if 3000 is taken
- Change `SMTP_PORT` if 2525 is taken
- Check Docker ports in `compose.dev.yml`

### Database connection errors
- Verify `DATABASE_URL` is correct
- Check if PostgreSQL is running: `docker compose -f compose.dev.yml ps postgres`
- Ensure database exists: `docker compose -f compose.dev.yml exec postgres psql -U postgres -c "CREATE DATABASE mymail;"`
