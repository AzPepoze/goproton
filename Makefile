BINARY_NAME=goproton
INSTANCE_BIN=goproton-instance
UI_BIN=goproton-ui
BIN_DIR=bin
PROJECT_DIR=ui
WAILS_BIN=$(shell go env GOPATH)/bin/wails

build: build-instance build-ui build-wrapper

build-instance:
	go build -o $(BIN_DIR)/$(INSTANCE_BIN) cmd/instance/*.go

build-ui: generate-bindings
	cd $(PROJECT_DIR) && $(WAILS_BIN) build
	mv $(PROJECT_DIR)/build/bin/$(UI_BIN) $(BIN_DIR)/$(UI_BIN)

generate-bindings:
	cd $(PROJECT_DIR) && $(WAILS_BIN) generate bindings

build-wrapper:
	go build -o $(BIN_DIR)/$(BINARY_NAME) cmd/goproton/*.go

dev: build
	cd $(PROJECT_DIR) && $(WAILS_BIN) dev

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(PROJECT_DIR)/build/bin
