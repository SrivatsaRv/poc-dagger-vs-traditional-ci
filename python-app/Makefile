# Variables
TAG ?= 1.1.0  # Default semantic version, can be overridden via command line
REGISTRY := docker.io
USERNAME := srivatsarv21
REPOSITORY := docker-dagger-poc
IMAGE_NAME := $(REGISTRY)/$(USERNAME)/$(REPOSITORY)

CONTAINER_NAME := python-app-container
DEPLOYMENT_FILE := infrastructure/deployment.yml

# Install dependencies
requirements:
	pip3 install -r requirements.txt

# Run tests
test:
	pytest test-app.py

# Run linting with flake8
lint:
	flake8 app.py

# Docker login
docker-login:
	@if [ -z "$(DOCKERHUB_USERNAME)" ] || [ -z "$(DOCKERHUB_TOKEN)" ]; then \
		echo "Error: DOCKERHUB_USERNAME or DOCKERHUB_TOKEN is not set"; \
		exit 1; \
	fi
	echo $(DOCKERHUB_TOKEN) | docker login -u $(DOCKERHUB_USERNAME) --password-stdin

# Build Docker image with semver tag
docker-build:
	docker build -t $(IMAGE_NAME):$(TAG) .

# Push Docker image to registry
docker-push:
	docker push $(IMAGE_NAME):$(TAG)

# Update version in deployment.yml
update-deployment:
	@if ! command -v yq &>/dev/null; then \
		echo "Error: yq is not installed. Please install it before running this target."; \
		exit 1; \
	fi
	yq -i '.spec.template.spec.containers[0].image = "$(IMAGE_NAME):$(TAG)"' $(DEPLOYMENT_FILE)

# Chain: Build, push, and update deployment
ci-pipeline: docker-login docker-build docker-push update-deployment

# Clean up Docker images and containers
clean:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true
	docker rmi $(IMAGE_NAME):$(TAG) || true
