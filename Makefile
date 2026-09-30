.PHONY: all build build-local crawl publish validate fmt lint vuln check clean serve help

# Default target
all: help

## Build targets
build: ## Build static registry catalog and manifests
	go run ./scripts/build

build-local: ## Build static registry with local download links for offline preview
	LOCAL_DOWNLOADS=1 go run ./scripts/build

crawl: ## Crawl and validate decentralized community bundles
	go run ./scripts/crawl

publish: ## Publish release assets to GitHub Releases (requires gh CLI)
	go run ./scripts/publish

## Verification & Quality
validate: ## Validate all official core bundles with okf CLI (--strict --drift)
	@for d in bundles/*; do \
		if [ -d "$$d" ]; then \
			okf validate "$$d" --strict --drift --stale || exit 1; \
		fi; \
	done

fmt: ## Format Go code
	go fmt ./...

lint: ## Run golangci-lint
	golangci-lint run

vuln: ## Run Go vulnerability check
	govulncheck ./...

check: fmt lint vuln validate ## Run all checks (fmt, lint, vuln, validate)

## Utilities
serve: build-local ## Build locally and serve public/ at http://localhost:8080
	@echo "Serving static registry at http://localhost:8080..."
	python3 -m http.server 8080 -d public

clean: ## Remove generated public directory
	rm -rf public

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
