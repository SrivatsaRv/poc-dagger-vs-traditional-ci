# Variables
IMAGE_NAME=dagger-poc
VERSION=1.0.0

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v ./...


## Run when app is up
.PHONY: test-e2e
test-e2e:
	go test -v -tags e2e -count 1 ./tests/e2e

.PHONY: build
build:
	CGO_ENABLED=0 go build -o server ./cmd/server

.PHONY: docker-build
docker-build:
	docker build -t ttl.sh/$(IMAGE_NAME):$(VERSION) .

.PHONY: docker-push
docker-push:
	docker push ttl.sh/$(IMAGE_NAME):$(VERSION)
	@echo "Image ttl.sh/$(IMAGE_NAME):$(VERSION) is available for 1 hour"

.PHONY: up
up:
	docker compose up --build -d
