package middleware

import (
	"github.com/gin-gonic/gin"
)

func RequireActiveOutlet() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		// super admin & owner bebas
		if role == "super_admin" || role == "owner" {
			c.Next()
			return
		}

		outletID, exists := c.Get("active_outlet_id")
		if !exists || outletID == nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "please select outlet first",
			})
			return
		}

		c.Next()
	}
}
