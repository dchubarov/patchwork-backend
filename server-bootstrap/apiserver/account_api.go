package apiserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twowls.org/patchwork/server/bootstrap/services"
)

func registerEndpointsAccount(r gin.IRoutes) {
	r.GET("users/:login", func(c *gin.Context) {
		if user, err := services.Account().FindUser(c.Param("login"), false, retrieveAuth(c)); err != nil {
			handleStandardError(err, c)
		} else {
			c.JSON(http.StatusOK, user)
		}
	})
}
