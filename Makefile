COMMIT         := $(shell git describe --dirty --long --always)
VERSION        := $(shell cat $(CURDIR)/VERSION)-$(COMMIT)
LDFLAGS_COMMON := -X main.version=$(VERSION)

.PHONY: build test test-unit test-functional test-all clean bench

build:
	CGO_ENABLED=0 go build -a -ldflags="$(LDFLAGS_COMMON) -s -w -extldflags=-static" -trimpath -o $(CURDIR)/dist/echo main.go

# Run unit tests only
test-unit:
	go test -v -run "^Test.*" -bench= .

# Run functional tests only  
test-functional:
	go test -v -run "^TestFunctional.*" .

# Run all tests
test-all: test-unit test-functional

# Default test target (all tests)
test: test-all

# Run benchmarks
bench:
	go test -bench=. -benchmem .

# Test with coverage
test-coverage:
	go test -coverprofile=coverage.out .
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

# Clean build artifacts
clean:
	rm -rf $(CURDIR)/dist/
	rm -f coverage.out coverage.html

# Cross-compile for multiple platforms
build-all: clean
	mkdir -p $(CURDIR)/dist
	# Linux AMD64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="$(LDFLAGS_COMMON) -s -w" -trimpath -o $(CURDIR)/dist/echo-linux-amd64 main.go
	# Linux ARM64  
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags="$(LDFLAGS_COMMON) -s -w" -trimpath -o $(CURDIR)/dist/echo-linux-arm64 main.go
	
# Show binary sizes
size:
	@echo "Binary sizes:"
	@ls -lh $(CURDIR)/dist/