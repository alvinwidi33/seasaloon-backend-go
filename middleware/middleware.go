package middleware

import (
	"net/http"
	"seasaloon-backend-go/repository"
	"strings"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		userID, userRole, err := repository.ParseToken(parts[1]) 
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if len(requiredRoles) > 0 {
			allowed := false
			for _, role := range requiredRoles {
				if userRole == role {
					allowed = true
					break
				}
			}

			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}
		}

		c.Set("userID", userID)
		c.Set("userRole", userRole)
		c.Next()
	}
}