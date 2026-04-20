BACKEND_DIR  := backend
FRONTEND_DIR := frontend
DOCKER_COMPOSE := docker compose
GO := go

.PHONY: all dev build test lint clean migrate generate \
        dev-backend dev-user dev-admin \
        test-backend test-frontend \
        lint-backend lint-frontend \
        build-backend build-user build-admin

# ── Default ─────────────────────────────────────────────────────────────────
all: build

# ── Development ─────────────────────────────────────────────────────────────
dev:
	$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.dev.yml up

dev-backend:
	cd $(BACKEND_DIR) && air

dev-user:
	cd $(FRONTEND_DIR) && npm run dev:user

dev-admin:
	cd $(FRONTEND_DIR) && npm run dev:admin

# ── Build ────────────────────────────────────────────────────────────────────
build: build-backend build-user build-admin

build-backend:
	cd $(BACKEND_DIR) && $(GO) build -o bin/server ./cmd/server

build-user:
	cd $(FRONTEND_DIR) && npm run build:user

build-admin:
	cd $(FRONTEND_DIR) && npm run build:admin

# ── Test ─────────────────────────────────────────────────────────────────────
test: test-backend test-frontend

test-backend:
	cd $(BACKEND_DIR) && $(GO) test -race -coverprofile=coverage.out ./...

test-frontend:
	cd $(FRONTEND_DIR) && npm run test

# ── Lint ─────────────────────────────────────────────────────────────────────
lint: lint-backend lint-frontend

lint-backend:
	cd $(BACKEND_DIR) && golangci-lint run ./...

lint-frontend:
	cd $(FRONTEND_DIR) && npm run lint

# ── Database ─────────────────────────────────────────────────────────────────
migrate:
	cd $(BACKEND_DIR) && $(GO) run ./cmd/server -migrate

migrate-down:
	cd $(BACKEND_DIR) && migrate -path internal/infrastructure/database/migrations -database "$$DATABASE_URL" down 1

# ── Generate ─────────────────────────────────────────────────────────────────
generate:
	cd $(BACKEND_DIR) && $(GO) generate ./...

# ── Docker ───────────────────────────────────────────────────────────────────
docker-build:
	$(DOCKER_COMPOSE) build

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

# ── Setup ────────────────────────────────────────────────────────────────────
setup:
	cd $(BACKEND_DIR) && sqlc generate
	cd $(BACKEND_DIR) && $(GO) mod tidy && $(GO) mod download
	cd $(FRONTEND_DIR) && npm install

# ── Clean ────────────────────────────────────────────────────────────────────
clean:
	cd $(BACKEND_DIR) && rm -f bin/server coverage.out coverage.html
	cd $(FRONTEND_DIR) && rm -rf apps/user/dist apps/admin/dist
