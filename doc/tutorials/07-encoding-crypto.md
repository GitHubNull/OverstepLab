# 教程 07: 编码加密挑战

## 概念

**编码加密越权** 是在传统 IDOR 漏洞的基础上，增加了参数编码/加密层。攻击者不仅要发现越权漏洞，还需要先识别和解码参数编码方式，才能成功实施越权攻击。

**难度范围**: ★☆☆ ~ ★★★★ (入门到专家级)

## 编码 vs 加密

| 类型 | 特点 | 示例 |
|------|------|------|
| **编码 (Encoding)** | 可逆变换，无密钥，仅改变数据表示形式 | Base64, Base32, Base58, 凯撒密码 |
| **加密 (Encryption)** | 需要密钥才能解密，提供机密性保护 | AES-256-GCM, RSA-OAEP, SM4-CBC |
| **签名 (Signature)** | 使用密钥生成验证码，确保完整性 | HMAC-SHA256 |

## 工具准备

### 1. 后端编码API

系统提供了三个辅助接口：

```bash
# 获取加密密钥
curl http://localhost:5000/api/v1/crypto/keys \
  -H "Authorization: Bearer <your-token>"

# 编码/加密参数
curl -X POST http://localhost:5000/api/v1/crypto/encode \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"value": "2", "encoding_type": "base64"}'

# 解码/解密参数
curl -X POST http://localhost:5000/api/v1/crypto/decode \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"value": "Mg==", "encoding_type": "base64"}'
```

### 2. 浏览器控制台工具

打开浏览器 Console（F12 → Console），可使用 `CryptoUtils` 对象：

```javascript
CryptoUtils.base64Encode("2")    // → "Mg=="
CryptoUtils.base64Decode("Mg==") // → "2"
CryptoUtils.base32Encode("3")    // → "GM"
CryptoUtils.caesarEncode("4", 3) // → "7"
```

## 全局编码机制

从 v1.1 版本开始，编码挑战采用**全局透明编码中间件**架构：

- **前端**：axios 全局请求拦截器自动处理编码。当在挑战页面激活任意编码挑战后，所有 API 请求的 body 和 query 参数中的字符串值会被递归编码，并自动附加 `X-Encoding-Type` 请求头。
- **后端**：Gin 中间件 `EncodingMiddleware` 在认证之后透明解码请求数据，所有 Handler 无感知，直接处理明文数据。
- **排除路径**：`/auth/*`（登录/注册）、`/crypto/*`（编码工具端点）不受编码影响。
- **数字类型保护**：`vpsId: 1` 等数字类型字段保持原样，不会被编码。

这意味着在浏览器中正常操作（创建 API Key、查看 VPS、提交工单等）时，前端会自动编码参数，后端会自动解码，用户无需手动构造编码请求。

## 挑战详解

### E-01: Base64 编码 VPS ID 越权 ⭐

**目标**: 通过 Base64 编码的 VPS ID，查看其他企业的 VPS。

**步骤**:
1. 使用 `acme_ops` / `pass123` 登录
2. 将目标 VPS ID 进行 Base64 编码，例如 ID=4 → `NA==`
3. 发送请求：

```bash
curl "http://localhost:5000/api/v1/encoded/vps/NA==" \
  -H "Authorization: Bearer <your-token>" \
  -H "X-Encoding-Type: base64"
```

**提示**: Base64 编码的 ID 可以直接用 `echo -n "4" | base64` 获得。

---

### E-02: Base32 编码用户 ID 越权 ⭐

**目标**: 通过 Base32 编码的用户 ID，查看其他用户的个人信息。

**步骤**:
1. 使用 `acme_ops` / `pass123` 登录
2. 将目标用户 ID 进行 Base32 编码，例如 ID=8 →无填充格式
3. 发送请求：

```bash
curl "http://localhost:5000/api/v1/encoded/users/<base32-encoded-id>" \
  -H "Authorization: Bearer <your-token>" \
  -H "X-Encoding-Type: base32"
```

**提示**: Base32 仅包含 A-Z 和 2-7 字符，无 `=` 填充，注意与 Base64 区分。

---

### E-03: 凯撒密码订单 ID 越权 ⭐

**目标**: 识别凯撒密码编码规律后查看他人订单。

**步骤**:
1. 凯撒密码默认移位量为 3（shift=3）
2. 数字 ID 移位规律：1→4, 2→5, 3→6, 4→7...
3. 要查询 ID=3 的订单，编码值为 6

```bash
curl "http://localhost:5000/api/v1/encoded/orders/6" \
  -H "Authorization: Bearer <your-token>" \
  -H "X-Encoding-Type: caesar"
```

---

### E-04: 自定义 Base64 编码表绕过 ⭐⭐

**目标**: 识别并使用自定义 Base64 字符表进行编码。

**步骤**:
1. 自定义 Base64 字符表规则：A↔Z、a↔z、0↔9、+↔/ 互换
2. 使用 `/api/v1/crypto/encode` 接口进行编码测试
3. 对比标准 Base64 和自定义 Base64 的输出差异
4. 使用识别出的自定义编码表构造请求

```bash
curl "http://localhost:5000/api/v1/encoded/vps/<custom-encoded-id>" \
  -H "Authorization: Bearer <your-token>" \
  -H "X-Encoding-Type: custom_base64"
```

---

### E-05: 多层嵌套编码工单越权 ⭐⭐

**目标**: 逐层解码 Base64 → Base32 嵌套的参数。

**步骤**:
1. 参数经过 Base64 编码后再进行 Base32 编码
2. 解码顺序：Base32 解码 → Base64 解码 → 得到原始 ID
3. 需要手动逐层解码，没有单一 X-Encoding-Type

```bash
# 先 Base32 解码
curl -X POST http://localhost:5000/api/v1/crypto/decode \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"value": "<encoded-id>", "encoding_type": "base32"}'

# 再 Base64 解码获得真实 ID
```

---

### E-06: AES 加密 VPS ID 操作 ⭐⭐⭐

**目标**: 获取 AES 密钥后解密参数，操作其他企业 VPS。

**步骤**:
1. 获取 AES 密钥：

```bash
curl http://localhost:5000/api/v1/crypto/keys \
  -H "Authorization: Bearer <your-token>"
# 响应中 aes_key_base64 字段即为密钥
```

2. 使用密钥加密 VPS ID（加密模式 AES-256-GCM）
3. 发送请求：

```bash
curl "http://localhost:5000/api/v1/encoded/vps/<aes-encrypted-id>/start" \
  -H "Authorization: Bearer <your-token>" \
  -H "X-Encoding-Type: aes"
```

**提示**: 可直接使用 `/api/v1/crypto/encode` 接口传入 `"encoding_type": "aes"` 生成加密值。

---

### E-07: HMAC 签名验证绕过 ⭐⭐⭐

**目标**: 获取 HMAC 密钥后伪造签名，实现越权操作。

**步骤**:
1. 获取 HMAC 密钥（`hmac_key_base64` 字段）
2. 参数格式：`Base64(data).HMAC签名Hex`
3. 签名算法：HMAC-SHA256，签名的原始值为未编码的数据
4. 伪造签名后发送请求：

```bash
curl "http://localhost:5000/api/v1/encoded/vps/<signed-param>/stop" \
  -H "Authorization: Bearer <your-token>" \
  -H "X-Encoding-Type: signed"
```

**提示**: 可使用 `/api/v1/crypto/encode` 接口传入 `"encoding_type": "signed"` 生成签名。

---

### E-08: SM4 国密加密跨企业操作 ⭐⭐⭐⭐

**目标**: 破解 SM4 加密 + 跨企业上下文越权双重组合。

**步骤**:
1. 获取 SM4 密钥（`sm4_key_base64` 字段）
2. SM4 加密模式：CBC，IV 为密文前 16 字节，PKCS7 填充
3. 使用 acme_ops 登录，将 Globex Inc 的 VPS ID (4-5) 进行 SM4 加密
4. 利用跨企业越权（C-01）操作其他企业 VPS：

```bash
# 注意：复杂加密参数建议通过 ?v= 查询参数传递
curl "http://localhost:5000/api/v1/encoded/vps/ignored?v=<sm4-encrypted-id>" \
  -H "Authorization: Bearer <your-token>" \
  -H "X-Encoding-Type: sm4"
```

---

## 编码类型速查

| X-Encoding-Type | 类型 | 难度 |
|-----------------|------|------|
| `base64` | 标准 Base64 | ⭐ |
| `base32` | RFC 4648 Base32 (无填充) | ⭐ |
| `caesar` | 凯撒密码 (shift=3) | ⭐ |
| `custom_base64` | 自定义字符表 Base64 | ⭐⭐ |
| `base58` | Bitcoin Base58 | ⭐⭐ |
| `base85` | Ascii85 变体 | ⭐⭐ |
| `custom_base32` | 自定义字符表 Base32 | ⭐⭐ |
| `aes` | AES-256-GCM | ⭐⭐⭐ |
| `rsa` | RSA-OAEP | ⭐⭐⭐ |
| `signed` | HMAC-SHA256 签名 | ⭐⭐⭐ |
| `sm4` | SM4-CBC 国密 | ⭐⭐⭐⭐ |

## 请求格式示例

### 全局编码模式（推荐）

激活挑战后，正常调用 API 即可，编码由前后端自动处理：

```bash
# 创建 API Key（body 中的字符串字段会自动 Base64 编码）
curl -X POST http://localhost:5000/api/v1/apikeys \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "X-Encoding-Type: base64" \
  -d '{"name":"dGVzdGtleQ==","permissions":"cmVhZA=="}'

# 查询 VPS 详情（query 参数会自动编码）
curl "http://localhost:5000/api/v1/vps/detail?vpsId=MQ==" \
  -H "Authorization: Bearer <token>" \
  -H "X-Encoding-Type: base64"
```

### 传统编码端点（向后兼容）

`/api/v1/encoded/*` 路由仍然可用：

```bash
# 路径参数（适用于简单编码）
curl "/api/v1/encoded/vps/Mg==" -H "X-Encoding-Type: base64"

# 查询参数（适用于 AES/SM4/RSA 等复杂编码）
curl "/api/v1/encoded/vps/ignored?v=<url-encoded-value>" -H "X-Encoding-Type: sm4"
```

## 修复方案

在安全模式下，编码包装端点也会启用权限校验。编码仅增加了参数传输的复杂度，真正的安全防护仍需在后端进行权限验证。

## 下一步

完成编码加密挑战后，你已经掌握了：
- 多种编码和加密算法的识别与解码
- 使用辅助工具加速攻击流程
- 理解编码 ≠ 加密 ≠ 安全的深刻含义

→ 返回 [教程首页](index.md)
