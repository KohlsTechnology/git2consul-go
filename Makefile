BINARY = git2consul
COMMIT := $(shell git rev-parse HEAD)
BRANCH := $(shell git symbolic-ref --short -q HEAD || echo HEAD)
DATE := $(shell date -u +%Y%m%d-%H:%M:%S)
VERSION_PKG = github.com/KohlsTechnology/git2consul-go/pkg/version
LDFLAGS := "-X ${VERSION_PKG}.Branch=${BRANCH} -X ${VERSION_PKG}.BuildDate=${DATE} \
	-X ${VERSION_PKG}.GitSHA1=${COMMIT}"
TAG?=""

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -rf $(BINARY) dist/

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(BINARY) -ldflags $(LDFLAGS)

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: test
test: fmt vet test-unit

.PHONY: test-unit
test-unit:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-e2e
test-e2e: build
	go test -v ./pkg/e2e/e2e_test.go

# Make sure go.mod and go.sum are not modified
.PHONY: test-dirty
test-dirty: vendor build
	go mod tidy
	git diff --exit-code
	# TODO: also check that there are no untracked files, e.g. extra .go

# Make sure goreleaser is working
.PHONY: test-release
test-release:
	BRANCH=$(BRANCH) COMMIT=$(COMMIT) DATE=$(DATE) VERSION_PKG=$(VERSION_PKG) goreleaser --snapshot --skip-publish --rm-dist

.PHONY: fmt
fmt:
	test -z "$(shell gofmt -l . | grep -v ^vendor)"

.PHONY: lint
lint:
	LINT_INPUT="$(shell go list ./...)"; golint -set_exit_status $$LINT_INPUT

.PHONY: vet
vet:
	VET_INPUT="$(shell go list ./...)"; go vet $$VET_INPUT

.PHONY: tag
tag:
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)

# Requires GITHUB_TOKEN environment variable to be set
.PHONY: release
release:
	BRANCH=$(BRANCH) COMMIT=$(COMMIT) DATE=$(DATE) VERSION_PKG=$(VERSION_PKG) goreleaser
