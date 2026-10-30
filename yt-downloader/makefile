#!/bin/bash
# yt-downloader-go: consola en Go (sin docker). Se lanza desde la raiz del repo.
CONSOLE := build/console$(if $(filter Windows_NT,$(OS)),.exe,)
# version unica: la constante APP_VERSION del codigo (la muestra `console version`)
VERSION := $(shell sed -n 's/.*APP_VERSION *= *"\([^"]*\)".*/\1/p' src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go)
DIST_WIN   := build/dist/yt-win-$(VERSION).exe
DIST_LINUX := build/dist/yt-linux-$(VERSION)
DIST_MACOS := build/dist/yt-macos-$(VERSION)
LDFLAGS    := -ldflags "-s -w"

.PHONY: help env build console dist dist-win dist-linux dist-macos version test fmt vet gitpush

help: ## Show this help message · make help
	@echo "usage: make [target]"
	@echo
	@echo "targets:"
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env: ## crea .env desde .env.example (no pisa el existente) · make env
	@test -f .env || cp .env.example .env

build: env ## compila build/console · make build
	go build -o $(CONSOLE) ./cmd/console

console: build ## comando de consola · make console c="help"  |  make console c='video-info --url=https://youtu.be/xxx'
	./$(CONSOLE) $(c)

test: ## tests · make test  |  make test f="VideoFieldsValidator"
	go test -count=1 $(if $(f),-run $(f),) ./...

fmt: ## gofmt + go vet · make fmt
	gofmt -l -w ./cmd ./console ./src ./tests
	go vet ./...

gitpush: ## git push · make gitpush m="mensaje"
	clear;
	git pull;
	git add .; git commit -m "$(m)"; git push;
dist: dist-win dist-linux dist-macos ## los tres binarios para distribuir (build/dist/yt-{win,linux,macos}-<version>) · make dist
	@echo "subirlos a https://github.com/eacevedof/prj_go/releases (tag yt-downloader-v$(VERSION))"

dist-win: ## Windows amd64, build/dist/yt-win-<version>.exe · make dist-win
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_WIN) ./cmd/console

dist-linux: ## Linux amd64, build/dist/yt-linux-<version> · make dist-linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_LINUX) ./cmd/console

dist-macos: ## macOS Apple Silicon (arm64), build/dist/yt-macos-<version> · make dist-macos
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_MACOS) ./cmd/console

version: ## version actual (APP_VERSION) · make version
	@echo $(VERSION)

