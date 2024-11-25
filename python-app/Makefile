# Variables
TAG ?= 1.0.0  # Default semantic version, can be overridden via command line
REGISTRY := ttl.sh
IMAGE_NAME := python-app
CONTAINER_NAME := python-app-container

# Install dependencies
requirements:
	pip3 install -r python-app/requirements.txt

# Run tests
test:
	pytest python-app/test-app.py

# Run linting with flake8
lint:
	flake8 python-app/app.py

# Build Docker image with semver tag
docker-build:
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(TAG) python-app/

# Push Docker image to registry
docker-push:
	docker push $(REGISTRY)/$(IMAGE_NAME):$(TAG)

# Run Flask app using the last pushed Docker image
run:
	docker run --rm -d -p 5000:5000 --name $(CONTAINER_NAME) $(REGISTRY)/$(IMAGE_NAME):$(TAG)

# Stop and remove the running container (cleanup)
stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)
