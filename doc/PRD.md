# OverstepLab 产品需求文档（PRD）

| 属性     | 值                        |
| -------- | ------------------------- |
| 产品名称 | OverstepLab 越权测试靶场  |
| 版本     | v1.1.0                    |
| 创建日期 | 2026-04-24                |
| 更新日期 | 2026-07-12                |
| 技术栈   | Go + Vue3 + PrimeVue + TailwindCSS + SQLite3 |
| 包管理   | pnpm（前端）/ Go Modules（后端） |

---

## 1. 产品概述

### 1.1 产品背景

越权漏洞（Broken Access Control）长期位居 OWASP Top 10 首位，是 Web 应用中最常见且危害最大的安全问题之一。目前业界缺乏一个贴近真实业务场景、可供安全研究人员和开发者系统性学习与实践越权测试的靶场平台。

### 1.2 产品定位

OverstepLab 是一个**以越权漏洞为核心的安全测试靶场**，通过模拟一个多用户 VPS 云管理平台的业务场景，内置多种典型越权漏洞（水平越权、垂直越权、上下文越权），供安全从业者进行漏洞挖掘练习与安全开发学习。

### 1.3 目标用户

| 用户类型     | 说明                                       |
| ------------ | ------------------------------------------ |
| 安全研究员   | 利用靶场练习越权漏洞的发现与利用           |
| 渗透测试工程师 | 在靶场中磨练实战技巧                     |
| 开发工程师   | 通过反面案例学习安全编码规范               |
| 安全培训讲师 | 将靶场作为教学演示工具                     |

### 1.4 产品目标

- 提供一个真实感强、覆盖面广的越权漏洞练习环境
- 内置 **3 大类、22 个越权漏洞场景**（13 个传统越权 + 9 个编码加密挑战）
- 支持一键部署、开箱即用
- 提供漏洞提示（Hint）与漏洞解析（WriteUp）辅助学习
- 内置编码加密工具库（`utils/crypto.ts`），支持编码挑战的自动请求编码，同时暴露到浏览器控制台供学习调用

---

## 2. 业务场景建模

### 2.1 模拟平台：多用户 VPS 管理平台

OverstepLab 模拟一个名为 **"CloudNest"** 的虚拟 VPS 管理平台，该平台为企业和个人提供虚拟服务器的购买、管理与监控服务。

### 2.2 用户体系

系统包含以下用户类型与角色层级：

```
平台管理员 (Platform Admin)
│
├── 企业用户 (Company)
│   ├── 企业管理员 (Company Admin)
│   ├── 企业运维 (Company Operator)
│   ├── 企业财务 (Company Finance)
│   └── 企业只读成员 (Company Viewer)
│
└── 个人用户 (Individual)
```

#### 2.2.1 平台管理员（Platform Admin）

| 项目     | 说明                                       |
| -------- | ------------------------------------------ |
| 数量     | 系统内置，不可注册                         |
| 职责     | 管理所有企业和个人账户、平台配置、审计日志 |
| 特殊权限 | 可查看和操作平台内所有资源                 |

#### 2.2.2 企业用户（Company）

一个企业账户下可包含多个子用户，各子用户拥有不同角色与权限。

| 角色             | 权限说明                                                   |
| ---------------- | ---------------------------------------------------------- |
| 企业管理员       | 管理企业内所有用户、分配角色与权限、管理所有 VPS 实例、查看账单 |
| 企业运维         | 管理被分配的 VPS 实例（启动/停止/重启/控制台）、查看监控   |
| 企业财务         | 查看企业账单与消费记录、充值、管理支付方式                 |
| 企业只读成员     | 只读查看被授权的 VPS 信息                                  |

#### 2.2.3 个人用户（Individual）

| 项目     | 说明                                         |
| -------- | -------------------------------------------- |
| 数量     | 单账户单用户，无子用户概念                   |
| 职责     | 管理自己名下的 VPS 实例                      |
| 权限     | 等同企业管理员对自身资源的权限，但无用户管理能力 |

### 2.3 核心业务对象

| 对象         | 说明                                         |
| ------------ | -------------------------------------------- |
| VPS 实例     | 虚拟服务器，包含配置、状态、IP、系统镜像等   |
| 订单（Order）| VPS 的购买/续费/升级订单                     |
| 账单（Bill） | 企业或个人的消费记录                         |
| 工单（Ticket）| 用户提交的技术支持请求                      |
| 操作日志     | 用户对资源的所有操作记录                     |
| API Key      | 用户生成的 API 访问凭证                      |

---

## 3. 功能模块

### 3.1 认证与会话管理

| 功能         | 说明                                       |
| ------------ | ------------------------------------------ |
| 用户注册     | 支持企业注册和个人注册两种入口             |
| 用户登录     | 账号密码登录，返回 JWT Token               |
| 会话管理     | Token 刷新、登出、多端登录控制             |
| 密码管理     | 修改密码、重置密码                         |

### 3.2 用户管理

| 功能               | 可用角色                     | 说明                               |
| ------------------ | ---------------------------- | ---------------------------------- |
| 查看个人信息       | 所有用户                     | 查看/编辑自己的基本信息            |
| 企业成员列表       | 企业管理员                   | 查看企业下所有成员                 |
| 添加企业成员       | 企业管理员                   | 邀请新用户加入企业                 |
| 编辑企业成员       | 企业管理员                   | 修改成员角色与权限                 |
| 删除企业成员       | 企业管理员                   | 移除企业成员                       |
| 查看所有用户       | 平台管理员                   | 管理平台内所有账户                 |

### 3.3 VPS 实例管理

| 功能             | 可用角色                              | 说明                               |
| ---------------- | ------------------------------------- | ---------------------------------- |
| 购买 VPS         | 企业管理员、个人用户                  | 选择配置并下单购买                 |
| 查看 VPS 列表    | 所有用户（按权限过滤）               | 查看有权限访问的 VPS               |
| 查看 VPS 详情    | VPS 所有者/被授权者                   | 查看配置、状态、IP、流量等         |
| 启动/停止/重启   | 企业管理员、企业运维、个人用户        | 控制 VPS 电源状态                  |
| 重装系统         | 企业管理员、个人用户                  | 选择镜像重装操作系统               |
| 控制台（VNC）    | 企业管理员、企业运维、个人用户        | 模拟 Web 终端访问                  |
| 删除/退订 VPS    | 企业管理员、个人用户                  | 销毁 VPS 实例                      |

### 3.4 财务管理

| 功能             | 可用角色                     | 说明                               |
| ---------------- | ---------------------------- | ---------------------------------- |
| 查看账单         | 企业管理员、企业财务、个人用户 | 查看消费明细                      |
| 余额充值         | 企业管理员、企业财务、个人用户 | 模拟充值操作                      |
| 查看订单         | 企业管理员、企业财务、个人用户 | 查看历史订单记录                  |
| 导出账单         | 企业管理员、企业财务          | 导出 CSV 格式账单                 |

### 3.5 工单系统

| 功能             | 可用角色                     | 说明                               |
| ---------------- | ---------------------------- | ---------------------------------- |
| 创建工单         | 所有用户                     | 提交技术支持请求                   |
| 查看工单列表     | 工单创建者、企业管理员、平台管理员 | 查看相关工单               |
| 回复工单         | 工单创建者、平台管理员       | 在工单中追加回复                   |
| 关闭工单         | 工单创建者、平台管理员       | 关闭已解决的工单                   |

### 3.6 API Key 管理

| 功能             | 可用角色                     | 说明                               |
| ---------------- | ---------------------------- | ---------------------------------- |
| 生成 API Key     | 企业管理员、个人用户         | 创建 API 访问凭证                  |
| 查看 API Key     | Key 所有者                   | 查看已创建的 Key 列表              |
| 删除 API Key     | Key 所有者                   | 吊销指定 Key                       |

### 3.7 操作审计日志

| 功能             | 可用角色                     | 说明                               |
| ---------------- | ---------------------------- | ---------------------------------- |
| 查看个人操作日志 | 所有用户                     | 查看自己的操作记录                 |
| 查看企业操作日志 | 企业管理员                   | 查看企业内所有成员的操作记录       |
| 查看平台操作日志 | 平台管理员                   | 查看全平台操作记录                 |

### 3.8 平台管理（Platform Admin）

| 功能             | 说明                                       |
| ---------------- | ------------------------------------------ |
| 用户管理         | 查看/禁用/启用所有账户                     |
| 企业管理         | 查看/审核企业信息                          |
| VPS 全局管理     | 查看所有 VPS 实例状态                      |
| 系统配置         | 管理平台基础配置                           |
| 公告管理         | 发布/编辑/删除平台公告                     |

---

## 4. 越权漏洞设计

> **核心价值：** 以下漏洞均为**有意植入**的安全缺陷，供学习者发现和练习。

### 4.1 漏洞分类

| 类别         | 编码 | 说明                                       |
| ------------ | ---- | ------------------------------------------ |
| 水平越权 (IDOR) | H-xx | 同权限等级用户之间的资源越权访问           |
| 垂直越权       | V-xx | 低权限用户执行高权限操作                   |
| 上下文越权     | C-xx | 跨业务上下文的逻辑越权                     |

### 4.2 漏洞清单

| 编号  | 类别     | 漏洞场景                                     | 涉及模块       |
| ----- | -------- | -------------------------------------------- | -------------- |
| H-01  | 水平越权 | 通过修改 VPS ID 查看他人 VPS 详情            | VPS 管理       |
| H-02  | 水平越权 | 通过修改用户 ID 查看他人个人信息             | 用户管理       |
| H-03  | 水平越权 | 通过修改订单 ID 查看他人订单详情             | 财务管理       |
| H-04  | 水平越权 | 通过修改工单 ID 查看/回复他人工单            | 工单系统       |
| H-05  | 水平越权 | 通过修改 API Key ID 删除他人 API Key         | API Key 管理   |
| V-01  | 垂直越权 | 企业只读成员调用 VPS 启停接口                | VPS 管理       |
| V-02  | 垂直越权 | 企业运维调用用户管理接口添加成员             | 用户管理       |
| V-03  | 垂直越权 | 企业财务调用 VPS 重装系统接口                | VPS 管理       |
| V-04  | 垂直越权 | 个人用户调用平台管理员接口                   | 平台管理       |
| V-05  | 垂直越权 | 企业运维修改自身角色为企业管理员             | 用户管理       |
| C-01  | 上下文越权 | A 企业用户操作 B 企业的 VPS 实例           | VPS 管理       |
| C-02  | 上下文越权 | 个人用户通过接口创建企业成员               | 用户管理       |
| C-03  | 上下文越权 | 利用过期/被吊销的 API Key 仍可访问         | API Key 管理   |

### 4.3 编码加密挑战（E-01 ~ E-09）

编码加密挑战是在传统越权漏洞基础上增加参数编码/加密层，攻击者需先识别编码方式并解码，才能实施越权攻击。

**架构设计：** 采用全局透明编码中间件架构
- 前端：axios 全局请求拦截器，当挑战激活时对所有请求的 body/query 中的字符串值递归编码，自动添加 `X-Encoding-Type` 请求头
- 后端：Gin 全局中间件 `EncodingMiddleware`，在 `AuthMiddleware` 之后执行，透明解码请求数据，Handler 无感知
- 排除路径：`/auth/*`（登录注册）、`/crypto/*`（编码工具端点本身）不受编码影响

| 挑战 ID | 编码类型 | 前端编码 | 后端解码 | 说明 |
|---------|---------|---------|---------|------|
| E-01 | base64 | `base64Encode` | `Base64Decode` | Base64 编码 VPS ID |
| E-02 | base32 | `base32Encode` | `Base32Decode` | Base32 编码用户 ID（无填充） |
| E-03 | caesar | `caesarEncode(shift=3)` | `CaesarDecode` | 凯撒密码编码订单 ID |
| E-04 | custom_base64 | `customBase64Encode` | `CustomBase64Decode` | 自定义 Base64 字符表 |
| E-05 | multi | `base32(base64(x))` | `Base32Decode → Base64Decode` | 双层嵌套编码 |
| E-06 | aes | `cryptoApi.cryptoEncode` | `AESDecryptFromBase64` | AES-256-GCM 加密 |
| E-07 | signed | `encodeSignedParam` | `DecodeSignedParam` | HMAC-SHA256 签名 |
| E-08 | sm4 | `cryptoApi.cryptoEncode` | `SM4DecryptFromBase64` | SM4-CBC 国密加密 |
| E-09 | md5_hash | `cryptoApi.cryptoEncode` | `DecodeHashSignedParam` | MD5 哈希签名（可预测盐值） |

### 4.4 漏洞难度分级

| 等级   | 星级 | 说明                                       |
| ------ | ---- | ------------------------------------------ |
| 入门   | ★☆☆  | 直接修改 URL/参数即可触发                  |
| 进阶   | ★★☆  | 需要构造特定请求或绕过前端校验             |
| 高级   | ★★★  | 需要组合多个漏洞或深入理解业务逻辑         |

### 4.4 学习辅助

- **漏洞提示（Hint）：** 每个漏洞提供 3 级渐进提示
- **漏洞解析（WriteUp）：** 提供完整的漏洞分析、利用步骤与修复方案
- **安全模式开关：** 可在"安全模式"与"漏洞模式"之间切换，对比学习

---

## 5. 系统架构

### 5.1 整体架构

```
┌─────────────────────────────────────────────┐
│                 浏览器 (Browser)              │
│          Vue3 + PrimeVue + Pinia            │
└──────────────────┬──────────────────────────┘
                   │ HTTP/REST API
┌──────────────────▼──────────────────────────┐
│              Go HTTP Server                  │
│  ┌─────────┐ ┌──────────┐ ┌──────────────┐  │
│  │  Router  │ │Middleware│ │  Controller  │  │
│  │ (Gin)   │ │(Auth/RBAC│ │  (Handler)   │  │
│  └─────────┘ └──────────┘ └──────┬───────┘  │
│                                  │           │
│  ┌───────────┐  ┌────────────────▼────────┐  │
│  │  Service   │  │       Repository       │  │
│  │  (业务层)  │  │     (数据访问层)        │  │
│  └───────────┘  └────────────┬────────────┘  │
│                              │               │
│                    ┌─────────▼─────────┐     │
│                    │     SQLite3       │     │
│                    │   (数据库文件)     │     │
│                    └───────────────────┘     │
└──────────────────────────────────────────────┘
```

### 5.2 技术选型

| 层级     | 技术                 | 说明                             |
| -------- | -------------------- | -------------------------------- |
| 前端框架 | Vue 3 (Composition API) | 响应式 UI 框架                |
| UI 组件库 | PrimeVue 4          | 企业级 Vue 组件库               |
| CSS 框架 | TailwindCSS 4          | 原子化 CSS 框架                 |
| 状态管理 | Pinia                | Vue 官方推荐状态管理            |
| 路由     | Vue Router 4         | 前端路由管理                    |
| HTTP 客户端 | Axios             | 前端 HTTP 请求                  |
| 前端构建 | Vite                 | 快速开发构建工具                |
| 前端加密 | crypto-js, sm-crypto, jsencrypt, base-x | 编码/加密/国密/签名库 |
| 包管理   | pnpm                 | 高性能包管理器                  |
| 后端框架 | Gin                  | 高性能 Go Web 框架              |
| ORM      | GORM                 | Go ORM 库                       |
| 数据库   | SQLite3              | 轻量级嵌入式数据库              |
| 认证     | JWT (golang-jwt)     | JSON Web Token 认证             |
| API 文档 | Swagger (swaggo)     | 自动生成 API 文档               |

### 5.3 项目目录结构

```
OverstepLab/
├── doc/                          # 文档目录
│   └── PRD.md
├── src/
│   ├── backend/                  # Go 后端
│   │   ├── cmd/                  # 入口
│   │   │   └── server/
│   │   │       └── main.go
│   │   ├── internal/
│   │   │   ├── config/           # 配置管理
│   │   │   ├── middleware/       # 中间件（认证、RBAC、日志）
│   │   │   ├── handler/          # HTTP Handler（Controller）
│   │   │   ├── service/          # 业务逻辑层
│   │   │   ├── repository/       # 数据访问层
│   │   │   ├── model/            # 数据模型
│   │   │   ├── crypto/           # 编码/加密工具库（古典、现代、国密、签名、大整数）
│   │   │   ├── vuln/             # 漏洞模式控制与挑战元数据
│   │   │   ├── web/              # 嵌入式前端静态资源
│   │   ├── database/
│   │   │   ├── migration/        # 数据库迁移
│   │   │   └── seed/             # 初始数据填充
│   │   ├── go.mod
│   │   └── go.sum
│   │
│   └── frontend/                 # Vue3 前端
│       ├── src/
│       │   ├── api/              # API 请求封装（含编码/加密专用客户端）
│       │   ├── assets/           # 静态资源
│       │   ├── components/       # 公共组件
│       │   ├── composables/      # 组合式函数
│       │   ├── layouts/          # 布局组件
│       │   ├── router/           # 路由配置
│       │   ├── stores/           # Pinia 状态管理（含编码挑战状态）
│       │   ├── utils/            # 工具函数（含 crypto.ts 编码工具库）
│       │   ├── views/            # 页面视图
│       │   ├── App.vue
│       │   └── main.ts
│       ├── index.html
│       ├── vite.config.ts
│       ├── tsconfig.json
│       └── package.json
│
└── tmp/                          # 临时文件
```

---

## 6. 数据模型设计

### 6.1 ER 关系概览

```
Company (1) ──── (N) User
User    (1) ──── (N) VPS
User    (1) ──── (N) Order
User    (1) ──── (N) Ticket
User    (1) ──── (N) ApiKey
User    (1) ──── (N) AuditLog
Company (1) ──── (N) Bill
```

### 6.2 核心数据表

#### users 用户表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | 用户 ID                                  |
| username        | TEXT         | 用户名（唯一）                           |
| password_hash   | TEXT         | 密码哈希（bcrypt）                       |
| email           | TEXT         | 邮箱                                     |
| phone           | TEXT         | 手机号                                   |
| avatar          | TEXT         | 头像 URL                                 |
| user_type       | TEXT         | 用户类型：platform_admin / company / individual |
| company_id      | INTEGER (FK) | 所属企业 ID（个人用户为 NULL）           |
| role            | TEXT         | 角色：admin / operator / finance / viewer |
| status          | TEXT         | 状态：active / disabled                  |
| created_at      | DATETIME     | 创建时间                                 |
| updated_at      | DATETIME     | 更新时间                                 |

#### companies 企业表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | 企业 ID                                  |
| name            | TEXT         | 企业名称                                 |
| license_no      | TEXT         | 营业执照号                               |
| contact_name    | TEXT         | 联系人姓名                               |
| contact_phone   | TEXT         | 联系人电话                               |
| balance         | REAL         | 账户余额                                 |
| status          | TEXT         | 状态：active / suspended                 |
| created_at      | DATETIME     | 创建时间                                 |
| updated_at      | DATETIME     | 更新时间                                 |

#### vps_instances VPS 实例表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | VPS ID                                   |
| name            | TEXT         | 实例名称                                 |
| owner_id        | INTEGER (FK) | 所有者用户 ID                            |
| company_id      | INTEGER (FK) | 所属企业 ID（个人用户为 NULL）           |
| cpu             | INTEGER      | CPU 核数                                 |
| memory          | INTEGER      | 内存（MB）                               |
| disk            | INTEGER      | 磁盘（GB）                               |
| bandwidth       | INTEGER      | 带宽（Mbps）                             |
| ip_address      | TEXT         | IP 地址                                  |
| os_image        | TEXT         | 操作系统镜像                             |
| status          | TEXT         | 状态：running / stopped / creating / error |
| expire_at       | DATETIME     | 到期时间                                 |
| created_at      | DATETIME     | 创建时间                                 |
| updated_at      | DATETIME     | 更新时间                                 |

#### orders 订单表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | 订单 ID                                  |
| order_no        | TEXT         | 订单编号（唯一）                         |
| user_id         | INTEGER (FK) | 下单用户 ID                              |
| company_id      | INTEGER (FK) | 所属企业 ID                              |
| vps_id          | INTEGER (FK) | 关联 VPS ID                              |
| type            | TEXT         | 类型：purchase / renew / upgrade         |
| amount          | REAL         | 订单金额                                 |
| status          | TEXT         | 状态：pending / paid / cancelled         |
| created_at      | DATETIME     | 创建时间                                 |
| updated_at      | DATETIME     | 更新时间                                 |

#### tickets 工单表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | 工单 ID                                  |
| title           | TEXT         | 工单标题                                 |
| content         | TEXT         | 工单内容                                 |
| user_id         | INTEGER (FK) | 创建者用户 ID                            |
| company_id      | INTEGER (FK) | 所属企业 ID                              |
| status          | TEXT         | 状态：open / replied / closed            |
| created_at      | DATETIME     | 创建时间                                 |
| updated_at      | DATETIME     | 更新时间                                 |

#### ticket_replies 工单回复表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | 回复 ID                                  |
| ticket_id       | INTEGER (FK) | 关联工单 ID                              |
| user_id         | INTEGER (FK) | 回复者用户 ID                            |
| content         | TEXT         | 回复内容                                 |
| created_at      | DATETIME     | 创建时间                                 |

#### api_keys API Key 表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | Key ID                                   |
| user_id         | INTEGER (FK) | 所属用户 ID                              |
| name            | TEXT         | Key 名称                                 |
| key_value       | TEXT         | Key 值（哈希存储）                       |
| key_prefix      | TEXT         | Key 前缀（用于展示）                     |
| permissions     | TEXT         | 权限范围（JSON）                         |
| status          | TEXT         | 状态：active / revoked                   |
| last_used_at    | DATETIME     | 最后使用时间                             |
| expire_at       | DATETIME     | 过期时间                                 |
| created_at      | DATETIME     | 创建时间                                 |

#### audit_logs 审计日志表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | 日志 ID                                  |
| user_id         | INTEGER (FK) | 操作者用户 ID                            |
| company_id      | INTEGER (FK) | 所属企业 ID                              |
| action          | TEXT         | 操作类型（如 vps.start, user.create）    |
| resource_type   | TEXT         | 资源类型                                 |
| resource_id     | INTEGER      | 资源 ID                                  |
| detail          | TEXT         | 操作详情（JSON）                         |
| ip_address      | TEXT         | 操作者 IP                                |
| created_at      | DATETIME     | 操作时间                                 |

#### bills 账单表

| 字段            | 类型         | 说明                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INTEGER (PK) | 账单 ID                                  |
| company_id      | INTEGER (FK) | 所属企业 ID（个人用户关联 user_id）      |
| user_id         | INTEGER (FK) | 关联用户 ID                              |
| type            | TEXT         | 类型：expense / recharge                 |
| amount          | REAL         | 金额                                     |
| balance_after   | REAL         | 操作后余额                               |
| description     | TEXT         | 描述                                     |
| created_at      | DATETIME     | 创建时间                                 |

---

## 7. API 设计

### 7.1 API 规范

- 基础路径：`/api/v1`
- 认证方式：`Authorization: Bearer <JWT Token>`
- 响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 7.2 核心 API 列表

#### 认证模块

| Method | Path                      | 说明         |
| ------ | ------------------------- | ------------ |
| POST   | /api/v1/auth/register     | 用户注册     |
| POST   | /api/v1/auth/login        | 用户登录     |
| POST   | /api/v1/auth/refresh      | 刷新 Token   |
| POST   | /api/v1/auth/logout       | 用户登出     |

#### 用户模块

| Method | Path                            | 说明               |
| ------ | ------------------------------- | ------------------ |
| GET    | /api/v1/user/profile            | 获取个人信息       |
| PUT    | /api/v1/user/profile            | 修改个人信息       |
| PUT    | /api/v1/user/password           | 修改密码           |
| GET    | /api/v1/users/:id               | 获取指定用户信息   |

#### 企业成员管理

| Method | Path                                | 说明             |
| ------ | ----------------------------------- | ---------------- |
| GET    | /api/v1/company/members             | 获取企业成员列表 |
| POST   | /api/v1/company/members             | 添加企业成员     |
| PUT    | /api/v1/company/members/:id         | 编辑企业成员     |
| DELETE | /api/v1/company/members/:id         | 删除企业成员     |
| PUT    | /api/v1/company/members/:id/role    | 修改成员角色     |

#### VPS 管理

| Method | Path                              | 说明             |
| ------ | --------------------------------- | ---------------- |
| GET    | /api/v1/vps                       | 获取 VPS 列表    |
| POST   | /api/v1/vps                       | 购买 VPS         |
| GET    | /api/v1/vps/:id                   | 获取 VPS 详情    |
| POST   | /api/v1/vps/:id/start             | 启动 VPS         |
| POST   | /api/v1/vps/:id/stop              | 停止 VPS         |
| POST   | /api/v1/vps/:id/restart           | 重启 VPS         |
| POST   | /api/v1/vps/:id/reinstall         | 重装系统         |
| DELETE | /api/v1/vps/:id                   | 删除 VPS         |
| GET    | /api/v1/vps/:id/console           | 获取控制台连接   |

#### 财务模块

| Method | Path                              | 说明             |
| ------ | --------------------------------- | ---------------- |
| GET    | /api/v1/orders                    | 获取订单列表     |
| GET    | /api/v1/orders/:id                | 获取订单详情     |
| GET    | /api/v1/bills                     | 获取账单列表     |
| POST   | /api/v1/bills/recharge            | 余额充值         |
| GET    | /api/v1/bills/export              | 导出账单         |

#### 工单模块

| Method | Path                              | 说明             |
| ------ | --------------------------------- | ---------------- |
| GET    | /api/v1/tickets                   | 获取工单列表     |
| POST   | /api/v1/tickets                   | 创建工单         |
| GET    | /api/v1/tickets/:id               | 获取工单详情     |
| POST   | /api/v1/tickets/:id/reply         | 回复工单         |
| PUT    | /api/v1/tickets/:id/close         | 关闭工单         |

#### API Key 模块

| Method | Path                              | 说明             |
| ------ | --------------------------------- | ---------------- |
| GET    | /api/v1/apikeys                   | 获取 Key 列表    |
| POST   | /api/v1/apikeys                   | 创建 Key         |
| DELETE | /api/v1/apikeys/:id               | 删除 Key         |

#### 审计日志

| Method | Path                              | 说明             |
| ------ | --------------------------------- | ---------------- |
| GET    | /api/v1/audit-logs                | 获取操作日志     |

#### 编码挑战状态

| Method | Path                              | 说明             |
| ------ | --------------------------------- | ---------------- |
| GET    | /api/v1/encoding-challenge-state  | 获取当前编码挑战状态 |
| GET    | /api/v1/encoding-challenge-state/:type | 获取指定类型挑战状态 |

#### 编码/加密端点

| Method | Path                              | 说明             |
| ------ | --------------------------------- | ---------------- |
| POST   | /encoded/:type                    | 编码参数越权测试端点 |
| POST   | /crypto/:type                     | 加密参数越权测试端点 |
| POST   | /api/v1/crypto/encode             | 前端编码工具 API   |
| POST   | /api/v1/crypto/decode             | 前端解码工具 API   |

#### 平台管理

| Method | Path                              | 说明               |
| ------ | --------------------------------- | ------------------ |
| GET    | /api/v1/admin/users               | 获取所有用户       |
| PUT    | /api/v1/admin/users/:id/status    | 启用/禁用用户      |
| GET    | /api/v1/admin/companies           | 获取所有企业       |
| GET    | /api/v1/admin/vps                 | 获取所有 VPS       |
| GET    | /api/v1/admin/audit-logs          | 获取全平台日志     |

---

## 8. 权限模型

### 8.1 RBAC 权限矩阵

| 资源 \ 角色        | Platform Admin | Company Admin | Company Operator | Company Finance | Company Viewer | Individual |
| ------------------- | :---: | :---: | :---: | :---: | :---: | :---: |
| 平台管理            | ✅    | ❌    | ❌    | ❌    | ❌    | ❌     |
| 管理企业成员        | ✅    | ✅    | ❌    | ❌    | ❌    | ❌     |
| 购买 VPS            | ✅    | ✅    | ❌    | ❌    | ❌    | ✅     |
| VPS 启停重启        | ✅    | ✅    | ✅    | ❌    | ❌    | ✅     |
| VPS 重装/删除       | ✅    | ✅    | ❌    | ❌    | ❌    | ✅     |
| VPS 控制台          | ✅    | ✅    | ✅    | ❌    | ❌    | ✅     |
| 查看 VPS（已授权）  | ✅    | ✅    | ✅    | ❌    | ✅    | ✅     |
| 查看账单            | ✅    | ✅    | ❌    | ✅    | ❌    | ✅     |
| 充值                | ✅    | ✅    | ❌    | ✅    | ❌    | ✅     |
| 导出账单            | ✅    | ✅    | ❌    | ✅    | ❌    | ❌     |
| 创建工单            | ✅    | ✅    | ✅    | ✅    | ✅    | ✅     |
| 管理 API Key        | ✅    | ✅    | ❌    | ❌    | ❌    | ✅     |
| 查看审计日志（企业）| ✅    | ✅    | ❌    | ❌    | ❌    | ❌     |

---

## 9. 初始化数据

系统首次启动时自动填充以下测试数据：

### 9.1 预置账户

| 用户名          | 密码         | 类型         | 角色             | 所属企业       |
| --------------- | ------------ | ------------ | ---------------- | -------------- |
| admin           | admin123     | 平台管理员   | platform_admin   | -              |
| acme_admin      | pass123      | 企业用户     | admin            | Acme Corp      |
| acme_ops        | pass123      | 企业用户     | operator         | Acme Corp      |
| acme_finance    | pass123      | 企业用户     | finance          | Acme Corp      |
| acme_viewer     | pass123      | 企业用户     | viewer           | Acme Corp      |
| globex_admin    | pass123      | 企业用户     | admin            | Globex Inc     |
| globex_ops      | pass123      | 企业用户     | operator         | Globex Inc     |
| globex_finance  | pass123      | 企业用户     | finance          | Globex Inc     |
| globex_viewer   | pass123      | 企业用户     | viewer           | Globex Inc     |
| alice           | pass123      | 个人用户     | individual       | -              |
| bob             | pass123      | 个人用户     | individual       | -              |

### 9.2 预置 VPS 实例

- Acme Corp：3 台 VPS（分配给不同成员管理）
- Globex Inc：2 台 VPS
- alice：1 台 VPS
- bob：1 台 VPS

---

## 10. 非功能性需求

### 10.1 部署要求

- 支持单二进制部署（Go 编译前端静态资源嵌入）
- 支持 Docker 容器部署
- 首次启动自动初始化数据库与种子数据

### 10.2 性能要求

- 作为靶场项目，无高并发要求
- 单机部署，支持 10+ 并发用户即可

### 10.3 兼容性

- 前端：Chrome 90+、Firefox 90+、Edge 90+
- 后端：Go 1.25+
- 数据库：SQLite 3.35+

---

## 11. 靶场辅助功能

### 11.1 漏洞挑战面板

- 展示所有漏洞挑战列表及完成状态
- 每个挑战提供：描述、难度、分类标签
- 支持标记已完成的挑战

### 11.2 提示系统

- 每个漏洞提供 3 级渐进式提示：
  - Level 1：方向性提示（如"注意请求中的 ID 参数"）
  - Level 2：具体提示（如"尝试修改 GET /api/v1/vps/:id 中的 id 值"）
  - Level 3：完整解析步骤

### 11.3 安全模式切换

- 漏洞模式（默认）：所有漏洞均可触发
- 安全模式：启用完整的权限校验，所有漏洞被修复
- 通过接口或环境变量切换

### 11.4 重置功能

- 一键重置数据库到初始状态
- 清除所有用户操作痕迹

---

## 12. 里程碑规划

| 阶段   | 内容                                       | 预计周期 |
| ------ | ------------------------------------------ | -------- |
| M1     | 项目初始化、用户认证、用户管理模块         | 1 周     |
| M2     | VPS 管理、财务管理模块                     | 1 周     |
| M3     | 工单系统、API Key、审计日志模块            | 1 周     |
| M4     | 越权漏洞植入与漏洞挑战面板                 | 1 周     |
| M5     | 提示系统、安全模式、数据重置、整体联调     | 1 周     |
| M6     | Docker 部署、文档完善、发布                | 3 天     |
