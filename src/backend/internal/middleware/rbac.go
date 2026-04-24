package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

const (
	PermAdmin         = "admin"
	PermVPSView       = "vps:view"
	PermVPSManage     = "vps:manage"
	PermVPSConsole    = "vps:console"
	PermUserManage    = "user:manage"
	PermFinanceView   = "finance:view"
	PermFinanceExport = "finance:export"
	PermAPIKeyManage  = "apikey:manage"
	PermAuditLogView  = "audit:view"
)

var rolePermissions = map[string]map[string]bool{
	"platform_admin": {
		PermAdmin: true, PermVPSView: true, PermVPSManage: true, PermVPSConsole: true,
		PermUserManage: true, PermFinanceView: true, PermFinanceExport: true,
		PermAPIKeyManage: true, PermAuditLogView: true,
	},
	"admin": {
		PermVPSView: true, PermVPSManage: true, PermVPSConsole: true,
		PermUserManage: true, PermFinanceView: true, PermFinanceExport: true,
		PermAPIKeyManage: true, PermAuditLogView: true,
	},
	"operator": {
		PermVPSView: true, PermVPSManage: true, PermVPSConsole: true,
		PermFinanceView: true,
	},
	"finance": {
		PermFinanceView: true, PermFinanceExport: true,
	},
	"viewer": {
		PermVPSView: true,
	},
}

func RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			common.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// Platform admin has all permissions
		if user.IsPlatformAdmin() {
			c.Next()
			return
		}

		if vuln.IsSecureMode() {
			// Secure mode: enforce permissions
			perms := rolePermissions[user.Role]
			if !perms[perm] {
				common.Forbidden(c, "Insufficient permissions")
				c.Abort()
				return
			}
		}
		// Vulnerable mode: allow through (vulnerability is in service layer for granularity)

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			common.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		if vuln.IsSecureMode() && !user.IsPlatformAdmin() {
			common.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}
		// Vulnerable mode: allow non-admin to access admin endpoints (V-04)

		c.Next()
	}
}

func RequireCompanyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			common.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		if vuln.IsSecureMode() && !user.IsCompanyAdmin() && !user.IsPlatformAdmin() {
			common.Forbidden(c, "Company admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
