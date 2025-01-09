# Go parameters
BINARY_NAME=cleric
MAIN_PACKAGE=cmd/cleric/*.go
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

bundle:
	fyne bundle -o internal/ui/bundled.go -pkg ui assets/Icon.png

package:
	fyne package -sourceDir cmd/cleric

.PHONY: build run clean test test-coverage lint
