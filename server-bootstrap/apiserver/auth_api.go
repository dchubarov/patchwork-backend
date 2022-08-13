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
	Token   string               `json:"token,omitempty"`
	Expires int64                `json:"expires"`
	User    *service.UserAccount `json:"user"`
}

// registry

func registerEndpointsAuth(r gin.IRoutes) {
	r.GET("/login", loginEndpoint)
	r.GET("/logout", logoutEndpoint)

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

// endpoints

func loginEndpoint(c *gin.Context) {
	var aac *service.AuthContext
	if aac = services.GetAuthFromContext(c); aac == nil {
		loginAac, err := services.Auth().LoginWithCredentials(c, c.GetHeader(headers.Authorization), service.AuthServiceHeaderCredentials)

		if err != nil {
			handleStandardError(err, c)
		} else {
			aac = loginAac
		}
	}

	if aac != nil {
		httpCode := http.StatusOK
		if aac.Token != "" {
			httpCode = http.StatusCreated
		}

		c.JSON(httpCode, loginResponse{
			Token:   aac.Token,
			Expires: aac.Session.Expires.Unix(),
			User:    aac.User,
		})
	}
}

func logoutEndpoint(c *gin.Context) {
	if err := services.Auth().Logout(c); err != nil {
		handleStandardError(err, c)
	} else {
		c.Status(http.StatusNoContent)
	}
}
