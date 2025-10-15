.PHONY: build test clean fmt lint install check

# Build the bento binary
build:
	go build -o bin/bento ./cmd/bento

# Run all tests with race detector
test:
	@echo "Testing all packages..."
	@go test -v -race bento/pkg/neta
	@go test -v -race bento/pkg/itamae
	@go test -v -race bento/pkg/pantry

# Clean build artifacts
clean:
	rm -rf bin/

# Format all Go files
fmt:
	@gofmt -s -w pkg/ cmd/

# Run linter (if golangci-lint is installed)
lint:
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping lint"; \
	fi

# Install to GOPATH
install:
	go install ./cmd/bento

# Run all quality checks (Karen's requirements)
check: fmt test build
	@echo "All quality checks passed!"
