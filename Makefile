BINARY = git2consul
COMMIT := $(shell git rev-parse HEAD)
BRANCH := $(shell git symbolic-ref --short -q HEAD || echo HEAD)
DATE := $(shell date -u +%Y%m%d-%H:%M:%S)
VERSION_PKG = github.com/KohlsTechnology/git2consul-go/pkg/version
LDFLAGS := "-X ${VERSION_PKG}.Branch=${BRANCH} -X ${VERSION_PKG}.BuildDate=${DATE} \
	-X ${VERSION_PKG}.GitSHA1=${COMMIT}"

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -rf $(BINARY)

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(BINARY) -ldflags $(LDFLAGS)

.PHONY: test
test: fmt vet test-unit

.PHONY: test-unit
test-unit:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

# Make sure go.mod and go.sum are not modified
.PHONY: test-dirty
test-dirty: build
	go mod tidy
	git diff --exit-code

.PHONY: fmt
fmt:
	test -z "$(shell gofmt -l .)"

.PHONY: lint
lint:
	LINT_INPUT="$(shell go list ./...)"; golint -set_exit_status $$LINT_INPUT

.PHONY: vet
vet:
	VET_INPUT="$(shell go list ./...)"; go vet $$VET_INPUT
