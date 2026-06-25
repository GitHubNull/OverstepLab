package seed

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil // Already seeded
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Create companies
		comp1 := model.Company{
			ID: 1, Name: "Acme Corp", LicenseNo: "LIC-2024-001",
			ContactName: "John Smith", ContactPhone: "13800001001",
			Balance: 50000, Status: "active",
		}
		comp2 := model.Company{
			ID: 2, Name: "Globex Inc", LicenseNo: "LIC-2024-002",
			ContactName: "Jane Doe", ContactPhone: "13800001002",
			Balance: 30000, Status: "active",
		}
		tx.Create(&comp1)
		tx.Create(&comp2)

		// Create users
		users := []model.User{
			{ID: 1, Username: "admin", PasswordHash: mustHash("admin123"), Email: "admin@cloudnest.com", UserType: "platform_admin", Status: "active"},
			{ID: 2, Username: "acme_admin", PasswordHash: mustHash("pass123"), Email: "admin@acme.com", UserType: "company", CompanyID: ptrUint(1), Role: "admin", Status: "active"},
			{ID: 3, Username: "acme_ops", PasswordHash: mustHash("pass123"), Email: "ops@acme.com", UserType: "company", CompanyID: ptrUint(1), Role: "operator", Status: "active"},
			{ID: 4, Username: "acme_finance", PasswordHash: mustHash("pass123"), Email: "finance@acme.com", UserType: "company", CompanyID: ptrUint(1), Role: "finance", Status: "active"},
			{ID: 5, Username: "acme_viewer", PasswordHash: mustHash("pass123"), Email: "viewer@acme.com", UserType: "company", CompanyID: ptrUint(1), Role: "viewer", Status: "active"},
			{ID: 6, Username: "globex_admin", PasswordHash: mustHash("pass123"), Email: "admin@globex.com", UserType: "company", CompanyID: ptrUint(2), Role: "admin", Status: "active"},
			{ID: 7, Username: "globex_ops", PasswordHash: mustHash("pass123"), Email: "ops@globex.com", UserType: "company", CompanyID: ptrUint(2), Role: "operator", Status: "active"},
			{ID: 8, Username: "alice", PasswordHash: mustHash("pass123"), Email: "alice@example.com", UserType: "individual", Role: "individual", Status: "active"},
			{ID: 9, Username: "bob", PasswordHash: mustHash("pass123"), Email: "bob@example.com", UserType: "individual", Role: "individual", Status: "active"},
			{ID: 10, Username: "globex_finance", PasswordHash: mustHash("pass123"), Email: "finance@globex.com", UserType: "company", CompanyID: ptrUint(2), Role: "finance", Status: "active"},
			{ID: 11, Username: "globex_viewer", PasswordHash: mustHash("pass123"), Email: "viewer@globex.com", UserType: "company", CompanyID: ptrUint(2), Role: "viewer", Status: "active"},
		}
		for _, u := range users {
			tx.Create(&u)
		}

		// Create VPS instances
		expireAt := time.Now().AddDate(1, 0, 0)
		vpsList := []model.VPSInstance{
			{ID: 1, Name: "acme-web-01", OwnerID: 2, CompanyID: ptrUint(1), CPU: 4, Memory: 8192, Disk: 100, Bandwidth: 100, IPAddress: "10.0.1.10", OSImage: "ubuntu-22.04", Status: "running", ExpireAt: expireAt},
			{ID: 2, Name: "acme-db-01", OwnerID: 2, CompanyID: ptrUint(1), CPU: 8, Memory: 16384, Disk: 500, Bandwidth: 200, IPAddress: "10.0.1.11", OSImage: "centos-8", Status: "running", ExpireAt: expireAt},
			{ID: 3, Name: "acme-app-01", OwnerID: 3, CompanyID: ptrUint(1), CPU: 2, Memory: 4096, Disk: 50, Bandwidth: 50, IPAddress: "10.0.1.12", OSImage: "ubuntu-22.04", Status: "stopped", ExpireAt: expireAt},
			{ID: 4, Name: "globex-web-01", OwnerID: 6, CompanyID: ptrUint(2), CPU: 4, Memory: 8192, Disk: 200, Bandwidth: 100, IPAddress: "10.0.2.10", OSImage: "debian-12", Status: "running", ExpireAt: expireAt},
			{ID: 5, Name: "globex-api-01", OwnerID: 7, CompanyID: ptrUint(2), CPU: 2, Memory: 4096, Disk: 100, Bandwidth: 50, IPAddress: "10.0.2.11", OSImage: "ubuntu-22.04", Status: "running", ExpireAt: expireAt},
			{ID: 6, Name: "alice-personal-01", OwnerID: 8, CPU: 2, Memory: 4096, Disk: 50, Bandwidth: 50, IPAddress: "10.1.1.10", OSImage: "ubuntu-22.04", Status: "running", ExpireAt: expireAt},
			{ID: 7, Name: "bob-personal-01", OwnerID: 9, CPU: 1, Memory: 2048, Disk: 30, Bandwidth: 30, IPAddress: "10.1.2.10", OSImage: "centos-8", Status: "stopped", ExpireAt: expireAt},
			{ID: 8, Name: "globex-finance-01", OwnerID: 10, CompanyID: ptrUint(2), CPU: 2, Memory: 4096, Disk: 80, Bandwidth: 50, IPAddress: "10.0.2.20", OSImage: "ubuntu-22.04", Status: "running", ExpireAt: expireAt},
			{ID: 9, Name: "globex-viewer-01", OwnerID: 11, CompanyID: ptrUint(2), CPU: 1, Memory: 2048, Disk: 40, Bandwidth: 30, IPAddress: "10.0.2.21", OSImage: "debian-12", Status: "running", ExpireAt: expireAt},
		}
		for _, v := range vpsList {
			tx.Create(&v)
		}

		// Create orders
		orders := []model.Order{
			{ID: 1, OrderNo: "ORD20240101001", UserID: 2, CompanyID: ptrUint(1), VPSID: ptrUint(1), Type: "purchase", Amount: 299.99, Status: "paid"},
			{ID: 2, OrderNo: "ORD20240101002", UserID: 2, CompanyID: ptrUint(1), VPSID: ptrUint(2), Type: "purchase", Amount: 599.99, Status: "paid"},
			{ID: 3, OrderNo: "ORD20240105001", UserID: 6, CompanyID: ptrUint(2), VPSID: ptrUint(4), Type: "purchase", Amount: 399.99, Status: "paid"},
			{ID: 4, OrderNo: "ORD20240105002", UserID: 8, VPSID: ptrUint(6), Type: "purchase", Amount: 149.99, Status: "paid"},
			{ID: 5, OrderNo: "ORD20240110001", UserID: 10, CompanyID: ptrUint(2), VPSID: ptrUint(8), Type: "purchase", Amount: 249.99, Status: "paid"},
			{ID: 6, OrderNo: "ORD20240110002", UserID: 11, CompanyID: ptrUint(2), VPSID: ptrUint(9), Type: "purchase", Amount: 129.99, Status: "paid"},
		}
		for _, o := range orders {
			tx.Create(&o)
		}

		// Create bills
		bills := []model.Bill{
			{ID: 1, CompanyID: ptrUint(1), UserID: 2, Type: "expense", Amount: -299.99, BalanceAfter: 49700.01, Description: "Purchase VPS acme-web-01"},
			{ID: 2, CompanyID: ptrUint(1), UserID: 2, Type: "expense", Amount: -599.99, BalanceAfter: 49100.02, Description: "Purchase VPS acme-db-01"},
			{ID: 3, UserID: 8, Type: "expense", Amount: -149.99, BalanceAfter: 850.01, Description: "Purchase VPS alice-personal-01"},
			{ID: 4, CompanyID: ptrUint(2), UserID: 10, Type: "expense", Amount: -249.99, BalanceAfter: 29750.01, Description: "Purchase VPS globex-finance-01"},
			{ID: 5, CompanyID: ptrUint(2), UserID: 11, Type: "expense", Amount: -129.99, BalanceAfter: 29620.02, Description: "Purchase VPS globex-viewer-01"},
		}
		for _, b := range bills {
			tx.Create(&b)
		}

		// Create tickets
		tickets := []model.Ticket{
			{ID: 1, Title: "Cannot access VPS console", Content: "When I try to access the console for acme-web-01, it shows a connection error.", UserID: 2, CompanyID: ptrUint(1), Status: "open"},
			{ID: 2, Title: "Billing question", Content: "I see an unexpected charge on my bill for last month.", UserID: 8, Status: "replied"},
			{ID: 3, Title: "VNC console not responding", Content: "The VNC console for globex-finance-01 is not responding.", UserID: 10, CompanyID: ptrUint(2), Status: "open"},
			{ID: 4, Title: "Access request", Content: "I need access to the monitoring dashboard for globex-viewer-01.", UserID: 11, CompanyID: ptrUint(2), Status: "open"},
		}
		for _, t := range tickets {
			tx.Create(&t)
		}

		// Create ticket replies
		replies := []model.TicketReply{
			{ID: 1, TicketID: 2, UserID: 1, Content: "Could you please provide your order number for reference?"},
		}
		for _, r := range replies {
			tx.Create(&r)
		}

		// Create API keys
		apiKeys := []model.APIKey{
			{ID: 1, UserID: 2, Name: "Acme Production Key", KeyValue: hashKey("sk_acme_prod_abc123def456"), KeyPrefix: "sk_acme_p", Permissions: `["vps:read","vps:manage"]`, Status: "active"},
			{ID: 2, UserID: 8, Name: "Alice Test Key", KeyValue: hashKey("sk_alice_test_789xyz"), KeyPrefix: "sk_alice_t", Permissions: `["vps:read"]`, Status: "active"},
			{ID: 3, UserID: 6, Name: "Globex Revoked Key", KeyValue: hashKey("sk_globex_revoked_old"), KeyPrefix: "sk_globex_r", Permissions: `["vps:read","vps:manage"]`, Status: "revoked"},
			{ID: 4, UserID: 11, Name: "Globex Viewer Key", KeyValue: hashKey("sk_globex_viewer_abc789"), KeyPrefix: "sk_globex_v", Permissions: `["vps:read"]`, Status: "active"},
		}
		for _, k := range apiKeys {
			tx.Create(&k)
		}

		// Create audit logs
		logs := []model.AuditLog{
			{UserID: 2, CompanyID: ptrUint(1), Action: "vps.create", ResourceType: "vps", ResourceID: 1, Detail: `{"name":"acme-web-01"}`, IPAddress: "192.168.1.100"},
			{UserID: 8, Action: "vps.create", ResourceType: "vps", ResourceID: 6, Detail: `{"name":"alice-personal-01"}`, IPAddress: "192.168.2.50"},
			{UserID: 10, CompanyID: ptrUint(2), Action: "vps.create", ResourceType: "vps", ResourceID: 8, Detail: `{"name":"globex-finance-01"}`, IPAddress: "192.168.3.10"},
			{UserID: 11, CompanyID: ptrUint(2), Action: "vps.create", ResourceType: "vps", ResourceID: 9, Detail: `{"name":"globex-viewer-01"}`, IPAddress: "192.168.3.11"},
		}
		for _, l := range logs {
			tx.Create(&l)
		}

		// Create announcements
		announcements := []model.Announcement{
			{Title: "欢迎使用 CloudNest VPS 管理平台", Content: "这是 OverstepLab 越权测试靶场的模拟平台。请在挑战页面查看所有漏洞挑战，祝您学习愉快！", UserID: 1, IsPinned: true},
		}
		for _, a := range announcements {
			tx.Create(&a)
		}

		// Create system configs (defaults)
		configs := []model.SystemConfig{
			{Key: "site_name", Value: "CloudNest"},
			{Key: "allow_registration", Value: "true"},
			{Key: "default_vps_expire_days", Value: "30"},
			{Key: "max_vps_per_user", Value: "10"},
		}
		for _, c := range configs {
			tx.Create(&c)
		}

		return nil
	})
}

func mustHash(password string) string {
	hash, err := common.HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hash
}

func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func ptrUint(v uint) *uint {
	return &v
}

func ptrString(v string) *string {
	return &v
}

func ptrFloat(v float64) *float64 {
	return &v
}

// Expose ptr functions for use in other packages
var PtrUint = ptrUint
var PtrString = ptrString
var PtrFloat = ptrFloat

func MustHash(password string) string {
	return mustHash(password)
}

func HashKey(key string) string {
	return hashKey(key)
}
