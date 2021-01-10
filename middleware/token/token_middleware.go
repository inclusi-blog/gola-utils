package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/constants"
	"github.com/gola-glitch/gola-utils/http/util"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/model"
	"github.com/gola-glitch/gola-utils/oauth"
	"net/http"
)

type TokenMiddleware interface {
	DecryptIdToken() gin.HandlerFunc
}

type tokenMiddleware struct {
	oauthUtils oauth.Utils
}

func NewTokenMiddleware(oauthUtils oauth.Utils) TokenMiddleware {
	return tokenMiddleware{oauthUtils: oauthUtils}
}

func respondWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, model.ErrorResponse{
		ErrorCode:    message,
		ErrorMessage: message,
	})
}

func (t tokenMiddleware) DecryptIdToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logging.GetLogger(c)
		logger.Info("Started decryption")
		_, accessTokenErr := util.GetAccessToken(c)
		_, idTokenErr := util.GetEncryptedIDToken(c)

		if accessTokenErr != nil || idTokenErr != nil {
			logger.Error("Access or id token not available in tokenMiddleware::DecryptIdToken")
			respondWithError(c, http.StatusUnauthorized, "Invalid access token")
			return
		}

		decodedIdToken, decodeTokenError := t.oauthUtils.DecodeEncryptedIdToken(c)

		if decodeTokenError != nil {
			logger.Error("Unable to decode id token : ", decodeTokenError)
			respondWithError(c, http.StatusInternalServerError, "Unable to decode id token")
			return
		}
		logger.Info("Decryption successful")
		c.Set(constants.CONTEXT_ENC_ID_TOKEN, decodedIdToken)
		c.Next()
	}
}
