.SERVER_DIR = $(PWD)/cmd/api
.BUILD_DIR = $(PWD)/build
.SCRIPTS_DIR = $(PWD)/scripts
.API_DOCS_DIR = $(PWD)/docs/api
.SWAGGER_MODEL_DIR = $(PWD)/pkg/model/swagger
.LOCAL_DOCKER_COMPOSE_FILE = $(PWD)/docker-compose.local.yml
.INTEGRATION_TEST_DOCKER_COMPOSE_FILE = $(PWD)/docker-compose.integration-test.yml
.E2E_TEST_DOCKER_COMPOSE_FILE = $(PWD)/docker-compose.e2e-test.yml
.TEST_DIR = $(PWD)/tests
.TEST_RESULTS_DIR = $(PWD)/test-results
.LOGS_DIR = $(PWD)/logs
.CONFIG_PATH = $(PWD)/config/config.yaml
.MOCKERY_CONFIG_PATH = $(PWD)/.mockery.yaml
.LINTER_CONFIG_PATH = $(PWD)/.golangci.yaml
.ENV_FILE = $(PWD)/.env
.SERVER_TARGET = server
.DOCKER_IMAGE_PREFIX = backend

.GO = go
.NPM = npm
.AIR = air
.SWAG = swag
.SWAG2OAPI = swagger2openapi
.FORMATTER = gofmt
.LINTER = golangci-lint
.MOCKERY = mockery
.DOCKER = docker
.DOCKER_COMPOSE = docker compose

.GREEN_COLOR = \033[0;32m
.RED_COLOR = \033[0;31m
.NO_COLOR = \033[0m

.PHONY: clean
clean:
	@rm -rf $(.BUILD_DIR) $(.LOGS_DIR)

.PHONY: hooks
hooks:
	$(SHELL) $(.SCRIPTS_DIR)/setup-git-hooks.sh

.PHONY: install
install:
	$(.GO) mod tidy
	$(.GO) install github.com/air-verse/air@latest
	$(.GO) install github.com/swaggo/swag/cmd/swag@latest
	$(.GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(.GO) install github.com/vektra/mockery/v2@latest
	$(.NPM) i -g swagger2openapi

.PHONY: swagger.gen
swagger.gen:
	$(.GO) install github.com/swaggo/swag/cmd/swag@latest
	$(.SWAG) init --parseDependency --generalInfo internal/transport/http/router.go --outputTypes go,yaml,json --output $(.SWAGGER_MODEL_DIR)
	sed -i '' 's/github_com_mandarine-io_backend_pkg_model_//g' $(.SWAGGER_MODEL_DIR)/docs.go
	sed -i '' 's/github_com_mandarine-io_backend_pkg_model_//g' $(.SWAGGER_MODEL_DIR)/swagger.yaml
	sed -i '' 's/github_com_mandarine-io_backend_pkg_model_//g' $(.SWAGGER_MODEL_DIR)/swagger.json

.PHONY: mock.gen
mock.gen:
	$(.GO) install github.com/vektra/mockery/v2@latest
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

.PHONY: build
build:
	$(.GO) mod download
	$(.GO) build \
	-o $(.BUILD_DIR)/$(.SERVER_TARGET) \
	-installsuffix "static" \
	-tags "" \
	-ldflags " \
	-X main.Version=${ARTIFACT_VERSION} \
	-X main.GoVersion=$(shell go version | cut -d " " -f 3) \
	-X main.Compiler=$(shell go env CC) \
	-X main.Platform=$(shell go env GOOS)/$(shell go env GOARCH)" \
	$(.SERVER_DIR)

.PHONY: build.docker
build.docker:
	$(.DOCKER) build -f Dockerfile -t $(.DOCKER_IMAGE_PREFIX):latest .

.PHONY: start
start: build
	$(.BUILD_DIR)/$(.SERVER_TARGET)

.PHONY: start.dev
start.dev: build
	$(.AIR)

.PHONY: deploy.local
deploy.local:
	$(.DOCKER_COMPOSE) -f $(.LOCAL_DOCKER_COMPOSE_FILE) up -d

.PHONY: deploy.integration-test
deploy.integration-test:
	$(.DOCKER_COMPOSE) -f $(.INTEGRATION_TEST_DOCKER_COMPOSE_FILE) up -d

.PHONY: deploy.e2e-test
deploy.e2e-test:
	$(.DOCKER_COMPOSE) -f $(.E2E_TEST_DOCKER_COMPOSE_FILE) up -d

.PHONY: test.unit
test.unit:
	ALLURE_OUTPUT_PATH=$(.TEST_RESULTS_DIR) \
	ALLURE_OUTPUT_FOLDER=unit \
	$(.GO) test $(.TEST_DIR)/unit/... -v -shuffle on

.PHONY: test.integration
test.integration:
	ALLURE_OUTPUT_PATH=$(.TEST_RESULTS_DIR) \
	ALLURE_OUTPUT_FOLDER=integration \
	$(.GO) test $(.TEST_DIR)/integration/... -v -shuffle on

.PHONY: test.e2e
test.e2e:
	ALLURE_OUTPUT_PATH=$(.TEST_RESULTS_DIR) \
	ALLURE_OUTPUT_FOLDER=e2e \
	$(.GO) test $(.TEST_DIR)/e2e/... -v -shuffle on

.PHONY: help
help:
	@echo "Available commands:"
	@echo "	make help			${.GREEN_COLOR}Display this message${.NO_COLOR}"
	@echo "	make clean			${.GREEN_COLOR}Clean build and logs directories${.NO_COLOR}"
	@echo "	make hooks			${.GREEN_COLOR}Run pre-commit, pre-push Git hooks${.NO_COLOR}"
	@echo "	make install			${.GREEN_COLOR}Install necessary utilities${.NO_COLOR}"
	@echo "	make swagger.gen		${.GREEN_COLOR}Generate Swagger 2 and OpenAPI 3 specification${.NO_COLOR}"
	@echo "	make redoc.gen			${.GREEN_COLOR}Generate Redoc from OpenAPI 3${.NO_COLOR}"
	@echo "	make mock.gen			${.GREEN_COLOR}Generate mocks${.NO_COLOR}"
	@echo "	make format			${.GREEN_COLOR}Run formatting${.NO_COLOR}"
	@echo "	make format.fix			${.GREEN_COLOR}Run formatting and simplify code${.NO_COLOR}"
	@echo "	make lint			${.GREEN_COLOR}Run linters${.NO_COLOR}"
	@echo "	make lint.fix 			${.GREEN_COLOR}Run linters and fix found issues (if its supported by the linter)${.NO_COLOR}"
	@echo "	make build			${.GREEN_COLOR}Build the server executable file${.NO_COLOR}"
	@echo "	make start			${.GREEN_COLOR}Run the built server executable file${.NO_COLOR}"
	@echo "	make start.dev			${.GREEN_COLOR}Run the built server executable file with hot reload${.NO_COLOR}"
	@echo "	make deploy.local		${.GREEN_COLOR}Deploy to local environment${.NO_COLOR}"
	@echo "	make deploy.integration-test	${.GREEN_COLOR}Deploy to integration test environment${.NO_COLOR}"
	@echo "	make deploy.e2e-test		${.GREEN_COLOR}Deploy to E2E test environment${.NO_COLOR}"
	@echo "	make test.unit			${.GREEN_COLOR}Run unit tests${.NO_COLOR}"
	@echo "	make test.integration		${.GREEN_COLOR}Run integration tests${.NO_COLOR}"
	@echo "	make test.e2e			${.GREEN_COLOR}Run E2E tests${.NO_COLOR}"

.DEFAULT_GOAL := help