package request_response_trace

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HttpRequestResponseMiddlewareTest struct {
	suite.Suite
	recorder  *httptest.ResponseRecorder
	ginEngine *gin.Engine
	context   *gin.Context
}

func TestHttpRequestResponseMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(HttpRequestResponseMiddlewareTest))
}

func (suite *HttpRequestResponseMiddlewareTest) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.ginEngine = gin.Default()
}

type spanExporter struct {
	spans []*trace.SpanData
}

func (s *spanExporter) ExportSpan(d *trace.SpanData) {
	s.spans = append(s.spans, d)
}

type transactionsRequest struct {
	BatchSize  int `json:"batch_size"`
	PageNumber int `json:"page_number"`
}

type transactionsResponse struct {
	TransactionId   int    `json:"transaction_id"`
	TransactionType string `json:"transaction_type"`
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddRequestResponseAnnotationsInNewChildSpan() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
		{TransactionId: 454, TransactionType: "Send money"},
	}
	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddleware,
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddRequestResponseAnnotationsInNewChildSpanWithRequestResponseHooks() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
		{TransactionId: 454, TransactionType: "Send money"},
	}

	reqHook := func(*gin.Context, []byte) (string, error) {
		marshall, _ := json.Marshal(request)
		return string(marshall), nil
	}

	resHook := func(*gin.Context, []byte) (string, error) {
		marshall, _ := json.Marshal(response)
		return string(marshall), nil
	}

	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddlewareWithHooks(reqHook, resHook),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddRequestResponseAnnotationsInNewChildSpanWithOnlyResponseHooks() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
		{TransactionId: 454, TransactionType: "Send money"},
	}

	resHook := func(*gin.Context, []byte) (string, error) {
		marshall, _ := json.Marshal(response)
		return string(marshall), nil
	}

	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddlewareWithHooks(nil, resHook),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddRequestResponseAnnotationsInNewChildSpanWithResponseHookFailed() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
		{TransactionId: 454, TransactionType: "Send money"},
	}

	resHook := func(*gin.Context, []byte) (string, error) {
		return "", errors.New("resHook failed")
	}

	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddlewareWithHooks(nil, resHook),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddRequestResponseAnnotationsInNewChildSpanWithOnlyRequestHooks() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
		{TransactionId: 454, TransactionType: "Send money"},
	}

	reqHook := func(*gin.Context, []byte) (string, error) {
		marshall, _ := json.Marshal(request)
		return string(marshall), nil
	}

	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddlewareWithHooks(reqHook, nil),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddRequestResponseAnnotationsInNewChildSpanWithRequestHookFailed() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
		{TransactionId: 454, TransactionType: "Send money"},
	}

	reqHook := func(*gin.Context, []byte) (string, error) {
		return "", errors.New("req hook failed")
	}

	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddlewareWithHooks(reqHook, nil),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddNOBODYCONTENTAsRequestAnnotationIfRequestBodyIsEmpty() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
		{TransactionId: 454, TransactionType: "Send money"},
	}
	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddleware,
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, nil)
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal("NO BODY CONTENT", childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddNOBODYCONTENTAsRequestAnnotationIfRequestBodyIsEmptyWithHooks() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
		{TransactionId: 454, TransactionType: "Send money"},
	}

	reqHook := func(*gin.Context, []byte) (string, error) {
		return "NO BODY CONTENT", nil
	}

	resHook := func(*gin.Context, []byte) (string, error) {
		marshall, _ := json.Marshal(response)
		return string(marshall), nil
	}

	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddlewareWithHooks(reqHook, resHook),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, nil)
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal("NO BODY CONTENT", childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddBase64EncodedTextAsResponseAnnotationIfResponseBodyContainsInvalidChars() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []byte{1, 128, 129, 200}
	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddleware,
		func(c *gin.Context) {
			c.Data(http.StatusOK, "application/pdf", response)
			c.Next()
		})
	requestBody, _ := json.Marshal(request)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal("BASE64_ENCODED_CONTENT: AYCByA==", childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldAddNOBODYCONTENTAsResponseAnnotationsIfResponseBodyIsEmpty() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddleware,
		func(c *gin.Context) {
			c.Status(http.StatusOK)
			c.Next()
		})
	requestBody, _ := json.Marshal(request)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 2)
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusOK), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Len(childSpan.Annotations, 2)

	suite.Equal(url, childSpan.Annotations[0].Message)
	suite.Equal(url, childSpan.Annotations[1].Message)

	suite.Len(childSpan.Annotations[0].Attributes, 1)
	suite.Len(childSpan.Annotations[1].Attributes, 1)

	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal("NO BODY CONTENT", childSpan.Annotations[1].Attributes["response"])

	suite.Nil(childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldSetErrorAttributeInChildSpanIfResponseIsNotSuccess() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	suite.ginEngine.POST(url,
		HttpRequestResponseTracingAllMiddleware,
		func(c *gin.Context) {
			c.Status(http.StatusForbidden)
		})
	requestBody, _ := json.Marshal(request)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	childSpan := t.spans[0]
	suite.Equal(int64(http.StatusForbidden), childSpan.Attributes[ochttp.StatusCodeAttribute])
	suite.Equal(true, childSpan.Attributes["error"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldExcludeTracingForHealthEndpoint() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/v1/healthz"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	suite.ginEngine.GET(url,
		HttpRequestResponseTracingAllMiddleware,
		func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
	requestBody, _ := json.Marshal(request)
	r, _ := http.NewRequest("GET", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 1)
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldExcludeTracingForCustomHealthEndpoint() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	healthEndpoint := "/api/v1/info"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	suite.ginEngine.GET(healthEndpoint,
		HttpRequestResponseTracingAllMiddlewareWithCustomHealthEndpoint(healthEndpoint),
		func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
	requestBody, _ := json.Marshal(request)
	r, _ := http.NewRequest("GET", healthEndpoint, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	suite.Len(t.spans, 1)
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldNotLogRequestAndResponseByDefaultIfApiIsIncludedInIgnoredApisList() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
	}
	expectedRequestLog := "Request payload not logged for security reasons"
	expectedResponseLog := "Response body not logged for security reasons"
	suite.ginEngine.POST(url,
		HttpRequestResponseTracingMiddleware([]IgnoreRequestResponseLogs{
			{PartialApiPath: "TRANSACTIONS"},
			{PartialApiPath: "filters"},
		}, "/healthz", nil, nil),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	childSpan := t.spans[0]
	suite.Equal(expectedRequestLog, childSpan.Annotations[0].Attributes["request"])
	suite.Equal(expectedResponseLog, childSpan.Annotations[1].Attributes["response"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldLogRequestOnlyIfApiIsIncludedInIgnoredApisListAndAllowedRequestIsSetToTrue() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
	}
	expectedResponseLog := "Response body not logged for security reasons"
	suite.ginEngine.POST(url,
		HttpRequestResponseTracingMiddleware([]IgnoreRequestResponseLogs{
			{PartialApiPath: "transactionS", IsRequestLogAllowed: true},
			{PartialApiPath: "filters"},
		}, "/healthz", nil, nil),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	childSpan := t.spans[0]
	suite.Equal(string(requestBody), childSpan.Annotations[0].Attributes["request"])
	suite.Equal(expectedResponseLog, childSpan.Annotations[1].Attributes["response"])
}

func (suite *HttpRequestResponseMiddlewareTest) TestShouldLogResponseOnlyIfApiIsIncludedInIgnoredApisListAndAllowedResponseIsSetToTrue() {
	t := spanExporter{}
	trace.RegisterExporter(&t)
	url := "/api/transactions"
	request := transactionsRequest{BatchSize: 20, PageNumber: 4}
	response := []transactionsResponse{
		{TransactionId: 111, TransactionType: "Transfer funds"},
	}
	expectedRequestLog := "Request payload not logged for security reasons"
	suite.ginEngine.POST(url,
		HttpRequestResponseTracingMiddleware([]IgnoreRequestResponseLogs{
			{PartialApiPath: "Transactions", IsResponseLogAllowed: true},
			{PartialApiPath: "filters"},
		}, "/healthz", nil, nil),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, response)
		})
	requestBody, _ := json.Marshal(request)
	responseBody, _ := json.Marshal(response)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	ctx, s := trace.StartSpan(r.Context(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	suite.ginEngine.ServeHTTP(suite.recorder, r.WithContext(ctx))

	s.End()
	childSpan := t.spans[0]
	suite.Equal(expectedRequestLog, childSpan.Annotations[0].Attributes["request"])
	suite.Equal(string(responseBody), childSpan.Annotations[1].Attributes["response"])
}
