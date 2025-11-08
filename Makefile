include .env



APP_NAME=connect
MAIN_PATH=cmd/api/server.go
BINARY_PATH=bin/$(APP_NAME).exe

ifndef DATABASE_URL
$(warning WARNING: DATABASE_URL not set. Migration commands will fail. Set it in your environment or .env file)
endif

postgres:
	@echo "Starting PostgreSQL Docker container..."
	@if [ -z "$(POSTGRES_PASSWORD)" ]; then \
		echo "ERROR: POSTGRES_PASSWORD environment variable must be set"; \
		echo "Example: export POSTGRES_PASSWORD='your-secure-password'"; \
		exit 1; \
	fi
	docker run --name $(APP_NAME)-postgres -p 5432:5432 \
	-e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(APP_NAME) -d postgres:15-alpine

createdb:
	@echo "Creating database $(APP_NAME)..."
	docker exec -it $(APP_NAME)-postgres createdb --username=postgres --owner=postgres $(APP_NAME)

dropdb:
	@echo "Dropping database $(APP_NAME)..."
	docker exec -it $(APP_NAME)-postgres dropdb $(APP_NAME)

db-up:
	@echo "Starting PostgreSQL container..."
	docker start $(APP_NAME)-postgres || make postgres

db-down:
	@echo "Stopping PostgreSQL container..."
	-docker stop $(APP_NAME)-postgres
	-docker rm $(APP_NAME)-postgres

migrate-up:
	@echo "Running all migrations..."
	migrate -path migration -database "$(DATABASE_URL)" -verbose up

migrate-up1:
	@echo "Running one migration up..."
	migrate -path migration -database "$(DATABASE_URL)" -verbose up 1

migrate-down:
	@echo "Rolling back all migrations..."
	migrate -path migration -database "$(DATABASE_URL)" -verbose down

migrate-down1:
	@echo "Rolling back one migration..."
	migrate -path migration -database "$(DATABASE_URL)" -verbose down 1

migrate-create:
	@echo "Creating new migration: $(NAME)"
	migrate create -ext sql -dir migration -seq $(NAME)

migrate-to:
	@echo "Migrating to version $(VERSION)..."
	migrate -path migration -database "$(DATABASE_URL)" goto $(VERSION)

sqlc:
	@echo "Generating SQLC code..."
	sqlc generate

mock:
	@echo "Generating mock database interfaces..."
	mockgen -package mockdb -destination db/mock/store.go github.com/kelvinmhacwilson/connect/db/sqlc Store

proto:
	@echo "Generating protobuf files..."
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	       --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	       proto/*.proto

build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

run: build
	@echo "Running $(APP_NAME)..."
	@$(BINARY_PATH)

server:
	@echo "Starting $(APP_NAME) API server..."
	go run $(MAIN_PATH)

run-dev:
	@echo "Starting $(APP_NAME) in dev mode with hot reload..."
	@export PATH="$$PATH:$(shell go env GOPATH)/bin" && air -c .air.toml || (echo "Air not installed. Run 'make install-tools' or 'go install github.com/cosmtrek/air@v1.49.0'" && exit 1)

test:
	@echo "Running all integration tests against PostgreSQL database..."
	@echo "Database: $(TEST_DATABASE_URL)"
	go test -v -cover ./test/api/... -timeout 30m

test-is-error-present:
	@echo "Running tests (quiet mode - only errors shown)..."
	@echo "Database: $(TEST_DATABASE_URL)"
	@go test ./test/api/... -timeout 30m || exit 1
	@echo "✓ All tests passed"

test-integration:
	@echo "Running integration tests..."
	go test -v -cover ./test/api/... -timeout 30m

test-coverage:
	@echo "Generating test coverage report..."
	go test -v -coverprofile=coverage.out ./test/api/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit:
	@echo "Running unit tests (if any)..."
	go test -v -cover -short ./...

lint:
	@echo "Running linter..."
	golangci-lint run

fmt:
	@echo "Formatting code..."
	go fmt ./...

tidy:
	@echo "Tidying modules..."
	go mod tidy

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

clean:
	@echo "Cleaning up build artifacts..."
	go clean
	rm -rf bin/
	@echo "Clean complete."

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

install-tools:
	@echo "Installing dev tools..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@v1.49.0

help:
	@echo "Available Makefile Commands:"
	@echo "  make postgres         - Run PostgreSQL container"
	@echo "  make createdb         - Create database"
	@echo "  make dropdb           - Drop database"
	@echo "  make migrate-up       - Apply all migrations"
	@echo "  make migrate-down     - Roll back migrations"
	@echo "  make migrate-create NAME=<name> - Create a new migration"
	@echo "  make sqlc             - Generate SQLC code"
	@echo "  make build            - Build Go binary"
	@echo "  make run              - Build & run binary"
	@echo "  make server           - Run API server directly"
	@echo "  make run-dev          - Run with Air (hot reload)"
	@echo "  make test             - Run all tests (verbose)"
	@echo "  make test-is-error-present - Run tests (quiet, errors only)"
	@echo "  make lint             - Run code linter"
	@echo "  make fmt              - Format code"
	@echo "  make deps             - Install dependencies"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make tidy            - Tidy up Modules"

pull:
	git pull

.PHONY: postgres createdb dropdb migrate-up migrate-down migrate-create \
sqlc mock proto build run server run-dev test test-is-error-present lint fmt tidy deps clean \
db_docs db_schema install-tools help pull
