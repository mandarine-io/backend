.TEST_DIR = $(PWD)/tests
.TEST_RESULTS_DIR = $(PWD)/test-results
.UNIT_TEST_DIR = $(.TEST_DIR)/unit
.E2E_TEST_DIR = $(.TEST_DIR)/e2e
.LOGS_DIR = $(PWD)/logs
.MOCKERY_CONFIG_PATH = $(PWD)/.mockery.yaml
.LINTER_CONFIG_PATH = $(PWD)/golangci.yaml

.GO = go
.FORMATTER = gofmt
.LINTER = golangci-lint
.MOCKERY = mockery

.GREEN_COLOR = \033[0;32m
.RED_COLOR = \033[0;31m
.NO_COLOR = \033[0m

.PHONY: clean
clean:
	@rm -rf $(.LOGS_DIR)

.PHONY: install
install:
	$(.GO) mod download
	$(.GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(.GO) install github.com/vektra/mockery/v2@latest

.PHONY: mock.gen
mock.gen:
	$(.MOCKERY) --config $(.MOCKERY_CONFIG_PATH)

.PHONY: format
format:
	$(.FORMATTER) -w .

.PHONY: format.fix
format.fix:
	$(.FORMATTER) -s -w .

.PHONY: lint
lint:
	$(.GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(.LINTER) run -c $(.LINTER_CONFIG_PATH)

.PHONY: lint.fix
lint.fix:
	$(.GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(.LINTER) run --fix -c $(.LINTER_CONFIG_PATH)

.PHONY: test.unit
test.unit:
	ALLURE_OUTPUT_PATH=$(.TEST_RESULTS_DIR) \
	ALLURE_OUTPUT_FOLDER=unit \
	$(.GO) test $(.UNIT_TEST_DIR)/... -v -shuffle on

.PHONY: test.e2e
test.e2e:
	ALLURE_OUTPUT_PATH=$(.TEST_RESULTS_DIR) \
	ALLURE_OUTPUT_FOLDER=e2e \
	APP_TEST_DOCKER_COMPOSE_PATH=$(.DOCKER_COMPOSE_TEST_PATH) \
	$(.GO) test $(.E2E_TEST_DIR)/... -v -shuffle on

.PHONY: help
help:
	@echo "Available commands:"
	@echo "	make help			${.GREEN_COLOR}Display this message${.NO_COLOR}"
	@echo "	make clean			${.GREEN_COLOR}Clean build and logs directories${.NO_COLOR}"
	@echo "	make install			${.GREEN_COLOR}Install necessary utilities${.NO_COLOR}"
	@echo "	make mock.gen			${.GREEN_COLOR}Generate mocks${.NO_COLOR}"
	@echo "	make format			${.GREEN_COLOR}Run formatting${.NO_COLOR}"
	@echo "	make format.fix			${.GREEN_COLOR}Run formatting and simplify code${.NO_COLOR}"
	@echo "	make lint			${.GREEN_COLOR}Run linters${.NO_COLOR}"
	@echo "	make lint.fix 			${.GREEN_COLOR}Run linters and fix found issues (if its supported by the linter)${.NO_COLOR}"
	@echo "	make test.unit			${.GREEN_COLOR}Run unit tests${.NO_COLOR}"
	@echo "	make test.e2e			${.GREEN_COLOR}Run e2e tests${.NO_COLOR}"

.DEFAULT_GOAL := help