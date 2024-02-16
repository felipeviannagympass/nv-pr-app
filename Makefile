GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
BINARY_DIR=bin
BINARY_NAME=app
CONFIG_DIR=configs
MIGRATIONS_DIR=db
GOBUILD_VARS = CGO_ENABLED=0 GOOS=linux
GOBUILD_BASE_PARAMS = -mod vendor -ldflags="-X 'main.buildVersion=$(BUILD_TAG)' -X 'main.buildTime=$(BUILD_TIME)'-X 'main.goVersion=$(GOVERSION)'"
OS=linux
ARCH=amd64
UNAME_S = $(shell uname -s)
UNAME_M = $(shell uname -m)

ifeq ($(UNAME_S),Darwin)
	OS=darwin
	PKG_CFG_PATH = PKG_CONFIG_PATH="$(shell find /opt/homebrew/Cellar/openssl@1.1 -type d -d 1 )/lib/pkgconfig/"
	GOBUILD_PARAMS = $(GOBUILD_BASE_PARAMS) -tags dynamic
	GOBUILD_VARS =  GOOS=${OS}
endif

ifeq ($(UNAME_M),arm64)
	ARCH=arm64
endif

.PHONY: bindir
bindir:
	mkdir -p $(BINARY_DIR)
	mkdir -p $(BINARY_DIR)/$(CONFIG_DIR)
	cp $(CONFIG_DIR)/sample.json $(BINARY_DIR)/$(CONFIG_DIR)/envs.json
	cp -R ${MIGRATIONS_DIR} $(BINARY_DIR)/${MIGRATIONS_DIR}

.PHONY: build
build: bindir
	$(PKG_CFG_PATH) $(GOBUILD_VARS) $(GOBUILD) $(GOBUILD_PARAMS) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/app

.PHONY: deps
deps:
	$(GOMOD) tidy
	$(GOMOD) vendor