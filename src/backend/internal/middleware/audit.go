package middleware

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
)

// actionRules 用于根据请求方法和路径自动生成审计动作。
// 格式：method:pathRegex -> action
var actionRules = []struct {
	method  string
	pattern *regexp.Regexp
	action  string
	rt      string
}{
	{"POST", regexp.MustCompile(`^/api/v1/auth/login$`), "auth.login", "auth"},
	{"POST", regexp.MustCompile(`^/api/v1/auth/logout$`), "auth.logout", "auth"},
	{"POST", regexp.MustCompile(`^/api/v1/auth/register$`), "auth.register", "auth"},
	{"POST", regexp.MustCompile(`^/api/v1/vps$`), "vps.create", "vps"},
	{"POST", regexp.MustCompile(`^/api/v1/vps/([^/]+)/start$`), "vps.start", "vps"},
	{"POST", regexp.MustCompile(`^/api/v1/vps/([^/]+)/stop$`), "vps.stop", "vps"},
	{"POST", regexp.MustCompile(`^/api/v1/vps/([^/]+)/restart$`), "vps.restart", "vps"},
	{"POST", regexp.MustCompile(`^/api/v1/vps/([^/]+)/reinstall$`), "vps.reinstall", "vps"},
	{"DELETE", regexp.MustCompile(`^/api/v1/vps/([^/]+)$`), "vps.delete", "vps"},
	{"POST", regexp.MustCompile(`^/api/v1/company/members$`), "company.member.add", "company"},
	{"PUT", regexp.MustCompile(`^/api/v1/company/members/([^/]+)$`), "company.member.update", "company"},
	{"DELETE", regexp.MustCompile(`^/api/v1/company/members/([^/]+)$`), "company.member.delete", "company"},
	{"PUT", regexp.MustCompile(`^/api/v1/company/members/([^/]+)/role$`), "company.member.role", "company"},
	{"POST", regexp.MustCompile(`^/api/v1/tickets$`), "ticket.create", "ticket"},
	{"POST", regexp.MustCompile(`^/api/v1/tickets/([^/]+)/reply$`), "ticket.reply", "ticket"},
	{"PUT", regexp.MustCompile(`^/api/v1/tickets/([^/]+)/close$`), "ticket.close", "ticket"},
	{"POST", regexp.MustCompile(`^/api/v1/apikeys$`), "apikey.create", "apikey"},
	{"DELETE", regexp.MustCompile(`^/api/v1/apikeys/([^/]+)$`), "apikey.delete", "apikey"},
	{"POST", regexp.MustCompile(`^/api/v1/bills/recharge$`), "bill.recharge", "bill"},
	{"PUT", regexp.MustCompile(`^/api/v1/user/profile$`), "user.profile.update", "user"},
	{"PUT", regexp.MustCompile(`^/api/v1/user/password$`), "user.password.change", "user"},
	{"PUT", regexp.MustCompile(`^/api/v1/admin/users/([^/]+)/status$`), "admin.user.status", "user"},
	{"PUT", regexp.MustCompile(`^/api/v1/admin/users/([^/]+)/password$`), "admin.user.password", "user"},
	{"POST", regexp.MustCompile(`^/api/v1/admin/announcements$`), "admin.announcement.create", "announcement"},
	{"PUT", regexp.MustCompile(`^/api/v1/admin/announcements/([^/]+)$`), "admin.announcement.update", "announcement"},
	{"DELETE", regexp.MustCompile(`^/api/v1/admin/announcements/([^/]+)$`), "admin.announcement.delete", "announcement"},
	{"PUT", regexp.MustCompile(`^/api/v1/admin/config$`), "admin.config.update", "config"},
	{"POST", regexp.MustCompile(`^/api/v1/admin/reset$`), "admin.db.reset", "system"},
	{"PUT", regexp.MustCompile(`^/api/v1/security-mode$`), "security.mode.change", "system"},
}

func resolveAuditAction(method, path string) (action, resourceType string) {
	for _, rule := range actionRules {
		if strings.EqualFold(rule.method, method) && rule.pattern.MatchString(path) {
			return rule.action, rule.rt
		}
	}
	return "", ""
}

func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		user := GetCurrentUser(c)
		if user == nil {
			return
		}

		action := c.GetString("auditAction")
		resourceType := c.GetString("auditResourceType")
		if action == "" {
			action, resourceType = resolveAuditAction(c.Request.Method, c.Request.URL.Path)
		}
		if action == "" {
			return
		}
		if resourceType == "" {
			resourceType = c.GetString("auditResourceType")
		}

		// 尝试从路径中解析资源 ID
		resourceID := c.GetUint("auditResourceID")
		if resourceID == 0 {
			matches := regexp.MustCompile(`/([^/]+)$`).FindStringSubmatch(c.Request.URL.Path)
			if matches != nil {
				fmt.Sscanf(matches[1], "%d", &resourceID)
			}
		}

		log := model.AuditLog{
			UserID:       user.ID,
			CompanyID:    user.CompanyID,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Detail:       c.GetString("auditDetail"),
			IPAddress:    c.ClientIP(),
		}

		auditRepo := repository.GetAuditLogRepo()
		go func() {
			auditRepo.Create(&log)
		}()
	}
}
