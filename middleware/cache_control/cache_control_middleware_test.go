package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CacheControlMiddlewareTestSuite struct {
	suite.Suite
	recorder   *httptest.ResponseRecorder
	context    *gin.Context
	ginEngine  *gin.Engine
	middleware CacheControl
}

func TestCacheControlMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(CacheControlMiddlewareTestSuite))
}

func (suite *CacheControlMiddlewareTestSuite) SetupTest() {
	suite.middleware = NewCacheControlMiddleware([]string{"/allow-cache"})
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.ginEngine = gin.Default()
}

func (suite *CacheControlMiddlewareTestSuite) TearDownTest() {

}

func (suite *CacheControlMiddlewareTestSuite) TestShouldNotSetCacheHeadersForAllowedUrl() {
	httpRequest, _ := http.NewRequest("GET", "/allow-cache", nil)
	suite.ginEngine.Use(suite.middleware.StopCaching())
	suite.ginEngine.GET("/allow-cache", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	suite.ginEngine.ServeHTTP(suite.recorder, httpRequest)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(suite.recorder.Header(), http.Header{"Content-Type": []string{"application/json; charset=utf-8"}})
}

func (suite *CacheControlMiddlewareTestSuite) TestShouldSetCacheHeadersForAllUrls() {
	httpRequest, _ := http.NewRequest("GET", "/should-not-cache", nil)
	suite.ginEngine.Use(suite.middleware.StopCaching())
	suite.ginEngine.GET("/should-not-cache", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	suite.ginEngine.ServeHTTP(suite.recorder, httpRequest)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(suite.recorder.Header(), http.Header{"Cache-Control": []string{"no-cache, no-store, must-revalidate"}, "Content-Type": []string{"application/json; charset=utf-8"}, "Expires": []string{"0"}, "Pragma": []string{"no-cache"}})
}

func (suite *CacheControlMiddlewareTestSuite) TestShouldNotSetCacheHeadersForAllUrlsNonGETRequestType() {
	httpRequest, _ := http.NewRequest("POST", "/should-not-cache", nil)
	suite.ginEngine.Use(suite.middleware.StopCaching())
	suite.ginEngine.POST("/should-not-cache", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	suite.ginEngine.ServeHTTP(suite.recorder, httpRequest)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(suite.recorder.Header(), http.Header{"Content-Type": []string{"application/json; charset=utf-8"}})
}
