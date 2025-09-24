package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no role"})
			return
		}
		role := roleVal.(string)

		allowed := slices.Contains(roles, role)
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not allowed"})
			return
		}
		c.Next()
	}
}
