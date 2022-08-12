package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"twowls.org/patchwork/commons/service"
	"twowls.org/patchwork/server/bootstrap/services"
)

type loginResponse struct {
	Token  string               `json:"token"`
	Expire int64                `json:"expires"`
	User   *service.AccountUser `json:"user"`
}

func registerEndpointsAuth(r gin.IRoutes) {
	r.GET("/login", func(c *gin.Context) {
		var aac *service.AuthContext
		if aac = retrieveAuth(c); aac == nil {
			loginAac, err := services.Auth().LoginWithCredentials(
				c.GetHeader(headers.Authorization),
				service.AuthServiceHeaderCredentials)

			if err != nil {
				handleStandardError(err, c)
			} else {
				aac = loginAac
			}
		}

		if aac != nil {
			c.JSON(http.StatusOK, loginResponse{
				Expire: aac.Session.Expires.Unix(),
				User:   aac.User,
				Token:  aac.Token,
			})
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		if err := services.Auth().Logout(retrieveAuth(c)); err != nil {
			handleStandardError(err, c)
		} else {
			c.Status(http.StatusNoContent)
		}
	})

	// TODO to be removed
	r.GET("/password/hash", func(c *gin.Context) {
		password := c.Query("p")
		if hash, err := bcrypt.GenerateFromPassword([]byte(password), 8); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.String(http.StatusOK, string(hash))
		}
	})
}
