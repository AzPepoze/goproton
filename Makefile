BINARY_NAME=goproton
INSTANCE_BIN=goproton-instance
UI_BIN=goproton-ui
PROJECT_DIR=ui
WAILS_BIN=$(shell go env GOPATH)/bin/wails

build: build-instance build-ui build-wrapper

build-instance:
	go build -o $(INSTANCE_BIN) cmd/instance/main.go

build-ui:
	cd $(PROJECT_DIR) && $(WAILS_BIN) build
	mv $(PROJECT_DIR)/build/bin/$(UI_BIN) $(UI_BIN)

build-wrapper:
	go build -o $(BINARY_NAME) cmd/goproton/main.go

run: build
	cd $(PROJECT_DIR) && $(WAILS_BIN) dev

clean:
	rm -f $(BINARY_NAME) $(INSTANCE_BIN) $(UI_BIN)
	rm -rf $(PROJECT_DIR)/build/bin
