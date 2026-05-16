.PHONY: build run clean test tidy generate

APP_NAME := git-sync-service
BUILD_DIR := ./output

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .

run:
	@go run main.go

tidy:
	@go mod tidy

test:
	@go test ./... -v

clean:
	@rm -rf $(BUILD_DIR) data/

generate:
	@echo "Generating code from IDL..."
	@cd idl && thriftgo --out ../biz --go --go-recurse 10 git_sync.thrift

docker-build:
	@docker build -t $(APP_NAME):latest .
