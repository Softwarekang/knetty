include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk

.PHONY: build
build:
	@$(MAKE) go.build

.PHONY: cover
cover:
	@$(MAKE) go.cover

.PHONY: test
test:
	@$(MAKE) go.test

.PHONY: lint
lint:
	@$(MAKE) go.lint