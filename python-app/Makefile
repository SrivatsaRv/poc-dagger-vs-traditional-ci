# Variables
TAG ?= 1.0.0  # Default semantic version, can be overridden via command line
REGISTRY := docker.io
IMAGE_NAME := python-app
CONTAINER_NAME := python-app-container

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
	docker build -t $(REGISTRY)/$(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG) .

# Push Docker image to registry
docker-push:
	docker push $(REGISTRY)/$(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)

# Run Flask app using the last pushed Docker image
run:
	docker run --rm -d -p 5000:5000 --name $(CONTAINER_NAME) $(REGISTRY)/$(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)

# Stop and remove the running container (cleanup)
stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)
