package service

//go:generate mockgen -source=introspection_service.go -destination=./../mocks/mock_introspection_service.go -package=mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/http"
)

type IntrospectionService interface {
	Introspect(context *gin.Context, accessToken string) (response []byte, err error)
}

type introspectionService struct {
	introspectionHttpClient     http.IntrospectionHttpClient
	introspectionRequestBuilder http.IntrospectionRequestBuilder
	hydraAdminUrl               string
}

func NewIntrospectionService(client http.IntrospectionHttpClient, builder http.IntrospectionRequestBuilder, hydraAdminUrl string) IntrospectionService {
	return introspectionService{introspectionHttpClient: client, introspectionRequestBuilder: builder, hydraAdminUrl: hydraAdminUrl}
}

func (introspectionService introspectionService) Introspect(context *gin.Context, accessToken string) ([]byte, error) {
	logger := logging.GetLogger(context)
	logger.Info("IntrospectionService.Introspect : Making introspection request ")
	request, requestError := introspectionService.introspectionRequestBuilder.Build(introspectionService.hydraAdminUrl, accessToken)

	if requestError != nil {
		logger.Error("IntrospectionService.Introspect: Error in introspection Request builder ", requestError)
		return nil, requestError
	}

	return introspectionService.introspectionHttpClient.Introspect(context, request)
}
