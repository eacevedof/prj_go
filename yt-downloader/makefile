#!/bin/bash
# yt-downloader-go: consola en Go (sin docker). Se lanza desde la raiz del repo.
CONSOLE := build/console$(if $(filter Windows_NT,$(OS)),.exe,)

.PHONY: help env build console test fmt vet gitpush

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
