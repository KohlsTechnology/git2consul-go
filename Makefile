TEST?=$(shell go list ./... | grep -v /vendor/)


default: test

test: generate
	@echo " ==> Running tests..."
	@go list $(TEST) \
		| grep -v "/vendor/" \
		| xargs -n1 go test -timeout=60s $(TESTARGS)

generate:
	@echo " ==> Generating..."
	@find . -type f -name '.DS_Store' -delete
	@go list ./... \
		| grep -v "/vendor/" \
		| xargs -n1 go generate $(PACKAGES)
