package service

//go:generate mockgen -source=token_service.go -destination=./../mocks/mock_token_service.go -package=mocks

import (
	"github.com/gin-gonic/gin"
	errorUtils "github.com/gola-glitch/gola-utils/golang_error"
	"github.com/gola-glitch/gola-utils/http/util"
	loggingUtil "github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/constants"
	error2 "github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/error"
	"github.com/gola-glitch/gola-utils/oauth"
)

type TokenService interface {
	Validate(ctx *gin.Context) *error2.OAuthMiddlewareError
	Decrypt(ctx *gin.Context) *error2.OAuthMiddlewareError
}

type tokenService struct {
	introspectionService IntrospectionService
	oauthUtils           oauth.Utils
}

func NewTokenService(introspectionService IntrospectionService, oauthUtils oauth.Utils) TokenService {
	return tokenService{
		introspectionService: introspectionService,
		oauthUtils:           oauthUtils,
	}
}

func (tokenService tokenService) Validate(ctx *gin.Context) *error2.OAuthMiddlewareError {
	logger := loggingUtil.GetLogger(ctx)

	logger.Info("TokenService.Validate: Starting token validation")

	accessToken, err := util.GetAccessToken(ctx)
	if err != nil {
		for key, value := range ctx.Request.Header {
			logger.Info("Header Key: ", key)
			logger.Info("Header Key: ", value)
		}
		logger.Error("TokenService.Validate: Could not fetch access token, reason: ", err)
		return error2.InvalidAccessTokenError
	}

	if accessToken == "" {
		logger.Error("TokenService.Validate: Access token is empty")
		return error2.InvalidAccessTokenError
	}

	_, introspectionError := tokenService.introspectionService.Introspect(ctx, accessToken)
	if introspectionError != nil {
		logger.Error("TokenService.Validate: Token introspection failed with error ", introspectionError)
		switch introspectionError.(type) {
		case error2.HydraInternalServerError:
			logger.Error("TokenService.Validate: server error while introspecting token")
			return error2.InternalServerErrorFunc(introspectionError.Error())
		default:
			logger.Error("TokenService.Validate: Access token is invalid")
			return error2.InvalidAccessTokenError
		}
	}

	logger.Info("TokenService.Validate: Finished successfully")

	return nil
}

func (tokenService tokenService) Decrypt(ctx *gin.Context) *error2.OAuthMiddlewareError {
	logger := loggingUtil.GetLogger(ctx)
	logger.Info("TokenService.Decrypt: Starting id token decryption")

	idToken, decodeTokenError := tokenService.oauthUtils.DecodeEncryptedIdToken(ctx)
	if decodeTokenError != nil {
		logger.Error("TokenService.Decrypt: Id Token decryption error ", decodeTokenError.Error())
		switch decodeTokenError.(type) {
		case errorUtils.InternalServerError:
			logger.Error("TokenService.Decrypt: server error while decrypting Id Token")
			return error2.InternalServerErrorFunc("error in token decryption")
		default:
			logger.Error("TokenService.Decrypt: invalid Id Token")
			return error2.InvalidIdTokenError
		}
	}
	ctx.Set(constants.ContextDecryptedIdTokenKey, idToken)
	logger.Debug("TokenService.Decrypt: Id Token decrypted successfully and set to context")
	return nil
}
