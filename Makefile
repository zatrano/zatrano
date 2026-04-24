# ZATRANO — optional Makefile for POSIX shells (Git Bash, macOS, Linux).
# On Windows PowerShell, use: go run ./cmd/zatrano serve

# Changelog: install git-cliff — https://github.com/orhun/git-cliff#installation
#   (e.g. cargo install git-cliff, WinGet, or GitHub release binary)

.PHONY: build test fmt vet lint run doctor air gen-example openapi-export verify verify-race config-validate \
	changelog next-version
APP := ./cmd/zatrano

build:
	go build -o bin/zatrano $(APP)

test:
	go test ./... -count=1

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

run:
	go run $(APP) serve

doctor:
	go run $(APP) doctor

air:
	air

gen-example:
	go run $(APP) gen module example_item --dry-run

openapi-export:
	go run $(APP) openapi export --output api/openapi.merged.yaml

# Offline config check (no DB/Redis); same flags as CLI
config-validate:
	go run $(APP) config validate

# Same checks as: go run ./cmd/zatrano verify
verify:
	go run $(APP) verify

verify-race:
	go run $(APP) verify --race

# Regenerate CHANGELOG.md (Conventional Commits; requires git-cliff on PATH, full `git` history)
changelog:
	@command -v git-cliff >/dev/null 2>&1 || { echo "install git-cliff: https://github.com/orhun/git-cliff#installation"; exit 1; }
	git-cliff -c cliff.toml -o CHANGELOG.md

# Print next semver from commits since the last v* tag (no file changes)
next-version:
	@command -v git-cliff >/dev/null 2>&1 || { echo "install git-cliff: https://github.com/orhun/git-cliff#installation"; exit 1; }
	@git-cliff --bumped-version
