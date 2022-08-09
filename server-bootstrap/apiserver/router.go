package apiserver

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"twowls.org/patchwork/commons/logging"
	"twowls.org/patchwork/server/bootstrap/apiserver/auth"
	"twowls.org/patchwork/server/bootstrap/apiserver/health"
	"twowls.org/patchwork/server/bootstrap/config"
)

func Router(log logging.Facade) http.Handler {
	if !log.IsDebugEnabled() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware(log))
	router.Use(cors.New(cors.Config{
		// TODO development only
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Authorization"},
		AllowCredentials: true,
		AllowWildcard:    true,
	}))

	api := router.Group("/api")
	{
		health.RegisterEndpoints(api.Group("/health"))
		auth.RegisterEndpoints(api.Group("/auth"))
	}

	return router
}

func loggingMiddleware(log logging.Facade) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		fields := map[string]any{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"from":   c.ClientIP(),
			"status": c.Writer.Status(),
		}

		log.InfoFields(fields, "%s %-7s %-30s",
			coloredHttpStatus(fields["status"].(int)), fields["method"], fields["path"])
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
