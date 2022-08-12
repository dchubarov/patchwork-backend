package apiserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/server/bootstrap/database"
)

func registerEndpointsAccount(r gin.IRoutes) {
	r.GET("users/:login", func(c *gin.Context) {
		login := c.Param("login")
		aac := retrieveAuth(c)
		if aac != nil {
			if aac != nil && (aac.User.IsPrivileged() || aac.User.Is(login)) {
				accountRepo := database.Client().(repos.AccountRepository)
				if account, found := accountRepo.AccountFindUser(login, false); found {
					c.JSON(http.StatusOK, account)
				} else {
					c.Status(http.StatusNotFound)
				}
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	})
}
