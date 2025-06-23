CMD = ./cmd/web/main.go 
OUT_DIR = ./build

.PHONY: help up down build generate dev test clean logs shell

help: ## This help dialog
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

run-local-server: ## Run the app locally
	air -c .air.toml

requirements: ## Generate go.mod & go.sum files
	go mod tidy

clean-packages: ## Clean packages
	go clean -modcache

build-production:
	go build -o $(OUT_DIR)/bookings $(CMD) -production=true -dbname=bookingsdb -dbuser=nelson -dbpassword=nelson9199 -dbhost=localhost -dbport=5432 -dbssl=disable

build-development:
	go build -o $(OUT_DIR)/bookings-dev $(CMD) -production=false -dbname=bookingsdb -dbuser=nelson -dbpassword=nelson9199 -dbhost=localhost -dbport=5432 -dbssl=disable

test:
	go test -v ./...
