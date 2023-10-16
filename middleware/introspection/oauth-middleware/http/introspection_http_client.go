package http

//go:generate mockgen -source=introspection_http_client.go -destination=./../mocks/mock_introspection_http_client.go -package=mocks
import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/http/request"
	"github.com/inclusi-blog/gola-utils/http/util"
	"github.com/inclusi-blog/gola-utils/logging"
	middlewareError "github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/error"
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/model"
	"net/http"
)

type IntrospectionHttpClient interface {
	Introspect(context *gin.Context, request *http.Request) (response []byte, error error)
}

type introspectionHttpClient struct {
	httpRequestBuilder request.HttpRequestBuilder
}

func NewIntrospectionClient(httpRequestBuilder request.HttpRequestBuilder) IntrospectionHttpClient {
	return introspectionHttpClient{httpRequestBuilder: httpRequestBuilder}
}

func (introspectionHttpClient introspectionHttpClient) Introspect(context *gin.Context, request *http.Request) ([]byte, error) {
	logger := logging.GetLogger(context)

	logger.Info("IntrospectionHttpClient.Introspect: Making introspect call to hydra. Url ", request.URL.String())

	introspectionResponse := model.IntrospectionResponse{}

	introspectionRequest := make(map[string]interface{})
	accessToken, _ := util.GetAccessToken(context)
	introspectionRequest["token"] = accessToken

	httpError := introspectionHttpClient.httpRequestBuilder.NewRequestWithContext(context).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		AddHeader("Accept", "application/json").
		WithFormURLEncoded(introspectionRequest).
		ResponseAs(&introspectionResponse).
		Post(request.URL.String())

	responseBytes, marshalError := json.Marshal(introspectionResponse)
	if marshalError != nil {
		return nil, middlewareError.HydraInternalServerError{}
	}

	if httpError != nil {
		if apiError, ok := httpError.(golaerror.HttpError); ok {
			if apiError.StatusCode >= 500 && apiError.StatusCode <= 505 {
				logger.Error("IntrospectionHttpClient.Introspect: Status code from hydra", apiError.StatusCode)
				return responseBytes, middlewareError.HydraInternalServerError{}
			}
		}
		logger.Error("IntrospectionHttpClient.Introspect: Error while making call to hydra : ", httpError)
		return responseBytes, middlewareError.AuthenticationError{}
	}

	if introspectionResponse.Active {
		logger.Info("IntrospectionHttpClient.Introspect: Valid & Active access token!")
		return responseBytes, nil
	}

	logger.Error("IntrospectionHttpClient.Introspect: Access Token is not valid. Response ", introspectionResponse)

	return responseBytes, middlewareError.AuthenticationError{}
}
