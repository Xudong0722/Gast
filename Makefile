BINARY_NAME=gast
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

.PHONY: all build clean install uninstall test fmt vet

# 默认目标
all: build

# 编译
build:
	@echo "正在编译 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 安装到系统
install: build
	@echo "正在安装 $(BINARY_NAME) 到 $(INSTALL_DIR)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "安装完成，现在可以在任何位置使用 '$(BINARY_NAME)' 命令"

# 卸载
uninstall:
	@echo "正在卸载 $(BINARY_NAME)..."
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "卸载完成"

# 清理
clean:
	@echo "正在清理..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	@echo "清理完成"

# 运行测试
test:
	@echo "正在运行测试..."
	go test -v ./...

# 格式化代码
fmt:
	@echo "正在格式化代码..."
	go fmt ./...

# 静态检查
vet:
	@echo "正在进行静态检查..."
	go vet ./...

# 发布版本 (交叉编译)
release:
	@echo "正在构建发布版本..."
	@mkdir -p $(BUILD_DIR)/release
	# Linux amd64
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 .
	# Linux arm64
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 .
	# Windows amd64
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe .
	# macOS amd64
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 .
	# macOS arm64
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 .
	@echo "发布版本构建完成，文件位于 $(BUILD_DIR)/release/"

# 显示帮助
help:
	@echo "可用的make目标:"
	@echo "  build     - 编译程序"
	@echo "  install   - 安装到系统"
	@echo "  uninstall - 从系统卸载"
	@echo "  clean     - 清理编译文件"
	@echo "  test      - 运行测试"
	@echo "  fmt       - 格式化代码"
	@echo "  vet       - 静态检查"
	@echo "  release   - 构建发布版本"
	@echo "  help      - 显示此帮助" 