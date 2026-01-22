# Quick Start Guide

## Development Setup (5 minutes)

### 1. Install Air for Go Hot Reload

```bash
make install-air
# Or manually: go install github.com/air-verse/air@latest
```

### 2. Start Infrastructure

```bash
make dev-infra
```

This starts PostgreSQL, Redis, and MinIO in Docker.

### 3. Run Migrations

```bash
make migrate-dev
```

### 4. Start Development Servers

Open 4 terminal windows/tabs:

**Terminal 1 - API:**
```bash
make dev-api
# Runs on http://localhost:3000
```

**Terminal 2 - SMTP:**
```bash
make dev-smtp
# Runs on localhost:2525
```

**Terminal 3 - Worker:**
```bash
make dev-worker
# Processes email queue
```

**Terminal 4 - UI:**
```bash
make dev-ui
# Runs on http://localhost:5173
```

### 5. Test It!

1. Open UI: http://localhost:5173
2. Register a user
3. Send a test email to your mailbox using:
   ```bash
   swaks --to yourmailbox@mymail.com \
         --from test@example.com \
         --server localhost:2525
   ```

## Production Setup

See [DEPLOYMENT.md](./DEPLOYMENT.md) for production deployment.

## Troubleshooting

- **Air not found**: Run `make install-air`
- **Port in use**: Change ports in Makefile or kill existing processes
- **Database errors**: Ensure `make dev-infra` is running
- **Go build errors**: Run `go mod tidy` in smtp/ and worker/ directories

For more details, see [DEVELOPMENT.md](./DEVELOPMENT.md).
