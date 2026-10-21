package handlers

import (
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "http.MethodOptions" || c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		next(c)
	}
}
