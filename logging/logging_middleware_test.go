package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LoggingMiddlewareTestSuite struct {
	suite.Suite
	context  *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestLoggingMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingMiddlewareTestSuite))
}

func (suite *LoggingMiddlewareTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *LoggingMiddlewareTestSuite) TestContextShouldHaveLoggerWithCorrelationIdField() {
	suite.context.Request, _ = http.NewRequest("GET", "/customer-profile?user-id=12345", nil)

	loggingMiddleware := LoggingMiddleware(NewLoggerEntry())

	loggingMiddleware(suite.context)

	logger, ok := suite.context.Keys["logger"]
	suite.True(ok)
	_, ok = logger.(*golaLoggerEntry).Data["traceID"]
	suite.True(ok)
	// Default traceID
	u1 := "no-trace-id"

	// Default traceID if not present in header
	loggingMiddleware(suite.context)

	logger, ok = suite.context.Keys["logger"]
	suite.True(ok)
	_, ok = logger.(*golaLoggerEntry).Data["traceID"]
	suite.True(ok)
	u2 := "no-trace-id"

	suite.Equal(u1, u2)
}

func (suite *LoggingMiddlewareTestSuite) TestContextShouldHaveLoggerWithSameCorrelationIdIfPassedInRequestHeader() {
	suite.context.Request, _ = http.NewRequest("GET", "/customer-profile?user-id=12345", nil)
	suite.context.Request.Header.Set("X-B3-Traceid", "sampleID")

	loggingMiddleware := LoggingMiddleware(NewLoggerEntry())

	loggingMiddleware(suite.context)

	logger, ok := suite.context.Keys["logger"]
	suite.True(ok)
	traceID, ok := logger.(*golaLoggerEntry).Data["traceID"]
	suite.True(ok)
	suite.Equal("sampleID", traceID)
}

func (suite *LoggingMiddlewareTestSuite) TestContextShouldHaveSameCorrelationIdIfPassedInRequestHeader() {
	suite.context.Request, _ = http.NewRequest("GET", "/customer-profile?user-id=12345", nil)
	suite.context.Request.Header.Set("X-B3-Traceid", "sampleID")

	loggingMiddleware := LoggingMiddleware(NewLoggerEntry())

	loggingMiddleware(suite.context)

	traceID, ok := suite.context.Keys["traceID"]
	suite.True(ok)
	suite.Equal("sampleID", traceID)
}

func (suite *LoggingMiddlewareTestSuite) TestContextShouldHaveNewCorrelationIdIfNotPassedInRequestHeader() {
	suite.context.Request, _ = http.NewRequest("GET", "/customer-profile?user-id=12345", nil)

	loggingMiddleware := LoggingMiddleware(NewLoggerEntry())

	loggingMiddleware(suite.context)

	traceID, ok := suite.context.Keys["traceID"]
	suite.True(ok)
	// traceID is Default one
	id := "no-trace-id"
	suite.Equal(traceID, id)
}
