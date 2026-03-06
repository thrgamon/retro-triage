default:
    @just --list

# --- Development ---

# Start all services (postgres, backend, frontend)
dev:
    docker compose up --build

# Stop all services
dev-down:
    docker compose down

# View backend logs
logs:
    docker compose logs -f backend

# Start with monitoring stack
dev-monitoring:
    docker compose --profile monitoring up --build

# Run Go backend locally (outside docker)
backend:
    go run ./cmd/server

# Run frontend locally (outside docker)
frontend:
    yarn install && yarn run dev

# Install git hooks
install-hooks:
    cp scripts/pre-push .git/hooks/pre-push
    chmod +x .git/hooks/pre-push

# --- Code Generation ---

# Regenerate sqlc Go types from SQL queries
sqlc:
    sqlc generate

# Regenerate swagger.json from Go annotations
api-docs:
    swag init -g cmd/server/main.go -o docs --parseInternal

# Regenerate frontend TypeScript client from swagger
api-types:
    npx orval

# Full sync: sqlc + type check
sync: sqlc
    yarn run check
    go vet ./...

# --- Quality ---

# Run Go linters
lint:
    golangci-lint run ./...
    go mod tidy

# Run all tests
test:
    go test -race ./...

# Format Go code
fmt:
    gofmt -s -w .

# Frontend lint
fe-lint:
    yarn run lint

# Frontend lint with auto-fix
fe-lint-fix:
    yarn run lint:fix

# Frontend format
fe-fmt:
    yarn run format

# Run all checks (lint + test + type-check)
check: lint test
    yarn run check
    yarn run lint

# --- E2E Tests ---

# Run Playwright e2e tests (requires server at localhost:3000)
e2e:
    npx playwright test

# Run Playwright with UI mode
e2e-ui:
    npx playwright test --ui

# --- Database ---

# Run migrations up
migrate:
    goose -dir migrations postgres "${DATABASE_URL}" up

# Roll back one migration
migrate-down:
    goose -dir migrations postgres "${DATABASE_URL}" down

# Reset all migrations
migrate-reset:
    goose -dir migrations postgres "${DATABASE_URL}" reset

# Show migration status
migrate-status:
    goose -dir migrations postgres "${DATABASE_URL}" status

# Create a new migration
migrate-create NAME:
    goose -dir migrations create {{NAME}} sql

# Reset database (destroy volume and recreate)
db-reset:
    docker compose down -v
    docker compose up -d postgres
    @echo "Waiting for postgres..."
    @sleep 3
    just migrate

# --- Build ---

# Build production Docker image
build:
    docker build -t retro-triage -f Dockerfile .

# Build Dokku Docker image
dokku-build:
    docker build -t retro-triage-dokku -f Dockerfile.dokku .

# Clean build artifacts
clean:
    rm -rf bin/ tmp/ coverage.out docs/
    rm -rf .next build node_modules/.cache test-results playwright-report

# --- Dokku Deployment ---

_dokku-host:
    @echo "${DOKKU_HOST}"

# Deploy to Dokku
dokku-deploy:
    git push dokku main

# View Dokku app logs
dokku-logs:
    ssh dokku@$(just _dokku-host) logs retro-triage -t

# View Dokku app config
dokku-config:
    ssh dokku@$(just _dokku-host) config:show retro-triage

# View Dokku app process status
dokku-ps:
    ssh dokku@$(just _dokku-host) ps:report retro-triage

# Connect to Dokku database
dokku-db-connect:
    ssh dokku@$(just _dokku-host) postgres:connect retro-triage-db

# Backup Dokku database locally
dokku-db-backup:
    ssh dokku@$(just _dokku-host) postgres:export retro-triage-db > retro-triage-db-backup.sql

# --- Monitoring ---

# Start monitoring stack
monitoring-up:
    docker compose -f monitoring/docker-compose.yml up -d

# Stop monitoring stack
monitoring-down:
    docker compose -f monitoring/docker-compose.yml down

# View monitoring logs
monitoring-logs:
    docker compose -f monitoring/docker-compose.yml logs -f

# Restart monitoring stack
monitoring-restart:
    docker compose -f monitoring/docker-compose.yml restart
