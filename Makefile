.PHONY: all build dev run test clean mock

# Variables
APP_NAME=flight-api
MAIN_PATH=cmd/http/main.go
BIN_DIR=bin

all: build

# Build binary aplikasi
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)

dev:
	@echo "Running $(APP_NAME) locally with .env..."
	@export $$(cat .env | xargs) && go run $(MAIN_PATH)

run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BIN_DIR)/

mock:
	@echo "Generating mocks..."
	@go run github.com/vektra/mockery/v2@latest --dir internal --all --output ./tests/mock --outpkg mock