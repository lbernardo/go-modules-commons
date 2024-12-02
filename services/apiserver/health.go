package apiserver

import "github.com/gin-gonic/gin"

func Health(r *gin.Engine) {
	r.GET("/api/commons/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
