package middleware

import "github.com/gin-gonic/gin"

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		// super admin can access everything
		if role == "super_admin" {
			c.Next()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{
			"message": "forbidden",
		})
	}
}
