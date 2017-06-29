
# the product we're building
NAME := instaunit
# the product's main package
MAIN := ./src/main
# fix our gopath
GOPATH := $(GOPATH):$(PWD)

# build and packaging
TARGETS := $(PWD)/bin
PRODUCT := $(TARGETS)/$(NAME)

# build and packaging for release
VERSION 				?= $(shell git log --pretty=format:'%h' -n 1)
RELEASE_TARGETS	 = $(PWD)/target/$(GOOS)_$(GOARCH)
RELEASE_PRODUCT	 = instaunit-$(VERSION)
RELEASE_ARCHIVE	 = $(RELEASE_PRODUCT).tgz
RELEASE_PACKAGE	 = $(RELEASE_TARGETS)/$(RELEASE_ARCHIVE)
RELEASE_BINARY   = $(RELEASE_TARGETS)/$(RELEASE_PRODUCT)/bin/$(NAME)

# sources
SRC = $(shell find src -name \*.go -print)

# tests
TEST_PKGS = hunit \
						hunit/expr \
						hunit/text/slug

.PHONY: all test clean release build build_darwin_amd64 build_linux_amd64 build_freebsd_amd64

all: build

$(PRODUCT): $(SRC)
	go build -o $@ $(MAIN)

build: $(PRODUCT) ## Build the product

$(RELEASE_PACKAGE): $(SRC)
	go build -o $(RELEASE_BINARY) $(MAIN)
	(cd $(RELEASE_TARGETS) && tar -zcf $(RELEASE_ARCHIVE) $(RELEASE_PRODUCT))

build_release: $(RELEASE_PACKAGE)

release: ## Build for all supported architectures
	make build_release GOOS=darwin GOARCH=amd64
	make build_release GOOS=linux GOARCH=amd64
	make build_release GOOS=freebsd GOARCH=amd64

test: ## Run tests
	go test -test.v $(TEST_PKGS)

clean: ## Delete the built product and any generated files
	rm -rf $(TARGETS)
