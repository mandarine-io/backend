.SERVER_DIR = $(PWD)/cmd/api
.BUILD_DIR = $(PWD)/build
.TOOLS_DIR = $(PWD)/tools
.API_DOCS_DIR = $(PWD)/docs/api
.DOCKER_COMPOSE_DIR = $(PWD)/deploy/docker
.LOCAL_DOCKER_COMPOSE_FILE = $(.DOCKER_COMPOSE_DIR)/docker-compose.local.yml
.DEV_DOCKER_COMPOSE_FILE = $(.DOCKER_COMPOSE_DIR)/docker-compose.dev.yml
.TEST_DOCKER_COMPOSE_FILE = $(.DOCKER_COMPOSE_DIR)/docker-compose.test.yml
.TEST_DIR = $(PWD)/tests
.UNIT_TEST_DIR = $(.TEST_DIR)/unit
.E2E_TEST_DIR = $(.TEST_DIR)/e2e
.LOAD_TEST_DIR = $(.TEST_DIR)/load
.LOGS_DIR = $(PWD)/logs
.FORMATER_LOG_DIR = $(.LOGS_DIR)/format
.LINTER_LOG_DIR = $(.LOGS_DIR)/lint
.UNIT_TEST_LOG_DIR = $(.LOGS_DIR)/unit-tests
.E2E_TEST_LOG_DIR = $(.LOGS_DIR)/e2e-tests
.LOAD_TEST_LOG_DIR = $(.LOGS_DIR)/load-tests
.CONFIG_PATH = $(PWD)/config/config.yaml
.APP_ENV_FILE = $(PWD)/.env
.LOCAL_ENV_FILE = $(.DOCKER_COMPOSE_DIR)/.env.local
.DEV_ENV_FILE = $(.DOCKER_COMPOSE_DIR)/.env.dev
.TEST_ENV_FILE = $(.DOCKER_COMPOSE_DIR)/.env.test
.SERVER_TARGET = server
.TIMESTAMP = $(shell date +%s)

.GO = go
.NPM = npm
.AIR = air
.SWAG = swag
.SWAG2OP = swagger2openapi
.REDOC = redocly
.FORMATTER = gofmt
.LINTER = golangci-lint
.DOCKER = docker
.DOCKER_COMPOSE = docker compose
.K6 = k6
.K6_IMAGE = grafana/k6:0.54.0

.GREEN_COLOR = \033[0;32m
.RED_COLOR = \033[0;31m
.NO_COLOR = \033[0m

.PHONY: clean
clean:
	@rm -rf $(.BUILD_DIR) $(.LOGS_DIR)

.PHONY: hooks
hooks:
	$(SHELL) $(.TOOLS_DIR)/setup-git-hooks.sh

.PHONY: install
install:
	$(.GO) mod download
	$(.GO) install github.com/air-verse/air@latest
	$(.GO) install github.com/swaggo/swag/cmd/swag@latest
	$(.GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: swagger.gen
swagger.gen:
	$(.GO) install github.com/swaggo/swag/cmd/swag@latest
	$(.SWAG) init --generalInfo ./internal/api/transport/http/router.go --outputTypes go,yaml,json --output $(.API_DOCS_DIR)

.PHONY: openapi.gen
openapi.gen:
	$(.NPM) i -g swagger2openapi
	$(.SWAG2OP) --yaml --outfile $(.API_DOCS_DIR)/openapi.yaml $(.API_DOCS_DIR)/swagger.yaml

.PHONY: redoc.gen
redoc.gen:
	$(.NPM) i -g @redocly/cli
	$(.REDOC) build-docs --output $(.API_DOCS_DIR)/redoc.html $(.API_DOCS_DIR)/swagger.yaml

.PHONY: format
format:
	@mkdir -p $(.LOGS_DIR)
	@mkdir -p $(.FORMATER_LOG_DIR)
	$(.FORMATTER) -w . | tee $(.FORMATER_LOG_DIR)/output-$(.TIMESTAMP).log

.PHONY: format.fix
format.fix:
	@mkdir -p $(.LOGS_DIR)
	@mkdir -p $(.FORMATER_LOG_DIR)
	$(.FORMATTER) -s -w . | tee $(.FORMATER_LOG_DIR)/output-$(.TIMESTAMP).log

.PHONY: lint
lint:
	$(.GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@mkdir -p $(.LOGS_DIR)
	@mkdir -p $(.LINTER_LOG_DIR)
	$(.LINTER) run | tee $(.LINTER_LOG_DIR)/output-$(.TIMESTAMP).log

.PHONY: lint.fix
lint.fix:
	$(.GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@mkdir -p $(.LOGS_DIR)
	@mkdir -p $(.LINTER_LOG_DIR)
	$(.LINTER) run --fix | tee $(.LINTER_LOG_DIR)/output-$(.TIMESTAMP).log

.PHONY: build
build:
	$(.GO) mod tidy
	$(.GO) build -o $(.BUILD_DIR)/$(.SERVER_TARGET) $(.SERVER_DIR)

.PHONY: start
start: build
	$(.GO) mod tidy
	if [ ! -f $(.CONFIG_PATH) ] && [ ! -f $(.APP_ENV_FILE) ]; \
		then $(.BUILD_DIR)/$(.SERVER_TARGET); \
	elif [ ! -f $(.CONFIG_PATH) ]; \
		then $(.BUILD_DIR)/$(.SERVER_TARGET) --env $(.APP_ENV_FILE); \
	elif [ ! -f $(.APP_ENV_FILE) ]; \
		then $(.BUILD_DIR)/$(.SERVER_TARGET) --config $(.CONFIG_PATH); \
	else \
	  $(.BUILD_DIR)/$(.SERVER_TARGET) --config $(.CONFIG_PATH) --env $(.APP_ENV_FILE); \
	fi

.PHONY: start.dev
start.dev: build
	$(.GO) mod tidy
	$(.GO) install github.com/air-verse/air@latest
	$(.AIR)

.PHONY: test.unit
test.unit:
	@mkdir -p $(.LOGS_DIR)
	@mkdir -p $(.UNIT_TEST_LOG_DIR)
	$(.GO) test $(.UNIT_TEST_DIR)/... -v -shuffle on -covermode atomic -coverprofile $(.UNIT_TEST_LOG_DIR)/cover.out | tee $(.UNIT_TEST_LOG_DIR)/output-$(.TIMESTAMP).log
	$(.GO) tool cover -html $(.UNIT_TEST_LOG_DIR)/cover.out -o $(.UNIT_TEST_LOG_DIR)/cover.html

.PHONY: test.e2e
test.e2e:
	@mkdir -p $(.LOGS_DIR)
	@mkdir -p $(.E2E_TEST_LOG_DIR)
	$(.GO) test $(.E2E_TEST_DIR)/... -v -shuffle on -covermode atomic -coverprofile $(.E2E_TEST_LOG_DIR)/cover.out | tee $(.E2E_TEST_LOG_DIR)/output-$(.TIMESTAMP).log
	$(.GO) tool cover -html $(.E2E_TEST_LOG_DIR)/cover.out -o $(.E2E_TEST_LOG_DIR)/cover.html

.PHONY: test.load
test.load:
	if [ ! -f $(.LOAD_TEST_DIR)/$(LOAD_TEST_NAME) ]; then echo "${.RED_COLOR}Load test file $(.LOAD_TEST_DIR)/$(LOAD_TEST_NAME) not found${.NO_COLOR}"; exit 1; fi
	@mkdir -p $(.LOGS_DIR)
	@mkdir -p $(.LOAD_TEST_LOG_DIR)
	@mkdir -p $(.LOAD_TEST_LOG_DIR)/$(LOAD_TEST_NAME)
	K6_WEB_DASHBOARD=true \
	K6_WEB_DASHBOARD_EXPORT=$(.LOAD_TEST_LOG_DIR)/$(LOAD_TEST_NAME)/report-$(.TIMESTAMP).html \
	$(.K6) run \
	--profiling-enabled \
	$(.LOAD_TEST_DIR)/$(LOAD_TEST_NAME) | tee $(.LOAD_TEST_LOG_DIR)/$(LOAD_TEST_NAME)/output-$(.TIMESTAMP).log

.PHONY: deploy.local
deploy.local:
	if [ ! -f $(.LOCAL_ENV_FILE) ]; then echo "${.RED_COLOR}Local environment file $(.LOCAL_ENV_FILE) not found${.NO_COLOR}"; exit 1; fi
	$(.DOCKER_COMPOSE) -f $(.LOCAL_DOCKER_COMPOSE_FILE) --env-file $(.LOCAL_ENV_FILE) up -d

.PHONY: deploy.dev
deploy.dev:
	if [ ! -f $(.DEV_ENV_FILE) ]; then echo "${.RED_COLOR}Development environment file $(.DEV_ENV_FILE) not found${.NO_COLOR}"; exit 1; fi
	$(.DOCKER_COMPOSE) -f $(.DEV_DOCKER_COMPOSE_FILE) --env-file $(.DEV_ENV_FILE) up -d

.PHONY: deploy.test
deploy.test:
	if [ ! -f $(.TEST_ENV_FILE) ]; then echo "${.RED_COLOR}Test environment file $(.TEST_ENV_FILE) not found${.NO_COLOR}"; exit 1; fi
	$(.DOCKER_COMPOSE) -f $(.TEST_DOCKER_COMPOSE_FILE) --env-file $(.TEST_ENV_FILE) up -d

.PHONY: help
help:
	@echo "Available commands:"
	@echo "	make help			${.GREEN_COLOR}Display this message${.NO_COLOR}"
	@echo "	make clean			${.GREEN_COLOR}Clean build and logs directories${.NO_COLOR}"
	@echo "	make hooks			${.GREEN_COLOR}Run pre-commit, pre-push Git hooks${.NO_COLOR}"
	@echo "	make install			${.GREEN_COLOR}Install necessary utilities${.NO_COLOR}"
	@echo "	make swagger.gen		${.GREEN_COLOR}Generate Swagger 2 specification${.NO_COLOR}"
	@echo "	make openapi.gen		${.GREEN_COLOR}Generate OpenAPI 3 specification from Swagger 2 specification${.NO_COLOR}"
	@echo "	make redoc.gen			${.GREEN_COLOR}Generate Redoc from OpenAPI 3${.NO_COLOR}"
	@echo "	make format			${.GREEN_COLOR}Run formatting${.NO_COLOR}"
	@echo "	make format.fix			${.GREEN_COLOR}Run formatting and simplify code${.NO_COLOR}"
	@echo "	make lint			${.GREEN_COLOR}Run linters${.NO_COLOR}"
	@echo "	make lint.fix 			${.GREEN_COLOR}Run linters and fix found issues (if it's supported by the linter)${.NO_COLOR}"
	@echo "	make build			${.GREEN_COLOR}Build the server executable file${.NO_COLOR}"
	@echo "	make start			${.GREEN_COLOR}Run the built server executable file${.NO_COLOR}"
	@echo "	make start.dev			${.GREEN_COLOR}Run the built server executable file with hot reload${.NO_COLOR}"
	@echo "	make test.unit			${.GREEN_COLOR}Run unit tests${.NO_COLOR}"
	@echo "	make test.e2e			${.GREEN_COLOR}Run e2e tests${.NO_COLOR}"
	@echo "	make test.load LOAD_TEST_NAME=<name>	${.GREEN_COLOR}Run load tests (set env vars used in test)${.NO_COLOR}"
	@echo "	make deploy.local		${.GREEN_COLOR}Deploy to local environment${.NO_COLOR}"
	@echo "	make deploy.dev			${.GREEN_COLOR}Deploy to development environment${.NO_COLOR}"
	@echo "	make deploy.test		${.GREEN_COLOR}Deploy to test environment${.NO_COLOR}"

.DEFAULT_GOAL := help