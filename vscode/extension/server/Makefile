# Kro LSP Server Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run

# Binary name
BINARY_NAME=kro-lsp

# Build information
GIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date +%FT%T%z)
VERSION?=0.1.0

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.BuildHash=$(GIT_HASH)"

# Build directories
BUILD_DIR=./build

.PHONY: all build clean test coverage deps tidy

all: clean build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build the binary in the current directory for direct usage
build-here:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

deps:
	$(GOGET) -u ./...

tidy:
	$(GOMOD) tidy

run:
	$(GORUN) main.go
