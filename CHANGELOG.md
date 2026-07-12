# Changelog

本项目遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/) 格式，版本号基于 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

## [0.3.1] - 2026-07-12

### 修复
- 重写前端 MD5 算法实现（`crypto.ts`），替换为更标准可靠的版本，修复哈希计算正确性问题

## [0.3.0] - 2026-07-12

### 新增
- 编码加密越权挑战扩展：新增 6 个编码/加密相关漏洞挑战（E-07 至 E-12）
  - E-07 Base64 编码绕过权限校验
  - E-08 URL 编码绕过参数过滤
  - E-09 哈希签名验证绕过
  - E-10 JWT 弱密钥破解
  - E-11 AES-CBC 位翻转攻击
  - E-12 经典密码（凯撒/维吉尼亚）编码绕过
- 后端新增 `internal/crypto/` 密码学工具模块（经典密码、现代加密、编码、国密 SM2/SM3/SM4、大整数运算、签名验证）
- 新增 `encoded_handler.go` 处理编码挑战 API 端点
- 新增 `encoding.go` 中间件支持多种编码格式自动检测与解码
- 前端新增编码挑战状态管理 (`encodingChallenge` store)
- 前端新增 `tools/` 视图（编码转换、加解密工具箱）
- 前端新增 `crypto.ts` / `encoded.ts` API 客户端
- 新增教程文档 `07-encoding-crypto.md`

### 变更
- API 客户端 (`client.ts`) 重构以支持编码挑战的请求拦截与自动编码转换
- 路由 (`router.go`) 扩展编码挑战相关端点
- 挑战元数据 (`challenges.go`) 大幅扩展以支持编码/加密类挑战
- 前端导航栏新增"工具箱"入口
- PRD 文档和教程索引同步更新

### 修复
- 前端多处视图细节修正（VPS 详情、工单列表/详情、Dashboard、挑战页）

## [0.2.2] - 2026-06-29

### 修复
- 修复 Dockerfile 前端构建阶段缺少 pnpm-lock.yaml 导致 `pnpm install --frozen-lockfile` 失败
- 修复 gen-release-notes.sh 脚本在 GitHub Actions 环境下的兼容性问题（替换关联数组为普通变量）
- 添加 .gitattributes 强制 shell 脚本使用 LF 换行符

## [0.2.1] - 2026-06-29

### 修复
- 修复 Dockerfile 中 Go 版本与 go.mod 不匹配（1.21 → 1.25）

## [0.2.0] - 2026-06-29

### 新增
- GitHub Actions 自动发布工作流，支持 tag 推送和手动触发
- 自动构建多平台可执行文件：Linux amd64、Linux arm64、Windows amd64
- Docker 镜像自动构建并推送至 GitHub Container Registry (ghcr.io)
- 源码 zip 压缩包和 SHA256SUMS 校验文件自动生成
- 基于 Conventional Commits 的中文 Release Notes 自动生成脚本

### 变更
- 统一所有文档、配置、脚本中的服务端口为 `5000`（原为 `8080`）
  - 涉及文件：Dockerfile、docker-compose.yml、README、start.bat/sh/ps1、教程文档

## [0.1.2] - 2026-06-28

### 修复
- 修复前后端联调中的多项功能问题（API 路由、数据展示、权限校验）

## [0.1.1] - 2026-06-27

### 修复
- 修复后端 API 和服务层多处逻辑问题（认证、VPS 管理、工单模块）

## [0.1.0] - 2026-06-26

### 新增
- 项目初始化，搭建 Go + Vue 3 全栈架构
- 13 个预置越权漏洞挑战（含 horizontal/vertical/context escalation）
- 管理后台完整功能：用户管理、公司管理、系统配置
- 公告和公告管理模块
- VPS 管理、订单管理、账单管理、工单系统
- JWT 认证与 RBAC 权限控制（含可绕过设计）
- 安全模式开关（`OVERSTEPLAB_SAFE_MODE` 环境变量）
- SQLite 数据库与 GORM 自动迁移
- TailwindCSS + PrimeVue 4 前端 UI 框架
- 一键启动脚本（Windows CMD/PowerShell、Linux/macOS Shell）

[0.2.2]: https://github.com/GitHubNull/OverstepLab/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/GitHubNull/OverstepLab/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/GitHubNull/OverstepLab/compare/v0.1.2...v0.2.0
[0.1.2]: https://github.com/GitHubNull/OverstepLab/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/GitHubNull/OverstepLab/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/GitHubNull/OverstepLab/releases/tag/v0.1.0
