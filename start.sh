#!/bin/bash

# OverstepLab - 靶场服务启动脚本 (Linux/macOS)
# 支持双击运行或通过终端执行

set -euo pipefail

# 设置 UTF-8 编码
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 清屏
clear 2>/dev/null || printf '\033c'

printf "${CYAN}============================================${NC}\n"
printf "${CYAN}   OverstepLab - 靶场服务启动脚本${NC}\n"
printf "${CYAN}============================================${NC}\n"
printf "\n"

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
FRONTEND_DIR="${SCRIPT_DIR}/src/frontend"
BACKEND_DIR="${SCRIPT_DIR}/src/backend"
WEB_DIST_DIR="${BACKEND_DIR}/internal/web/dist"
EXE_NAME="oversteplab"
EXE_PATH="${BACKEND_DIR}/${EXE_NAME}"

# 检测 Go 环境
printf "${YELLOW}[1/5] 正在检测 Go 编译环境...${NC}\n"
if ! command -v go &> /dev/null; then
    printf "${RED}[错误] 未检测到 Go 环境，请安装 Go 1.25+ 并添加到 PATH。${NC}\n"
    printf "${RED}下载地址: https://go.dev/dl/${NC}\n"
    printf "\n"
    printf "按 Enter 键退出..."
    read -r
    exit 1
fi

GO_VERSION=$(go version 2>/dev/null | awk '{print $3}')
printf "${GREEN}[信息] 检测到 Go 版本: ${GO_VERSION}${NC}\n"

# 检测 Node.js / pnpm 环境
printf "\n"
printf "${YELLOW}[2/5] 正在检测前端构建环境...${NC}\n"
HAS_PNPM=0
HAS_NPM=0
if command -v pnpm &> /dev/null; then
    HAS_PNPM=1
fi
if command -v npm &> /dev/null; then
    HAS_NPM=1
fi

if [ "$HAS_PNPM" -eq 0 ] && [ "$HAS_NPM" -eq 0 ]; then
    printf "${RED}[错误] 未检测到 pnpm 或 npm，请安装 Node.js 和 pnpm。${NC}\n"
    printf "${RED}下载地址: https://nodejs.org/  安装后执行: npm install -g pnpm${NC}\n"
    printf "\n"
    printf "按 Enter 键退出..."
    read -r
    exit 1
fi

# 编译前端
printf "\n"
printf "${YELLOW}[3/5] 正在编译前端...${NC}\n"
if [ ! -d "$FRONTEND_DIR" ]; then
    printf "${RED}[错误] 无法找到前端目录: ${FRONTEND_DIR}${NC}\n"
    printf "\n"
    printf "按 Enter 键退出..."
    read -r
    exit 1
fi

cd "$FRONTEND_DIR"

if [ ! -d "node_modules" ]; then
    printf "${CYAN}[信息] 首次运行，正在安装前端依赖...${NC}\n"
    if [ "$HAS_PNPM" -eq 1 ]; then
        pnpm install
    else
        npm install
    fi
    if [ $? -ne 0 ]; then
        printf "${RED}[错误] 前端依赖安装失败。${NC}\n"
        printf "\n"
        printf "按 Enter 键退出..."
        read -r
        exit 1
    fi
fi

if [ "$HAS_PNPM" -eq 1 ]; then
    pnpm build
else
    npm run build
fi
if [ $? -ne 0 ]; then
    printf "${RED}[错误] 前端编译失败。${NC}\n"
    printf "\n"
    printf "按 Enter 键退出..."
    read -r
    exit 1
fi

# 复制前端产物到后端嵌入目录
if [ ! -d "$WEB_DIST_DIR" ]; then
    mkdir -p "$WEB_DIST_DIR"
fi
if [ -d "dist" ]; then
    cp -r dist/* "$WEB_DIST_DIR/" 2>/dev/null || printf "${YELLOW}[警告] 复制前端产物到后端目录失败，尝试继续...${NC}\n"
fi

printf "${GREEN}[信息] 前端编译完成。${NC}\n"

# 进入后端目录
if [ ! -d "$BACKEND_DIR" ]; then
    printf "${RED}[错误] 无法找到后端目录: ${BACKEND_DIR}${NC}\n"
    printf "\n"
    printf "按 Enter 键退出..."
    read -r
    exit 1
fi

cd "$BACKEND_DIR"

# 编译后端
printf "\n"
printf "${YELLOW}[4/5] 正在编译后端服务...${NC}\n"
if [ -f "$EXE_NAME" ]; then
    rm -f "$EXE_NAME"
fi

if ! go build -ldflags "-s -w" -o "$EXE_NAME" ./cmd/server/; then
    printf "${RED}[错误] 编译失败，请检查代码或依赖。${NC}\n"
    printf "\n"
    printf "按 Enter 键退出..."
    read -r
    exit 1
fi

printf "${GREEN}[信息] 编译成功: ${EXE_NAME}${NC}\n"

# 启动服务
printf "\n"
printf "${YELLOW}[5/5] 正在启动 OverstepLab 服务...${NC}\n"
printf "${GREEN}[信息] 默认访问地址: http://localhost:8080${NC}\n"
printf "${GREEN}[信息] 按 Ctrl+C 停止服务${NC}\n"
printf "\n"
printf "${CYAN}============================================${NC}\n"
printf "${CYAN}  靶场已启动！请在浏览器中访问:${NC}\n"
printf "${CYAN}  http://localhost:8080${NC}\n"
printf "${CYAN}============================================${NC}\n"
printf "\n"

export GIN_MODE="release"

# 启动服务并等待
"${EXE_PATH}" &
SERVER_PID=$!

# 捕获 Ctrl+C 信号进行优雅退出
cleanup() {
    printf "\n"
    printf "${CYAN}============================================${NC}\n"
    printf "${YELLOW}[信息] 正在停止服务...${NC}\n"
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
    printf "${YELLOW}[信息] 服务已停止。${NC}\n"
    printf "${CYAN}============================================${NC}\n"
    printf "\n"
    printf "按 Enter 键退出..."
    read -r
    exit 0
}

trap cleanup SIGINT SIGTERM

# 等待服务进程
wait "$SERVER_PID"
EXIT_CODE=$?

# 服务停止后的提示
printf "\n"
printf "${CYAN}============================================${NC}\n"
printf "${YELLOW}[信息] 服务已停止 (退出码: ${EXIT_CODE})。${NC}\n"
printf "${CYAN}============================================${NC}\n"
printf "\n"
printf "按 Enter 键退出..."
read -r
