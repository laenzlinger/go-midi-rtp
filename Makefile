GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=midibox
MIDIBOX=pi@midibox

.DEFAULT_GOAL := help

.PHONY: help

all: test build

build: ## build
	$(GOBUILD)

test: ## run unit tests
	$(GOTEST) -v ./...

clean: ## clean all temporary files
	$(GOCLEAN)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
