# 教程 06: 漏洞 WriteUp 参考

## 水平越权 (IDOR)

### H-01: 查看他人 VPS
- **端点**: `GET /api/v1/vps/:id`
- **原因**: 缺少所有权校验
- **修复**: 校验 `vps.OwnerID == user.ID || vps.CompanyID == user.CompanyID`

### H-02: 查看他人信息
- **端点**: `GET /api/v1/users/:id`
- **原因**: 允许查看任意用户
- **修复**: 仅允许查看自己（平台管理员除外）

### H-03: 查看他人订单
- **端点**: `GET /api/v1/orders/:id`
- **原因**: 缺少订单归属校验
- **修复**: 校验 `order.UserID == user.ID`

### H-04: 查看/回复他人工单
- **端点**: `GET /api/v1/tickets/:id`, `POST /api/v1/tickets/:id/reply`
- **原因**: 缺少工单创建者校验
- **修复**: 校验 `ticket.UserID == user.ID`

### H-05: 删除他人 API Key
- **端点**: `DELETE /api/v1/apikeys/:id`
- **原因**: 缺少 Key 归属校验
- **修复**: 校验 `key.UserID == user.ID`

## 垂直越权

### V-01: 只读成员控制 VPS
- **端点**: `POST /api/v1/vps/:id/start`
- **原因**: 服务层没有角色校验
- **修复**: 要求 `operator` 或 `admin` 角色

### V-02: 运维添加成员
- **端点**: `POST /api/v1/company/members`
- **原因**: 中间件在漏洞模式跳过角色检查
- **修复**: `RequirePermission(PermUserManage)`

### V-03: 财务重装 VPS
- **端点**: `POST /api/v1/vps/:id/reinstall`
- **原因**: 服务层没有角色校验
- **修复**: 要求 `admin` 或 `individual`

### V-04: 个人用户调用管理接口
- **端点**: `GET /api/v1/admin/users`
- **原因**: `RequireAdmin` 在漏洞模式跳过检查
- **修复**: 强制校验 `platform_admin`

### V-05: 自升管理员
- **端点**: `PUT /api/v1/company/members/:id/role`
- **原因**: 缺少调用者权限和目标关系校验
- **修复**: 验证调用者是管理员，且禁止修改自身角色

## 上下文越权

### C-01: 跨企业操作
- **端点**: `POST /api/v1/vps/:id/stop`
- **原因**: 缺少企业边界校验
- **修复**: 校验 `vps.CompanyID == user.CompanyID`

### C-02: 个人创建企业成员
- **端点**: `POST /api/v1/company/members`
- **原因**: 缺少 userType 校验
- **修复**: 验证 `user.UserType == "company"`

### C-03: 吊销 Key 仍可访问
- **端点**: 所有 API Key 认证的接口
- **原因**: 认证中间件不检查 key status 和 expire_at
- **修复**: 增加 `key.Status == "active" && (key.ExpireAt == nil || key.ExpireAt > now)`
