package vuln

type Challenge struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Difficulty  int      `json:"difficulty"`
	Description string   `json:"description"`
	Hints       []string `json:"hints"`
	WriteUp     string   `json:"writeup"`
	Endpoint    string   `json:"endpoint"`
	Method      string   `json:"method"`
}

var Challenges = []Challenge{
	{
		ID: "H-01", Title: "查看他人 VPS 详情", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改 VPS ID 参数，查看其他用户或企业的 VPS 详细信息。",
		Hints: []string{
			"注意请求中的 VPS ID 参数",
			"尝试修改 GET /api/v1/vps/:id 中的 id 值，使用其他用户的 VPS ID",
			"Acme Corp 的 VPS ID 为 1-3，Globex Inc 的 VPS ID 为 4-5",
		},
		WriteUp:   "后端 GetDetail 方法在 vulnerable 模式下没有检查 VPS 所有权或企业归属，仅验证 VPS 是否存在。攻击者可通过遍历 ID 查看任意 VPS 详情。修复方案：在 secure 模式下增加 ownership/company 校验。",
		Endpoint:  "/api/v1/vps/:id", Method: "GET",
	},
	{
		ID: "H-02", Title: "查看他人个人信息", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改用户 ID 参数，查看其他用户的个人资料。",
		Hints: []string{
			"注意获取用户信息的接口使用了路径参数",
			"尝试访问 GET /api/v1/users/:id 并修改 id 值",
			"系统中预置了 9 个用户，ID 为 1-9",
		},
		WriteUp:   "GetUserByID 接口在 vulnerable 模式下没有校验当前用户与目标用户的关系，允许查看任意用户资料。修复方案：限制用户只能查看自己的资料（平台管理员除外）。",
		Endpoint:  "/api/v1/users/:id", Method: "GET",
	},
	{
		ID: "H-03", Title: "查看他人订单", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改订单 ID，查看其他用户的购买/续费订单详情。",
		Hints: []string{
			"订单详情接口使用了订单 ID 作为路径参数",
			"尝试访问 GET /api/v1/orders/:id 并修改 id",
			"系统中有 4 个预置订单，ID 为 1-4",
		},
		WriteUp:   "订单详情查询没有校验订单归属用户，在 vulnerable 模式下允许查看任意订单。修复方案：增加 user_id 匹配校验。",
		Endpoint:  "/api/v1/orders/:id", Method: "GET",
	},
	{
		ID: "H-04", Title: "查看/回复他人工单", Category: "水平越权 (IDOR)", Difficulty: 1,
		Description: "通过修改工单 ID，查看或回复其他用户的工单。",
		Hints: []string{
			"工单详情和回复接口都使用了工单 ID",
			"尝试访问 GET /api/v1/tickets/:id 和 POST /api/v1/tickets/:id/reply",
			"系统中有 2 个预置工单，ID 为 1-2",
		},
		WriteUp:   "工单详情和回复接口没有校验工单创建者，允许查看和回复他人工单。修复方案：增加 user_id 匹配校验。",
		Endpoint:  "/api/v1/tickets/:id", Method: "GET",
	},
	{
		ID: "H-05", Title: "删除他人 API Key", Category: "水平越权 (IDOR)", Difficulty: 2,
		Description: "通过修改 API Key ID，删除其他用户创建的 API Key。",
		Hints: []string{
			"删除 API Key 的接口使用了 Key ID",
			"尝试发送 DELETE /api/v1/apikeys/:id 请求，使用其他用户的 Key ID",
			"系统中预置了 3 个 API Key，ID 为 1-3",
		},
		WriteUp:   "API Key 删除接口没有校验 Key 归属，在 vulnerable 模式下允许删除任意 Key。修复方案：增加 user_id 匹配校验。",
		Endpoint:  "/api/v1/apikeys/:id", Method: "DELETE",
	},
	{
		ID: "V-01", Title: "只读成员控制 VPS 启停", Category: "垂直越权", Difficulty: 1,
		Description: "使用只读成员 (viewer) 账户，通过 API 直接调用 VPS 启停接口。",
		Hints: []string{
			"前端界面隐藏了 viewer 角色的操作按钮",
			"尝试用 acme_viewer 登录后，直接调用 POST /api/v1/vps/:id/start",
			"前端隐藏不等于后端限制",
		},
		WriteUp:   "VPS 控制接口在 vulnerable 模式下没有校验用户角色，viewer 角色本应只有查看权限，但可以执行启停操作。修复方案：在 secure 模式下增加角色权限检查（require operator 或 admin）。",
		Endpoint:  "/api/v1/vps/:id/start", Method: "POST",
	},
	{
		ID: "V-02", Title: "运维人员添加企业成员", Category: "垂直越权", Difficulty: 2,
		Description: "使用运维 (operator) 角色，调用添加成员接口创建新企业用户。",
		Hints: []string{
			"添加成员接口应该只有企业管理员可用",
			"尝试用 acme_ops 登录后，调用 POST /api/v1/company/members",
			"检查 RBAC 中间件和服务层双重校验",
		},
		WriteUp:   "添加企业成员接口在 vulnerable 模式下没有严格校验用户角色，operator 角色可以执行管理员操作。修复方案：在 middleware 和 service 层同时增加角色检查。",
		Endpoint:  "/api/v1/company/members", Method: "POST",
	},
	{
		ID: "V-03", Title: "财务人员重装 VPS 系统", Category: "垂直越权", Difficulty: 2,
		Description: "使用财务 (finance) 角色，调用 VPS 重装系统接口。",
		Hints: []string{
			"财务人员不应该有 VPS 操作权限",
			"尝试用 acme_finance 登录后，调用 POST /api/v1/vps/:id/reinstall",
			"检查服务层的角色校验逻辑",
		},
		WriteUp:   "VPS 重装接口在 vulnerable 模式下没有检查用户角色，财务角色可以执行管理员才能做的操作。修复方案：增加角色权限校验，仅允许 admin 或 individual 操作。",
		Endpoint:  "/api/v1/vps/:id/reinstall", Method: "POST",
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
			"尝试用 acme_ops 登录后，调用 PUT /api/v1/company/members/:id/role（:id 为自身 ID）",
			"将 role 字段改为 admin",
		},
		WriteUp:   "角色修改接口在 vulnerable 模式下没有校验调用者权限和目标用户关系，用户可以修改自身角色实现权限提升。修复方案：增加管理员权限校验，且禁止普通用户修改自身角色。",
		Endpoint:  "/api/v1/company/members/:id/role", Method: "PUT",
	},
	{
		ID: "C-01", Title: "跨企业操作 VPS", Category: "上下文越权", Difficulty: 2,
		Description: "使用 A 企业的用户，操作 B 企业的 VPS 实例。",
		Hints: []string{
			"不同企业的资源应该完全隔离",
			"尝试用 acme_ops 登录后，操作 Globex Inc 的 VPS（ID 4-5）",
			"检查企业边界校验逻辑",
		},
		WriteUp:   "VPS 操作接口在 vulnerable 模式下没有严格校验企业归属边界，一个企业的用户可以操作其他企业的 VPS。修复方案：增加 company_id 匹配校验。",
		Endpoint:  "/api/v1/vps/:id/stop", Method: "POST",
	},
	{
		ID: "C-02", Title: "个人用户创建企业成员", Category: "上下文越权", Difficulty: 3,
		Description: "使用个人用户身份，调用添加企业成员接口创建企业用户。",
		Hints: []string{
			"个人用户不属于任何企业",
			"尝试用 alice 登录后，调用 POST /api/v1/company/members",
			"检查 userType 和业务上下文校验",
		},
		WriteUp:   "添加企业成员接口在 vulnerable 模式下没有校验用户的 userType，个人账户可以创建企业成员。修复方案：增加 userType == 'company' 校验。",
		Endpoint:  "/api/v1/company/members", Method: "POST",
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
