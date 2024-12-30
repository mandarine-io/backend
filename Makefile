.TOOLS_DIR = $(PWD)/tools

.GREEN_COLOR = \033[0;32m
.RED_COLOR = \033[0;31m
.NO_COLOR = \033[0m

.PHONY: hooks
hooks:
	$(SHELL) $(.TOOLS_DIR)/setup-git-hooks.sh

.PHONY: help
help:
	@echo "Available commands:"
	@echo "	make help			${.GREEN_COLOR}Display this message${.NO_COLOR}"
	@echo "	make hooks			${.GREEN_COLOR}Run pre-commit, pre-push Git hooks${.NO_COLOR}"

.DEFAULT_GOAL := help