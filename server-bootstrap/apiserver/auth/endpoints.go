package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
	"time"
)

type authenticatedSession struct {
	Expire  int64      `json:"expires"`
	Refresh int64      `json:"refresh"`
	Host    string     `json:"host"`
	User    *userEntry `json:"user"`
}

const (
	sessionTTL             = 3600
	maxGenerateSidAttempts = 10
)

func RegisterEndpoints(r gin.IRoutes) {
	sessionStore := make(map[string]*authenticatedSession)
	userStore := []userEntry{
		{
			"dime",
			"dime@twowls.org",
			"Dmitry Chubarov",
			"$2a$08$0uFMX8KVPrxTSnpH1LL.pesu5/JvYnEbHbQhFWbF5xK/squZGPL7e", // bcrypt hash for 'xxx'
			[]userMembership{
				{"devs", "", "admin"},
				{"top12", "Top 1-2", "contributor"},
			}},
	}

	r.GET("/login", func(c *gin.Context) {
		a := c.GetHeader(headers.Authorization)
		if strings.HasPrefix(a, "Basic") {
			buf, err := base64.StdEncoding.DecodeString(a[6:])
			if err == nil {
				credentials := strings.Split(string(buf), ":")
				if len(credentials) > 1 {
					idx := slices.IndexFunc(userStore, func(u userEntry) bool {
						return strings.Compare(credentials[0], u.Login) == 0 ||
							strings.Compare(credentials[0], u.Email) == 0
					})

					if idx < 0 || bcrypt.CompareHashAndPassword([]byte(userStore[idx].PasswordHash), []byte(credentials[1])) != nil {
						c.JSON(http.StatusUnauthorized, gin.H{
							"error": "Invalid user name and (or) password",
						})
					} else {
						for i := 0; i < maxGenerateSidAttempts; i++ {
							sid := xid.New().String()
							if _, found := sessionStore[sid]; !found {
								timestamp := time.Now().Unix()
								session := &authenticatedSession{
									Expire:  timestamp + sessionTTL,
									Refresh: timestamp,
									User:    &userStore[idx],
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

						// could not generate unique session id
						c.AbortWithStatus(http.StatusInternalServerError)
					}

					return
				}
			}
		}

		c.AbortWithStatus(http.StatusBadRequest)
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
