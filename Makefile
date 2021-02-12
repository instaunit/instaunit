
# the product we're building
NAME   := instaunit
MODULE := github.com/instaunit/instaunit
MAIN   := $(MODULE)/main

# defaults
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# build and packaging
GITHASH   := $(shell git log --pretty=format:'%h' -n 1)
BRANCH    := $(shell git rev-parse --abbrev-ref HEAD)
VERSION   ?= $(GITHASH)
LATEST    ?= latest

BUILD_DIR  := $(PWD)/target
PRODUCT    ?= $(NAME)-local
TARGET_DIR := $(BUILD_DIR)/$(PRODUCT)
ARCHIVE    := $(PRODUCT).tgz
PACKAGE    := $(BUILD_DIR)/$(ARCHIVE)
ARTIFACTS  := s3://instaunit/releases

# build and install
PREFIX ?= /usr/local

# sources
SRC = $(shell find src -name \*.go -print)
# tests
TEST_PKGS = $(MODULE)/hunit/...

.PHONY: all test clean install build package

all: build

$(TARGET_DIR)/bin/$(NAME): $(SRC)
	(cd src && go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $@ $(MAIN))

build: $(TARGET_DIR)/bin/$(NAME) ## Build the product

$(PACKAGE): $(TARGET_DIR)/bin/$(NAME)
	(cd $(BUILD_DIR) && tar -zcf $(ARCHIVE) $(PRODUCT))

package: $(PACKAGE)

install: build ## Build and install
	install -m 0755 $(TARGET_DIR)/bin/$(NAME) $(PREFIX)/bin/

test: ## Run tests
	(cd src && go test $(FLAGS) $(TEST_PKGS))

clean: ## Delete the built product and any generated files
	rm -rf $(BUILD_DIR)
