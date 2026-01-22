# Deployment Guide

## Prerequisites

1. **Domain Setup**
   - Configure MX records pointing to your server
   - Set up SPF, DKIM, and DMARC records
   - Ensure port 25 is open (SMTP)

2. **Server Requirements**
   - Docker and Docker Compose (or Docker Swarm)
   - Minimum 4GB RAM
   - 50GB+ storage for emails
   - CPU: 2+ cores recommended

## Quick Start

### 1. Clone and Configure

```bash
git clone <your-repo>
cd mymail
cp .env.example .env
# Edit .env with your configuration
```

### 2. Start Services

```bash
docker-compose up -d
```

### 3. Run Migrations

```bash
docker-compose exec api bun run migrate
```

### 4. Verify Services

```bash
docker-compose ps
```

## Docker Swarm Deployment

### Initialize Swarm

```bash
docker swarm init
```

### Deploy Stack

```bash
docker stack deploy -c docker-compose.yml mymail
```

### Scale Services

```bash
docker service scale mymail_worker=5
docker service scale mymail_api=3
```

### View Services

```bash
docker service ls
docker service ps mymail_api
```

## DNS Configuration

### MX Record
```
Type: MX
Name: @
Value: 10 mail.yourdomain.com
Priority: 10
```

### A Record
```
Type: A
Name: mail
Value: YOUR_SERVER_IP
```

### SPF Record
```
Type: TXT
Name: @
Value: v=spf1 mx ~all
```

### DKIM Record
```
Type: TXT
Name: default._domainkey
Value: v=DKIM1; k=rsa; p=YOUR_PUBLIC_KEY
```

### DMARC Record
```
Type: TXT
Name: _dmarc
Value: v=DMARC1; p=none; rua=mailto:admin@yourdomain.com
```

## TLS/SSL Setup

### Generate Certificates

```bash
# Using Let's Encrypt
certbot certonly --standalone -d mail.yourdomain.com

# Copy certificates
cp /etc/letsencrypt/live/mail.yourdomain.com/fullchain.pem ./certs/
cp /etc/letsencrypt/live/mail.yourdomain.com/privkey.pem ./certs/
```

### Update Environment

```env
TLS_ENABLED=true
TLS_CERT_FILE=/certs/fullchain.pem
TLS_KEY_FILE=/certs/privkey.pem
```

## Monitoring

### Health Checks

```bash
# API
curl http://localhost:3000/health

# Services
docker-compose ps
```

### Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f smtp
docker-compose logs -f worker
```

## Backup

### Database Backup

```bash
docker-compose exec postgres pg_dump -U postgres mymail > backup.sql
```

### MinIO Backup

```bash
# MinIO data is in docker volume
docker volume inspect mymail_minio-data
```

## Troubleshooting

### SMTP Not Receiving Emails

1. Check MX records: `dig MX yourdomain.com`
2. Verify port 25 is open: `telnet your-server-ip 25`
3. Check SMTP logs: `docker-compose logs smtp`

### Database Connection Issues

1. Check postgres is running: `docker-compose ps postgres`
2. Verify connection string in .env
3. Check logs: `docker-compose logs postgres`

### Worker Not Processing

1. Check worker logs: `docker-compose logs worker`
2. Verify Redis connection
3. Check queue_jobs table in database

## Performance Tuning

### Database

```env
DB_MAX_CONNECTIONS=50
```

### Worker

```env
WORKER_CONCURRENCY=20
WORKER_BATCH_SIZE=200
```

### Redis

- Increase maxmemory in redis.conf
- Set appropriate eviction policy

## Security Checklist

- [ ] Change default JWT_SECRET
- [ ] Use strong database passwords
- [ ] Enable TLS for SMTP
- [ ] Configure firewall rules
- [ ] Set up rate limiting
- [ ] Enable DKIM signing
- [ ] Configure SPF/DMARC
- [ ] Regular backups
- [ ] Monitor logs for abuse
