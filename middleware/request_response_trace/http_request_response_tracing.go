package request_response_trace

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/logging"
	span "github.com/inclusi-blog/gola-utils/trace"
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

type traceHookFunc func(*gin.Context, []byte) (string, error)

func HttpRequestResponseTracingAllMiddleware(ctx *gin.Context) {
	HttpRequestResponseTracingMiddleware(nil, "/healthz", nil, nil)(ctx)
}

func HttpRequestResponseTracingAllMiddlewareWithHooks(requestHook traceHookFunc, responseHook traceHookFunc) gin.HandlerFunc {
	return HttpRequestResponseTracingMiddleware(nil, "/healthz", requestHook, responseHook)
}

func HttpRequestResponseTracingAllMiddlewareWithCustomHealthEndpoint(healthEndpoint string) gin.HandlerFunc {
	return HttpRequestResponseTracingMiddleware(nil, healthEndpoint, nil, nil)
}

func HttpRequestResponseTracingMiddleware(apisToBeIgnored []IgnoreRequestResponseLogs, healthEndpoint string, requestHook traceHookFunc, responseHook traceHookFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logging.GetLogger(ctx)
		logger.Info("Initiating http request response tracing")
		httpRequest := ctx.Request
		requestURI := httpRequest.URL.RequestURI()
		if !strings.Contains(requestURI, healthEndpoint) {
			dataSpan := createNewSpan(ctx, httpRequest)
			logHttpRequest(ctx, dataSpan, httpRequest, apisToBeIgnored, requestHook)

			responseWriter := addAndGetCustomResponseWriter(ctx)
			defer func() {
				logHttpResponse(ctx, dataSpan, responseWriter, apisToBeIgnored, responseHook)
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

func logHttpRequest(ctx *gin.Context, dataSpan *trace.Span, httpRequest *http.Request, apisToBeIgnored []IgnoreRequestResponseLogs, requestHook traceHookFunc) {
	logger := logging.GetLogger(ctx)
	if !isRequestLogAllowed(httpRequest.URL.RequestURI(), apisToBeIgnored) {
		logger.Debug("request payload for api is not allowed to log")
		addAnnotation(dataSpan, REQUEST_SEQURITY_MSG, REQUEST, httpRequest.URL.RequestURI())
		return
	}
	var requestBodyBytes = []byte(span.NO_BODY_CONTENT)
	var readErr error
	if httpRequest.Body != nil {
		requestBodyBytes, readErr = ioutil.ReadAll(httpRequest.Body)
		if readErr != nil {
			logger.Errorf("Error occurred while reading request body. Error: %+v", readErr)
			return
		}
		httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
	}
	if requestHook != nil {
		newRequestBody, err := requestHook(ctx, requestBodyBytes)
		if err != nil {
			logger.Debug("logHttpRequest: Request Hook function failed")
			addAnnotation(dataSpan, span.EscapeSpecialChar(requestBodyBytes), REQUEST, httpRequest.URL.RequestURI())
		} else {
			logger.Debug("logHttpRequest: Added request payload in newly created span")
			addAnnotation(dataSpan, span.EscapeSpecialChar([]byte(newRequestBody)), REQUEST, httpRequest.URL.RequestURI())
		}
	} else {
		addAnnotation(dataSpan, span.EscapeSpecialChar(requestBodyBytes), REQUEST, httpRequest.URL.RequestURI())
		logger.Debug("Added request payload in newly created span")
	}
}

func addAndGetCustomResponseWriter(ctx *gin.Context) *responseBodyWriter {
	responseWriter := NewCustomResponseWriter(ctx.Writer)
	ctx.Writer = responseWriter
	return responseWriter
}

func logHttpResponse(ctx *gin.Context, dataSpan *trace.Span, responseWriter *responseBodyWriter, apisToBeIgnored []IgnoreRequestResponseLogs, responseHook traceHookFunc) {
	logger := logging.GetLogger(ctx)
	if !isResponseLogAllowed(ctx.Request.URL.RequestURI(), apisToBeIgnored) {
		logger.Debug("response for api is not allowed to log")
		addAnnotation(dataSpan, RESPONSE_SEQURITY_MSG, RESPONSE, ctx.Request.URL.RequestURI())
		return
	}

	var requestBodyBytes = []byte(span.NO_BODY_CONTENT)
	if responseWriter.body != nil {
		requestBodyBytes = responseWriter.body.Bytes()
	}

	if responseHook != nil {
		responseBodyFromHook, err := responseHook(ctx, requestBodyBytes)
		if err != nil {
			logger.Debug("logHttpResponse: Response Hook function failed")
			addAnnotation(dataSpan, span.EscapeSpecialChar(requestBodyBytes), RESPONSE, ctx.Request.URL.RequestURI())
		} else {
			logger.Debug("logHttpResponse: Added response body in newly created span. Closing the span.")
			addAnnotation(dataSpan, span.EscapeSpecialChar([]byte(responseBodyFromHook)), RESPONSE, ctx.Request.URL.RequestURI())
		}
	} else {
		addAnnotation(dataSpan, span.EscapeSpecialChar(requestBodyBytes), RESPONSE, ctx.Request.URL.RequestURI())
		logger.Debug("Added response body in newly created span. Closing the span.")
	}
	addResponseTags(responseWriter.Status(), dataSpan)
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
