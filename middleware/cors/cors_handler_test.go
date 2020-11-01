package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	gohttp "net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type CorsMiddlewareTest struct {
	suite.Suite
	mockCtrl  *gomock.Controller
	ginEngine *gin.Engine
	recorder  *httptest.ResponseRecorder
}

func TestCorsMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(CorsMiddlewareTest))
}

func (suite *CorsMiddlewareTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
}

func (suite *CorsMiddlewareTest) TestCorsMiddleware() {
	config := model.CorsConfig{
		AllowedOrigins: []string{"http://app.gola.xyz", "http://localhost:3000"},
	}
	suite.ginEngine = suite.setupGinServer(config)

	httpRequest, _ := gohttp.NewRequest("GET", "api/idp/v1/token/validate", nil)
	httpRequest.Header.Add("Origin", "http://app.gola.xyz")
	httpRequest.URL, _ = url.Parse("http://app.gola.xyz/api/idp/v1/token/validate")

	suite.ginEngine.ServeHTTP(suite.recorder, httpRequest)

	suite.NotEqual(gohttp.StatusNoContent, suite.recorder.Code)

}

func (suite *CorsMiddlewareTest) TestCorsMiddleware_ShouldContinueTheFlowWhenOriginHostAndRequestedUrlHostAreSame() {
	config := model.CorsConfig{
		AllowedOrigins: []string{"http://app.gola.xyz"},
	}
	suite.ginEngine = suite.setupGinServer(config)

	httpRequest, _ := gohttp.NewRequest("GET", "api/idp/v1/token/validate", nil)
	httpRequest.Header.Add("Origin", "http://app.gola.xyz")
	httpRequest.URL, _ = url.Parse("http://app.gola.xyz/api/idp/v1/token/validate")

	suite.ginEngine.ServeHTTP(suite.recorder, httpRequest)

	suite.NotEqual(gohttp.StatusNotAcceptable, suite.recorder.Code)

}

func (suite *CorsMiddlewareTest) TestCorsMiddleware_WhenOriginNotAllowed() {
	config := model.CorsConfig{
		AllowedOrigins: []string{"http://app.gola.xyz", "http://localhost:3000"},
	}
	suite.ginEngine = suite.setupGinServer(config)

	httpRequest, _ := gohttp.NewRequest("GET", "api/idp/v1/token/validate", nil)
	httpRequest.Header.Add("Origin", "http://some-other.domain.com")
	httpRequest.URL, _ = url.Parse("http://app.gola.xyz/api/idp/v1/token/validate")

	suite.ginEngine.ServeHTTP(suite.recorder, httpRequest)

	suite.Equal(gohttp.StatusNotAcceptable, suite.recorder.Code)
}

func (suite *CorsMiddlewareTest) setupGinServer(config model.CorsConfig) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware(config))
	r.POST("api/idp/v1/token/validate",
		func(c *gin.Context) {
			c.Status(gohttp.StatusOK)
		})

	r.POST("api/idp/v1/user/profile",
		func(c *gin.Context) {
			c.Status(gohttp.StatusOK)
		})

	return r
}
