## Makefile: Thin, explicit wrappers for tools
## - One responsibility per target
## - No dynamic variables or shell logic
## - Real logic lives in scripts (TypeScript/Bun, bash)

.PHONY: lint format test test-unit cov coverage coverage-go coverage-critical coverage-threshold build build-dev build-dist typecheck e2e release release-gh clean complexity sec dup perf-smoke test-race contract-snapshots release-check doc-design review thoth-meta thoth-meta-go thoth-meta-go-test thoth-meta-ts-e2e check-tools install-tools-help repo-sync repo-audit help

BIOME := npx @biomejs/biome
BUN := bun
GO := go
GOLINT := golangci-lint
GO_ENV := GOCACHE=$(PWD)/.gocache GOMODCACHE=$(PWD)/.gomodcache
GOLINT_ENV := $(GO_ENV) GOLANGCI_LINT_CACHE=$(PWD)/.golangci-lint-cache
THOTH := thoth
TMP_DIR := ./temp
GO_PACKAGES := ./...
COVER_PROFILE := $(TMP_DIR)/test-unit.coverage.out
COVER_HTML := $(TMP_DIR)/test-unit.coverage.html
CRITICAL_COVER_PROFILE := $(TMP_DIR)/critical.coverage.out
CRITICAL_PACKAGES := . ./internal/schema ./internal/argv ./internal/validators
COVERAGE_MIN := 90.0


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

test-unit:
	mkdir -p $(TMP_DIR)
	$(GO_ENV) $(GO) test -v -coverprofile=$(COVER_PROFILE) -covermode=count $(GO_PACKAGES)
	$(GO_ENV) $(GO) tool cover -func=$(COVER_PROFILE)

test:
	$(GO_ENV) $(GO) test -coverprofile=coverage.out ./...
	$(GO_ENV) $(GO) tool cover -func=coverage.out

coverage: coverage-go

coverage-go: test-unit
	$(GO_ENV) $(GO) tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@printf "Coverage HTML: %s\n" "$(COVER_HTML)"

coverage-critical:
	mkdir -p $(TMP_DIR)
	$(GO_ENV) $(GO) test -coverprofile=$(CRITICAL_COVER_PROFILE) -covermode=count $(CRITICAL_PACKAGES)
	$(GO_ENV) $(GO) tool cover -func=$(CRITICAL_COVER_PROFILE)
	@printf "Critical coverage profile: %s\n" "$(CRITICAL_COVER_PROFILE)"

coverage-threshold:
	mkdir -p $(TMP_DIR)
	$(GO_ENV) $(GO) test -coverprofile=$(COVER_PROFILE) -covermode=count $(GO_PACKAGES)
	@report="$$( $(GO_ENV) $(GO) tool cover -func=$(COVER_PROFILE) )"; \
	printf "%s\n" "$$report"; \
	below="$$(printf "%s\n" "$$report" | awk -v min=$(COVERAGE_MIN) ' \
	/^[^[:space:]].*:[0-9]+:/ { \
		pct=$$NF; gsub(/%/, "", pct); \
		if (pct+0 < min+0) print $$0; \
	} \
	')"; \
	if [ -n "$$below" ]; then \
		printf "\nFunctions below %.1f%% coverage:\n%s\n" $(COVERAGE_MIN) "$$below"; \
		exit 1; \
	fi

cov:
	npm run test:cov

build:
	$(BUN) run build-go.ts

build-dist:
	gh flarebyte build

build-dev:
	mkdir -p .e2e-bin
	$(GO_ENV) CGO_ENABLED=0 $(GO) build -o .e2e-bin/flyb ./cmd/flyb

typecheck:
	npm run typecheck

e2e:
	mkdir -p .e2e-bin
	$(GO_ENV) $(GO) build -o .e2e-bin/flyb ./cmd/flyb
	$(BUN) test script/e2e

doc-design: build-dev
	mkdir -p doc/design
	./.e2e-bin/flyb validate --config doc/design-meta/app.cue
	./.e2e-bin/flyb generate markdown --config doc/design-meta/app.cue

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

release-gh:
	gh flarebyte release

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

review: format test e2e lint

thoth-meta: thoth-meta-go thoth-meta-go-test thoth-meta-ts-e2e

thoth-meta-go:
	$(THOTH) run --config ./pipeline-go-maat.thoth.cue

thoth-meta-go-test:
	$(THOTH) run --config ./pipeline-go-test-maat.thoth.cue

thoth-meta-ts-e2e:
	$(THOTH) run --config ./pipeline-ts-e2e-maat.thoth.cue

repo-sync:
	gh flarebyte repo update --repo flarebyte/baldrick-flying-buttress

repo-audit:
	gh flarebyte repo audit --repo flarebyte/baldrick-flying-buttress

check-tools:
	@printf "go=%s\n" "$$(command -v $(GO) >/dev/null 2>&1 && echo true || echo false)"
	@printf "bun=%s\n" "$$(command -v $(BUN) >/dev/null 2>&1 && echo true || echo false)"
	@printf "npx=%s\n" "$$(command -v npx >/dev/null 2>&1 && echo true || echo false)"
	@printf "biome=%s\n" "$$(command -v npx >/dev/null 2>&1 && npx --yes @biomejs/biome --version >/dev/null 2>&1 && echo true || echo false)"
	@printf "golangci-lint=%s\n" "$$(command -v $(GOLINT) >/dev/null 2>&1 && echo true || echo false)"
	@printf "thoth=%s\n" "$$(command -v $(THOTH) >/dev/null 2>&1 && echo true || echo false)"

install-tools-help:
	@printf "Install hints:\n"
	@printf "  go: https://go.dev/doc/install\n"
	@printf "  bun: https://bun.sh/docs/installation\n"
	@printf "  npx/npm (Node.js): https://nodejs.org/\n"
	@printf "  biome: npm i -D @biomejs/biome\n"
	@printf "  golangci-lint: https://golangci-lint.run/welcome/install/\n"
	@printf "  thoth: project-specific CLI; ensure it is installed and on PATH\n"

help:
	@printf "Targets:\n"
	@printf "  lint       Run Biome checks.\n"
	@printf "  format     Format code with Biome and apply safe fixes.\n"
	@printf "  test       Run Go unit tests and coverage summary.\n"
	@printf "  cov        Run unit tests with coverage report (text-summary + lcov).\n"
	@printf "  build      Build Go release binaries into build/.\n"
	@printf "  build-dist Build distribution artifacts via gh flarebyte build.\n"
	@printf "  typecheck  Run TypeScript type-check only.\n"
	@printf "  e2e        Run Bun-powered end-to-end tests.\n"
	@printf "  doc-design Generate design docs in doc/design from doc/design-meta CUE.\n"
	@printf "  perf-smoke Run deterministic moderate-size Go smoke tests.\n"
	@printf "  test-race  Run Go tests with the race detector.\n"
	@printf "  contract-snapshots  Run contract snapshot and contract invariants.\n"
	@printf "  release-check  Run deterministic release gates in fixed order.\n"
	@printf "  release    Prepare release artifacts (depends on build).\n"
	@printf "  release-gh Publish GitHub release via gh flarebyte release.\n"
	@printf "  clean      Remove dist/ artifacts.\n"
	@printf "  complexity Show top TypeScript files by complexity via scc.\n"
	@printf "  sec        Run Semgrep security scan.\n"
	@printf "  dup        Run duplicate code detection (jscpd).\n"
	@printf "  check-tools Report required tool availability as key=value lines.\n"
	@printf "  install-tools-help  Show installation hints for required tools.\n"
	@printf "  repo-sync  Apply .gh-flarebyte.cue settings to GitHub repo.\n"
	@printf "  repo-audit Audit .gh-flarebyte.cue drift against GitHub repo.\n"
	@printf "  help       Show this help message.\n"
