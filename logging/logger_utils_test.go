package logging

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type LoggerUtilsTestSuite struct {
	suite.Suite
	context  *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestLoggerUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerUtilsTestSuite))
}

func (suite *LoggerUtilsTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *LoggerUtilsTestSuite) TestShouldGetDefaultLoggerFromContext() {
	logger := GetLogger(suite.context)
	traceID, ok := logger.Data["traceID"]
	suite.True(ok)
	suite.Equal("no-trace-id", traceID)
}

func (suite *LoggerUtilsTestSuite) TestShouldGetLoggerFromGinContextWithCorrelationID() {
	suite.context.Set("logger", NewLoggerEntry().WithContext(suite.context).WithField("traceID", "sampleID"))

	logger := GetLogger(suite.context)

	traceID, ok := logger.Data["traceID"]
	suite.True(ok)
	suite.Equal("sampleID", traceID)
}

func (suite *LoggerUtilsTestSuite) TestShouldGetLoggerFromGoContextWithNoTraceId() {
	logger := GetLogger(context.TODO())
	traceID, ok := logger.Data["traceID"]
	suite.True(ok)
	suite.Equal("no-trace-id", traceID)
}

func (suite *LoggerUtilsTestSuite) TestShouldGetLoggerFromGoContextWithCorrelationID() {
	ctx := context.WithValue(context.Background(), "logger", NewLoggerEntry().WithField("traceID", "sampleID"))
	logger := GetLogger(ctx)
	traceID, ok := logger.Data["traceID"]
	suite.True(ok)
	suite.Equal("sampleID", traceID)
}

func (suite *LoggerUtilsTestSuite) TestShouldGetDefaultLoggerWhenContextIsNil() {
	// ctx := context.WithValue(context.Background(), "logger", NewLoggerEntry().WithField("traceID", "sampleID"))
	logger := GetLogger(nil)
	traceID, ok := logger.Data["traceID"]
	suite.True(ok)
	suite.Equal("no-trace-id", traceID)
}
