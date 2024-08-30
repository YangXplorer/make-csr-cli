# 定义变量
SRC_DIR := ./src
BIN_DIR := ./bin
APP_NAME := makeCsr

# 目标操作系统和架构
OS_ARCH := linux/amd64 darwin/amd64 windows/amd64

# 默认目标
.PHONY: all
all: $(OS_ARCH)

# 为每个平台生成二进制文件
$(OS_ARCH):
	@mkdir -p $(BIN_DIR)
	GOOS=$(word 1, $(subst /, ,$@)) GOARCH=$(word 2, $(subst /, ,$@)) \
		go build -o $(BIN_DIR)/$(APP_NAME)-$(word 1, $(subst /, ,$@))-$(word 2, $(subst /, ,$@))$(if $(findstring windows,$(word 1, $(subst /, ,$@))),.exe) $(SRC_DIR)/main.go

# 清理生成的文件
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/*
