# Base Makefile — default goal prints help (supports inline "##" and preceding "##" styles)
.DEFAULT_GOAL := default

.PHONY: default help build test fmt lint run setup clean test-fuzz coverage-int coverage-unit coverage-integration-html coverage-integration-total coverage-unit-html coverage-unit-total build-docker run-docker test-bench test-race coverage-total coverage-html coverage monitoring-up monitoring-down

VERSION := $(shell git describe --tags 2>/dev/null || echo dev)

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
	  | awk -F': ' '{printf "  \033[1;36m%-15s\033[0m %s\n", $$1, $$2}'

# -----------------------
# Example targets (annotate with "##" either inline or on the previous line)
# -----------------------
#

setup: ## Setup environment
	@echo "==> setup"
	@[ -x ${PWD}/bin/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ${PWD}/bin v2.5.0
	@[ -x ${PWD}/bin/air ] || curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b ${PWD}/bin

fmt: ## Format code
	@echo "==> format"
	gofmt -l -w .

lint: setup ## Lint code
	@echo "==> lint"
	${PWD}/bin/golangci-lint run

test-fuzz: ## Fuzz code
	@echo "==> test-fuzz"
	go test ./pkg/core -fuzz=FuzzParser -fuzztime=60s
	go test ./pkg/twitch/ -fuzz=FuzzParseChat -fuzztime=60s

test-pprof:
	@echo "==> test-pprof"
	go test -bench=. -cpuprofile=cpu.prof -benchmem ./integration/
	go tool pprof -top -nodecount=10 cpu.prof
	rm cpu.prof integration.test

coverage-unit-total:
	@echo "==> coverage-unit-total"
	@go list ./... | grep -v integration | xargs go test -coverprofile=cover.out -covermode=count > /dev/null
	@go tool cover -func=cover.out -o=cover.out
	@tail -n1 cover.out | tr -s '\t'
	@rm cover.out

coverage-unit-html:
	@echo "==> coverage-unit-html"
	@go list ./... | grep -v integration | xargs go test -coverprofile=cover.out
	@go tool cover -html=cover.out
	@rm cover.out

coverage-unit: coverage-unit-total coverage-unit-html ## Coverage of unit tests
	@echo "==> coverage-unit"

coverage-integration-total:
	@echo "==> coverage-integration-total"
	@go test ./integration -coverprofile=cover.out -coverpkg=./pkg/hub,./pkg/room,./pkg/state,./pkg/core -covermode=count > /dev/null
	@go tool cover -func=cover.out -o=cover.out
	@tail -n1 cover.out | tr -s '\t'
	@rm cover.out

coverage-integration-html:
	@echo "==> coverage-integration-html"
	@go test ./integration -coverprofile=cover.out -coverpkg=./pkg/hub,./pkg/room,./pkg/state,./pkg/core
	@go tool cover -html=cover.out
	@rm cover.out

coverage-int: coverage-integration-total coverage-integration-html ## Coverage of integration tests
	@echo "==> coverage-integration"

coverage-html:
	@echo "==> coverage-html"
	@go test ./integration/ -coverprofile=integration.tmp.out -coverpkg=./pkg/hub,./pkg/room,./pkg/state,./pkg/core > /dev/null
	@go list ./... | grep -v integration | xargs go test -coverprofile=main.tmp.out > /dev/null
	@echo "mode: set" > cover.out
	@grep -h -v mode *.tmp.out >> cover.out
	@go tool cover -html=cover.out
	@rm *.out

coverage-total:
	@echo "==> coverage-total"
	@go test ./integration/ -coverprofile=integration.tmp.out -coverpkg=./pkg/hub,./pkg/room,./pkg/state,./pkg/core -covermode=count > /dev/null
	@go list ./... | grep -v integration | xargs go test -coverprofile=main.tmp.out -covermode=count > /dev/null
	@echo "mode: count" > cover.out
	@grep -h -v mode *.tmp.out >> cover.out
	@go tool cover -func=cover.out -o=cover.out
	@tail -n1 cover.out | tr -s '\t'
	@rm cover.out

coverage: coverage-total coverage-html ## Total coverage
	@echo "==> coverage"

test-race: ## Test for data races
	@echo "==> test-race"
	go test -race ./...

test: ## Run tests
	@echo "==> test"
	go test ./...

build: ## Build the binary
	@echo "==> build"
	mkdir -p build
	go build -o build/main cmd/*.go

run: setup ## Build and run with air (live reloading)
	@echo "==> run"
	mkdir -p build
	rm -f build/* 2> /dev/null
	./bin/air --build.cmd 'go build -ldflags "-X main.version=$(VERSION)" -o build/main cmd/*.go' --build.bin "./build/main"

run-memory: setup ## Build and run with in-memory loader
	@echo "==> run-memory"
	mkdir -p build
	rm -f build/* 2> /dev/null
	./bin/air --build.cmd 'go build -ldflags "-X main.version=$(VERSION)" -o build/main cmd/*.go' --build.bin "./build/main" --build.args_bin "-f config/config-memory.yaml"

clean: ## Remove build artifacts
	@echo "==> clean"

build-docker: ## Build docker container
	@echo "==> build-docker"
	docker build --build-arg VERSION=$(VERSION) -t board .

run-docker: build-docker ## Run docker container
	@echo "==> run-docker"
	docker run -p 8080:8080 board

test-bench: ## Run benchmarks
	@echo "==> test-bench"
	go test -bench=. -benchmem ./integration/

monitoring-up: ## Run app and monitoring
	mkdir -p ./logs
	GRAFANA_ROOT_URL="" GRAFANA_SUB_PATH="" LOG_PATH=./logs docker compose --profile monitoring -f docker-compose.yaml up -d
	@docker logs -f board > ./logs/board.log &
	@echo "Remember to shut everything down with 'make monitoring-down'"

monitoring-down: ## Close app and monitoring
	GRAFANA_ROOT_URL="" GRAFANA_SUB_PATH="" LOG_PATH=./logs docker compose --profile monitoring -f docker-compose.yaml down

