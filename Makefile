
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
BINARY     := $(TARGET_DIR)/bin/$(NAME)

# build and install
PREFIX ?= /usr/local

SRC       := $(shell find src -name \*.go -print)
TEST_PKGS := $(MODULE)/hunit/...
FIXTURES  := $(PWD)/fixtures
GRPC      := $(FIXTURES)/grpc

# utils
PSCTL ?= psctl

# platform-specific config
ifeq ($(shell uname),Linux)
	BASE64    ?= base64 -w 0
	ECHO      ?= echo -e
	INSTALL   ?= install -D
else
	BASE64    ?= base64
	ECHO      ?= echo
	INSTALL   ?= install
endif

.PHONY: all
all: build

.PHONY: tools
tools:
	@test -n "$$VILLAINS_SKIP_DEVEL_TOOL_CHECKS" || $(ECHO) "✔ Checking development tools; to disable checks, set: VILLAINS_SKIP_DEVEL_TOOL_CHECKS=true"
	@test -n "$$VILLAINS_SKIP_DEVEL_TOOL_CHECKS" || which $(PSCTL) &> /dev/null || ($(ECHO) "You must install Process Control; try something like:\n\t$$ brew install bww/stable/psctl\nor download it from:\n\t➡ https://github.com/bww/psctl/releases" && exit 1)

$(TARGET_DIR)/bin/$(NAME): $(SRC)
	(cd src && go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $@ $(MAIN))

.PHONY: build
build: $(TARGET_DIR)/bin/$(NAME) ## Build the product

$(PACKAGE): $(TARGET_DIR)/bin/$(NAME)
	(cd $(BUILD_DIR) && tar -zcf $(ARCHIVE) $(PRODUCT))

.PHONY: package
package: $(PACKAGE)

install: build ## Build and install
	@echo "Using sudo to install; you may be prompted for a password..."
	sudo $(INSTALL) -m 0755 $(TARGET_DIR)/bin/$(NAME) $(PREFIX)/bin/

.PHONY: ci
ci: export INSTAUNIT = $(BINARY)
ci: export TEST_SUITE := $(PWD)/test/grpc
ci: export GRPC := $(GRPC)
ci: tools build ## Run integration tests
	(cd $(GRPC) && make build) && $(PSCTL) --file test/grpc/test.yml

.PHONY: test
test: ## Run tests
	(cd src && go test $(FLAGS) $(TEST_PKGS))

.PHONY: clean
clean: ## Delete the built product and any generated files
	rm -rf $(BUILD_DIR)
