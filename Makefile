# Base Makefile — default goal prints help (supports inline "##" and preceding "##" styles)
.DEFAULT_GOAL := default

.PHONY: default help build test fmt lint run setup clean fuzz coverage

default: help

help: ## Show this help.
	@echo "Available targets:"
	@awk 'BEGIN{ORS="";} \
	     /^## / {desc=substr($$0,4); next} \
	     /^[a-zA-Z0-9._-]+:/ { t=$$1; sub(/:$$/,"",t); d=""; \
	         if(match($$0,/## /)) { sub(/.*## /,"",$$0); d=$$0 } \
	         else if(desc) { d=desc; desc="" } \
	         if(d) print t": "d"\n" }' $(MAKEFILE_LIST) \
	  | sort -u \
	  | awk -F': ' '{printf "  \033[1;36m%-12s\033[0m %s\n", $$1, $$2}'

# -----------------------
# Example targets (annotate with "##" either inline or on the previous line)
# -----------------------
#

setup: ## Setup environment
	@echo "==> setup"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ${PWD}/bin v2.5.0
	curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b ${PWD}/bin

fmt: ## Format code
	@echo "==> format"
	gofmt -l -w .

lint: ## Lint code
	@echo "==> lint"
	${PWD}/bin/golangci-lint run

fuzz: ## Fuzz code
	@echo "==> fuzz"
	go test ./pkg/core -fuzz=FuzzParser -fuzztime=60s

coverage: ## Test coverage
	@echo "==> coverage"
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out
	rm cover.out

test: ## Run tests
	@echo "==> test"
	go test ./...

run: build ## Build and run
	@echo "==> run"
	mkdir -p build
	rm -f build/* 2> /dev/null
	./bin/air --build.cmd 'go build -o build cmd/*' --build.bin "./build/main"

clean: ## Remove build artifacts
	@echo "==> clean"

