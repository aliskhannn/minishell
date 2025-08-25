APP_NAME = minishell
BIN = bin/$(APP_NAME)

.PHONY: all build run test clean

all: build

## Build the minishell binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p bin
	@go build -o $(BIN) ./cmd/minishell

## Run minishell
run: build
	@$(BIN)

## Run integration tests
test: build
	@echo "Running integration tests..."
	@go test -v ./...

# Format Go code using goimports
format:
	goimports -local github.com/aliskhannn/minishell -w .

# Run linters: vet + golangci-lint
lint:
	go vet ./... && golangci-lint run ./...

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin
