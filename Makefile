.PHONY: clean build test

all: assemble

lint:
	@echo "\nApplying golint\n"
	@golint ./...

fmt:
	@echo "\nFormatting go files\n"
	@go fmt ./...

vet:
	@echo "\nApplying go vet\n"
	@go vet ./...

assemble: fmt lint vet
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
	@echo "\nStopping fake account API\n"
	@docker-compose stop

clean: uninstall
	@echo "\nCleaning all images"
	@cd scripts  && sh clean.sh