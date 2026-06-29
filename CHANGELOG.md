# Changelog

本项目遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/) 格式，版本号基于 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

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

[0.2.0]: https://github.com/GitHubNull/OverstepLab/compare/v0.1.2...v0.2.0
[0.1.2]: https://github.com/GitHubNull/OverstepLab/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/GitHubNull/OverstepLab/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/GitHubNull/OverstepLab/releases/tag/v0.1.0
