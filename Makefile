# Define variables
IMAGE_NAME = quantum-isp365-webapp
DOCKER_COMPOSE = docker-compose
DOCKERFILE = Dockerfile
GO_VERSION = 1.23
BUILD_DIR = build

# Build Docker Image
build: generate-dockerfile
	@echo "Building the Docker image..."
	@docker build --no-cache -t $(IMAGE_NAME) -f $(DOCKERFILE) .

# Clean the build directory (if needed)
clean:
	@echo "Cleaning the build directory..."
	@rm -rf $(BUILD_DIR)

# Run the application with Docker Compose
up:
	@echo "Starting the application with Docker Compose..."
	@$(DOCKER_COMPOSE) up

# Pull Go dependencies
get-deps:
	@echo "Fetching Go dependencies..."
	@go mod tidy
	@go mod download

# Run tests (if any)
test:
	@echo "Running tests..."
	@go test ./...

# Clean and rebuild everything
rebuild: clean build

# Build with no cache
build-no-cache:
	@echo "Building with no cache..."
	@docker build --no-cache -t $(IMAGE_NAME) -f $(DOCKERFILE) .

# Down (stop and remove containers)
down:
	@echo "Stopping containers..."
	@$(DOCKER_COMPOSE) down

.PHONY: generate-dockerfile build clean go-build get-deps up test rebuild build-no-cache down
