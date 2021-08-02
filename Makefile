.PHONY: clean build test

all: assemble

lint:
	@echo "\nApplying golint\n"
	@golint ./...

fmt:
	@echo "\nFormatting go files\n"
	@go fmt ./...

assemble: fmt lint
	@echo "\nBuilding application\n"
	@go build

unit-test: assemble
	@echo "\nRunning unit tests\n"
	@go test -cover -short ./...

it-test: assemble install
	@echo "\nRunning integration tests\n"
	@go test -cover -run Integration ./...

test: fmt unit-test install it-test
	@echo "\nRunning tests\n"

install:
	@echo "\nStarting fake account API\n"
	@docker-compose up -d --force-recreate

uninstall:
	@echo "\nStopping fake account API"
	@docker-compose stop
