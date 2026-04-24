# 教程 04: 上下文越权

## 概念

**上下文越权**是指在业务逻辑层面存在缺陷，允许用户跨越正常的业务上下文边界进行操作。

**难度**: ★★★ (高级)

## 漏洞列表

### C-01: 跨企业操作 VPS

使用 A 企业的用户操作 B 企业的 VPS。

```bash
# 使用 Acme Corp 用户操作 Globex 的 VPS
curl -X POST http://localhost:8080/api/v1/vps/4/stop \
  -H "Authorization: Bearer <acme-token>"
```

### C-02: 个人用户创建企业成员

个人账户不属于任何企业，但在漏洞模式下可以调用企业成员管理接口。

```bash
curl -X POST http://localhost:8080/api/v1/company/members \
  -H "Authorization: Bearer <alice-token>" \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"pass123","email":"n@n.com","role":"viewer"}'
```

### C-03: 已吊销 API Key 仍可访问

系统中有一个已被吊销的 API Key（Globex 用户创建）。在漏洞模式下，认证中间件只检查 Key 是否存在，不检查状态。

## 修复方案

上下文越权的修复需要更严格的业务校验：
- 企业操作必须验证 userType = 'company'
- API Key 认证必须检查 status 和 expire_at
- 跨企业操作必须验证 company_id 匹配

→ [教程 05: 高级组合攻击](05-advanced-combo.md)
