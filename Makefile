# Project variables
GOCMD ?= go
PKG ?= ./...
MAIN_PKG ?= cmd/main.go
BINARY_NAME ?= pharmacy-api
BIN_DIR ?= bin
BIN_PATH := $(BIN_DIR)/$(BINARY_NAME)

# Docker/Compose
DOCKER ?= docker
DOCKERFILE ?= Dockerfile
IMAGE_NAME ?= pharmacy-api
VERSION ?= latest
IMAGE_TAG ?= $(VERSION)
IMAGE := $(IMAGE_NAME):$(IMAGE_TAG)
COMPOSE ?= docker compose
COMPOSE_FILE ?= docker-compose.yaml
SERVICE ?= pharmacy-api
PORT ?= 8080

# Tools
GOTEST := $(GOCMD) test
GOVET := $(GOCMD) vet

# Default help
.PHONY: help
help: ## Show this help
	@echo "Usage: make <target>"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Targets:\n"} /^[a-zA-Z0-9_.-]+:.*##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Run
.PHONY: server run
server run: ## Run the application (go run cmd/main.go)
	$(GOCMD) run $(MAIN_PKG)

# Build
.PHONY: build
build: ## Build local binary for current OS/ARCH into bin/
	@mkdir -p $(BIN_DIR)
	$(GOCMD) build -ldflags "-s -w" -o $(BIN_PATH) $(MAIN_PKG)

.PHONY: build-linux
build-linux: ## Build linux/amd64 binary into bin/
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GOCMD) build -ldflags "-s -w" -o $(BIN_PATH)-linux-amd64 $(MAIN_PKG)

# Test & Coverage
.PHONY: test
test: ## Run tests with verbose output
	$(GOTEST) -v $(PKG)

.PHONY: cover
cover: ## Run tests with coverage to coverage.out
	$(GOTEST) -coverprofile=coverage.out $(PKG)

.PHONY: cover-html
cover-html: cover ## Generate coverage HTML report coverage.html
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser."

# Quality
.PHONY: fmt
fmt: ## Format code
	$(GOCMD) fmt $(PKG)

.PHONY: fmt-check
fmt-check: ## Check formatting (fails if files need formatting)
	@files=$$(gofmt -l .); if [ -n "$$files" ]; then echo "Unformatted files:"; echo "$$files"; exit 1; else echo "All files formatted"; fi

.PHONY: vet
vet: ## Run go vet
	$(GOVET) $(PKG)

.PHONY: tidy
tidy: ## Run go mod tidy
	$(GOCMD) mod tidy

.PHONY: mod-download
mod-download: ## Download Go modules
	$(GOCMD) mod download

.PHONY: mod-verify
mod-verify: ## Verify dependencies
	$(GOCMD) mod verify

.PHONY: lint
lint: fmt-check vet ## Basic lint (fmt-check + vet)
	@echo "Lint passed"

# Clean
.PHONY: clean
clean: ## Clean build artifacts and coverage files
	rm -rf $(BIN_DIR) coverage.out coverage.html

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image (IMAGE_NAME, IMAGE_TAG overrideable)
	$(DOCKER) build -f $(DOCKERFILE) -t $(IMAGE) .

.PHONY: docker-run
docker-run: ## Run Docker image (maps PORT)
	$(DOCKER) run --rm -d --name $(IMAGE_NAME) -p $(PORT):$(PORT) -e PORT=$(PORT) $(IMAGE)

.PHONY: docker-stop
docker-stop: ## Stop running Docker container
	-$(DOCKER) stop $(IMAGE_NAME) || true

.PHONY: docker-rm
docker-rm: ## Remove stopped Docker container
	-$(DOCKER) rm $(IMAGE_NAME) || true

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	$(DOCKER) push $(IMAGE)

# Docker Compose targets
.PHONY: compose-up
compose-up: ## Start services with Docker Compose (detached)
	$(COMPOSE) -f $(COMPOSE_FILE) up -d --build

.PHONY: compose-down
compose-down: ## Stop and remove services with Docker Compose
	$(COMPOSE) -f $(COMPOSE_FILE) down

.PHONY: compose-logs
compose-logs: ## Tail logs for the service
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f $(SERVICE)

.PHONY: compose-restart
compose-restart: ## Restart the service
	$(COMPOSE) -f $(COMPOSE_FILE) restart $(SERVICE)