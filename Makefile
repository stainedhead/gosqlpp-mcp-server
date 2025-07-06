.PHONY: build test lint clean run docker-build docker-run deploy-dev deploy-prod help

# Variables
BINARY_NAME=gosqlpp-mcp-server
DOCKER_IMAGE=gosqlpp-mcp-server
GO_VERSION=1.23

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and generate coverage report
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

run-http: build ## Build and run the application in HTTP mode
	@echo "Running $(BINARY_NAME) in HTTP mode..."
	./$(BINARY_NAME) --transport http --port 8080

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -f deployment/docker/Dockerfile -t $(DOCKER_IMAGE) .

docker-run: docker-build ## Build and run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(DOCKER_IMAGE) --transport http --host 0.0.0.0

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# CDK deployment targets
cdk-bootstrap: ## Bootstrap CDK (run once per account/region)
	@echo "Bootstrapping CDK..."
	cd deployment/cdk && cdk bootstrap

deploy-dev: ## Deploy to development environment
	@echo "Deploying to development..."
	cd deployment/cdk && ENVIRONMENT=development cdk deploy

deploy-prod: ## Deploy to production environment
	@echo "Deploying to production..."
	cd deployment/cdk && ENVIRONMENT=production cdk deploy

destroy-dev: ## Destroy development environment
	@echo "Destroying development environment..."
	cd deployment/cdk && ENVIRONMENT=development cdk destroy

destroy-prod: ## Destroy production environment
	@echo "Destroying production environment..."
	cd deployment/cdk && ENVIRONMENT=production cdk destroy

# Development helpers
dev-setup: deps install-tools ## Set up development environment
	@echo "Development environment setup complete!"

check: lint test ## Run all checks (lint and test)
	@echo "All checks passed!"

# Docker helpers
docker-shell: ## Run shell in Docker container
	docker run -it --entrypoint /bin/sh $(DOCKER_IMAGE)

docker-logs: ## Show Docker container logs
	docker logs $(shell docker ps -q --filter ancestor=$(DOCKER_IMAGE))

# AWS helpers
aws-login: ## Login to AWS ECR
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(shell aws sts get-caller-identity --query Account --output text).dkr.ecr.us-east-1.amazonaws.com

# Release helpers
version: ## Show current version
	@echo "Version: 1.0.0"
	@echo "Go version: $(GO_VERSION)"
	@echo "Git commit: $(shell git rev-parse --short HEAD)"

# Quick development cycle
dev: clean build test lint ## Full development cycle (clean, build, test, lint)
	@echo "Development cycle complete!"
