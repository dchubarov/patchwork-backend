package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterEndpoints(r gin.IRoutes) {
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
}
