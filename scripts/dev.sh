#!/bin/bash

# Development script to run all services with hot reload
# Usage: ./scripts/dev.sh [service]
# Services: api, smtp, worker, ui, all, infra

set -e

SERVICE=${1:-all}

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

echo_success() {
    echo -e "${GREEN}✓${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

start_infra() {
    echo_info "Starting infrastructure services..."
    docker compose -f compose.dev.yml up -d
    echo_success "Infrastructure running"
    echo "  PostgreSQL: localhost:5432"
    echo "  Redis: localhost:6379"
    echo "  MinIO: http://localhost:9000"
}

stop_infra() {
    echo_info "Stopping infrastructure services..."
    docker compose -f compose.dev.yml down
    echo_success "Infrastructure stopped"
}

start_api() {
    echo_info "Starting API server..."
    cd api
    DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail \
    REDIS_URL=redis://localhost:6379 \
    MINIO_ENDPOINT=localhost:9000 \
    MINIO_ACCESS_KEY=minioadmin \
    MINIO_SECRET_KEY=minioadmin \
    MINIO_BUCKET=mails \
    JWT_SECRET=dev-secret-key \
    API_PORT=3000 \
    API_HOST=0.0.0.0 \
    bun run dev
}

start_smtp() {
    # Find air binary
    AIR_BIN=$(which air 2>/dev/null || echo "$(go env GOPATH)/bin/air")
    
    if [ ! -f "$AIR_BIN" ]; then
        echo_warn "Air not found. Installing..."
        go install github.com/air-verse/air@latest
        AIR_BIN="$(go env GOPATH)/bin/air"
    fi
    
    echo_info "Starting SMTP server..."
    cd smtp
    PATH="$(go env GOPATH)/bin:$PATH" \
    DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail?sslmode=disable \
    REDIS_URL=redis://localhost:6379 \
    MINIO_ENDPOINT=localhost:9000 \
    MINIO_ACCESS_KEY=minioadmin \
    MINIO_SECRET_KEY=minioadmin \
    MINIO_BUCKET=mails \
    SMTP_HOST=0.0.0.0 \
    SMTP_PORT=25 \
    SMTP_DOMAIN=mymail.com \
    SMTP_MAX_SIZE=10485760 \
    "$AIR_BIN"
}

start_worker() {
    # Find air binary
    AIR_BIN=$(which air 2>/dev/null || echo "$(go env GOPATH)/bin/air")
    
    if [ ! -f "$AIR_BIN" ]; then
        echo_warn "Air not found. Installing..."
        go install github.com/air-verse/air@latest
        AIR_BIN="$(go env GOPATH)/bin/air"
    fi
    
    echo_info "Starting worker..."
    cd worker
    PATH="$(go env GOPATH)/bin:$PATH" \
    DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail?sslmode=disable \
    REDIS_URL=redis://localhost:6379 \
    WORKER_CONCURRENCY=10 \
    WORKER_BATCH_SIZE=100 \
    "$AIR_BIN"
}

start_ui() {
    echo_info "Starting UI..."
    cd ui
    bun run dev
}

case "$SERVICE" in
    infra)
        start_infra
        ;;
    api)
        start_infra
        sleep 2
        start_api
        ;;
    smtp)
        start_infra
        sleep 2
        start_smtp
        ;;
    worker)
        start_infra
        sleep 2
        start_worker
        ;;
    ui)
        start_ui
        ;;
    all)
        echo_info "Starting all services..."
        echo_warn "This will start services in separate processes."
        echo_warn "Press Ctrl+C to stop all services."
        echo ""
        start_infra
        sleep 3
        
        # Start services in background
        start_api &
        API_PID=$!
        start_smtp &
        SMTP_PID=$!
        start_worker &
        WORKER_PID=$!
        start_ui &
        UI_PID=$!
        
        echo_success "All services started!"
        echo "  API: http://localhost:3000 (PID: $API_PID)"
        echo "  SMTP: localhost:2525 (PID: $SMTP_PID)"
        echo "  Worker: running (PID: $WORKER_PID)"
        echo "  UI: http://localhost:5173 (PID: $UI_PID)"
        echo ""
        echo "Press Ctrl+C to stop all services"
        
        # Wait for interrupt
        trap "kill $API_PID $SMTP_PID $WORKER_PID $UI_PID 2>/dev/null; stop_infra; exit" INT TERM
        wait
        ;;
    stop)
        stop_infra
        pkill -f "bun run dev" || true
        pkill -f "air" || true
        pkill -f "vite" || true
        echo_success "All services stopped"
        ;;
    *)
        echo "Usage: $0 [infra|api|smtp|worker|ui|all|stop]"
        exit 1
        ;;
esac
