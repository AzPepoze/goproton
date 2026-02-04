BINARY_NAME=goproton
INSTANCE_BIN=goproton-instance
BIN_DIR=bin
PROJECT_DIR=src/ui
WAILS_BIN=$(shell go env GOPATH)/bin/wails3

build: build-instance build-ui

build-instance:
	cd src && go build -o ../$(BIN_DIR)/$(INSTANCE_BIN) instance/*.go

build-ui: generate-bindings
	cd $(PROJECT_DIR) && PATH="$(shell go env GOPATH)/bin:$(PATH)" $(WAILS_BIN) build
	mv $(PROJECT_DIR)/bin/$(BINARY_NAME) $(BIN_DIR)/$(BINARY_NAME)

generate-bindings:
	cd $(PROJECT_DIR) && PATH="$(shell go env GOPATH)/bin:$(PATH)" $(WAILS_BIN) generate bindings

dev: build-instance
	cd $(PROJECT_DIR) && PATH="$(shell go env GOPATH)/bin:$(PATH)" $(WAILS_BIN) dev

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(PROJECT_DIR)/build/bin
