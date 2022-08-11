package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/server/bootstrap/database"
)

func RegisterEndpoints(r gin.IRoutes) {
	r.GET("users/:login", func(c *gin.Context) {
		accountRepo := database.Client().(repos.AccountUserRepository)
		if account, found := accountRepo.AccountUserFind(c.Param("login")); found {
			c.JSON(http.StatusOK, account)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
}
