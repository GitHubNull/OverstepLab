package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
)

func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		user := GetCurrentUser(c)
		if user == nil {
			return
		}

		action := c.GetString("auditAction")
		if action == "" {
			return
		}

		log := model.AuditLog{
			UserID:       user.ID,
			CompanyID:    user.CompanyID,
			Action:       action,
			ResourceType: c.GetString("auditResourceType"),
			ResourceID:   c.GetUint("auditResourceID"),
			Detail:       c.GetString("auditDetail"),
			IPAddress:    c.ClientIP(),
		}

		auditRepo := repository.GetAuditLogRepo()
		go func() {
			auditRepo.Create(&log)
		}()
	}
}
