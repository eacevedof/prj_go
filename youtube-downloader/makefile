#!/bin/bash
TODAY := $(shell date +'%Y-%m-%d')
CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

BUILD_DIR := bin
DIST_DIR  := dist
BINARY    := $(BUILD_DIR)/extractor.exe
DIST_ZIP  := $(DIST_DIR)/extractor.zip

build: ## Compile the extractor binary to bin/extractor.exe
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) ./cmd/extractor

dist: build ## Build and package extractor.zip (extractor.exe + yt-dlp + ffmpeg + ffprobe)
	@mkdir -p $(DIST_DIR)
	powershell -Command "Compress-Archive -Force \
		-Path '$(BINARY)','$(BUILD_DIR)/yt-dlp.exe','$(BUILD_DIR)/ffmpeg.exe','$(BUILD_DIR)/ffprobe.exe' \
		-DestinationPath '$(DIST_ZIP)'"
	@echo ">> $(DIST_ZIP)"

help: ## Show this help message
	@echo "usage: make [target]"
	@echo
	@echo "targets:"
	@egrep "^(.+)\:\ ##\ (.+)" ${MAKEFILE_LIST} | column -t -c 2 -s ":#"

gitpush: ## git push m=any message
	clear;
	git add .; git commit -m "yt-donwloader: $(m)"; git push;

gp-def: ## git push with default yt-donwloader message
	clear;
	git add .; git commit -m "yt-donwloader"; git push;

CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
update-branch:  ## update master branch
	git fetch --all;
	git checkout master; git reset --hard origin/master;

	git checkout $(CURRENT_BRANCH);
	git diff --name-only develop master

