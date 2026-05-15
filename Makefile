include .env
export

DB_URL := $(DATABASE_URL)
MIGRATIONS_DIR := ./migrations

env-export:
	export $(cat .env | xargs)

# ── Helpers ────────────────────────────────────────────────────────────────────


help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ── Migrate ────────────────────────────────────────────────────────────────────


migrate-up: ## Apply all pending migrations
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up


migrate-down: ## Rollback the last migration
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1


migrate-down-all: ## Rollback ALL migrations
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down


migrate-reset: migrate-down-all migrate-up ## Full reset (down all + up)


migrate-status: ## Show current migration version
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version


migrate-force: ## Force a specific version (usage: make migrate-force V=1)
	@test -n "$(V)" || (echo "❌  Usage: make migrate-force V=<version>"; exit 1)
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(V)


migrate-create: ## Create a new migration (usage: make migrate-create NAME=create_something)
	@test -n "$(NAME)" || (echo "❌  Usage: make migrate-create NAME=<migration_name>"; exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

sqlc: ## Regenerate type-safe Go code from SQL queries
		sqlc generate -f internal/user/sqlc.yaml
		sqlc generate -f internal/auth/sqlc.yaml

# ── Docker ─────────────────────────────────────────────────────────────────────
 
.PHONY: docker-up
docker-up: ## Start all containers
	docker compose up -d
 
.PHONY: docker-down
docker-down: ## Stop all containers
	docker compose down
 
.PHONY: docker-build
docker-build: ## Rebuild the api image
	docker compose build api
 
.PHONY: docker-logs
docker-logs: ## Follow api logs
	docker compose logs -f api
 
.PHONY: docker-migrate
docker-migrate: ## Run migrations inside Docker (uses migrate profile)
	docker compose --profile migrate run --rm migrate
 

.PHONY: run
run: ## Run the API locally (requires .env file)
	go run cmd/main.go

