package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type cacheControl struct {
	AllowCache []string
}

type CacheControl interface {
	StopCaching() gin.HandlerFunc
	ShouldNotAllowCache(context *gin.Context) bool
}

func NewCacheControlMiddleware(allowCache []string) CacheControl {
	return cacheControl{AllowCache: allowCache}
}

func (middleware cacheControl) StopCaching() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.Method == http.MethodGet && middleware.ShouldNotAllowCache(context) {
			addNoCacheHeaders(context)
		}
		context.Next()
	}
}

func (middleware cacheControl) ShouldNotAllowCache(context *gin.Context) bool {
	for _, url := range middleware.AllowCache {
		if strings.Contains(context.Request.URL.RequestURI(), url) {
			return false
		}
	}
	return true
}

func addNoCacheHeaders(context *gin.Context) {
	context.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	context.Writer.Header().Set("Pragma", "no-cache")
	context.Writer.Header().Set("Expires", "0")
}
