
# the product we're building
NAME   := instaunit
MODULE := github.com/instaunit/instaunit
MAIN   := $(MODULE)/main

# defaults
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# build and packaging
GITHASH   := $(shell git log --pretty=format:'%h' -n 1)
VERSION   ?= $(GITHASH)

BUILD_DIR  := $(PWD)/target
PRODUCT    ?= $(NAME)-local
TARGET_DIR := $(BUILD_DIR)/$(PRODUCT)
ARCHIVE    := $(PRODUCT).tgz
PACKAGE    := $(BUILD_DIR)/$(ARCHIVE)

# build and install
PREFIX ?= /usr/local

# sources
SRC = $(shell find src -name \*.go -print)
# tests
TEST_PKGS = $(MODULE)/hunit/...

.PHONY: all test clean install release build package formula

all: build

$(TARGET_DIR)/bin/$(NAME): $(SRC)
	(cd src && go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $@ $(MAIN))

build: $(TARGET_DIR)/bin/$(NAME) ## Build the product

$(PACKAGE): $(TARGET_DIR)/bin/$(NAME)
	(cd $(BUILD_DIR) && tar -zcf $(ARCHIVE) $(PRODUCT))

package: $(PACKAGE)

formula: package
	$(PWD)/build/update-formula -v $(VERSION) -o $(PWD)/formula/instaunit.rb $(PACKAGE)

release: test ## Build for all supported architectures
	make package PRODUCT=$(NAME)-$(VERSION)-linux-amd64 GOOS=linux GOARCH=amd64
	make package PRODUCT=$(NAME)-$(VERSION)-freebsd-amd64 GOOS=freebsd GOARCH=amd64
	make formula PRODUCT=$(NAME)-$(VERSION)-darwin-amd64 GOOS=darwin GOARCH=amd64

install: build ## Build and install
	install -m 0755 $(TARGET_DIR)/bin/$(NAME) $(PREFIX)/bin/

test: ## Run tests
	(cd src && go test $(TEST_PKGS))

clean: ## Delete the built product and any generated files
	rm -rf $(BUILD_DIR)
