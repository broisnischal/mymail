# MyMail - Self-Hosted SMTP Mail Server

A modern, self-hosted email server built with Bun, Go, React, PostgreSQL, Redis, and MinIO. Designed to handle 100k users and 2M emails per day with low memory usage and fast response times.

## Architecture

- **API Server**: Bun + TypeScript + Drizzle ORM
- **SMTP Server**: Go (with DKIM, SPF, DMARC, TLS support)
- **Worker/Queue**: Go (processes emails asynchronously)
- **Frontend**: React + React Router + Tailwind CSS
- **Database**: PostgreSQL
- **Cache**: Redis
- **Storage**: MinIO (for email storage)

## Features

- ✅ Self-hosted SMTP server
- ✅ User authentication
- ✅ Mailbox management (aliases, temp addresses)
- ✅ Auto-create temp mailboxes (e.g., abc@mymail.com)
- ✅ Email parsing and storage
- ✅ Rate limiting and anti-abuse measures
- ✅ Modern UI for email management
- ✅ Docker Swarm ready
- ✅ Low memory footprint
- ✅ High performance

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Domain with MX records configured
- (Optional) TLS certificates for SMTP

### Environment Variables

Create a `.env` file:

```env
JWT_SECRET=your-secret-key-here
SMTP_DOMAIN=mymail.com
DATABASE_URL=postgresql://postgres:postgres@postgres:5432/mymail
REDIS_URL=redis://redis:6379
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
```

### Start Services

```bash
docker-compose up -d
```

### Run Migrations

```bash
docker-compose exec api bun run migrate
```

### Access Services

- **UI**: http://localhost
- **API**: http://localhost:3000
- **SMTP**: localhost:25
- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)

## Configuration

### SMTP Configuration

Configure your domain's MX records to point to your server:

```
MX Record: 10 mail.mymail.com
A Record: mail.mymail.com -> YOUR_SERVER_IP
```

### DKIM Setup

1. Generate DKIM keys:
```bash
openssl genrsa -out dkim-private.pem 2048
openssl rsa -in dkim-private.pem -pubout -out dkim-public.pem
```

2. Add DNS TXT record:
```
default._domainkey.mymail.com TXT "v=DKIM1; k=rsa; p=YOUR_PUBLIC_KEY"
```

3. Set environment variables:
```env
DKIM_ENABLED=true
DKIM_PRIVATE_KEY=...
DKIM_SELECTOR=default
DKIM_DOMAIN=mymail.com
```

### SPF Record

Add to your DNS:
```
mymail.com TXT "v=spf1 mx ~all"
```

### DMARC Record

Add to your DNS:
```
_dmarc.mymail.com TXT "v=DMARC1; p=none; rua=mailto:admin@mymail.com"
```

## Development

### API Server

```bash
cd api
bun install
bun run dev
```

### SMTP Server

```bash
cd smtp
go mod download
go run main.go
```

### Worker

```bash
cd worker
go mod download
go run main.go
```

### UI

```bash
cd ui
npm install
npm run dev
```

## Deployment

### Docker Swarm

```bash
docker swarm init
docker stack deploy -c docker-compose.yml mymail
```

### Scaling

```bash
docker service scale mymail_worker=5
```

## Performance

- Handles 100k users per day
- Processes 2M emails per day (100k users × 20 emails)
- Low memory usage (< 512MB per service)
- Fast response times (< 100ms API, < 50ms cache)

## Security

- Rate limiting per user and IP
- JWT authentication
- TLS support for SMTP
- DKIM, SPF, DMARC support
- Input validation
- SQL injection protection (parameterized queries)

## License

MIT
