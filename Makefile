AVRDUDE ?= avrdude
PARTNO ?= atmega328p
TARGET ?= arduino

ifeq ($(SHELL),sh.exe)
BUILD_DIR ?= .build
else
# Assume a Linux system
TOP_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR ?= $(TOP_DIR)/.build

.PHONY: all
all: go ## All targets

# Help target (from https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html)
.PHONY: help
help: ## Display this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

TINYGO ?= tinygo
GO_VERSION := $(shell grep 'go ' go.mod | cut -f 2 -d ' ')
GO_COMMON_SOURCES := go.mod go.sum
GO_PROJECTS := $(shell find ./ -maxdepth 1 -type d -regex '\./[^.].*')
HEX_ARTIFACTS := $(patsubst %,$(BUILD_DIR)/hex/%.hex,$(GO_PROJECTS))

.PHONY: go
go: go-fmt $(HEX_ARTIFACTS) ## Build go code

.PHONY: go-mod-tidy
go-mod-tidy: ## Go mod tidy
	go mod tidy --compat=$(GO_VERSION)

GO_FMT_STAMP := $(BUILD_DIR)/stamps/go-fmt
.PHONY: go-fmt
go-fmt: $(GO_FMT_STAMP) ## Go format check
$(GO_FMT_STAMP): $(GO_SOURCES)
	mkdir -p $(@D)
	@fixformat=$$(gofmt -l ./); \
		[ -z "$$fixformat" ] && exit 0; \
		echo "Files with wrong formatting:"; for fn in $$fixformat; do echo "  $$fn: "; gofmt -d $$fn; done; exit 1
	touch $@

$(BUILD_DIR)/hex/%.hex: %/* $(GO_COMMON_SOURCES) $(GO_FMT_STAMP)
	mkdir -p $(@D)
	$(TINYGO) build -o $@ -target $(TARGET) ./$*

.PHONY: clean
clean: ## Clean build area
	rm -rf $(BUILD_DIR)/

endif # ifeq ($(SHELL),sh.exe)

.PHONY: flash-xyz
flash-xyz: ## Flash project <xyz> (env. var. PORT has to be defined)
	@echo $(AVRDUDE) -p $(PARTNO) -c $(TARGET) -P $$PORT -D -U flash:w:$(BUILD_DIR)/hex/xyz.hex:i

ifneq ($(PORT),)
flash-%: $(BUILD_DIR)/hex/%.hex
	$(AVRDUDE) -p $(PARTNO) -c $(TARGET) -P $(PORT) -D -U flash:w:$<:i
endif




