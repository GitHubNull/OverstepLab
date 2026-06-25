# 教程 00: 跨平台一键启动

## 概述

OverstepLab 提供跨平台的一键启动脚本，让你无需手动配置即可快速运行靶场环境。脚本会自动检测依赖、编译前端、构建后端并启动服务。

## 环境要求

### 必需环境

| 环境 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.21+ | 后端编译环境 |
| Node.js | 18+ | 前端构建环境 |
| pnpm 或 npm | - | 前端包管理器（脚本会自动检测） |

### 操作系统支持

- **Windows**: 支持 Windows 10/11，提供 `.bat` (CMD) 和 `.ps1` (PowerShell) 两种脚本
- **Linux**: 支持主流发行版（Ubuntu、CentOS、Debian 等）
- **macOS**: 支持 macOS 10.15+

## 使用方式

### Windows

#### 方式 A: 使用 CMD 脚本（推荐）

在项目根目录下，直接双击运行：

```
start.bat
```

或在 CMD 终端中执行：

```cmd
start.bat
```

#### 方式 B: 使用 PowerShell 脚本

在项目根目录下，右键选择"使用 PowerShell 运行"：

```powershell
start.ps1
```

或在 PowerShell 终端中执行：

```powershell
.\start.ps1
```

> **注意**: 如果 PowerShell 执行策略限制了脚本运行，可以临时调整：
> ```powershell
> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
> ```

### Linux / macOS

在项目根目录下，打开终端执行：

```bash
# 首次运行需要赋予执行权限
chmod +x start.sh

# 运行启动脚本
./start.sh
```

## 脚本执行流程

启动脚本会按以下步骤自动执行：

```
[1/5] 检测 Go 编译环境
[2/5] 检测前端构建环境 (pnpm / npm)
[3/5] 编译前端（自动安装依赖）
       └── 复制前端产物到后端嵌入目录
[4/5] 编译后端服务
[5/5] 启动 OverstepLab 服务
```

启动完成后，脚本会显示：

```
============================================
  靶场已启动！请在浏览器中访问:
  http://localhost:8080
============================================
```

## 停止服务

- **Windows**: 在运行窗口中按 `Ctrl + C`
- **Linux/macOS**: 在终端中按 `Ctrl + C`，脚本会优雅地停止服务

## 常见问题

### Q1: 提示"未检测到 Go 环境"

**原因**: 系统中未安装 Go 或 Go 未添加到 PATH。

**解决**:
1. 访问 https://go.dev/dl/ 下载对应系统的 Go 安装包
2. 安装后重新打开终端，再次运行脚本

### Q2: 提示"未检测到 pnpm 或 npm"

**原因**: 系统中未安装 Node.js。

**解决**:
1. 访问 https://nodejs.org/ 下载并安装 Node.js 18+
2. 安装完成后，执行 `npm install -g pnpm` 安装 pnpm（推荐）
3. 重新运行脚本

### Q3: 前端编译失败

**原因**: 可能是网络问题导致依赖安装失败，或前端代码存在语法错误。

**解决**:
1. 检查网络连接，确保能访问 npm  registry
2. 手动进入前端目录尝试构建：
   ```bash
   cd src/frontend
   pnpm install
   pnpm build
   ```
3. 查看具体错误信息并修复

### Q4: 端口 8080 被占用

**原因**: 系统中已有其他程序占用了 8080 端口。

**解决**:
- **Windows**: 脚本会自动尝试终止占用进程
- **Linux/macOS**: 手动查找并终止占用进程：
  ```bash
  lsof -i :8080
  kill -9 <PID>
  ```

### Q5: PowerShell 无法执行脚本

**原因**: Windows 默认的执行策略限制了 PowerShell 脚本运行。

**解决**:
```powershell
# 查看当前执行策略
Get-ExecutionPolicy

# 临时允许当前用户运行脚本
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## 与其他启动方式对比

| 启动方式 | 适用场景 | 优点 | 缺点 |
|---------|---------|------|------|
| 一键启动脚本 | 快速体验 | 全自动，零配置 | 不适合开发调试 |
| 开发模式 | 前后端开发 | 热重载，便于调试 | 需要两个终端 |
| Makefile 构建 | 生产部署 | 单二进制文件 | 需要手动执行多步 |
| Docker | 容器化部署 | 环境隔离 | 需要 Docker 环境 |

## 下一步

启动成功后，请阅读 [教程 01: 快速入门](01-quickstart.md) 了解系统功能和基本操作。
