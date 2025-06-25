COMMIT         := $(shell git describe --dirty --long --always)
VERSION        := $(shell cat $(CURDIR)/VERSION)-$(COMMIT)
LDFLAGS_COMMON   := -X main.version=$(VERSION)


build:
	CGO_ENABLED=0 go build -a -ldflags="$(LDFLAGS_COMMON) -s -w -extldflags=-static" -trimpath -o $(CURDIR)/dist/echo main.go

info:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"

