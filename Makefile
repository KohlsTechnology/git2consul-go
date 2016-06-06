TEST?=$(shell go list ./... | grep -v /vendor/)

.PHONY: build test

default: test

test: generate
	@echo " ==> Running tests..."
	@go list $(TEST) \
		| grep -v "/vendor/" \
		| xargs -n1 go test -v -timeout=60s $(TESTARGS)

generate:
	@echo " ==> Generating..."
	@find . -type f -name '.DS_Store' -delete
	@go list ./... \
		| grep -v "/vendor/" \
		| xargs -n1 go generate $(PACKAGES)

create_docker_image:
	@docker build -t cimpress/git2consul-builder $(CURDIR)/build/

build:
	@echo " ===> Building..."
	@docker run --rm --name git2consul-builder -v $(CURDIR):/app -v $(CURDIR)/build/bin:/build/bin cimpress/git2consul-builder
