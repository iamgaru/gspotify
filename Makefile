# Binary name
BINARY_NAME=gspotty

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard cmd/$(BINARY_NAME)/*.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean test help

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

## help: Display this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /' 