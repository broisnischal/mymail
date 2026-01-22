# Development Guide

This guide explains how to run and develop the MyMail project with hot reload.

## Prerequisites

- **Go 1.21+** - For SMTP server and worker
- **Bun** - For API server (or Node.js 20+)
- **Docker & Docker Compose** - For infrastructure services
- **Air** (optional) - For Go hot reload (`make install-air`)

## Quick Start

### 1. Start Infrastructure Services

Start PostgreSQL, Redis, and MinIO in Docker:

```bash
make dev-infra
```

This starts:
- PostgreSQL on `localhost:5432`
- Redis on `localhost:6379`
- MinIO on `http://localhost:9000` (Console: `http://localhost:9001`)

### 2. Run Database Migrations

```bash
make migrate-dev
```

### 3. Start Development Servers

You have two options:

#### Option A: Run Each Service Separately (Recommended)

Open separate terminal windows/tabs:

**Terminal 1 - API Server:**
```bash
make dev-api
```
- Runs on `http://localhost:3000`
- Hot reload with Bun's `--watch` flag
- Auto-restarts on file changes

**Terminal 2 - SMTP Server:**
```bash
make dev-smtp
```
- Runs on `localhost:2525` (using port 2525 to avoid requiring root)
- Hot reload with Air
- Auto-rebuilds and restarts on `.go` file changes

**Terminal 3 - Worker:**
```bash
make dev-worker
```
- Processes email queue jobs
- Hot reload with Air
- Auto-rebuilds and restarts on `.go` file changes

**Terminal 4 - UI:**
```bash
make dev-ui
```
- Runs on `http://localhost:5173`
- Hot reload with Vite
- Instant updates in browser

#### Option B: Use tmux/screen

You can use `tmux` or `screen` to run multiple services:

```bash
# Create a new tmux session
tmux new-session -d -s mymail

# Split into panes and run each service
tmux split-window -h
tmux split-window -v
tmux split-window -v

# Run services in each pane
tmux send-keys -t 0 "make dev-api" C-m
tmux send-keys -t 1 "make dev-smtp" C-m
tmux send-keys -t 2 "make dev-worker" C-m
tmux send-keys -t 3 "make dev-ui" C-m

# Attach to session
tmux attach -t mymail
```

## Hot Reload Details

### Go Services (SMTP & Worker)

**Using Air (Recommended):**
```bash
# Install Air first
make install-air

# Then run with hot reload
make dev-smtp
make dev-worker
```

Air watches for `.go` file changes and automatically:
- Rebuilds the Go binary
- Restarts the service
- Shows build errors in the terminal

**Manual Go Development:**
```bash
# SMTP Server
cd smtp
go run main.go

# Worker
cd worker
go run main.go
```

### API Server (Bun)

The API server uses Bun's built-in `--watch` flag:

```bash
cd api
bun run dev
```

Bun automatically:
- Restarts on TypeScript file changes
- Shows compilation errors
- Fast refresh (no full restart needed for most changes)

### UI (React + Vite)

The UI uses Vite's HMR (Hot Module Replacement):

```bash
cd ui
bun run dev
```

Vite provides:
- Instant updates in browser (no page refresh)
- Fast compilation
- Error overlay in browser

## Environment Variables

Development services use these default values:

```bash
# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail

# Redis
REDIS_URL=redis://localhost:6379

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=mails

# API
JWT_SECRET=dev-secret-key
API_PORT=3000

# SMTP (using port 2525 to avoid root requirement)
SMTP_PORT=2525
SMTP_DOMAIN=mymail.com
```

You can override these by setting environment variables before running the dev commands.

## Testing SMTP Server

### Using telnet:

```bash
telnet localhost 2525
```

Then:
```
EHLO test.com
MAIL FROM: <sender@example.com>
RCPT TO: <test@mymail.com>
DATA
Subject: Test Email

This is a test email body.
.
QUIT
```

### Using swaks (Swiss Army Knife for SMTP):

```bash
# Install swaks
brew install swaks  # macOS
# or
sudo apt-get install swaks  # Linux

# Send test email
swaks --to test@mymail.com \
      --from sender@example.com \
      --server localhost:2525 \
      --body "Test email body"
```

### Using mail command:

```bash
echo "Test body" | mail -s "Test Subject" test@mymail.com \
  -S smtp=localhost:2525 \
  -S from=sender@example.com
```

## Testing the API

```bash
# Register a user
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@mymail.com","password":"testpass123"}'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@mymail.com","password":"testpass123"}'

# Get mailboxes (use token from login)
curl http://localhost:3000/api/mailboxes \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Debugging

### View Logs

**Infrastructure logs:**
```bash
docker compose -f compose.dev.yml logs -f
```

**Go services:**
- Air shows logs in the terminal
- Build errors appear immediately

**API/Bun:**
- Logs appear in the terminal
- Use `console.log()` for debugging

**UI:**
- Browser console for runtime errors
- Terminal for build errors
- Vite error overlay in browser

### Database Access

```bash
# Connect to PostgreSQL
docker compose -f compose.dev.yml exec postgres psql -U postgres -d mymail

# Or use psql directly
psql postgresql://postgres:postgres@localhost:5432/mymail
```

### Redis Access

```bash
# Connect to Redis CLI
docker compose -f compose.dev.yml exec redis redis-cli

# Or use redis-cli directly
redis-cli -h localhost -p 6379
```

### MinIO Console

Access MinIO web console at: `http://localhost:9001`
- Username: `minioadmin`
- Password: `minioadmin`

## Stopping Services

```bash
# Stop infrastructure
make dev-infra-down

# Stop all dev services
make dev-stop
```

## Troubleshooting

### Port Already in Use

If a port is already in use:

```bash
# Find process using port
lsof -i :3000  # API
lsof -i :2525  # SMTP
lsof -i :5173  # UI

# Kill the process
kill -9 <PID>
```

### Air Not Found

```bash
# Install Air
make install-air

# Or manually
go install github.com/air-verse/air@latest

# Add to PATH (add to ~/.bashrc or ~/.zshrc for permanent)
export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
air -v
```

**Note:** Air moved from `github.com/cosmtrek/air` to `github.com/air-verse/air`. 
The Makefile and scripts have been updated to use the new repository.

### Database Connection Issues

1. Ensure infrastructure is running: `make dev-infra`
2. Check PostgreSQL is ready: `docker compose -f compose.dev.yml ps`
3. Verify connection string matches Docker service names

### Go Build Errors

```bash
# Clean Go modules
cd smtp && go mod tidy
cd ../worker && go mod tidy

# Rebuild
go clean -cache
```

## Production vs Development

| Feature | Development | Production |
|---------|------------|------------|
| Hot Reload | ✅ Yes | ❌ No |
| Source Maps | ✅ Yes | ❌ No |
| Minification | ❌ No | ✅ Yes |
| Error Details | ✅ Full | ❌ Limited |
| Port | Various | Standard (25, 3000, 80) |
| Infrastructure | Docker Compose | Docker Swarm/K8s |

## Next Steps

1. Read the [README.md](./README.md) for architecture overview
2. Check [DEPLOYMENT.md](./DEPLOYMENT.md) for production deployment
3. Review code structure in each service directory
4. Start developing! 🚀
