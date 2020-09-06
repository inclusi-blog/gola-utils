package trace

import (
	"context"
	"net/http"
	"net/http/httptrace"

	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp"
	openTrace "go.opencensus.io/trace"
)

type Trace interface {
	Continue(ctx context.Context, httpRequest *http.Request) (*openTrace.Span, *http.Request)
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

func New() Trace {
	return trace{}
}
