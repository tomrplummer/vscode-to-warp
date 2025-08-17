# VS Code to Warp Theme Converter Makefile

BINARY_NAME=vscode-to-warp
GO_FILES=$(shell find . -name "*.go" -not -name "*_test.go" -type f)

# Build the binary
build: $(BINARY_NAME)

$(BINARY_NAME): $(GO_FILES)
	go build -o $(BINARY_NAME) .

# Install to /usr/local/bin (requires sudo)
install: build
	sudo cp $(BINARY_NAME) /usr/local/bin/

# Install to user's local bin (no sudo required)
install-user: build
	mkdir -p ~/bin
	cp $(BINARY_NAME) ~/bin/
	@echo "Add ~/bin to your PATH if it's not already there:"
	@echo "export PATH=\"\$$HOME/bin:\$$PATH\""

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)

# Run the application
run: build
	./$(BINARY_NAME)

# Download dependencies
deps:
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  install       - Install to /usr/local/bin (requires sudo)"
	@echo "  install-user  - Install to ~/bin (no sudo required)"
	@echo "  clean         - Remove build artifacts"
	@echo "  run           - Build and run the application"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  help          - Show this help"

# Build for multiple platforms
build-all:
	@echo "Building for all platforms..."
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -o $(BINARY_NAME)-windows-arm64.exe .
	@echo "✅ Built binaries for all platforms (Windows, macOS, Linux)"

.PHONY: build install install-user clean run deps help build-all
