APP ?= cronweb
BUILD_DIR ?= runtime
CONF ?= ./conf/base.json

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        Build $(APP) into $(BUILD_DIR)"
	@echo "  run-node     Run cronnode with the local configuration"
	@echo "  run-web      Run cronweb with the local configuration"
	@echo "  tidy         Run gofmt on Go sources"

.PHONY: build
build:
        mkdir -p $(BUILD_DIR)
        GO111MODULE=on go build -o $(BUILD_DIR)/$(APP) ./cmd/$(APP)

.PHONY: run-node
run-node:
	GO111MODULE=on go run ./cmd/cronnode --conf=$(CONF)

.PHONY: run-web
run-web:
	GO111MODULE=on go run ./cmd/cronweb --conf=$(CONF)

.PHONY: tidy
tidy:
	gofmt -w $$(find . -name '*.go' -not -path './web/ui/node_modules/*')
