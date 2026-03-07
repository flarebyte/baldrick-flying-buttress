## Makefile: Thin, explicit wrappers for tools
## - One responsibility per target
## - No dynamic variables or shell logic
## - Real logic lives in scripts (TypeScript/Bun, bash)

.PHONY: lint format test cov build typecheck e2e release clean complexity sec dup perf-smoke test-race contract-snapshots release-check help

BIOME := npx @biomejs/biome
BUN := bun
GO := go
GOLINT := golangci-lint
GO_ENV := GOCACHE=$(PWD)/.gocache GOMODCACHE=$(PWD)/.gomodcache
GOLINT_ENV := $(GO_ENV) GOLANGCI_LINT_CACHE=$(PWD)/.golangci-lint-cache

lint:
	$(BIOME) check
	$(GO_ENV) $(GO) vet ./...
	$(GOLINT_ENV) $(GOLINT) run

format:
	find . -type f -name '*.go' \
		-not -path './.git/*' \
		-not -path './.gocache/*' \
		-not -path './.gomodcache/*' \
		-not -path './.e2e-bin/*' \
		-not -path './node_modules/*' \
		-print0 | xargs -0 -r gofmt -w
	$(BIOME) format --write .
	$(BIOME) check --unsafe --write

test:
	$(GO_ENV) $(GO) test -coverprofile=coverage.out ./...
	$(GO_ENV) $(GO) tool cover -func=coverage.out

cov:
	npm run test:cov

build:
	$(BUN) run build-go.ts

build-dev:
	mkdir -p .e2e-bin
	$(GO_ENV) CGO_ENABLED=0 $(GO) build -o .e2e-bin/flyb ./cmd/flyb

typecheck:
	npm run typecheck

e2e:
	mkdir -p .e2e-bin
	$(GO_ENV) $(GO) build -o .e2e-bin/flyb ./cmd/flyb
	$(BUN) test script/e2e

perf-smoke:
	$(GO_ENV) $(GO) test -run PerfSmoke ./internal/cli

test-race:
	$(GO_ENV) $(GO) test -race ./...

contract-snapshots:
	$(GO_ENV) $(GO) test -run 'TestContract|TestContractSnapshot' ./internal/...

release-check:
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) contract-snapshots
	$(MAKE) test-race
	$(MAKE) perf-smoke
	$(MAKE) e2e

release: release-check build
	$(BUN) run release-go.ts

clean:
	npm run clean

complexity:
	scc --sort complexity --by-file -i go . | head -n 15
	scc --sort complexity --by-file -i ts . | head -n 15

sec:
	semgrep scan --config auto

dup:
	npx jscpd --format go --min-lines 10 --ignore "**/.gomodcache/**,**/.gocache/**,**/.e2e-bin/**,**/node_modules/**,**/dist/**" --gitignore .
	npx jscpd --format typescript --min-lines 10 --gitignore .

help:
	@printf "Targets:\n"
	@printf "  lint       Run Biome checks.\n"
	@printf "  format     Format code with Biome and apply safe fixes.\n"
	@printf "  test       Run unit tests (Node test runner via tsx).\n"
	@printf "  cov        Run unit tests with coverage report (text-summary + lcov).\n"
	@printf "  build      Build Go release binaries into build/.\n"
	@printf "  typecheck  Run TypeScript type-check only.\n"
	@printf "  e2e        Run Bun-powered end-to-end tests.\n"
	@printf "  perf-smoke Run deterministic moderate-size Go smoke tests.\n"
	@printf "  test-race  Run Go tests with the race detector.\n"
	@printf "  contract-snapshots  Run contract snapshot and contract invariants.\n"
	@printf "  release-check  Run deterministic release gates in fixed order.\n"
	@printf "  release    Prepare release artifacts (depends on build).\n"
	@printf "  clean      Remove dist/ artifacts.\n"
	@printf "  complexity Show top TypeScript files by complexity via scc.\n"
	@printf "  sec        Run Semgrep security scan.\n"
	@printf "  dup        Run duplicate code detection (jscpd).\n"
	@printf "  help       Show this help message.\n"
