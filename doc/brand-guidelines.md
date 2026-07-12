# OverstepLab 品牌视觉识别系统使用规范

## 概述

本文档定义了 OverstepLab 越权测试靶场项目的视觉识别系统（VIS）使用规范，确保品牌在所有触点上的一致性和专业性。

---

## 1. Logo 使用规范

### 1.1 Logo 版本

| 版本 | 文件名 | 使用场景 |
|------|--------|----------|
| 图标版 | `logo-icon.svg` | 应用图标、头像、小尺寸展示 |
| 横排完整版 | `logo-full-horizontal.svg` | 网站头部、导航栏、名片 |
| 竖排完整版 | `logo-full-vertical.svg` | 社交媒体头像、竖向布局 |
| 纯文字版 | `logo-wordmark.svg` | 极简场景、水印 |
| 单色深 | `logo-icon-mono-dark.svg` | 浅色背景 |
| 单色浅 | `logo-icon-mono-light.svg` | 深色背景 |

### 1.2 清除空间

Logo 周围必须保持最小清除空间，等于图标高度的 1/4。任何文字、图形或其他元素不得侵入此区域。

```
+----------------------------+
|     [清除空间]             |
|  +--------------------+    |
|  |                    |    |
|  |      LOGO          |    |
|  |                    |    |
|  +--------------------+    |
|     [清除空间]             |
+----------------------------+
```

### 1.3 最小尺寸

| 场景 | 最小尺寸 |
|------|----------|
| 数字端图标 | 32x32 px |
| 数字端完整Logo | 120x40 px |
| 印刷端 | 20x20 mm |

### 1.4 正确/错误使用示例

**正确用法：**
- 在深色背景上使用全彩或单色浅版本
- 在浅色背景上使用单色深版本
- 保持原始比例，不拉伸变形

**错误用法：**
- 改变 Logo 颜色为品牌色以外的颜色
- 旋转、倾斜或拉伸 Logo
- 在复杂背景上使用无背景版本（缺乏对比度）
- 添加阴影、描边或其他效果
- 将图标与文字分离后单独使用文字部分

---

## 2. 色彩规范

### 2.1 品牌主色

| 角色 | 色值 (HEX) | 色值 (RGB) | 使用场景 |
|------|------------|------------|----------|
| Indigo 500 | `#4f46e5` | rgb(79, 70, 229) | 品牌主色、按钮、链接 |
| Indigo 400 | `#818cf8` | rgb(129, 140, 248) | 深色主题主色、高亮 |
| Indigo 600 | `#4338ca` | rgb(67, 56, 202) | 渐变终点、按下状态 |
| Indigo 300 | `#a5b4fc` | rgb(165, 180, 252) | 装饰、次要元素 |

### 2.2 功能色

| 角色 | 色值 (HEX) | 使用场景 |
|------|------------|----------|
| 成功绿 | `#10b981` / `#34d399` | 安全状态、通过、正常 |
| 警告黄 | `#f59e0b` / `#fbbf24` | 提示、注意、警告 |
| 危险红 | `#ef4444` / `#fb7185` | 漏洞、错误、越权 |
| 信息蓝 | `#3b82f6` / `#60a5fa` | 信息提示、IDOR标签 |

### 2.3 背景色

| 角色 | 色值 (HEX) | 使用场景 |
|------|------------|----------|
| 深色背景 | `#0a0a0a` / `#0c0c0c` | 深色模式主背景 |
| 浅色背景 | `#f8fafc` / `#ffffff` | 浅色模式主背景 |
| 表面色 | `#161616` | 深色模式卡片背景 |
| 边框色 | `#2a2a2a` / `#e2e8f0` | 深色/浅色模式边框 |

### 2.4 文字色

| 角色 | 色值 (HEX) | 使用场景 |
|------|------------|----------|
| 主文字（暗） | `#f0f0f0` | 深色模式标题正文 |
| 次文字（暗） | `#a0a0a0` | 深色模式次要文字 |
| 主文字（亮） | `#0f172a` | 浅色模式标题正文 |
| 次文字（亮） | `#64748b` | 浅色模式次要文字 |

---

## 3. 字体规范

| 层级 | 字体 | 字重 | 用途 |
|------|------|------|------|
| 标题 | Inter | 700 (Bold) | 页面标题、Logo文字 |
| 正文 | Inter | 400/500 | 正文、描述 |
| 代码/技术 | JetBrains Mono | 400/500/600 | 代码块、技术标签、数据 |

---

## 4. 图标使用指南

### 4.1 核心功能图标

所有图标位于 `src/frontend/src/assets/icons/`，为 24x24 像素的 SVG 格式。

| 图标 | 文件名 | 用途 | 默认颜色 |
|------|--------|------|----------|
| 安全模式 | `icon-shield-check.svg` | 安全模式指示 | `#34d399` |
| 漏洞模式 | `icon-shield-cross.svg` | 漏洞模式指示 | `#fb7185` |
| API Key | `icon-key.svg` | API Key管理 | `#818cf8` |
| 越权突破 | `icon-lock-open.svg` | 越权/突破 | `#fbbf24` |
| 权限层级 | `icon-hierarchy.svg` | 权限层级 | `#818cf8` |
| 漏洞标记 | `icon-bug.svg` | 漏洞标记 | `#fb7185` |
| 终端 | `icon-terminal.svg` | 命令/控制台 | `#818cf8` |
| 身份认证 | `icon-fingerprint.svg` | 身份/认证 | `#818cf8` |
| 隐藏/IDOR | `icon-eye-crossed.svg` | 隐藏/IDOR | `#818cf8` |
| 越权动作 | `icon-arrow-jump.svg` | 越权动作 | `#fbbf24` |

### 4.2 空状态插画

| 插画 | 文件名 | 用途 | 尺寸 |
|------|--------|------|------|
| 无数据 | `empty-state-no-data.svg` | 无数据状态 | 200x200 |
| 权限不足 | `empty-state-no-permission.svg` | 权限不足状态 | 200x200 |
| 搜索无结果 | `empty-state-search.svg` | 搜索无结果 | 200x200 |

---

## 5. 文件目录结构

```
OverstepLab/
├── src/frontend/public/
│   ├── logo/
│   │   ├── logo-icon.svg                    # 图标版 (512x512)
│   │   ├── logo-full-horizontal.svg         # 横排完整版
│   │   ├── logo-full-vertical.svg           # 竖排完整版
│   │   ├── logo-wordmark.svg              # 纯文字版
│   │   ├── logo-icon-mono-dark.svg        # 单色深 (浅色背景)
│   │   └── logo-icon-mono-light.svg       # 单色浅 (深色背景)
│   ├── favicon.ico                          # 浏览器图标
│   ├── favicon-32x32.png                  # 32x32 PNG
│   ├── apple-touch-icon.png               # 180x180 iOS
│   ├── android-chrome-192x192.png         # Android 192
│   ├── android-chrome-512x512.png         # Android 512
│   └── mstile-150x150.png                 # Windows磁贴
│
├── src/frontend/src/assets/icons/
│   ├── icon-shield-check.svg
│   ├── icon-shield-cross.svg
│   ├── icon-key.svg
│   ├── icon-lock-open.svg
│   ├── icon-hierarchy.svg
│   ├── icon-bug.svg
│   ├── icon-terminal.svg
│   ├── icon-fingerprint.svg
│   ├── icon-eye-crossed.svg
│   ├── icon-arrow-jump.svg
│   ├── empty-state-no-data.svg
│   ├── empty-state-no-permission.svg
│   └── empty-state-search.svg
│
├── assets/
│   ├── brand/
│   │   ├── svg/                          # 品牌SVG源文件备份
│   │   └── png/                          # 品牌PNG导出文件
│   ├── posters/
│   │   ├── poster-social-1200x630.svg    # 社交媒体分享图
│   │   ├── poster-banner-1500x500.svg    # GitHub横幅
│   │   └── poster-square-1080x1080.svg   # 方形海报
│   └── doc-images/
│       ├── readme-header.svg             # README头图
│       ├── bg-grid-dark.svg              # 深色网格背景
│       ├── bg-gradient-spot.svg          # 渐变光斑分隔
│       └── bg-shield-pattern.svg         # 盾牌水印图案
│
└── doc/
    └── brand-guidelines.md               # 本文档
```

---

## 6. 技术栈标签色

在 README 和文档中引用技术栈时，使用以下颜色：

| 技术 | 背景色 | 文字色 |
|------|--------|--------|
| Go | `#00ADD8` | `#ffffff` |
| Vue | `#4FC08D` | `#ffffff` |
| Gin | `#4f46e5` | `#ffffff` |
| JWT | `#4f46e5` | `#ffffff` |
| SQLite | `#003B57` | `#ffffff` |
| Docker | `#2496ED` | `#ffffff` |

---

## 7. 宣传海报使用

| 海报 | 尺寸 | 用途 |
|------|------|------|
| `poster-social-1200x630.svg` | 1200x630 | GitHub/Twitter/微博分享 |
| `poster-banner-1500x500.svg` | 1500x500 | GitHub仓库顶部横幅 |
| `poster-square-1080x1080.svg` | 1080x1080 | Instagram/小红书/公众号 |

---

## 8. 更新记录

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2026-07-12 | v1.0 | 初始版本，包含Logo、图标、海报、文档配图完整规范 |

---

*本规范由 OverstepLab 项目维护团队制定，所有项目贡献者应遵循本规范进行品牌相关设计工作。*
