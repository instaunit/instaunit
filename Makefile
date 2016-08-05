
# the product we're building
PRODUCT := hunit
# the product's main package
MAIN := ./src/main
# fix our gopath
GOPATH := $(GOPATH):$(PWD)

# build and packaging
TARGETS		:= $(PWD)/target
BUILD_DIR	:= $(TARGETS)/$(PRODUCT)
PRODUCT		:= $(BUILD_DIR)/bin/hunit

# sources
SRC = $(shell find src -name \*.go -print)

.PHONY: all build test clean

all: build

$(PRODUCT): $(SRC)
	mkdir -p $(BUILD_DIR)/bin
	go build -o $@ $(MAIN)

build: $(PRODUCT) ## Build the service

test: ## Run tests

clean: ## Delete the built product and any generated files
	rm -rf $(TARGETS)
