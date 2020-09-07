package tracing_middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	span "github.com/gola-glitch/gola-utils/trace"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"io/ioutil"
	"strings"
)

const (
	ERROR    = "error"
	REQUEST  = "request"
	RESPONSE = "response"
)

func HttpTracer(healthz string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logging.GetLogger(ctx)
		logger.Info("Initiating http request tracing")

		httpRequest := ctx.Request
		requestURI := httpRequest.URL.RequestURI()

		if !strings.Contains(requestURI, healthz) {
			requestBodyBytes, err := ioutil.ReadAll(httpRequest.Body)
			if err != nil {
				logger.Error("read error occurred, gracefully exiting the middleware: ", err)
				return
			}
			httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
			responseWriter := createCustomWriter(ctx)

			currentSpan, _ := span.New().Continue(ctx, httpRequest)
			dataSpan := logRequest(ctx, currentSpan, string(requestBodyBytes), REQUEST, requestURI)
			defer func() {
				logger.Info("Initiating http response tracing")

				dataSpan := logResponse(dataSpan, responseWriter.body.String(), RESPONSE, requestURI)
				addResponseTags(responseWriter.Status(), dataSpan)
				dataSpan.End()
			}()
		}
		ctx.Next()
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (writer responseBodyWriter) Write(bytesArray []byte) (int, error) {
	writer.body.Write(bytesArray)
	return writer.ResponseWriter.Write(bytesArray)
}

func createCustomWriter(ctx *gin.Context) *responseBodyWriter {
	responseWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
	ctx.Writer = responseWriter
	return responseWriter
}

func addResponseTags(statusCode int, span *trace.Span) {
	span.AddAttributes(trace.Int64Attribute(ochttp.StatusCodeAttribute, int64(statusCode)))
	if statusCode >= 400 {
		span.AddAttributes(trace.BoolAttribute(ERROR, true))
	}
}

func logRequest(ctx context.Context, span *trace.Span, request string, logKeyName string, logDescription string) *trace.Span {
	spanName := REQUEST + "/" + RESPONSE
	_, dataSpan := trace.StartSpanWithRemoteParent(ctx, logDescription+" | "+spanName, span.SpanContext())

	dataSpan.Annotate([]trace.Attribute{
		trace.StringAttribute(logKeyName, request),
	}, logDescription)
	return dataSpan
}

func logResponse(span *trace.Span, response string, logKeyName string, logDescription string) *trace.Span {
	span.Annotate([]trace.Attribute{
		trace.StringAttribute(logKeyName, response),
	}, logDescription)
	return span
}
