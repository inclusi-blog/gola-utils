package error

//go:generate mockgen -source=error_response_interceptor.go -destination=./../mocks/mock_error_response_interceptor.go -package=mocks

import (
	"github.com/gin-gonic/gin"
	loggingUtil "github.com/gola-glitch/gola-utils/logging"
)

type ResponseInterceptor interface {
	HandleServiceError(ctx *gin.Context, serviceError *OAuthMiddlewareError)
}

type errorResponseInterceptor struct {
}

func NewErrorResponseInterceptor() ResponseInterceptor {
	return errorResponseInterceptor{}
}

func (errorResponseInterceptor) HandleServiceError(ctx *gin.Context, serviceError *OAuthMiddlewareError) {
	loggingUtil.GetLogger(ctx).Error("ResponseInterceptor: Service Error: ", serviceError.ErrorResponse)
	if serviceError.HttpStatusCode >= 500 && serviceError.HttpStatusCode < 600 {
		loggingUtil.GetLogger(ctx).Error("ResponseInterceptor: Returning httpStatus 500", serviceError.ErrorResponse)
		ctx.AbortWithStatus(serviceError.HttpStatusCode)
		return
	}
	ctx.AbortWithStatusJSON(serviceError.HttpStatusCode,
		gin.H{"error": serviceError.ErrorResponse.ErrorMessage, "errorCode": serviceError.ErrorResponse.ErrorCode,
			"errorMessage": serviceError.ErrorResponse.ErrorMessage})
}
