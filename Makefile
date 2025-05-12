DOCKER_COMPOSE = docker-compose -f deploy/docker-compose.yml
APP_SERVICE = app
PROJECT_DIR = $(PWD)

.PHONY: help up down build generate dev test clean logs shell

help: ## This help dialog
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

run-local-server: ## Run the app locally
	air -c .air.toml

run-ui-desktop: ## Run the app locally
	go run ./cmd/kybermed_UI/ui.go

run-ui-mobile: ## Run the app locally
	go run -tags mobile ./cmd/kybermed_UI/ui.go

requirements: ## Generate go.mod & go.sum files
	go mod tidy

clean-packages: ## Clean packages
	go clean -modcache

up:
	$(DOCKER_COMPOSE) up $(APP_SERVICE) --build 

down: ## Stop all running containers
	$(DOCKER_COMPOSE) down

build: ## Build the production Docker image
	$(DOCKER_COMPOSE) build --no-cache

test:
	go test -v ./...

logs: ## View logs from the app container
	$(DOCKER_COMPOSE) logs -f $(APP_SERVICE)

shell: ## Open an interactive shell in the app container
	docker exec -it $(APP_SERVICE) /bin/sh

clean: ## Clean up Docker volumes and containers
	$(DOCKER_COMPOSE) down -v
