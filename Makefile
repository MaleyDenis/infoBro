.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down frontend-install frontend-build frontend-start run-all

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=infobro
MAIN_PKG=./cmd/server

# Node parameters
NPMCMD=npm
NPMINSTALL=$(NPMCMD) install
NPMBUILD=$(NPMCMD) run build
NPMSTART=$(NPMCMD) start

# Docker parameters
DOCKER_IMAGE=infobro
DOCKER_TAG=latest

all: test build

build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) $(MAIN_PKG)

run: build
	./bin/$(BINARY_NAME)

test:
	$(GOTEST) -v ./...

clean:
	rm -f ./bin/$(BINARY_NAME)
	rm -rf ./vendor
	rm -rf ./web/node_modules
	rm -rf ./web/build

deps:
	$(GOMOD) download

tidy:
	$(GOMOD) tidy

docker-build:
	docker build -t $(DOCKER_IMAGE)-api:$(DOCKER_TAG) .
	docker build -t $(DOCKER_IMAGE)-frontend:$(DOCKER_TAG) ./web

docker-run: docker-build
	docker run -p 8080:8080 $(DOCKER_IMAGE)-api:$(DOCKER_TAG)

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

# Frontend commands
frontend-install:
	cd web && $(NPMINSTALL)

frontend-build: frontend-install
	cd web && $(NPMBUILD)

frontend-start: frontend-install
	cd web && $(NPMSTART)

# Development helpers
dev-env:
	docker-compose up -d mongodb redis

dev-run: build
	./bin/$(BINARY_NAME) --mongo-uri mongodb://localhost:27017 --redis-addr localhost:6379

# Run all components (backend, frontend, mongodb, redis)
run-all:
	./scripts/run-dev.sh

# Helpers for connectors
run-reddit:
	./bin/$(BINARY_NAME) --run-connector reddit

run-all-connectors:
	./bin/$(BINARY_NAME) --run-all-connectors