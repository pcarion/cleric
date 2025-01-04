# Go parameters
BINARY_NAME=cleric
MAIN_PACKAGE=cmd/cleric/main.go
BUILD_DIR=build

# Build commands
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

run:
	go run $(MAIN_PACKAGE)

clean:
	go clean
	rm -rf $(BUILD_DIR)

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run

.PHONY: build run clean test test-coverage lint
