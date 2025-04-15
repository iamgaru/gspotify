# Binary name
BINARY_NAME=gspotty

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard cmd/$(BINARY_NAME)/*.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean test help deps fmt vet install uninstall

all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "Done!"

## clean: Clean build files
clean:
	@echo "Cleaning build files..."
	rm -f $(BINARY_NAME)
	go clean
	@echo "Done!"

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Done!"

## run: Build and run the application
run: build
	./$(BINARY_NAME)

## deps: Install dependencies
deps:
	go mod download

## fmt: Format code
fmt:
	go fmt ./...

## vet: Run go vet on the code
vet:
	go vet ./...

## install: Install the application to /usr/local/bin
install: build
	cp $(BINARY_NAME) /usr/local/bin/

## uninstall: Remove the application from /usr/local/bin
uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

## help: Display this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /' 