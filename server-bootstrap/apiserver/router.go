package apiserver

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"twowls.org/patchwork/commons/logging"
	"twowls.org/patchwork/commons/service"
	"twowls.org/patchwork/server/bootstrap/config"
	"twowls.org/patchwork/server/bootstrap/services"
)

func Router(log logging.Facade) http.Handler {
	if !log.IsDebugEnabled() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware(log))
	router.Use(tokenInterceptorMiddleware())
	router.Use(cors.New(cors.Config{
		// TODO development only
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Authorization"},
		AllowCredentials: true,
		AllowWildcard:    true,
	}))

	api := router.Group("/api")
	{
		registerEndpointsAccount(api.Group("/accounts"))
		registerEndpointsAuth(api.Group("/auth"))
		registerEndpointsHealth(api.Group("/health"))
	}

	return router
}

// private

var authContextKey = "$$aac." + strconv.Itoa(os.Getpid())

func handleStandardError(err error, c *gin.Context) {
	errStd, ok := err.(*service.E)
	if !ok {
		log.Warn("Service returned error which will not be forwarded to user: %v", err)
		errStd = service.ErrServiceUnspecific
	}

	httpCode := errStd.HttpCode
	if httpCode < 200 || httpCode >= 599 {
		httpCode = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(httpCode, errStd)
}

func retrieveAuth(c *gin.Context) *service.AuthContext {
	if raw, found := c.Get(authContextKey); found {
		if auth, ok := raw.(*service.AuthContext); ok {
			return auth
		}
	}
	return nil
}

func tokenInterceptorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(headers.Authorization)
		if strings.HasPrefix(authHeader, services.AuthSchemeBearer) {
			auth, err := services.Auth().LoginWithCredentials(authHeader, service.AuthServiceHeaderCredentials)
			if err != nil {
				handleStandardError(err, c)
			} else {
				c.Set(authContextKey, auth)
			}
		}
	}
}

func loggingMiddleware(log logging.Facade) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		fields := map[string]any{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"from":     c.ClientIP(),
			"status":   c.Writer.Status(),
			"duration": duration,
		}

		log.InfoFields(fields, "%s %-7s %-30s (%s)",
			coloredHttpStatus(fields["status"].(int)), fields["method"], fields["path"],
			duration.Round(time.Microsecond).String())
	}
}

func coloredHttpStatus(status int) string {
	result := strconv.Itoa(status)
	if config.Values().Logging.NoColor {
		return result
	}

	if status >= 200 && status < 300 {
		result = fmt.Sprintf("\033[1;30m\033[42m%3s\033[0m", result)
	} else if status >= 300 && status < 400 {
		result = fmt.Sprintf("\033[1;30m\033[47m%3s\033[0m", result)
	} else if status >= 400 && status < 500 {
		result = fmt.Sprintf("\033[1;30m\033[43m%3s\033[0m", result)
	} else if status >= 500 {
		result = fmt.Sprintf("\033[1;30m\033[41m%3s\033[0m", result)
	} else {
		result = fmt.Sprintf("\033[1;31m%3s\033[0m", result)
	}

	return result
}
