package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/commons/service"
	"twowls.org/patchwork/server/bootstrap/services"
)

type loginResponse struct {
	Expire int64              `json:"expires"`
	User   *repos.AccountUser `json:"user"`
	Token  string             `json:"token"`
}

type authenticatedSession struct {
	Expire  int64              `json:"expires"`
	Refresh int64              `json:"refresh"`
	Host    string             `json:"host"`
	User    *repos.AccountUser `json:"user"`
}

func registerEndpointsAuth(r gin.IRoutes) {
	sessionStore := make(map[string]*authenticatedSession)

	r.GET("/login", func(c *gin.Context) {
		var aac *service.AuthContext
		if aac = retrieveAuth(c); aac == nil {
			loginAac, err := services.Auth().LoginWithCredentials(
				c.GetHeader(headers.Authorization),
				service.AuthServiceHeaderCredentials)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			}

			aac = loginAac
		}

		// TODO temp
		sessionStore[aac.Session.Sid] = &authenticatedSession{
			Expire:  aac.Session.Expires.Unix(),
			Refresh: 0,
			User:    aac.User,
		}

		c.JSON(http.StatusOK, loginResponse{
			Expire: aac.Session.Expires.Unix(),
			User:   aac.User,
			Token:  aac.Token,
		})
		// TODO end temp
	})

	r.GET("/logout", func(c *gin.Context) {
		if aac := retrieveAuth(c); aac != nil {
			// TODO delete session
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
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

	// TODO to be removed
	r.GET("/dump/session", func(c *gin.Context) {
		c.JSON(http.StatusOK, sessionStore)
	})
}
