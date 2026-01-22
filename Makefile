.PHONY: build up down migrate clean

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	docker-compose exec api bun run migrate

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	docker system prune -f

setup:
	@echo "Setting up MyMail..."
	@cp .env.example .env
	@echo "Please edit .env file with your configuration"
	@docker-compose up -d postgres redis minio
	@sleep 5
	@docker-compose exec api bun run migrate
	@echo "Setup complete! Run 'make up' to start all services"
