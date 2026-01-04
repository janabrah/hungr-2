.PHONY: build run test test-v clean fmt vet lint help db-up db-down db-logs

build:
	go build -o bin/server .

run:
	go run .

test:
	go test ./...

test-v:
	go test -v ./...

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

db-up:
	docker compose up -d

db-down:
	docker compose down

db-logs:
	docker compose logs -f

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
	@echo "db-up      Start local postgres"
	@echo "db-down    Stop local postgres"
	@echo "db-logs    Tail postgres logs"
