SERVER_DIR = $(PWD)/cmd/api
BUILD_DIR = $(PWD)/build
LOGS_DIR = $(PWD)/logs
TOOLS_DIR = $(PWD)/tools
API_DOCS_DIR = $(PWD)/docs/api
FORMATER_LOG_DIR = $(LOGS_DIR)/format
LINTER_LOG_DIR = $(LOGS_DIR)/lint
UNIT_TEST_LOG_DIR = $(LOGS_DIR)/unit-tests
E2E_TEST_LOG_DIR = $(LOGS_DIR)/e2e-tests
ALL_TEST_LOG_DIR = $(LOGS_DIR)/all-tests
CONFIG_PATH = $(PWD)/config/config.yaml
ENV_FILE = $(PWD)/.env
SERVER_TARGET = server
TIMESTAMP = $(shell date +%s)

GO = go
NPM = npm
AIR = air
SWAG = swag
SWAG2OP = swagger2openapi
REDOC = redocly
FORMATTER = gofmt
LINTER = golangci-lint

GREEN_COLOR = \033[0;32m
NO_COLOR = \033[0m

ifneq (,$(wildcard $(ENV_FILE)))
	include $(ENV_FILE)
endif

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR) $(LOGS_DIR)

.PHONY: hooks
hooks:
	$(SHELL) $(TOOLS_DIR)/setup-git-hooks.sh

.PHONY: install
install:
	$(GO) mod download
	$(GO) install github.com/air-verse/air@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: swagger.gen
swagger.gen:
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(SWAG) init --generalInfo ./internal/api/rest/router.go --outputTypes go,yaml --output $(API_DOCS_DIR)

.PHONY: openapi.gen
openapi.gen:
	$(NPM) i -g swagger2openapi
	$(SWAG2OP) --yaml --outfile $(API_DOCS_DIR)/openapi.yaml $(API_DOCS_DIR)/swagger.yaml

.PHONY: redoc.gen
redoc.gen:
	$(NPM) i -g @redocly/cli
	$(REDOC) build-docs --output $(API_DOCS_DIR)/redoc.html $(API_DOCS_DIR)/swagger.yaml

.PHONY: format
format:
	@mkdir -p $(LOGS_DIR)
	@mkdir -p $(FORMATER_LOG_DIR)
	$(FORMATTER) -w . | tee $(FORMATER_LOG_DIR)/output-$(TIMESTAMP).log

.PHONY: format.fix
format.fix:
	@mkdir -p $(LOGS_DIR)
	@mkdir -p $(FORMATER_LOG_DIR)
	$(FORMATTER) -s -w . | tee $(FORMATER_LOG_DIR)/output-$(TIMESTAMP).log

.PHONY: lint
lint:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@mkdir -p $(LOGS_DIR)
	@mkdir -p $(LINTER_LOG_DIR)
	$(LINTER) run | tee $(LINTER_LOG_DIR)/output-$(TIMESTAMP).log

.PHONY: lint.fix
lint.fix:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@mkdir -p $(LOGS_DIR)
	@mkdir -p $(LINTER_LOG_DIR)
	$(LINTER) run --fix | tee $(LINTER_LOG_DIR)/output-$(TIMESTAMP).log

.PHONY: build
build:
	$(GO) mod tidy
	$(GO) build -o $(BUILD_DIR)/$(SERVER_TARGET) $(SERVER_DIR)

.PHONY: start
start: build
	$(GO) mod tidy
	$(BUILD_DIR)/$(SERVER_TARGET) --config $(CONFIG_PATH) --env $(ENV_FILE)

.PHONY: start.dev
start.dev: build
	$(GO) mod tidy
	$(GO) install github.com/air-verse/air@latest
	$(AIR)

.PHONY: test.unit
test.unit:
	@mkdir -p $(LOGS_DIR)
	@mkdir -p $(UNIT_TEST_LOG_DIR)
	$(GO) test ./tests/unit/... -v -shuffle on -covermode atomic -coverprofile $(UNIT_TEST_LOG_DIR)/cover.out | tee $(UNIT_TEST_LOG_DIR)/output-$(TIMESTAMP).log
	$(GO) tool cover -html $(UNIT_TEST_LOG_DIR)/cover.out -o $(UNIT_TEST_LOG_DIR)/cover.html

.PHONY: test.e2e
test.e2e:
	@mkdir -p $(LOGS_DIR)
	@mkdir -p $(E2E_TEST_LOG_DIR)
	$(GO) test ./tests/e2e/... -v -shuffle on -covermode atomic -coverprofile $(E2E_TEST_LOG_DIR)/cover.out | tee $(E2E_TEST_LOG_DIR)/output-$(TIMESTAMP).log
	$(GO) tool cover -html $(E2E_TEST_LOG_DIR)/cover.out -o $(E2E_TEST_LOG_DIR)/cover.html

.PHONY: test.load
test.load:
	@echo Not supported

.PHONY: test.all
test.all:
	@mkdir -p $(LOGS_DIR)
	@mkdir -p $(ALL_TEST_LOG_DIR)
	@touch $(ALL_TEST_LOG_DIR)/all-tests-output-$(TIMESTAMP).log
	$(GO) test ./... -v -shuffle on -covermode atomic -coverprofile $(ALL_TEST_LOG_DIR)/cover.out | tee $(ALL_TEST_LOG_DIR)/all-tests-output-$(TIMESTAMP).log
	$(GO) tool cover -html $(ALL_TEST_LOG_DIR)/cover.out -o $(ALL_TEST_LOG_DIR)/cover.html

.PHONY: help
help:
	@echo "Available commands:"
	@echo "	make help			${GREEN_COLOR}Display this message${NO_COLOR}"
	@echo "	make clean			${GREEN_COLOR}Clean build and logs directories${NO_COLOR}"
	@echo "	make hooks			${GREEN_COLOR}Run pre-commit, pre-push Git hooks${NO_COLOR}"
	@echo "	make install			${GREEN_COLOR}Install necessary utilities${NO_COLOR}"
	@echo "	make swagger.gen		${GREEN_COLOR}Generate Swagger 2 specification${NO_COLOR}"
	@echo "	make openapi.gen		${GREEN_COLOR}Generate OpenAPI 3 specification from Swagger 2 specification${NO_COLOR}"
	@echo "	make redoc.gen			${GREEN_COLOR}Generate Redoc from OpenAPI 3${NO_COLOR}"
	@echo "	make format			${GREEN_COLOR}Run formatting${NO_COLOR}"
	@echo "	make format.fix			${GREEN_COLOR}Run formatting and simplify code${NO_COLOR}"
	@echo "	make lint			${GREEN_COLOR}Run linters${NO_COLOR}"
	@echo "	make lint.fix 			${GREEN_COLOR}Run linters and fix found issues (if it's supported by the linter)${NO_COLOR}"
	@echo "	make build			${GREEN_COLOR}Build the server executable file${NO_COLOR}"
	@echo "	make start			${GREEN_COLOR}Run the built server executable file${NO_COLOR}"
	@echo "	make start.dev			${GREEN_COLOR}Run the built server executable file with hot reload${NO_COLOR}"
	@echo "	make test.unit			${GREEN_COLOR}Run unit tests${NO_COLOR}"
	@echo "	make test.e2e			${GREEN_COLOR}Run e2e tests${NO_COLOR}"
	@echo "	make test.load			${GREEN_COLOR}Run load tests${NO_COLOR}"
	@echo "	make test.all			${GREEN_COLOR}Run all tests (without load)${NO_COLOR}"

.DEFAULT_GOAL := help