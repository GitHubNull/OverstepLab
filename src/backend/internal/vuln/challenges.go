package vuln

type Challenge struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Category        string   `json:"category"`
	Difficulty      int      `json:"difficulty"`
	Description     string   `json:"description"`
	Hints           []string `json:"hints"`
	WriteUp         string   `json:"writeup"`
	Endpoint        string   `json:"endpoint"`
	Method          string   `json:"method"`
	EncodingType    string   `json:"encoding_type,omitempty"`
	EncodedEndpoint string   `json:"encoded_endpoint,omitempty"`
}

var Challenges = []Challenge{
	{
		ID: "H-01", Title: "查看他人 VPS 详情", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改 VPS ID 参数，查看其他用户或企业的 VPS 详细信息。",
		Hints: []string{
			"注意请求中的 VPS ID 参数",
			"尝试修改 GET /api/v1/vps?vpsId= 中的 vpsId 值，使用其他用户的 VPS ID",
			"Acme Corp 的 VPS ID 为 1-3，Globex Inc 的 VPS ID 为 4-5",
			"[编码提示] 也可使用 /api/v1/encoded/vps?v=，通过 X-Encoding-Type 头指定编码类型",
		},
		WriteUp:         "后端 GetDetail 方法在 vulnerable 模式下没有检查 VPS 所有权或企业归属，仅验证 VPS 是否存在。攻击者可通过遍历 ID 查看任意 VPS 详情。修复方案：在 secure 模式下增加 ownership/company 校验。",
		Endpoint:        "/api/v1/vps?vpsId=", Method: "GET",
		EncodedEndpoint: "/api/v1/encoded/vps?v=",
	},
	{
		ID: "H-02", Title: "查看他人个人信息", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改用户 ID 参数，查看其他用户的个人资料。",
		Hints: []string{
			"注意获取用户信息的接口使用了查询参数",
			"尝试访问 GET /api/v1/users?id= 并修改 id 值",
			"系统中预置了 9 个用户，ID 为 1-9",
			"[编码提示] 也可使用 /api/v1/encoded/users?v=，通过 X-Encoding-Type 头指定编码类型",
		},
		WriteUp:         "GetUserByID 接口在 vulnerable 模式下没有校验当前用户与目标用户的关系，允许查看任意用户资料。修复方案：限制用户只能查看自己的资料（平台管理员除外）。",
		Endpoint:        "/api/v1/users?id=", Method: "GET",
		EncodedEndpoint: "/api/v1/encoded/users?v=",
	},
	{
		ID: "H-03", Title: "查看他人订单", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改订单 ID，查看其他用户的购买/续费订单详情。",
		Hints: []string{
			"订单详情接口使用了订单 ID 作为查询参数",
			"尝试访问 GET /api/v1/orders?orderId= 并修改 orderId",
			"系统中有 4 个预置订单，ID 为 1-4",
			"[编码提示] 也可使用 /api/v1/encoded/orders?v=，通过 X-Encoding-Type 头指定编码类型",
		},
		WriteUp:         "订单详情查询没有校验订单归属用户，在 vulnerable 模式下允许查看任意订单。修复方案：增加 user_id 匹配校验。",
		Endpoint:        "/api/v1/orders?orderId=", Method: "GET",
		EncodedEndpoint: "/api/v1/encoded/orders?v=",
	},
	{
		ID: "H-04", Title: "查看/回复他人工单", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改工单 ID，查看或回复其他用户的工单。",
		Hints: []string{
			"工单详情和回复接口都使用了工单 ID",
			"尝试访问 GET /api/v1/tickets?ticketId= 和 POST /api/v1/tickets/reply（body 含 ticketId）",
			"系统中有 2 个预置工单，ID 为 1-2",
			"[编码提示] 也可使用 /api/v1/encoded/tickets?v=，通过 X-Encoding-Type 头指定编码类型",
		},
		WriteUp:         "工单详情和回复接口没有校验工单创建者，允许查看和回复他人工单。修复方案：增加 user_id 匹配校验。",
		Endpoint:        "/api/v1/tickets?ticketId=", Method: "GET",
		EncodedEndpoint: "/api/v1/encoded/tickets?v=",
	},
	{
		ID: "H-05", Title: "删除他人 API Key", Category: "水平越权 (IDOR)", Difficulty: 2,
		Description: "通过修改 API Key ID，删除其他用户创建的 API Key。",
		Hints: []string{
			"删除 API Key 的接口使用了 Key ID",
			"尝试发送 DELETE /api/v1/apikeys 请求（body 含 id），使用其他用户的 Key ID",
			"系统中预置了 3 个 API Key，ID 为 1-3",
			"[编码提示] 也可使用 /api/v1/encoded/apikeys（body 含 v），通过 X-Encoding-Type 头指定编码类型",
		},
		WriteUp:         "API Key 删除接口没有校验 Key 归属，在 vulnerable 模式下允许删除任意 Key。修复方案：增加 user_id 匹配校验。",
		Endpoint:        "/api/v1/apikeys（body 含 id）", Method: "DELETE",
		EncodedEndpoint: "/api/v1/encoded/apikeys（body 含 v）",
	},
	{
		ID: "V-01", Title: "只读成员控制 VPS 启停", Category: "垂直越权", Difficulty: 1,
		Description: "使用只读成员 (viewer) 账户，通过 API 直接调用 VPS 启停接口。",
		Hints: []string{
			"前端界面隐藏了 viewer 角色的操作按钮",
			"尝试用 acme_viewer 登录后，直接调用 POST /api/v1/vps/start（body 含 vpsId）",
			"前端隐藏不等于后端限制",
			"[编码提示] 也可使用 /api/v1/encoded/vps/start（body 含 v）",
		},
		WriteUp:         "VPS 控制接口在 vulnerable 模式下没有校验用户角色，viewer 角色本应只有查看权限，但可以执行启停操作。修复方案：在 secure 模式下增加角色权限检查（require operator 或 admin）。",
		Endpoint:        "/api/v1/vps/start（body 含 vpsId）", Method: "POST",
		EncodedEndpoint: "/api/v1/encoded/vps/start（body 含 v）",
	},
	{
		ID: "V-02", Title: "运维人员添加企业成员", Category: "垂直越权", Difficulty: 2,
		Description: "使用运维 (operator) 角色，调用添加成员接口创建新企业用户。",
		Hints: []string{
			"添加成员接口应该只有企业管理员可用",
			"尝试用 acme_ops 登录后，调用 POST /api/v1/company/members",
			"检查 RBAC 中间件和服务层双重校验",
			"[编码提示] 也可使用 /api/v1/encoded/company/members",
		},
		WriteUp:         "添加企业成员接口在 vulnerable 模式下没有严格校验用户角色，operator 角色可以执行管理员操作。修复方案：在 middleware 和 service 层同时增加角色检查。",
		Endpoint:        "/api/v1/company/members", Method: "POST",
		EncodedEndpoint: "/api/v1/encoded/company/members",
	},
	{
		ID: "V-03", Title: "财务人员重装 VPS 系统", Category: "垂直越权", Difficulty: 2,
		Description: "使用财务 (finance) 角色，调用 VPS 重装系统接口。",
		Hints: []string{
			"财务人员不应该有 VPS 操作权限",
			"尝试用 acme_finance 登录后，调用 POST /api/v1/vps/reinstall（body 含 vpsId）",
			"检查服务层的角色校验逻辑",
			"[编码提示] 也可使用 /api/v1/encoded/vps/reinstall（body 含 v）",
		},
		WriteUp:         "VPS 重装接口在 vulnerable 模式下没有检查用户角色，财务角色可以执行管理员才能做的操作。修复方案：增加角色权限校验，仅允许 admin 或 individual 操作。",
		Endpoint:        "/api/v1/vps/reinstall（body 含 vpsId）", Method: "POST",
		EncodedEndpoint: "/api/v1/encoded/vps/reinstall（body 含 v）",
	},
	{
		ID: "V-04", Title: "个人用户调用平台管理接口", Category: "垂直越权", Difficulty: 2,
		Description: "使用个人账户，调用平台管理员接口获取全平台数据。",
		Hints: []string{
			"平台管理接口应该仅限 platform_admin 访问",
			"尝试用 alice 登录后，调用 GET /api/v1/admin/users",
			"检查 admin 中间件的权限逻辑",
		},
		WriteUp:   "平台管理接口在 vulnerable 模式下，RequireAdmin 中间件跳过了权限检查，普通用户可以访问所有管理接口。修复方案：在 secure 模式下强制校验 platform_admin 角色。",
		Endpoint:  "/api/v1/admin/users", Method: "GET",
	},
	{
		ID: "V-05", Title: "运维人员自升为管理员", Category: "垂直越权", Difficulty: 3,
		Description: "使用运维角色，调用角色修改接口将自己的角色提升为企业管理员。",
		Hints: []string{
			"角色修改接口应该只有企业管理员可用",
			"尝试用 acme_ops 登录后，调用 PUT /api/v1/company/members/role（body 含 id 为自身 ID）",
			"将 role 字段改为 admin",
			"[编码提示] 也可使用 /api/v1/encoded/company/members/role（body 含 v）",
		},
		WriteUp:         "角色修改接口在 vulnerable 模式下没有校验调用者权限和目标用户关系，用户可以修改自身角色实现权限提升。修复方案：增加管理员权限校验，且禁止普通用户修改自身角色。",
		Endpoint:        "/api/v1/company/members/role（body 含 id）", Method: "PUT",
		EncodedEndpoint: "/api/v1/encoded/company/members/role（body 含 v）",
	},
	{
		ID: "C-01", Title: "跨企业操作 VPS", Category: "上下文越权", Difficulty: 2,
		Description: "使用 A 企业的用户，操作 B 企业的 VPS 实例。",
		Hints: []string{
			"不同企业的资源应该完全隔离",
			"尝试用 acme_ops 登录后，操作 Globex Inc 的 VPS（ID 4-5）",
			"检查企业边界校验逻辑",
			"[编码提示] 也可使用 /api/v1/encoded/vps/stop（body 含 v）",
		},
		WriteUp:         "VPS 操作接口在 vulnerable 模式下没有严格校验企业归属边界，一个企业的用户可以操作其他企业的 VPS。修复方案：增加 company_id 匹配校验。",
		Endpoint:        "/api/v1/vps/stop（body 含 vpsId）", Method: "POST",
		EncodedEndpoint: "/api/v1/encoded/vps/stop（body 含 v）",
	},
	{
		ID: "C-02", Title: "个人用户创建企业成员", Category: "上下文越权", Difficulty: 3,
		Description: "使用个人用户身份，调用添加企业成员接口创建企业用户。",
		Hints: []string{
			"个人用户不属于任何企业",
			"尝试用 alice 登录后，调用 POST /api/v1/company/members",
			"检查 userType 和业务上下文校验",
			"[编码提示] 也可使用 /api/v1/encoded/company/members",
		},
		WriteUp:         "添加企业成员接口在 vulnerable 模式下没有校验用户的 userType，个人账户可以创建企业成员。修复方案：增加 userType == 'company' 校验。",
		Endpoint:        "/api/v1/company/members", Method: "POST",
		EncodedEndpoint: "/api/v1/encoded/company/members",
	},
	{
		ID: "C-03", Title: "已吊销 API Key 仍可访问", Category: "上下文越权", Difficulty: 3,
		Description: "使用已被吊销或过期的 API Key，仍然能够通过认证访问 API。",
		Hints: []string{
			"API Key 有 status 和 expire_at 字段",
			"系统中有一个已被吊销的 API Key（ID 3，globex 用户创建）",
			"尝试使用已吊销的 key 进行认证",
		},
		WriteUp:   "认证中间件在 vulnerable 模式下只检查 API Key 是否存在，没有验证 status 是否为 active 以及 expire_at 是否过期。修复方案：在 secure 模式下增加状态和有效期校验。",
		Endpoint:  "/api/v1/vps", Method: "GET",
	},
	// ---- 编码加密挑战 (E 系列) ----
	{
		ID: "E-01", Title: "Base64 编码 VPS ID 越权", Category: "编码越权", Difficulty: 1,
		Description: "VPS ID 使用 Base64 编码传输，解码后即可利用 IDOR 漏洞查看他人 VPS 详情。",
		Hints: []string{
			"观察请求参数是否为 Base64 编码",
			"使用 X-Encoding-Type: base64 头调用 /api/v1/encoded/vps?v=",
			"例如：将 ID=2 编码为 Base64（'2' → 'Mg=='），然后请求 GET /api/v1/encoded/vps?v=Mg==",
			"可使用 /api/v1/crypto/encode 辅助接口进行编码",
		},
		WriteUp:         "Base64 是最常见的编码方式，仅用于数据表示而非加密。攻击者识别出 Base64 编码后即可解码参数并构造越权请求。",
		Endpoint:        "/api/v1/encoded/vps?v=", Method: "GET",
		EncodingType:    "base64",
		EncodedEndpoint: "/api/v1/encoded/vps?v=",
	},
	{
		ID: "E-02", Title: "Base32 编码用户 ID 越权", Category: "编码越权", Difficulty: 1,
		Description: "用户 ID 使用 Base32 编码（无填充），需要识别编码类型并解码后实施越权。",
		Hints: []string{
			"Base32 编码仅包含 A-Z 和 2-7 字符，注意与 Base64 区分",
			"使用 X-Encoding-Type: base32 头调用 /api/v1/encoded/users?v=",
			"例如：ID=5 → 'OBXW4===', 无填充格式为 'OBXW4'",
			"可使用 /api/v1/crypto/encode 辅助接口进行编码",
		},
		WriteUp:         "Base32 编码使用 32 个字符而非 64 个字符，常用于需要大小写不敏感的传输场景。攻击者需要先识别编码类型。",
		Endpoint:        "/api/v1/encoded/users?v=", Method: "GET",
		EncodingType:    "base32",
		EncodedEndpoint: "/api/v1/encoded/users?v=",
	},
	{
		ID: "E-03", Title: "凯撒密码订单 ID 越权", Category: "编码越权", Difficulty: 1,
		Description: "订单 ID 经过凯撒密码（shift=3）加密传输。需要发现移位规律并解码。",
		Hints: []string{
			"观察数字是否有规律性的偏移",
			"凯撒密码默认偏移量为 3（shift=3）",
			"使用 X-Encoding-Type: caesar 头调用 /api/v1/encoded/orders?v=",
			"例如：ID=4，凯撒编码(shift=3)后变为 7",
		},
		WriteUp:         "凯撒密码是最经典的单表替换密码。对数字应用凯撒移位后，凭肉眼即可发现规律（ID本身为数字，偏移后仍是数字）。",
		Endpoint:        "/api/v1/encoded/orders?v=", Method: "GET",
		EncodingType:    "caesar",
		EncodedEndpoint: "/api/v1/encoded/orders?v=",
	},
	{
		ID: "E-04", Title: "自定义 Base64 编码表绕过", Category: "编码越权", Difficulty: 2,
		Description: "API 使用了自定义字符表的 Base64（A↔Z 互换），需要逆向识别自定义编码表。",
		Hints: []string{
			"观察编码后字符是否与标准 Base64 相似但不同",
			"自定义 Base64 字符表：A↔Z、a↔z、0↔9、+↔/ 互换",
			"使用 X-Encoding-Type: custom_base64 头调用 /api/v1/encoded/vps?v=",
			"先尝试标准 Base64 解码，如失败则考虑自定义编码表",
		},
		WriteUp:         "自定义 Base64 通过修改标准字符表实现混淆。攻击者需要对比标准编码输出与自定义编码输出，发现字符表差异后构造解码器。",
		Endpoint:        "/api/v1/encoded/vps?v=", Method: "GET",
		EncodingType:    "custom_base64",
		EncodedEndpoint: "/api/v1/encoded/vps?v=",
	},
	{
		ID: "E-05", Title: "多层嵌套编码工单越权", Category: "编码越权", Difficulty: 2,
		Description: "工单 ID 经过 Base64 → Base32 双层嵌套编码，需要反向链式解码。",
		Hints: []string{
			"参数可能经过多次编码，尝试逐层解码",
			"解码链：Base64 解码 → Base32 解码 → 得到原始 ID",
			"没有对应的单一 X-Encoding-Type，需要手动逐层解码",
			"可使用 /api/v1/crypto/decode 接口逐步解码",
		},
		WriteUp:         "多层嵌套编码常用于绕过简单的安全过滤。攻击者需要逐层识别并解码，最终到达原始参数。",
		Endpoint:        "/api/v1/encoded/tickets?v=", Method: "GET",
		EncodingType:    "multi",
		EncodedEndpoint: "/api/v1/encoded/tickets?v=",
	},
	{
		ID: "E-06", Title: "AES 加密 VPS ID 操作", Category: "加密越权", Difficulty: 3,
		Description: "VPS ID 使用 AES-256-GCM 加密传输。需要获取加密密钥后解密，再进行越权操作。",
		Hints: []string{
			"AES 密钥可通过 /api/v1/crypto/keys 接口获取（aes_key_base64 字段）",
			"密钥编码方式为 Base64，加密模式为 AES-256-GCM",
			"加密数据格式：Base64(Nonce + Ciphertext)，Nonce 长度 12 字节",
			"使用 X-Encoding-Type: aes 头调用 /api/v1/encoded/vps/start（body 含 v）",
			"先解密得到原始 ID，再用原始 ID 构造请求或直接用加密后的值并设置正确的 X-Encoding-Type",
		},
		WriteUp:         "AES-GCM 提供认证加密。攻击者获取密钥后可以解密任意参数，这是现实攻击中常见的场景——密钥通过不安全的方式泄露。",
		Endpoint:        "/api/v1/encoded/vps/start（body 含 v）", Method: "POST",
		EncodingType:    "aes",
		EncodedEndpoint: "/api/v1/encoded/vps/start（body 含 v）",
	},
	{
		ID: "E-07", Title: "HMAC 签名验证绕过", Category: "加密越权", Difficulty: 3,
		Description: "参数附带 HMAC-SHA256 签名保护。需要获取签名密钥后伪造签名，实现越权。",
		Hints: []string{
			"参数格式：Base64(data).HMAC_Hex(data)，点为分隔符",
			"HMAC 密钥可通过 /api/v1/crypto/keys 接口获取（hmac_key_base64 字段）",
			"签名算法为 HMAC-SHA256，签名原始值为未编码的 data 字节",
			"使用 X-Encoding-Type: signed 头调用接口",
		},
		WriteUp:         "HMAC 签名可确保数据完整性和真实性。但如果签名密钥泄露，攻击者可以为任意参数生成有效签名，从而绕过完整性检查。",
		Endpoint:        "/api/v1/encoded/vps/stop（body 含 v）", Method: "POST",
		EncodingType:    "signed",
		EncodedEndpoint: "/api/v1/encoded/vps/stop（body 含 v）",
	},
	{
		ID: "E-08", Title: "SM4 国密加密跨企业操作", Category: "加密越权", Difficulty: 4,
		Description: "VPS 操作参数使用国密 SM4-CBC 加密，还需要跨企业进行上下文越权。组合加密破解与企业边界突破。",
		Hints: []string{
			"SM4 密钥可通过 /api/v1/crypto/keys 接口获取（sm4_key_base64 字段）",
			"SM4 加密模式为 CBC，需要正确处理 IV（前 16 字节）和 PKCS7 填充",
			"使用 X-Encoding-Type: sm4 头调用 /api/v1/encoded/vps/reinstall（body 含 v）",
			"破解加密后，还需要利用上下文越权（跨企业）才能完成挑战",
			"尝试用 acme_ops 登录，操作 Globex Inc 的 VPS（ID 4-5）",
		},
		WriteUp:         "SM4 是中国国家密码标准。挑战结合了国密算法破解和跨企业越权双重步骤。攻击者需要：1) 获取 SM4 密钥 2) 解密参数 3) 编码 Globex VPS ID 4) 利用上下文越权（C-01）操作其他企业的 VPS。",
		Endpoint:        "/api/v1/encoded/vps/reinstall（body 含 v）", Method: "POST",
		EncodingType:    "sm4",
		EncodedEndpoint: "/api/v1/encoded/vps/reinstall（body 含 v）",
	},
	{
		ID: "E-09", Title: "简单哈希签名绕过", Category: "加密越权", Difficulty: 2,
		Description: "所有请求参数附带 MD5 哈希签名，后端校验参数完整性。需要发现签名机制并伪造签名实现越权。",
		Hints: []string{
			"观察请求参数格式是否为 value:hash 形式",
			"哈希算法为 MD5，格式：原始值 + ':' + md5(原始值 + salt)",
			"salt 可通过源码分析或 /api/v1/crypto/keys 接口获取（hash_salt 字段）",
			"使用 X-Encoding-Type: hash 头调用任意业务接口",
			"尝试修改参数值后重新计算 MD5 哈希，保持签名一致",
		},
		WriteUp:         "简单哈希签名通过将参数值与其 MD5 哈希绑定，防止参数被篡改。但如果 salt 固定且可预测，攻击者可以轻易重新计算哈希并伪造任意参数值。修复方案：使用 HMAC 替代简单 MD5，或确保 salt 随机且不可预测。",
		Endpoint:        "/api/v1/vps/detail?vpsId=", Method: "GET",
		EncodingType:    "hash",
		EncodedEndpoint: "/api/v1/vps/detail?vpsId=",
	},
}

type ChallengeProgress struct {
	ChallengeID string `json:"challenge_id"`
	Completed   bool   `json:"completed"`
}

var progressStore = make(map[string]bool)

func MarkChallengeCompleted(id string) {
	progressStore[id] = true
}

func IsChallengeCompleted(id string) bool {
	return progressStore[id]
}

func ResetProgress() {
	progressStore = make(map[string]bool)
}
