package apiserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twowls.org/patchwork/server/bootstrap/services"
)

func registerEndpointsAccount(r gin.IRoutes) {
	r.GET("user/:login", userAccountEndpoint)
	r.GET("users/:login", userAccountEndpoint)
}

// endpoints

func userAccountEndpoint(c *gin.Context) {
	if user, err := services.Account().FindUser(c, c.Param("login"), false); err != nil {
		handleStandardError(err, c)
	} else {
		c.JSON(http.StatusOK, user)
	}
}
