TAGS := -tags 'netgo'
BINARY := sonic-configdb-utils

SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
# gnu date format iso-8601 is parsable with Go RFC3339
BUILDDATE := $(shell date --iso-8601=seconds)
VERSION := $(or ${VERSION},$(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD))

LINKMODE := $(LINKMODE) \
		 -X 'github.com/metal-stack/v.Version=$(VERSION)' \
		 -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
		 -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
		 -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'

.PHONY: all
all: test test-generate build

.PHONY: build
build:
	go build \
		$(TAGS) \
		-ldflags \
		"$(LINKMODE)" \
		-o bin/$(BINARY) \
		github.com/metal-stack/sonic-configdb-utils

.PHONY: test
test:
	go test ./... -count 1

.PHONY: test-generate
test-generate: docker-build
	./tests/test.sh

.PHONY: lint
lint:
	golangci-lint run -p bugs -p unused

.PHONY: docker-build
docker-build:
	docker build -t sonic-configdb-utils:local .
