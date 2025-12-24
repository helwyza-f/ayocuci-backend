package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OutletContext middleware
// - membaca outlet aktif dari header X-Outlet-ID
// - validasi user punya akses ke outlet tsb
// - set outlet_id ke gin context
func OutletContext(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		outletIDStr := c.GetHeader("X-Outlet-ID")
		if outletIDStr == "" {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "X-Outlet-ID header required",
			})
			return
		}

		outletID, err := strconv.Atoi(outletIDStr)
		if err != nil || outletID <= 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "invalid outlet id",
			})
			return
		}

		userID := c.GetUint("user_id")
		if userID == 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "unauthorized",
			})
			return
		}

		// validasi: user punya akses ke outlet ini?
		var count int64
		err = db.Table("user_outlets").
			Where("user_id = ? AND outlet_id = ? AND active = true", userID, outletID).
			Count(&count).Error

		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error": "database error",
			})
			return
		}

		if count == 0 {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "no access to this outlet",
			})
			return
		}

		// set outlet context
		c.Set("outlet_id", uint(outletID))

		c.Next()
	}
}
