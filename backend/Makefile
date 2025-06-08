.PHONY: build run clean test

# Binary name
BINARY_NAME=task-management-api

# Go build settings
GO=go
GOFLAGS=-v
BUILD_DIR=./bin

# Environment
ENV_FILE=.env

# Build the application
build:
	${GO} build ${GOFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/api

# Run the application
run:
	${GO} run ./cmd/api/main.go

# Clean build artifacts
clean:
	${GO} clean
	rm -rf ${BUILD_DIR}

# Run tests
test:
	${GO} test -v ./...

# Install dependencies
deps:
	${GO} mod tidy

# Build and run the application
dev: deps run

# Show help
help:
	@echo "make build - Build the application"
	@echo "make run - Run the application without building"
	@echo "make clean - Clean build artifacts"
	@echo "make test - Run tests"
	@echo "make deps - Install dependencies"
	@echo "make dev - Run the application in development mode"

# Default target
default: build 