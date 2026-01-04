include .env
export

.PHONY: build run test test-v clean fmt vet lint help migrate-up migrate-down migrate-status migrate-create db-create db-drop test-db-create test-db-drop seed types

build:
	go build -o bin/server .

run:
	go run .

test: test-db-create
	DATABASE_URL=$(DATABASE_URL_TEST) go test ./...

test-v: test-db-create
	DATABASE_URL=$(DATABASE_URL_TEST) go test -v ./...

test-pkg:
	go test -v ./$(PKG)/...

cover:
	go test -cover ./...

cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/ coverage.out coverage.html

deps:
	go mod download

dev: build
	./bin/server

check: fmt vet test

db-create:
	createdb hungr_db

db-drop:
	dropdb hungr_db

test-db-create:
	dropdb --if-exists hungr_db_test
	createdb hungr_db_test
	goose -dir migrations postgres "$(DATABASE_URL_TEST)" up

test-db-drop:
	dropdb --if-exists hungr_db_test

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DATABASE_URL)" status

migrate-create:
	goose -dir migrations create $(NAME) sql

seed:
	go run ./cmd/seed

types:
	tygo generate

help:
	@echo "build      Build the binary"
	@echo "run        Run the server"
	@echo "test       Run all tests"
	@echo "test-v     Verbose tests"
	@echo "test-pkg   Test a package (PKG=storage)"
	@echo "cover      Tests with coverage"
	@echo "cover-html HTML coverage report"
	@echo "fmt        Format code"
	@echo "vet        Run go vet"
	@echo "tidy       Tidy dependencies"
	@echo "clean      Clean build artifacts"
	@echo "deps       Download dependencies"
	@echo "dev        Build and run"
	@echo "check      fmt + vet + test"
	@echo "db-create       Create hungr_db database"
	@echo "db-drop         Drop hungr_db database"
	@echo "migrate-up     Run migrations"
	@echo "migrate-down   Rollback one migration"
	@echo "migrate-status Show migration status"
	@echo "migrate-create Create migration (NAME=name)"
	@echo "seed           Populate DB with test data"
	@echo "types          Generate TypeScript types from Go"
