# Makefile for kodi-renamer Go application

# Application name
APP_NAME := kodi-renamer
BUILD_DIR := build
CMD_DIR := cmd/kodi-renamer

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt

# Docker parameters
DOCKER := docker
DOCKER_REGISTRY ?=
DOCKER_IMAGE := kodi-renamer
DOCKER_TAG := latest
# Local image name (used for build and run)
LOCAL_IMAGE := $(DOCKER_IMAGE):$(DOCKER_TAG)
# Remote image name with registry (used for push)
ifdef DOCKER_REGISTRY
REMOTE_IMAGE := $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
else
REMOTE_IMAGE := $(LOCAL_IMAGE)
endif

# Build variables
BINARY_NAME := $(APP_NAME)
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

# Directory parameters (can be overridden via environment variables or make arguments)
MOVIE_TO_RENAME_DIR ?=
MOVIE_RENAMED_DIR ?=
SERIE_TO_RENAME_DIR ?=
SERIE_RENAMED_DIR ?=

# Default target - build and run
.PHONY: all
all: build run

# Build the application
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BINARY_PATH) ./$(CMD_DIR)
	@echo "Build complete: $(BINARY_PATH)"

# Run the application
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	@if [ ! -f $(BINARY_PATH) ]; then \
		echo "Binary not found. Building first..."; \
		$(MAKE) build; \
	fi
	@$(BINARY_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run dual API test
.PHONY: test-dual-api
test-dual-api:
	@echo "Running dual API test..."
	@./test_dual_api.sh

# Run integration test (requires real API keys)
.PHONY: test-integration
test-integration:
	@echo "Running integration test..."
	@./test_renaming.sh

# Run all tests
.PHONY: test-all
test-all: test test-dual-api
	@echo "All tests complete"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the application
.PHONY: install
install: build
	@echo "Installing $(APP_NAME)..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "Multi-platform build complete"

# Docker build - multi-stage with distroless
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	$(DOCKER) build -f docker/Dockerfile -t $(LOCAL_IMAGE) .
	@echo "Docker build complete: $(LOCAL_IMAGE)"

# Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	$(DOCKER) run --rm -it \
		-e TVDB_API_KEY="$(TVDB_API_KEY)" \
		-e TMDB_API_KEY="$(TMDB_API_KEY)" \
		-e MOVIE_TO_RENAME_DIR="/media/movie-to-rename" \
		-e MOVIE_RENAMED_DIR="/media/movie-renamed" \
		-e SERIE_TO_RENAME_DIR="/media/serie-to-rename" \
		-e SERIE_RENAMED_DIR="/media/serie-renamed" \
		$(if $(MOVIE_TO_RENAME_DIR),-v "$(MOVIE_TO_RENAME_DIR):/media/movie-to-rename",) \
		$(if $(MOVIE_RENAMED_DIR),-v "$(MOVIE_RENAMED_DIR):/media/movie-renamed",) \
		$(if $(SERIE_TO_RENAME_DIR),-v "$(SERIE_TO_RENAME_DIR):/media/serie-to-rename",) \
		$(if $(SERIE_RENAMED_DIR),-v "$(SERIE_RENAMED_DIR):/media/serie-renamed",) \
		$(LOCAL_IMAGE)

# Run Docker container in dry-run mode
.PHONY: docker-dry-run
docker-dry-run:
	@echo "Running Docker container in dry-run mode..."
	$(DOCKER) run --rm -it \
		-e TVDB_API_KEY="$(TVDB_API_KEY)" \
		-e TMDB_API_KEY="$(TMDB_API_KEY)" \
		-e MOVIE_TO_RENAME_DIR="/media/movie-to-rename" \
		-e MOVIE_RENAMED_DIR="/media/movie-renamed" \
		-e SERIE_TO_RENAME_DIR="/media/serie-to-rename" \
		-e SERIE_RENAMED_DIR="/media/serie-renamed" \
		$(if $(MOVIE_TO_RENAME_DIR),-v "$(MOVIE_TO_RENAME_DIR):/media/movie-to-rename",) \
		$(if $(MOVIE_RENAMED_DIR),-v "$(MOVIE_RENAMED_DIR):/media/movie-renamed",) \
		$(if $(SERIE_TO_RENAME_DIR),-v "$(SERIE_TO_RENAME_DIR):/media/serie-to-rename",) \
		$(if $(SERIE_RENAMED_DIR),-v "$(SERIE_RENAMED_DIR):/media/serie-renamed",) \
		$(LOCAL_IMAGE) -dry-run

# Show Docker image info
.PHONY: docker-info
docker-info:
	@echo "Docker image information: "
	@$(DOCKER) images --filter=reference='$(DOCKER_IMAGE):*' --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

# Remove Docker image
.PHONY: docker-clean
docker-clean:
	@echo "Removing Docker image $(LOCAL_IMAGE)..."
	$(DOCKER) rmi $(LOCAL_IMAGE) || true
	@echo "Docker clean complete"

# Tag and push Docker image to registry (requires docker login)
.PHONY: docker-push
docker-push:
ifdef DOCKER_REGISTRY
	@echo "Tagging $(LOCAL_IMAGE) as $(REMOTE_IMAGE)..."
	$(DOCKER) tag $(LOCAL_IMAGE) $(REMOTE_IMAGE)
endif
	@echo "Pushing Docker image $(REMOTE_IMAGE)..."
	$(DOCKER) push $(REMOTE_IMAGE)
	@echo "Docker push complete"

# Build and run Docker in one command
.PHONY: docker
docker: docker-build docker-run

# Run with dry-run mode
.PHONY: dry-run
dry-run: build
	@echo "Running in dry-run mode..."
	@export MOVIE_TO_RENAME_DIR="$(MOVIE_TO_RENAME_DIR)" \
		MOVIE_RENAMED_DIR="$(MOVIE_RENAMED_DIR)" \
		SERIE_TO_RENAME_DIR="$(SERIE_TO_RENAME_DIR)" \
		SERIE_RENAMED_DIR="$(SERIE_RENAMED_DIR)"; \
	$(BINARY_PATH) -dry-run

# Run with auto mode
.PHONY: auto
auto: build
	@echo "Running in auto mode..."
	@export MOVIE_TO_RENAME_DIR="$(MOVIE_TO_RENAME_DIR)" \
		MOVIE_RENAMED_DIR="$(MOVIE_RENAMED_DIR)" \
		SERIE_TO_RENAME_DIR="$(SERIE_TO_RENAME_DIR)" \
		SERIE_RENAMED_DIR="$(SERIE_RENAMED_DIR)"; \
	$(BINARY_PATH) -auto

# Run for movies only
.PHONY: run-movies
run-movies: build
	@echo "Running for movies only..."
	@if [ -z "$(MOVIE_TO_RENAME_DIR)" ]; then \
		echo "Error: MOVIE_TO_RENAME_DIR is required"; \
		echo "Usage: make run-movies MOVIE_TO_RENAME_DIR=./path/to/movies [MOVIE_RENAMED_DIR=./path/to/renamed]"; \
		exit 1; \
	fi
	@export MOVIE_TO_RENAME_DIR="$(MOVIE_TO_RENAME_DIR)" \
		MOVIE_RENAMED_DIR="$(MOVIE_RENAMED_DIR)"; \
	$(BINARY_PATH)

# Run for series only
.PHONY: run-series
run-series: build
	@echo "Running for series only..."
	@if [ -z "$(SERIE_TO_RENAME_DIR)" ]; then \
		echo "Error: SERIE_TO_RENAME_DIR is required"; \
		echo "Usage: make run-series SERIE_TO_RENAME_DIR=./path/to/series [SERIE_RENAMED_DIR=./path/to/renamed]"; \
		exit 1; \
	fi
	@export SERIE_TO_RENAME_DIR="$(SERIE_TO_RENAME_DIR)" \
		SERIE_RENAMED_DIR="$(SERIE_RENAMED_DIR)"; \
	$(BINARY_PATH)


# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  all            - Build and run the application (default)"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  clean          - Remove build artifacts"
	@echo "  build-all      - Build for multiple platforms"
	@echo ""
	@echo "Development targets:"
	@echo "  test              - Run Go unit tests"
	@echo "  test-dual-api     - Run dual API configuration tests"
	@echo "  test-integration  - Run integration tests (requires real API keys)"
	@echo "  test-all          - Run unit and dual API tests"
	@echo "  fmt               - Format source code"
	@echo "  deps              - Download and tidy dependencies"
	@echo "  install           - Install binary to GOPATH/bin"
	@echo ""
	@echo "Run modes:"
	@echo "  dry-run        - Run in dry-run mode"
	@echo "  auto           - Run in automatic mode"
	@echo "  run-movies     - Run for movies only"
	@echo "  run-series     - Run for series only"
	@echo ""
	@echo "Directory parameters:"
	@echo "  MOVIE_TO_RENAME_DIR - Directory containing movies to rename"
	@echo "  MOVIE_RENAMED_DIR   - Directory for renamed movies (optional)"
	@echo "  SERIE_TO_RENAME_DIR - Directory containing series to rename"
	@echo "  SERIE_RENAMED_DIR   - Directory for renamed series (optional)"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-dry-run - Run Docker in dry-run mode"
	@echo "  docker-info    - Show Docker image information"
	@echo "  docker-clean   - Remove Docker image"
	@echo "  docker-push    - Push Docker image to registry"
	@echo "  docker         - Build and run Docker in one command"
	@echo ""
	@echo "Docker variables:"
	@echo "  DOCKER_REGISTRY - Custom registry (e.g., ghcr.io/username, docker.io/username)"
	@echo "  DOCKER_TAG      - Image tag (default: latest)"
	@echo ""
	@echo "Examples:"
	@echo "  make docker-build DOCKER_TAG=v1.0.0"
	@echo "  make docker-push DOCKER_REGISTRY=ghcr.io/myuser DOCKER_TAG=v1.0.0"
	@echo "  make docker-push DOCKER_REGISTRY=docker.io/myuser"
	@echo "  make run-movies MOVIE_TO_RENAME_DIR=./media/movies MOVIE_RENAMED_DIR=./renamed/movies"
	@echo "  make run-series SERIE_TO_RENAME_DIR=./media/series SERIE_RENAMED_DIR=./renamed/series"
	@echo "  make dry-run MOVIE_TO_RENAME_DIR=./media/movies SERIE_TO_RENAME_DIR=./media/series"
	@echo "  make docker-run MOVIE_TO_RENAME_DIR=./media/movies SERIE_TO_RENAME_DIR=./media/series"
	@echo ""
	@echo "Help:"
	@echo "  help           - Show this help message"
