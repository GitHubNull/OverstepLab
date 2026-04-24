package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/common"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				common.InternalError(c, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
