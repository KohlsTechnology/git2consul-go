TEST?=$(shell go list ./... | grep -v /vendor/)

# Get git commit information
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

default: test

test: generate
	@echo " ==> Running tests..."
	@go list $(TEST) \
		| grep -v "/vendor/" \
		| xargs -n1 go test -v -timeout=60s $(TESTARGS)
.PHONY: test

generate:
	@echo " ==> Generating..."
	@find . -type f -name '.DS_Store' -delete
	@go list ./... \
		| grep -v "/vendor/" \
		| xargs -n1 go generate $(PACKAGES)
.PHONY: generate


build: generate
	@echo " ==> Building..."
	@go build -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY}" .
.PHONY: build

build-linux: create-build-image remove-dangling build-native
.PHONY: build-linux

create-build-image:
	@docker build -t cimpress/git2consul-builder $(CURDIR)/build/
.PHONY: create-build-image

remove-dangling:
	@docker images --quiet --filter dangling=true | grep . | xargs docker rmi
.PHONY: remove-dangling

run-build-image:
	@echo " ===> Building..."
	@docker run --rm --name git2consul-builder -v $(CURDIR):/app -v $(CURDIR)/build/bin:/build/bin --entrypoint /app/build/build.sh cimpress/git2consul-builder
.PHONY: run-build-image
