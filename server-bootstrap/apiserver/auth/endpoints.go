package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/server/bootstrap/services"
)

type authenticatedSession struct {
	Expire  int64              `json:"expires"`
	Refresh int64              `json:"refresh"`
	Host    string             `json:"host"`
	User    *repos.AccountUser `json:"user"`
}

const (
	sessionTTL             = 3600
	maxGenerateSidAttempts = 10
)

func RegisterEndpoints(r gin.IRoutes) {
	sessionStore := make(map[string]*authenticatedSession)

	r.GET("/login", func(c *gin.Context) {
		if aac, err := services.Auth().Login(c.GetHeader(headers.Authorization)); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			// TODO temp
			for i := 0; i < maxGenerateSidAttempts; i++ {
				sid := xid.New().String()
				if _, found := sessionStore[sid]; !found {
					timestamp := time.Now().Unix()
					session := &authenticatedSession{
						Expire:  timestamp + sessionTTL,
						Refresh: timestamp,
						User:    aac.User,
						Host:    c.ClientIP(),
					}

					sessionStore[sid] = session

					c.JSON(http.StatusOK, loginResponse{
						Session: sid,
						Expire:  session.Expire,
						User:    session.User,
					})

					return
				}
			}
			// TODO end temp
		}

		c.AbortWithStatus(http.StatusInternalServerError)
	})

	r.GET("/join", func(c *gin.Context) {
		sid := c.Query("s")
		if session, found := sessionStore[sid]; found {
			if session.Expire <= time.Now().Unix() {
				delete(sessionStore, sid)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Session already expired",
				})
			} else if strings.Compare(session.Host, c.ClientIP()) != 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Attempt to join an existing session from different host",
				})
			} else {
				session.Refresh = time.Now().Unix()
				c.JSON(http.StatusOK, loginResponse{
					Session: sid,
					Expire:  session.Expire,
					User:    session.User,
				})
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Session not found",
			})
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		sid := c.Query("s")
		if _, found := sessionStore[sid]; found {
			delete(sessionStore, sid)
			c.Status(http.StatusNoContent)
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Session not found",
			})
		}
	})

	r.GET("/password/hash", func(c *gin.Context) {
		password := c.Query("p")
		if hash, err := bcrypt.GenerateFromPassword([]byte(password), 8); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.String(http.StatusOK, string(hash))
		}
	})

	r.GET("/dump/session", func(c *gin.Context) {
		c.JSON(http.StatusOK, sessionStore)
	})
}
