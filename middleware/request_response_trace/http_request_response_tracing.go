package request_response_trace

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	span "github.com/gola-glitch/gola-utils/trace"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ERROR                 = "error"
	REQUEST               = "request"
	RESPONSE              = "response"
	REQUEST_SEQURITY_MSG  = "Request payload not logged for security reasons"
	RESPONSE_SEQURITY_MSG = "Response body not logged for security reasons"
)

func HttpRequestResponseTracingAllMiddleware(ctx *gin.Context) {
	HttpRequestResponseTracingMiddleware(nil)(ctx)
}

func HttpRequestResponseTracingMiddleware(apisToBeIgnored []IgnoreRequestResponseLogs) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logging.GetLogger(ctx)
		logger.Info("Initiating http request response tracing")
		httpRequest := ctx.Request
		requestURI := httpRequest.URL.RequestURI()
		if !strings.Contains(requestURI, "/healthz") {
			dataSpan := createNewSpan(ctx, httpRequest)
			logHttpRequest(ctx, dataSpan, httpRequest, apisToBeIgnored)

			responseWriter := addAndGetCustomResponseWriter(ctx)
			defer func() {
				logHttpResponse(ctx, dataSpan, responseWriter, apisToBeIgnored)
				dataSpan.End()
			}()
		}
		ctx.Next()
	}
}

func createNewSpan(ctx context.Context, httpRequest *http.Request) *trace.Span {
	logger := logging.GetLogger(ctx)
	currentSpan, _ := span.New().Continue(ctx, httpRequest)
	logDescription := httpRequest.URL.RequestURI()
	spanName := logDescription + " | " + REQUEST + "/" + RESPONSE
	logger.Debug("Creating new child span with span name ", spanName)
	_, dataSpan := trace.StartSpanWithRemoteParent(ctx, spanName, currentSpan.SpanContext())
	logger.Debug("Created new child span with span name ", spanName)
	return dataSpan
}

func logHttpRequest(ctx context.Context, dataSpan *trace.Span, httpRequest *http.Request, apisToBeIgnored []IgnoreRequestResponseLogs) {
	if !isRequestLogAllowed(httpRequest.URL.RequestURI(), apisToBeIgnored) {
		logging.GetLogger(ctx).Debug("request payload for api is not allowed to log")
		addAnnotation(dataSpan, REQUEST_SEQURITY_MSG, REQUEST, httpRequest.URL.RequestURI())
		return
	}
	var requestBody = "NO BODY CONTENT"
	if httpRequest.Body != nil {
		requestBodyBytes, readErr := ioutil.ReadAll(httpRequest.Body)
		if readErr != nil {
			logging.GetLogger(ctx).Errorf("Error occurred while reading request body. Error: %+v", readErr)
			return
		}
		httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
		requestBody = string(requestBodyBytes)
	}
	addAnnotation(dataSpan, requestBody, REQUEST, httpRequest.URL.RequestURI())
	logging.GetLogger(ctx).Debug("Added request payload in newly created span")
}

func addAndGetCustomResponseWriter(ctx *gin.Context) *responseBodyWriter {
	responseWriter := NewCustomResponseWriter(ctx.Writer)
	ctx.Writer = responseWriter
	return responseWriter
}

func logHttpResponse(ctx *gin.Context, dataSpan *trace.Span, responseWriter *responseBodyWriter, apisToBeIgnored []IgnoreRequestResponseLogs) {
	if !isResponseLogAllowed(ctx.Request.URL.RequestURI(), apisToBeIgnored) {
		logging.GetLogger(ctx).Debug("response for api is not allowed to log")
		addAnnotation(dataSpan, RESPONSE_SEQURITY_MSG, RESPONSE, ctx.Request.URL.RequestURI())
		return
	}
	var responseBody = "NO BODY CONTENT"
	if responseWriter.body != nil && responseWriter.body.Len() > 0 {
		responseBody = responseWriter.body.String()
	}
	addAnnotation(dataSpan, responseBody, RESPONSE, ctx.Request.URL.RequestURI())
	addResponseTags(responseWriter.Status(), dataSpan)
	logging.GetLogger(ctx).Debug("Added response body in newly created span. Closing the span.")
}

func isRequestLogAllowed(requestURI string, apisToBeIgnored []IgnoreRequestResponseLogs) bool {
	return isLogAllowed(requestURI, apisToBeIgnored, func(apiToBeIgnored IgnoreRequestResponseLogs) bool {
		return apiToBeIgnored.IsRequestLogAllowed
	})
}

func isResponseLogAllowed(requestURI string, apisToBeIgnored []IgnoreRequestResponseLogs) bool {
	return isLogAllowed(requestURI, apisToBeIgnored, func(apiToBeIgnored IgnoreRequestResponseLogs) bool {
		return apiToBeIgnored.IsResponseLogAllowed
	})
}

func isLogAllowed(requestURI string, apisToBeIgnored []IgnoreRequestResponseLogs, getIsAllowed func(IgnoreRequestResponseLogs) bool) bool {
	for _, api := range apisToBeIgnored {
		if strings.Contains(strings.ToLower(requestURI), strings.ToLower(api.PartialApiPath)) {
			return getIsAllowed(api)
		}
	}
	return true
}

func addAnnotation(dataSpan *trace.Span, logValue, logKey, logDescription string) {
	dataSpan.Annotate([]trace.Attribute{
		trace.StringAttribute(logKey, logValue),
	}, logDescription)
}

func addResponseTags(statusCode int, span *trace.Span) {
	span.AddAttributes(trace.Int64Attribute(ochttp.StatusCodeAttribute, int64(statusCode)))
	if statusCode >= 400 {
		span.AddAttributes(trace.BoolAttribute(ERROR, true))
	}
}
