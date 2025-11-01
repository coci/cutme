APP=cutme
APP_EXECUTABLE="./$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor)
SHELL := /bin/zsh

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

check-quality:
	make lint
	make fmt
	make vet

lint:
	golangci-lint run --enable-all

vet:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

test:
	make tidy
	make vendor
	go test -v -timeout 10m ./... -coverprofile=coverage.out -json > report.json

coverage:
	make test
	go tool cover -html=coverage.out


build:
	mkdir -p out/
	go build -o $(APP_EXECUTABLE)
	@echo "Build passed"

run:
	make build
	chmod +x $(APP_EXECUTABLE)
	$(APP_EXECUTABLE)

clean:
	go clean
	rm -rf out/
	rm -f coverage*.out

vendor:
	go mod vendor


.PHONY: all test build vendor


all:
	make check-quality
	make test
	make build