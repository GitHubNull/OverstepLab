> ⚠️ **安全警告**: OverstepLab 是一个**故意包含安全漏洞**的应用程序，仅供教育和学习目的使用。请勿将其部署到生产环境或暴露在公网。详见 [SECURITY.md](SECURITY.md)。

# OverstepLab - 越权测试靶场

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D.svg)](https://vuejs.org/)

**[English README](README_EN.md)**

## 项目简介

OverstepLab 是一个**以越权漏洞为核心的安全测试靶场**，通过模拟一个多用户 VPS 云管理平台（CloudNest），内置多种典型越权漏洞场景，供安全从业者进行漏洞挖掘练习与安全开发学习。

## 核心特性

- 🎯 **22 个漏洞挑战场景**：
  - 13 个传统越权漏洞（水平越权 IDOR、垂直越权、上下文越权）
  - 9 个编码加密挑战（Base64/Base32/凯撒/自定义编码/多层嵌套/AES/HMAC/SM4/MD5 签名）
- 🏢 **真实业务场景**：模拟多用户 VPS 管理平台，包含企业/个人两种用户类型
- 🔐 **编码加密挑战系统**：9 个编码/加密相关越权挑战（Base64/Base32/凯撒/自定义编码/多层嵌套/AES/HMAC/SM4/MD5 签名），前端自动编码请求参数
- 🔄 **安全/漏洞模式切换**：一键切换，对比学习与验证
- 💡 **渐进式提示系统**：每个漏洞提供 3 级提示
- 📋 **完整 WriteUp**：详细的漏洞分析和修复方案
- 🔧 **一键重置数据库**：随时恢复初始状态
- 📦 **单二进制部署**：Go 编译嵌入前端资源，无需额外部署
- 🚀 **跨平台一键启动脚本**：Windows (.bat/.ps1)、Linux/macOS (.sh) 双击即运行
- 🐳 **Docker 支持**：一键容器化部署
- 🌙 **主题切换**：支持浅色/深色/跟随系统三种模式

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.25+ / Gin / GORM / SQLite (纯 Go) / JWT |
| 前端 | Vue 3 / PrimeVue 4 / Pinia / Vite / TailwindCSS 4 |
| 数据库 | SQLite3 (嵌入式) |
| 密码学 | golang.org/x/crypto / emmansun/gmsm (SM2/SM3/SM4) |

## 快速开始

### 方式一：一键启动脚本（推荐）

最简单的启动方式，自动完成前端构建、后端编译和服务启动。

**Windows：**
```powershell
# 双击运行
start.bat
# 或
start.ps1
```

**Linux / macOS：**
```bash
# 赋予执行权限后运行
chmod +x start.sh
./start.sh
```

启动后访问 `http://localhost:5000`。

### 方式二：开发模式

需要两个终端分别运行前后端。

```bash
# 终端 1: 启动后端
cd src/backend && go mod tidy && go run ./cmd/server/main.go

# 终端 2: 启动前端
cd src/frontend && pnpm install && pnpm dev
```

### 方式三：生产模式 (单二进制)

```bash
# 构建
make build

# 运行
./bin/oversteplab
# 访问 http://localhost:5000
```

### 方式四：Docker 部署

```bash
docker-compose up --build
# 访问 http://localhost:5000
```

## 预置测试账户

| 用户名 | 密码 | 类型 | 角色 | 所属企业 |
|--------|------|------|------|----------|
| admin | admin123 | 平台管理员 | platform_admin | - |
| acme_admin | pass123 | 企业用户 | admin | Acme Corp |
| acme_ops | pass123 | 企业用户 | operator | Acme Corp |
| acme_finance | pass123 | 企业用户 | finance | Acme Corp |
| acme_viewer | pass123 | 企业用户 | viewer | Acme Corp |
| globex_admin | pass123 | 企业用户 | admin | Globex Inc |
| globex_ops | pass123 | 企业用户 | operator | Globex Inc |
| globex_finance | pass123 | 企业用户 | finance | Globex Inc |
| globex_viewer | pass123 | 企业用户 | viewer | Globex Inc |
| alice | pass123 | 个人用户 | - | - |
| bob | pass123 | 个人用户 | - | - |

## 目录结构

```
OverstepLab/
├── src/
│   ├── backend/          # Go 后端
│   │   ├── cmd/server/     # 入口
│   │   ├── internal/       # 业务代码 (handler/service/repository/model)
│   │   │   ├── crypto/     # 密码学工具模块（经典/现代/国密/签名）
│   │   │   ├── middleware/ # 中间件（JWT/RBAC/CORS/审计/编码）
│   │   │   └── vuln/       # 漏洞模式与挑战元数据
│   │   ├── database/       # 数据库层 (migration/seed)
│   │   └── router/         # 路由
│   └── frontend/           # Vue3 前端
│       └── src/
│           ├── api/        # API 封装
│           ├── stores/     # Pinia 状态管理
│           ├── views/      # 页面视图
│           └── router/     # 路由配置
├── doc/
│   ├── PRD.md              # 产品需求文档
│   └── tutorials/          # 使用教程（8 篇）
├── README.md               # 中文文档 (本文件)
├── README_EN.md            # 英文文档
├── SECURITY.md             # 安全声明
├── LICENSE                 # MIT 许可证
├── AGENTS.md               # AI Agent 指南
├── start.bat               # Windows 一键启动脚本 (CMD)
├── start.ps1               # Windows 一键启动脚本 (PowerShell)
├── start.sh                # Linux/macOS 一键启动脚本
├── Makefile
├── Dockerfile
└── docker-compose.yml
```

## 漏洞清单

### 传统越权漏洞（13 个）

| 编号 | 类别 | 场景 | 难度 |
|------|------|------|------|
| H-01 | 水平越权 | 修改 VPS ID 查看他人 VPS | ★☆☆ |
| H-02 | 水平越权 | 修改用户 ID 查看他人信息 | ★☆☆ |
| H-03 | 水平越权 | 修改订单 ID 查看他人订单 | ★☆☆ |
| H-04 | 水平越权 | 修改工单 ID 查看/回复他人工单 | ★☆☆ |
| H-05 | 水平越权 | 修改 API Key ID 删除他人 Key | ★★☆ |
| V-01 | 垂直越权 | 只读成员调用 VPS 启停接口 | ★☆☆ |
| V-02 | 垂直越权 | 运维调用用户管理接口添加成员 | ★★☆ |
| V-03 | 垂直越权 | 财务调用 VPS 重装系统接口 | ★★☆ |
| V-04 | 垂直越权 | 个人用户调用平台管理员接口 | ★★☆ |
| V-05 | 垂直越权 | 运维修改自身角色为管理员 | ★★★ |
| C-01 | 上下文越权 | A 企业用户操作 B 企业 VPS | ★★☆ |
| C-02 | 上下文越权 | 个人用户创建企业成员 | ★★★ |
| C-03 | 上下文越权 | 已吊销 API Key 仍可访问 | ★★★ |

### 编码加密挑战（9 个）

| 编号 | 类别 | 场景 | 难度 |
|------|------|------|------|
| E-01 | 编码 | Base64 编码参数绕过 | ★☆☆ |
| E-02 | 编码 | Base32 编码参数绕过 | ★☆☆ |
| E-03 | 编码 | 凯撒密码参数偏移绕过 | ★★☆ |
| E-04 | 编码 | 自定义 Base64 编码表绕过 | ★★☆ |
| E-05 | 编码 | 多层嵌套编码参数绕过 | ★★★ |
| E-06 | 加密 | AES-256-GCM 加密参数绕过 | ★★★ |
| E-07 | 签名 | HMAC-SHA256 签名验证绕过 | ★★★★ |
| E-08 | 加密 | SM4-CBC 国密加密参数绕过 | ★★★★ |
| E-09 | 签名 | MD5 哈希签名验证绕过 | ★★★ |

## 编码加密工具库

前端 `src/utils/crypto.ts` 内置密码学工具函数，主要用于编码挑战的自动请求编码，同时暴露到浏览器控制台供手动调用：

| 类别 | 支持算法 |
|------|----------|
| 编码 | Base64, Base32, Base58, Hex |
| 古典密码 | 凯撒密码 (ROT3) |
| 对称加密 | AES-256-GCM, SM4-CBC (简化实现) |
| 哈希/签名 | MD5, HMAC-SHA256 |
| 自定义编码 | 自定义 Base64 编码表, 多层嵌套编码 |

在浏览器控制台中可通过 `window.CryptoUtils` 手动调用这些工具函数。

## 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `OVERSTEPLAB_PORT` | `5000` | 服务监听端口 |
| `OVERSTEPLAB_DB_PATH` | `./oversteplab.db` | SQLite 数据库路径 |
| `OVERSTEPLAB_SAFE_MODE` | `false` | 安全模式开关（`true` 启用所有漏洞修复） |
| `OVERSTEPLAB_JWT_SECRET` | 随机生成 | JWT 签名密钥 |

## 贡献指南

欢迎提交 Issue 和 Pull Request。添加新漏洞场景时，请同时更新 PRD、教程和 WriteUp。

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE)。

## 安全声明

本项目包含**故意植入的安全漏洞**，仅供教育和学习目的使用。详情请见 [SECURITY.md](SECURITY.md)。
