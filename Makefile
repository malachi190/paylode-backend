BINARY        := paylode-app     # name of the produced binary
MAIN_PACKAGE  := ./cmd     # import path to main package
TABLE_NAME ?= users # default table name when creating migrations

.PHONY: build
build: ## Build the binary to ./bin/$(BINARY)
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) $(MAIN_PACKAGE)

.PHONY: run
run: ## Run the app with `go run`
	go run $(MAIN_PACKAGE)

# ---------- 3.  Development flow ----------
.PHONY: dev
dev: ## Start the app with air for live-reload (requires github.com/cosmtrek/air)
	@if ! command -v air >/dev/null 2>&1; then echo "Install air: go install github.com/cosmtrek/air@latest"; exit 1; fi
	air

.PHONY: tidy
tidy: ## Tidy go.mod & go.sum
	go mod tidy

.PHONY: clean
clean: ## Remove build artefacts
	rm -rf bin/

# ---------- 4.  Testing ----------
.PHONY: test
test: ## Run tests with race detector
	go test -race -coverprofile=cover.out ./...

# ---------- 5.  Set commands to create migration file and migrate ----------
.PHONY: migration
migration: 
	migrate create -ext sql -dir database/migrations -seq ${TABLE_NAME}