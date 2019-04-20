
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

.PHONY: all test clean install release build package formula

all: build

gate:
	@echo && echo "AWS Profile: $(AWS_PROFILE)" && echo "    Version: $(VERSION)" && echo "     Branch: $(BRANCH)"
	@echo && read -p "Release version $(VERSION)? [y/N] " -r continue && echo && [ "$${continue:-N}" = "y" ]

$(TARGET_DIR)/bin/$(NAME): $(SRC)
	(cd src && go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $@ $(MAIN))

build: $(TARGET_DIR)/bin/$(NAME) ## Build the product

$(PACKAGE): $(TARGET_DIR)/bin/$(NAME)
	(cd $(BUILD_DIR) && tar -zcf $(ARCHIVE) $(PRODUCT))

package: $(PACKAGE)

publish: package
	aws s3 cp --acl public-read $(PACKAGE) $(ARTIFACTS)/$(VERSION)/$(ARCHIVE)

formula: publish
	$(PWD)/build/update-formula -v $(VERSION) -o $(TARGET_DIR)/$(NAME).rb $(PACKAGE)
	aws s3 cp --acl public-read $(TARGET_DIR)/$(NAME).rb $(ARTIFACTS)/$(LATEST)/$(NAME).rb
	aws s3 cp --acl public-read $(TARGET_DIR)/$(NAME).rb $(ARTIFACTS)/$(VERSION)/$(NAME).rb
	@echo "\nHomebrew formula for version $(VERSION):\n\thttps://instaunit.s3.amazonaws.com/releases/$(LATEST)/$(NAME).rb"

release: gate test ## Build for all supported architectures
	make publish PRODUCT=$(NAME)-$(VERSION)-linux-amd64 GOOS=linux GOARCH=amd64
	make publish PRODUCT=$(NAME)-$(VERSION)-freebsd-amd64 GOOS=freebsd GOARCH=amd64
	make formula PRODUCT=$(NAME)-$(VERSION)-darwin-amd64 GOOS=darwin GOARCH=amd64
	@echo && echo "Tag this release:\n\t$ git commit -a -m \"Version $(VERSION)\" && git tag $(VERSION)" && echo

install: build ## Build and install
	install -m 0755 $(TARGET_DIR)/bin/$(NAME) $(PREFIX)/bin/

test: ## Run tests
	(cd src && go test $(TEST_PKGS))

clean: ## Delete the built product and any generated files
	rm -rf $(BUILD_DIR)
