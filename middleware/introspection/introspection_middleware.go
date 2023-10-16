package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/http/request"
	"github.com/inclusi-blog/gola-utils/http/util"
	loggingUtil "github.com/inclusi-blog/gola-utils/logging"
	error2 "github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/error"
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/http"
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/service"
	oauthUtils "github.com/inclusi-blog/gola-utils/oauth"
	"go.opencensus.io/plugin/ochttp"
	"net"
	gohttp "net/http"
	"time"
)

type introspectionMiddleware struct {
	tokenService             service.TokenService
	errorResponseInterceptor error2.ResponseInterceptor
	protectedUrlService      service.ProtectedUrlService
	skipTokenDecrypt         bool
}

type IntrospectionMiddleware interface {
	TokenValidationMiddleware() gin.HandlerFunc
	isSkipTokenDecrypt() bool
}

func NewIntrospectionMiddleware(protectedUrlService service.ProtectedUrlService, hydraAdminUrl string) IntrospectionMiddleware {
	return NewIntrospectionAndDecryptionMiddleware(protectedUrlService, hydraAdminUrl, nil)
}

func NewIntrospectionAndDecryptionMiddleware(protectedUrlService service.ProtectedUrlService, hydraAdminUrl string,
	oauthUtils oauthUtils.Utils) IntrospectionMiddleware {

	transport := &gohttp.Transport{
		DialContext: (&net.Dialer{
			Timeout: 50 * time.Second,
		}).DialContext,
	}

	httpClient := &gohttp.Client{Transport: &ochttp.Transport{Base: transport}}
	httpRequestBuilder := request.NewHttpRequestBuilder(httpClient)

	return introspectionMiddleware{
		tokenService: service.NewTokenService(service.NewIntrospectionService(http.NewIntrospectionClient(httpRequestBuilder),
			http.NewIntrospectionRequestBuilder(), hydraAdminUrl),
			oauthUtils),
		errorResponseInterceptor: error2.NewErrorResponseInterceptor(),
		protectedUrlService:      protectedUrlService,
		skipTokenDecrypt:         oauthUtils == nil,
	}
}

func (introspectionMiddleware introspectionMiddleware) TokenValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := loggingUtil.GetLogger(c)
		logger.Debug("IntrospectionMiddleware.TokenValidationMiddleware: Start introspection middleware execution.")
		url := c.Request.URL.Path

		if introspectionMiddleware.protectedUrlService.IsProtected(url) {
			logger.Debug("IntrospectionMiddleware.TokenValidationMiddleware: url is protected: ", url)

			validationError := introspectionMiddleware.tokenService.Validate(c)
			if validationError != nil {
				logger.Error("IntrospectionMiddleware.TokenValidationMiddleware: Validation service error: ", validationError)
				introspectionMiddleware.errorResponseInterceptor.HandleServiceError(c, validationError)
				return
			}

			_, tokenError := util.GetEncryptedIDToken(c)
			if !introspectionMiddleware.skipTokenDecrypt && tokenError == nil {
				decryptionError := introspectionMiddleware.tokenService.Decrypt(c)
				if decryptionError != nil {
					logger.Error("IntrospectionMiddleware.TokenValidationMiddleware: Decryption service error: ", decryptionError)
					introspectionMiddleware.errorResponseInterceptor.HandleServiceError(c, decryptionError)
					return
				}
			} else {
				logger.Debug("IntrospectionMiddleware.TokenValidationMiddleware: Skipping id token decryption")
			}
		} else {
			logger.Debug("IntrospectionMiddleware.TokenValidationMiddleware: url is not protected: ", url)
		}
		logger.Debug("IntrospectionMiddleware.TokenValidationMiddleware: Completed introspection middleware execution for the url ", url)
		c.Next()
	}
}

func (introspectionMiddleware introspectionMiddleware) isSkipTokenDecrypt() bool {
	return introspectionMiddleware.skipTokenDecrypt
}
