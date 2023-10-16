package trace

//go:generate mockgen -source=trace.go -destination=./mocks/mock_trace.go -package=mocks

import (
	"context"
	"github.com/inclusi-blog/gola-utils/constants"
	"net/http"
	"net/http/httptrace"

	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp"
	openTrace "go.opencensus.io/trace"
)

type Trace interface {
	Continue(ctx context.Context, httpRequest *http.Request) (*openTrace.Span, *http.Request)
	DeriveTracingFromRemoteParent(spanName string, traceID openTrace.TraceID, spanID openTrace.SpanID) (context.Context, *openTrace.Span)
	StartTracing(ctx context.Context, cmdName string) (context.Context, *openTrace.Span)
	CreateNewSpan(ctx context.Context, spanName string) (*openTrace.Span, context.Context)
}

type trace struct {
}

func (t trace) Continue(ctx context.Context, httpRequest *http.Request) (*openTrace.Span, *http.Request) {
	var span *openTrace.Span

	if ginContext, ok := ctx.(*gin.Context); ok {
		ctx = ginContext.Request.Context()
	}

	span = openTrace.FromContext(ctx)
	clientTrace := ochttp.NewSpanAnnotatingClientTrace(httpRequest, span)
	newContext := httptrace.WithClientTrace(ctx, clientTrace)
	return span, httpRequest.WithContext(newContext)
}

func (t trace) DeriveTracingFromRemoteParent(spanName string, traceID openTrace.TraceID, spanID openTrace.SpanID) (context.Context, *openTrace.Span) {
	spanContext := openTrace.SpanContext{
		TraceID:      traceID,
		SpanID:       spanID,
		TraceOptions: 0,
		Tracestate:   nil,
	}
	newContextWithTrace, span := openTrace.StartSpanWithRemoteParent(context.Background(), spanName, spanContext)
	newContextWithTrace = context.WithValue(newContextWithTrace, constants.TRACE_KEY, traceID.String())

	return newContextWithTrace, span
}

func (t trace) StartTracing(ctx context.Context, cmdName string) (context.Context, *openTrace.Span) {
	return openTrace.StartSpan(ctx, cmdName)
}

func (t trace) CreateNewSpan(ctx context.Context, spanName string) (*openTrace.Span, context.Context) {
	context, dataSpan := openTrace.StartSpanWithRemoteParent(ctx, spanName,
		openTrace.FromContext(ctx).SpanContext())
	return dataSpan, context
}

func New() Trace {
	return trace{}
}

