package middlewares

import (
	"net/http"
	"rentx/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	id, email, phone, role, err := utils.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.Set("userId", id)
	c.Set("email", email)
	c.Set("phone", phone)
	c.Set("role", role)

	c.Next()
}
