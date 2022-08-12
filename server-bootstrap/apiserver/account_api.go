package apiserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/server/bootstrap/database"
)

func registerEndpointsAccount(r gin.IRoutes) {
	r.GET("users/:login", func(c *gin.Context) {
		accountRepo := database.Client().(repos.AccountRepository)
		if account, found := accountRepo.AccountFindUser(c.Param("login"), false); found {
			c.JSON(http.StatusOK, account)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
}
