package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/constants"
	"github.com/stretchr/testify/suite"
	"go.opencensus.io/trace"
)

type SessionTracingMiddlewareTestSuite struct {
	suite.Suite
	recorder  *httptest.ResponseRecorder
	context   *gin.Context
	ginEngine *gin.Engine
}

func TestCacheControlMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(SessionTracingMiddlewareTestSuite))
}

func (suite *SessionTracingMiddlewareTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.ginEngine = gin.Default()
}

func (suite *SessionTracingMiddlewareTestSuite) TearDownTest() {

}

type testExporter struct {
	SpanData *trace.SpanData
}

func (t *testExporter) ExportSpan(s *trace.SpanData) {
	t.SpanData = s
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldSetSessionIdInTracingAttributeWhenTracingSessionHeaderKeyIsPresentInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	r.Header.Set(constants.TRACING_SESSION_HEADER_KEY, "12345")
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal("12345", t.SpanData.Attributes[constants.TRACING_SESSION_ID])
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldNotSetSessionIdInTracingAttributeWhenTracingSessionHeaderKeyIsAbsentFromRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal(nil, t.SpanData.Attributes[constants.TRACING_SESSION_ID])
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldSetClientIPInTracingAttributeWhenTracingClientPublicIpHeaderKeyIsPresentInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	r.Header.Set(constants.TRACING_CLIENT_PUBLIC_IP_HEADER, "1.0.0.0")
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal("1.0.0.0", t.SpanData.Attributes[constants.TRACING_CLIENT_PUBLIC_IP])
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldNotSetClientIPInTracingAttributeWhenTracingClientPublicIpHeaderKeyIsAbsentFromRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal(nil, t.SpanData.Attributes[constants.TRACING_CLIENT_PUBLIC_IP])
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldSetAppVersionAttributeWhenAppVersionAttributeExistsInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	r.Header.Set(constants.TRACING_APP_VERSION_HEADER_KEY, "version1")
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal("version1", t.SpanData.Attributes[constants.TRACING_APP_VERSION])
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldNotSetAppVersionAttributeWhenAppVersionHeaderKeyDoesNotExistInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal(nil, t.SpanData.Attributes[constants.TRACING_APP_VERSION])
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldSetDeviceInfoAttributeWhenDeviceInfoHeaderKeyExistsInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	r.Header.Set(constants.TRACING_DEVICE_INFO_HEADER_KEY, "DeviceInfo")
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal("DeviceInfo", t.SpanData.Attributes[constants.TRACING_DEVICE_INFO])
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldNotSetDeviceInfoAttributeWhenDeviceInfoHeaderKeyDoesNotExistInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal(nil, t.SpanData.Attributes[constants.TRACING_DEVICE_INFO])
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldSetTraceIdHttpHeaderWhenTraceIdHttpHeaderKeyExistsInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	r.Header.Set(constants.TRACE_ID_HTTP_HEADER, "someTraceId")

	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal("someTraceId", suite.context.Writer.Header().Get(constants.TRACE_ID_HTTP_HEADER))
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldNotSetTraceIdHttpHeaderWhenTraceIdHttpHeaderKeyDoesNotExistInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal("", suite.context.Writer.Header().Get(constants.TRACE_ID_HTTP_HEADER))
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldSetClientIPSessionIdAppVersionDeviceInfoAndTraceIdHttpHeaderInTracingAttributeWhenAllFiveHeaderKeysArePresentInRequestHeader() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	r.Header.Set(constants.TRACING_SESSION_HEADER_KEY, "12345")
	r.Header.Set(constants.TRACING_CLIENT_PUBLIC_IP_HEADER, "1.0.0.0")
	r.Header.Set(constants.TRACING_APP_VERSION_HEADER_KEY, "version1")
	r.Header.Set(constants.TRACING_DEVICE_INFO_HEADER_KEY, "deviceInfo")
	r.Header.Set(constants.TRACE_ID_HTTP_HEADER, "traceId")
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal("12345", t.SpanData.Attributes[constants.TRACING_SESSION_ID])
	suite.Equal("1.0.0.0", t.SpanData.Attributes[constants.TRACING_CLIENT_PUBLIC_IP])
	suite.Equal("version1", t.SpanData.Attributes[constants.TRACING_APP_VERSION])
	suite.Equal("deviceInfo", t.SpanData.Attributes[constants.TRACING_DEVICE_INFO])
	suite.Equal("traceId", suite.context.Writer.Header().Get(constants.TRACE_ID_HTTP_HEADER))
}

func (suite *SessionTracingMiddlewareTestSuite) TestShouldSetTraceIdHttpHeaderWhenTraceIdHttpHeaderKeyExistsInRequestHeaderInLowerCase() {
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	r.Header.Set(strings.ToLower(constants.TRACE_ID_HTTP_HEADER), "traceId")
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	SessionTracingMiddleware(suite.context)
	s.End()
	suite.Equal("traceId", suite.context.Writer.Header().Get(constants.TRACE_ID_HTTP_HEADER))
}
