package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, message string, data any) {
	c.JSON(200, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"message": message,
	})
}
