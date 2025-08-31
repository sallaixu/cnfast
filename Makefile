# Makefile for multi-platform Go builds

# 配置变量
APP_NAME := cnfast
VERSION := 1.0.0
BUILD_DIR := build
GO_PACKAGE := .

# 定义平台架构组合
PLATFORMS := linux-amd64 linux-arm64 windows-amd64 windows-arm64 darwin-amd64 darwin-arm64

# 默认目标
all: clean $(PLATFORMS)

# 清理构建目录
clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)

# 通用构建规则
$(PLATFORMS):
	@echo "Building for $@..."
	$(eval OS = $(firstword $(subst -, ,$@)))
	$(eval ARCH = $(lastword $(subst -, ,$@)))
	
	@mkdir -p $(BUILD_DIR)/$(OS)-$(ARCH)
	
	@if [ "$(OS)" = "windows" ]; then \
		CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_DIR)/$(OS)-$(ARCH)/$(APP_NAME).exe -ldflags="-s -w" $(GO_PACKAGE); \
		echo "Built: $(BUILD_DIR)/$(OS)-$(ARCH)/$(APP_NAME).exe"; \
	else \
		CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_DIR)/$(OS)-$(ARCH)/$(APP_NAME) -ldflags="-s -w" $(GO_PACKAGE); \
		echo "Built: $(BUILD_DIR)/$(OS)-$(ARCH)/$(APP_NAME)"; \
	fi

# 平台组目标
linux: linux-amd64 linux-arm64
win: windows-amd64 windows-arm64
arm: linux-arm64 windows-arm64
amd: linux-amd64 windows-amd64
darwin: darwin-amd64 darwin-arm64

# 构建当前本地平台
local:
	@echo "Building for local platform..."
	@mkdir -p $(BUILD_DIR)/local
	go build -o $(BUILD_DIR)/local/$(APP_NAME) $(GO_PACKAGE)
	@echo "Built: $(BUILD_DIR)/local/$(APP_NAME)"

# 创建压缩包
package: all
	@echo "Creating packages..."
	@for platform in $(PLATFORMS); do \
		cd $(BUILD_DIR)/$$platform && \
		if [ "$${platform%%-*}" = "windows" ]; then \
			zip -r ../$(APP_NAME)-$$platform-$(VERSION).zip .; \
			echo "Created: $(BUILD_DIR)/$(APP_NAME)-$$platform-$(VERSION).zip"; \
		else \
			tar -czf ../$(APP_NAME)-$$platform-$(VERSION).tar.gz .; \
			echo "Created: $(BUILD_DIR)/$(APP_NAME)-$$platform-$(VERSION).tar.gz"; \
		fi; \
		cd ../..; \
	done

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  all          - Build for all platforms: $(PLATFORMS)"
	@echo "  clean        - Clean build directory"
	@echo "  linux        - Build for Linux (amd64 + arm64)"
	@echo "  win          - Build for Windows (amd64 + arm64)"
	@echo "  arm          - Build for ARM architectures (linux-arm64 + windows-arm64)"
	@echo "  amd          - Build for AMD64 architectures (linux-amd64 + windows-amd64)"
	@echo "  local        - Build for current local platform"
	@echo "  package      - Build and create compressed packages"
	@echo "  help         - Show this help message"

.PHONY: all clean $(PLATFORMS) linux win arm amd local package help