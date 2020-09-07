package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/constants"
	"go.opencensus.io/trace"
)

func SessionTracingMiddleware(c *gin.Context) {
	span := trace.FromContext(c.Request.Context())

	sessionHeaderKey := c.Request.Header.Get(constants.TRACING_SESSION_HEADER_KEY)
	if sessionHeaderKey != "" {
		span.AddAttributes(trace.StringAttribute(constants.TRACING_SESSION_ID, sessionHeaderKey))
	}
	clientIp := c.Request.Header.Get(constants.TRACING_CLIENT_PUBLIC_IP_HEADER)
	if clientIp != "" {
		span.AddAttributes(trace.StringAttribute(constants.TRACING_CLIENT_PUBLIC_IP, clientIp))
	}
	appVersion := c.Request.Header.Get(constants.TRACING_APP_VERSION_HEADER_KEY)
	if appVersion != "" {
		span.AddAttributes(trace.StringAttribute(constants.TRACING_APP_VERSION, appVersion))
	}

	deviceInfo := c.Request.Header.Get(constants.TRACING_DEVICE_INFO_HEADER_KEY)
	if deviceInfo != "" {
		span.AddAttributes(trace.StringAttribute(constants.TRACING_DEVICE_INFO, deviceInfo))
	}
	traceId := c.Request.Header.Get(constants.TRACE_ID_HTTP_HEADER)
	if traceId != "" {
		c.Writer.Header().Add(constants.TRACE_ID_HTTP_HEADER, traceId)
	}
	c.Next()
}
