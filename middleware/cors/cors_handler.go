package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/model"
	"net/http"
	"strings"
)

func CORSMiddleware(config model.CorsConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logging.GetLogger(ctx)
		logger.Info("cors middleware - enter")
		origins := ctx.Request.Header["Origin"]
		host := ctx.Request.URL.Host
		logger.Infof("cors middleware - host: %s", host)
		logger.Infof("cors middleware - origin: %s", origins)
		if len(origins) > 0 && !isRequestedFromSameHost(origins[0], host) {
			if !isAllowedOrigin(origins[0], config.AllowedOrigins) {
				logger.Errorf("cors middleware - origin %s not allowed", origins[0])
				ctx.AbortWithStatus(http.StatusNotAcceptable)
				return
			}
		}
		logger.Info("cors middleware - exit")
		ctx.Next()
	}
}

func isRequestedFromSameHost(origin string, requestedHost string) bool {
	return strings.Contains(origin, requestedHost)
}

func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == origin {
			return true
		}
	}
	return false
}
