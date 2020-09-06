package trace

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type TraceUtilsTestSuite struct {
	suite.Suite
	context  *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestTraceUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(TraceUtilsTestSuite))
}

func (suite *TraceUtilsTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *TraceUtilsTestSuite) TestShouldReturnTraceIdFromGinContext() {
	expectedTraceId := "sample-trace-id"
	suite.context.Set("traceID", expectedTraceId)

	gotTraceId := GetTraceId(suite.context)
	suite.Equal(expectedTraceId, gotTraceId)
}

func (suite *TraceUtilsTestSuite) TestShouldReturnNoTraceIdFromGinContextIfNoTraceIdExists() {
	noTraceId := "no-trace-id"

	gotTraceId := GetTraceId(suite.context)
	suite.Equal(noTraceId, gotTraceId)
}

func (suite *TraceUtilsTestSuite) TestShouldReturnNoTraceIdFromGinContextIfTraceIdIsNotString() {
	noTraceId := "no-trace-id"
	nonStringTraceID := 12345
	suite.context.Set("traceID", nonStringTraceID)

	gotTraceId := GetTraceId(suite.context)
	suite.Equal(noTraceId, gotTraceId)
}

func (suite *TraceUtilsTestSuite) TestShouldReturnTraceIdFromGoContext() {
	expectedTraceId := "sample-trace-id"
	goCtx := context.WithValue(context.TODO(), "traceID", expectedTraceId)

	gotTraceId := GetTraceId(goCtx)
	suite.Equal(expectedTraceId, gotTraceId)
}

func (suite *TraceUtilsTestSuite) TestShouldReturnNoTraceIdFromGoContextIfNoTraceIdExists() {
	noTraceId := "no-trace-id"

	gotTraceId := GetTraceId(context.TODO())
	suite.Equal(noTraceId, gotTraceId)
}

func (suite *TraceUtilsTestSuite) TestShouldReturnNoTraceIdFromGoContextIfTraceIdIsNotString() {
	noTraceId := "no-trace-id"
	nonStringTraceID := 12345
	goCtx := context.WithValue(context.TODO(), "traceID", nonStringTraceID)

	gotTraceId := GetTraceId(goCtx)
	suite.Equal(noTraceId, gotTraceId)
}

func (suite *TraceUtilsTestSuite) TestShouldReturnTrimmedTraceIDFromGoContextIfTraceIDHasLeadingZeros() {
	goCtx := context.WithValue(context.TODO(), "traceID", "00000012345")

	gotTraceId := GetTraceId(goCtx)
	suite.Equal("12345", gotTraceId)
}
