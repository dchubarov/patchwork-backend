package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/server/bootstrap/database"
)

func RegisterEndpoints(r gin.IRoutes) {
	r.GET("users/:login", func(c *gin.Context) {
		accountRepo := database.Client().(repos.AccountRepository)
		account := accountRepo.AccountFindUser(c.Param("login"))
		if account.Exists {
			c.JSON(http.StatusOK, account)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
}
