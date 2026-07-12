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
			},
		WriteUp:         "后端 GetDetail 方法在 vulnerable 模式下没有检查 VPS 所有权或企业归属，仅验证 VPS 是否存在。攻击者可通过遍历 ID 查看任意 VPS 详情。修复方案：在 secure 模式下增加 ownership/company 校验。",
		Endpoint:        "/api/v1/vps?vpsId=", Method: "GET",
	},
	{
		ID: "H-02", Title: "查看他人个人信息", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改用户 ID 参数，查看其他用户的个人资料。",
		Hints: []string{
			"注意获取用户信息的接口使用了查询参数",
			"尝试访问 GET /api/v1/users?id= 并修改 id 值",
			"系统中预置了 9 个用户，ID 为 1-9",
		},
		WriteUp:         "GetUserByID 接口在 vulnerable 模式下没有校验当前用户与目标用户的关系，允许查看任意用户资料。修复方案：限制用户只能查看自己的资料（平台管理员除外）。",
		Endpoint:        "/api/v1/users?id=", Method: "GET",
	},
	{
		ID: "H-03", Title: "查看他人订单", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改订单 ID，查看其他用户的购买/续费订单详情。",
		Hints: []string{
			"订单详情接口使用了订单 ID 作为查询参数",
			"尝试访问 GET /api/v1/orders?orderId= 并修改 orderId",
			"系统中有 4 个预置订单，ID 为 1-4",
		},
		WriteUp:         "订单详情查询没有校验订单归属用户，在 vulnerable 模式下允许查看任意订单。修复方案：增加 user_id 匹配校验。",
		Endpoint:        "/api/v1/orders?orderId=", Method: "GET",
	},
	{
		ID: "H-04", Title: "查看/回复他人工单", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改工单 ID，查看或回复其他用户的工单。",
		Hints: []string{
			"工单详情和回复接口都使用了工单 ID",
			"尝试访问 GET /api/v1/tickets?ticketId= 和 POST /api/v1/tickets/reply（body 含 ticketId）",
			"系统中有 2 个预置工单，ID 为 1-2",
		},
		WriteUp:         "工单详情和回复接口没有校验工单创建者，允许查看和回复他人工单。修复方案：增加 user_id 匹配校验。",
		Endpoint:        "/api/v1/tickets?ticketId=", Method: "GET",
	},
	{
		ID: "H-05", Title: "删除他人 API Key", Category: "水平越权 (IDOR)", Difficulty: 2,
		Description: "通过修改 API Key ID，删除其他用户创建的 API Key。",
		Hints: []string{
			"删除 API Key 的接口使用了 Key ID",
			"尝试发送 DELETE /api/v1/apikeys 请求（body 含 id），使用其他用户的 Key ID",
			"系统中预置了 3 个 API Key，ID 为 1-3",
		},
		WriteUp:         "API Key 删除接口没有校验 Key 归属，在 vulnerable 模式下允许删除任意 Key。修复方案：增加 user_id 匹配校验。",
		Endpoint:        "/api/v1/apikeys（body 含 id）", Method: "DELETE",
	},
	{
		ID: "V-01", Title: "只读成员控制 VPS 启停", Category: "垂直越权", Difficulty: 1,
		Description: "使用只读成员 (viewer) 账户，通过 API 直接调用 VPS 启停接口。",
		Hints: []string{
			"前端界面隐藏了 viewer 角色的操作按钮",
			"尝试用 acme_viewer 登录后，直接调用 POST /api/v1/vps/start（body 含 vpsId）",
			"前端隐藏不等于后端限制",
		},
		WriteUp:         "VPS 控制接口在 vulnerable 模式下没有校验用户角色，viewer 角色本应只有查看权限，但可以执行启停操作。修复方案：在 secure 模式下增加角色权限检查（require operator 或 admin）。",
		Endpoint:        "/api/v1/vps/start（body 含 vpsId）", Method: "POST",
	},
	{
		ID: "V-02", Title: "运维人员添加企业成员", Category: "垂直越权", Difficulty: 2,
		Description: "使用运维 (operator) 角色，调用添加成员接口创建新企业用户。",
		Hints: []string{
			"添加成员接口应该只有企业管理员可用",
			"尝试用 acme_ops 登录后，调用 POST /api/v1/company/members",
			"检查 RBAC 中间件和服务层双重校验",
		},
		WriteUp:         "添加企业成员接口在 vulnerable 模式下没有严格校验用户角色，operator 角色可以执行管理员操作。修复方案：在 middleware 和 service 层同时增加角色检查。",
		Endpoint:        "/api/v1/company/members", Method: "POST",
	},
	{
		ID: "V-03", Title: "财务人员重装 VPS 系统", Category: "垂直越权", Difficulty: 2,
		Description: "使用财务 (finance) 角色，调用 VPS 重装系统接口。",
		Hints: []string{
			"财务人员不应该有 VPS 操作权限",
			"尝试用 acme_finance 登录后，调用 POST /api/v1/vps/reinstall（body 含 vpsId）",
			"检查服务层的角色校验逻辑",
		},
		WriteUp:         "VPS 重装接口在 vulnerable 模式下没有检查用户角色，财务角色可以执行管理员才能做的操作。修复方案：增加角色权限校验，仅允许 admin 或 individual 操作。",
		Endpoint:        "/api/v1/vps/reinstall（body 含 vpsId）", Method: "POST",
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
		},
		WriteUp:         "角色修改接口在 vulnerable 模式下没有校验调用者权限和目标用户关系，用户可以修改自身角色实现权限提升。修复方案：增加管理员权限校验，且禁止普通用户修改自身角色。",
		Endpoint:        "/api/v1/company/members/role（body 含 id）", Method: "PUT",
	},
	{
		ID: "C-01", Title: "跨企业操作 VPS", Category: "上下文越权", Difficulty: 2,
		Description: "使用 A 企业的用户，操作 B 企业的 VPS 实例。",
		Hints: []string{
			"不同企业的资源应该完全隔离",
			"尝试用 acme_ops 登录后，操作 Globex Inc 的 VPS（ID 4-5）",
			"检查企业边界校验逻辑",
		},
		WriteUp:         "VPS 操作接口在 vulnerable 模式下没有严格校验企业归属边界，一个企业的用户可以操作其他企业的 VPS。修复方案：增加 company_id 匹配校验。",
		Endpoint:        "/api/v1/vps/stop（body 含 vpsId）", Method: "POST",
	},
	{
		ID: "C-02", Title: "个人用户创建企业成员", Category: "上下文越权", Difficulty: 3,
		Description: "使用个人用户身份，调用添加企业成员接口创建企业用户。",
		Hints: []string{
			"个人用户不属于任何企业",
			"尝试用 alice 登录后，调用 POST /api/v1/company/members",
			"检查 userType 和业务上下文校验",
		},
		WriteUp:         "添加企业成员接口在 vulnerable 模式下没有校验用户的 userType，个人账户可以创建企业成员。修复方案：增加 userType == 'company' 校验。",
		Endpoint:        "/api/v1/company/members", Method: "POST",
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
	{
		ID: "E-01", Title: "Base64编码参数绕过", Category: "编码加密", Difficulty: 1,
		Description: "API参数使用标准Base64编码，识别编码特征并解码参数值，绕过前端校验。",
		Hints: []string{
			"注意参数值末尾是否有 = 或 == 填充字符",
			"Base64编码字符集：A-Z, a-z, 0-9, +, /",
			"激活此挑战后，前端所有请求的参数值会自动使用Base64编码",
			"使用 CyberChef 或 Burp Decoder 解码参数值",
		},
		WriteUp:         "API使用Base64编码保护参数传输，但攻击者可以通过简单的解码工具恢复原始参数值。修复方案：使用加密而非编码保护敏感参数，配合签名确保数据完整性。",
		Endpoint:        "/api/v1/vps/detail?vpsId=", Method: "GET",
		EncodingType:    "base64",
		EncodedEndpoint: "/api/v1/vps/detail?vpsId=MQ==",
	},
	{
		ID: "E-02", Title: "Base32编码参数绕过", Category: "编码加密", Difficulty: 1,
		Description: "API参数使用Base32编码，需识别编码类型并解码以发现越权漏洞。",
		Hints: []string{
			"Base32仅使用A-Z和2-7共32个字符",
			"注意参数值是否只包含大写字母和数字2-7",
			"激活此挑战后，前端请求参数值会自动使用Base32编码",
		},
		WriteUp:         "API使用Base32编码保护参数，但攻击者可识别编码类型并使用在线工具解码。修复方案：编码不是安全机制，需配合访问控制。",
		Endpoint:        "/api/v1/users?id=", Method: "GET",
		EncodingType:    "base32",
		EncodedEndpoint: "/api/v1/users?id=GE======",
	},
	{
		ID: "E-03", Title: "凯撒密码参数偏移绕过", Category: "编码加密", Difficulty: 2,
		Description: "API参数使用凯撒密码（ROT3）进行混淆，识别密码类型并解出原始参数值。",
		Hints: []string{
			"凯撒密码对字母和数字均按固定偏移量移位",
			"默认偏移量为3（ROT3），例如 1→4, A→D",
			"激活此挑战后，前端请求参数值会自动使用凯撒密码编码",
		},
		WriteUp:         "凯撒密码作为古典密码，安全性极低，可使用暴力枚举或频率分析破解。修复方案：使用现代加密算法替代古典密码。",
		Endpoint:        "/api/v1/orders/detail?orderId=", Method: "GET",
		EncodingType:    "caesar",
		EncodedEndpoint: "/api/v1/orders/detail?orderId=4",
	},
	{
		ID: "E-04", Title: "自定义Base64编码表参数绕过", Category: "编码加密", Difficulty: 2,
		Description: "API使用自定义字符表的Base64变体编码，通过对比已知明文和密文推断编码表映射关系。",
		Hints: []string{
			"自定义Base64使用不同的字符排列顺序",
			"尝试对比标准Base64输出与目标系统输出的差异",
			"已知自己的用户ID，发送请求观察编码结果，建立映射表",
			"激活此挑战后，前端请求参数值会自动使用自定义Base64编码",
		},
		WriteUp:         "自定义Base64通过改变字符表顺序进行混淆，但攻击者可通过对比已知明文与密文来恢复映射表。修复方案：使用加密和签名保护参数。",
		Endpoint:        "/api/v1/tickets/detail?ticketId=", Method: "GET",
		EncodingType:    "custom-base64",
		EncodedEndpoint: "/api/v1/tickets/detail?ticketId=MTE=",
	},
	{
		ID: "E-05", Title: "多层嵌套编码参数绕过", Category: "编码加密", Difficulty: 3,
		Description: "API参数经过多层编码（Base64→Base32），需逐层剥离解码以恢复原始参数值。",
		Hints: []string{
			"参数值可能经过两层或更多层编码",
			"从最外层开始逐层解码，每层可能有不同的编码类型",
			"常见组合：Base64→Base32, Hex→Base64",
			"激活此挑战后，前端请求参数值会自动使用多层嵌套编码",
		},
		WriteUp:         "多层嵌套编码增加了攻击难度，但仍可通过逐层解码工具恢复原始参数。修复方案：使用签名验证确保参数完整性，配合严格的访问控制。",
		Endpoint:        "/api/v1/vps/start", Method: "POST",
		EncodingType:    "multi",
		EncodedEndpoint: "/api/v1/vps/start (body: {\"vpsId\":\"base32(base64(1))\"})",
	},
	{
		ID: "E-06", Title: "AES加密参数绕过", Category: "编码加密", Difficulty: 3,
		Description: "API参数使用AES-256-GCM加密，需要获取密钥才能解密参数值以进行越权攻击。",
		Hints: []string{
			"AES加密结果通常以Base64编码传输",
			"检查前端源码、错误信息、debug接口寻找密钥",
			"AES-GCM模式包含认证标签，确保数据完整性和机密性",
			"提示：密钥可能硬编码在前端JavaScript中",
		},
		WriteUp:         "AES加密提供强机密性保护，但如果密钥泄露（如前端硬编码），攻击者仍可解密参数。修复方案：密钥托管在安全的后端环境，不在前端暴露。",
		Endpoint:        "/api/v1/vps/stop", Method: "POST",
		EncodingType:    "aes",
		EncodedEndpoint: "/api/v1/vps/stop (body: {\"vpsId\":\"<AES_encrypted>\"})",
	},
	{
		ID: "E-07", Title: "HMAC签名验证绕过", Category: "编码加密", Difficulty: 4,
		Description: "API使用HMAC-SHA256签名保护请求完整性，在X-HMAC-Sign头中传递签名值。攻击者需获取签名密钥才能伪造合法请求。",
		Hints: []string{
			"检查HTTP请求头中是否有 X-HMAC-Sign 字段",
			"HMAC使用密钥和哈希函数生成签名，密钥长度通常为32字节",
			"签名原始数据为请求参数排序后的规范字符串",
			"如果密钥泄露，攻击者可伪造任意参数的签名",
		},
		WriteUp:         "HMAC-SHA256提供高强度的数据完整性保护，但密钥安全存储是关键。修复方案：使用安全密钥管理服务，定期轮换密钥。",
		Endpoint:        "/api/v1/vps/stop", Method: "POST",
		EncodingType:    "hmac",
		EncodedEndpoint: "/api/v1/vps/stop (header: X-HMAC-Sign: <signature>)",
	},
	{
		ID: "E-08", Title: "SM4国密加密参数绕过", Category: "编码加密", Difficulty: 4,
		Description: "API参数使用SM4-CBC模式加密（国密标准），需获取密钥和IV才能解密参数值以进行越权测试。",
		Hints: []string{
			"SM4是中国国家密码标准，128位密钥，CBC模式需要IV",
			"加密结果以Base64编码传输",
			"检查前端代码中是否暴露SM4密钥或IV",
			"提示：SM4密钥可能在初始化配置中",
		},
		WriteUp:         "SM4是中国国家标准密码算法，安全性对标AES。但密钥泄露将导致加密保护失效。修复方案：密钥后端集中管理，不在前端传输。",
		Endpoint:        "/api/v1/vps/start", Method: "POST",
		EncodingType:    "sm4",
		EncodedEndpoint: "/api/v1/vps/start (body: {\"vpsId\":\"<SM4_encrypted>\"})",
	},
	{
		ID: "E-09", Title: "MD5哈希签名验证绕过", Category: "编码加密", Difficulty: 3,
		Description: "API使用value:md5(value|salt)格式的哈希签名保护参数完整性。攻击者需获取盐值(salt)才能伪造参数的签名。",
		Hints: []string{
			"参数格式为 value:32位hex哈希值",
			"MD5哈希不可逆，但可暴力破解简单盐值",
			"签名计算方式为：hash = md5(value + '|' + salt)",
			"如果盐值固定且可猜解，攻击者可以重新计算哈希伪造参数",
			"提示：检查前端源码中是否暴露了盐值",
		},
		WriteUp:         "MD5哈希签名提供基本的完整性校验，但MD5已被证明存在碰撞漏洞，且固定盐值可被暴力破解。修复方案：使用HMAC-SHA256替代MD5签名，盐值随机生成并后端存储。",
		Endpoint:        "/api/v1/vps/stop", Method: "POST",
		EncodingType:    "hash-sign",
		EncodedEndpoint: "/api/v1/vps/stop (body: {\"vpsId\":\"1:5d41402abc4b2a76b9719d911017c592\"})",
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
