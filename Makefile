.PHONY: build run restart clean test tidy generate

APP_NAME := git-sync-service
BUILD_DIR := ./output

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .

# 本地开发启动:自动加载 .env(ENCRYPTION_KEY 等本地密钥,已被 gitignore)。
# 注意必须 `go run .` 整包编译,不能 `go run main.go`(后者只编译单文件,会报 undefined: register)
run:
	@if [ -f .env ]; then set -a; . ./.env; set +a; fi; go run .

# 一键重启:先停掉占用 8890 的旧后端(连同 go run 包装进程),再以当前代码重新编译启动。
# 避免改完代码(如版本号)后旧进程仍跑旧二进制,导致行为/版本不更新
restart:
	@pid=$$(lsof -ti tcp:8890 2>/dev/null); \
	if [ -n "$$pid" ]; then \
		echo ">> stopping old backend (pid $$pid) ..."; \
		pkill -f "go run \.$$" 2>/dev/null || true; \
		kill $$pid 2>/dev/null || true; \
		sleep 2; \
		kill -9 $$pid 2>/dev/null || true; \
	fi
	@$(MAKE) run

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
