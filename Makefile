.PHONY: build up down migrate clean dev dev-infra dev-api dev-smtp dev-worker dev-ui install-air

# Production commands
build:
	docker compose -f compose.yml build

up:
	docker compose -f compose.yml up -d

down:
	docker compose -f compose.yml down

migrate:
	docker compose -f compose.yml exec api bun run migrate

logs:
	docker compose -f compose.yml logs -f

clean:
	docker compose -f compose.yml down -v
	docker system prune -f

setup:
	@echo "Setting up MyMail..."
	@cp .env.example .env
	@echo "Please edit .env file with your configuration"
	@docker compose -f compose.yml up -d postgres redis minio
	@sleep 5
	@docker compose -f compose.yml exec api bun run migrate
	@echo "Setup complete! Run 'make up' to start all services"

# Development commands
dev-infra:
	@echo "Starting development infrastructure (PostgreSQL, Redis, MinIO)..."
	@docker compose -f compose.dev.yml up -d
	@echo "✓ Infrastructure running"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis: localhost:6379"
	@echo "  MinIO: http://localhost:9000 (Console: http://localhost:9001)"
	@echo ""
	@echo "Run migrations: make migrate-dev"
	@echo "Stop infrastructure: make dev-infra-down"

dev-infra-down:
	@docker compose -f compose.dev.yml down

migrate-dev:
	@echo "Running migrations..."
	@cd api && DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail bun run migrate

migrate-generate:
	@echo "Generating migrations from schema..."
	@cd api && DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail bunx drizzle-kit generate:pg

migrate-push:
	@echo "Pushing schema directly to database (dev only)..."
	@cd api && DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail bunx drizzle-kit push:pg

dev-api:
	@echo "Starting API server with hot reload..."
	@cd api && DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail \
		REDIS_URL=redis://localhost:6379 \
		MINIO_ENDPOINT=localhost:9000 \
		MINIO_ACCESS_KEY=minioadmin \
		MINIO_SECRET_KEY=minioadmin \
		MINIO_BUCKET=mails \
		JWT_SECRET=dev-secret-key \
		API_PORT=3000 \
		API_HOST=0.0.0.0 \
		bun run dev

dev-smtp:
	@echo "Starting SMTP server with hot reload (requires Air)..."
	@AIR_BIN=$$(which air 2>/dev/null || echo "$(shell go env GOPATH)/bin/air"); \
	if [ ! -f "$$AIR_BIN" ]; then \
		echo "Air not found. Installing..."; \
		make install-air; \
		AIR_BIN="$(shell go env GOPATH)/bin/air"; \
	fi; \
	cd smtp && \
		PATH="$(shell go env GOPATH)/bin:$$PATH" \
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
		$$AIR_BIN

dev-worker:
	@echo "Starting worker with hot reload (requires Air)..."
	@AIR_BIN=$$(which air 2>/dev/null || echo "$(shell go env GOPATH)/bin/air"); \
	if [ ! -f "$$AIR_BIN" ]; then \
		echo "Air not found. Installing..."; \
		make install-air; \
		AIR_BIN="$(shell go env GOPATH)/bin/air"; \
	fi; \
	cd worker && \
		PATH="$(shell go env GOPATH)/bin:$$PATH" \
		DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mymail?sslmode=disable \
		REDIS_URL=redis://localhost:6379 \
		WORKER_CONCURRENCY=10 \
		WORKER_BATCH_SIZE=100 \
		$$AIR_BIN

dev-ui:
	@echo "Starting UI with hot reload..."
	@cd ui && bun run dev

# Install Air for Go hot reload
install-air:
	@echo "Installing Air for Go hot reload..."
	@go install github.com/air-verse/air@latest
	@echo "✓ Air installed!"
	@echo ""
	@echo "Note: If 'air' command is not found, add Go bin to your PATH:"
	@echo "  export PATH=\$$PATH:$$(go env GOPATH)/bin"
	@echo ""
	@echo "Or add to your ~/.bashrc or ~/.zshrc:"
	@echo "  export PATH=\$$PATH:$$(go env GOPATH)/bin"

# Run all dev services in separate terminals
dev-all:
	@echo "Starting all development services..."
	@echo "This will start services in the background. Use 'make dev-stop' to stop them."
	@echo ""
	@echo "Starting infrastructure..."
	@make dev-infra
	@sleep 3
	@echo "Starting API..."
	@make dev-api &
	@echo "Starting SMTP..."
	@make dev-smtp &
	@echo "Starting Worker..."
	@make dev-worker &
	@echo "Starting UI..."
	@make dev-ui &
	@echo ""
	@echo "All services starting! Check logs in each terminal."

dev-stop:
	@echo "Stopping all development services..."
	@make dev-infra-down
	@pkill -f "bun run dev" || true
	@pkill -f "air" || true
	@pkill -f "vite" || true
	@echo "✓ All services stopped"
