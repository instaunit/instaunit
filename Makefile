
# the product we're building
NAME := hunit
# the product's main package
MAIN := ./src/main
# fix our gopath
GOPATH := $(GOPATH):$(PWD)

# build and packaging
TARGETS		:= $(PWD)/bin
PRODUCT		:= $(TARGETS)/$(NAME)

# sources
SRC = $(shell find src -name \*.go -print)

# tests
TEST_PKGS = hunit hunit/text/slug

.PHONY: all build test clean

all: build

$(PRODUCT): $(SRC)
	go build -o $@ $(MAIN)

build: $(PRODUCT) ## Build the service

test: ## Run tests
	go test -test.v $(TEST_PKGS)

clean: ## Delete the built product and any generated files
	rm -rf $(TARGETS)
