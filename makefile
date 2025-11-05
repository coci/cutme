APP=cutme
APP_EXECUTABLE=out/$(APP)
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor)
SHELL := /bin/zsh

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

# Quiet mode by default (use `make V=1 …` for verbose)
ifeq ($(V),1)
	Q :=
else
	Q := @
endif

# Avoid noisy "Entering/Leaving directory" messages
MAKEFLAGS += --no-print-directory

check-quality:
	$(Q)$(MAKE) lint
	$(Q)$(MAKE) fmt
	$(Q)$(MAKE) vet

lint:
	$(Q)golangci-lint run

vet:
	$(Q)go vet ./...

fmt:
	$(Q)go fmt ./...

tidy:
	$(Q)go mod tidy

test:
	$(Q)$(MAKE) tidy
	$(Q)$(MAKE) vendor
	$(Q)go test ./...

coverage:
	$(Q)$(MAKE) test
	$(Q)go tool cover -html=coverage.out


build:
	$(Q)$(MAKE) test
	$(Q)mkdir -p out/
	$(Q)go build -o $(APP_EXECUTABLE) ./cmd/main.go
	$(Q)echo "Build passed"

run:
	$(Q)$(MAKE) build
	$(Q)chmod +x $(APP_EXECUTABLE)
	$(Q)$(APP_EXECUTABLE) -config ./configs/config.yaml

clean:
	$(Q)go clean
	$(Q)rm -rf out/
	$(Q)find . -maxdepth 1 -name 'coverage*.out' -delete || true

vendor:
	$(Q)go mod vendor

all:
	$(Q)$(MAKE) check-quality
	$(Q)$(MAKE) test
	$(Q)$(MAKE) build