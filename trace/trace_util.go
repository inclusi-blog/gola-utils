package trace

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/constants"
)

func GetTraceId(ctx context.Context) string {
	if ginContext, ok := ctx.(*gin.Context); ok {
		if traceId, ok := ginContext.Get(constants.TRACE_KEY); ok {
			if traceIdString, ok := traceId.(string); ok {
				return strings.TrimLeft(traceIdString, "0")
			}
		}
	}
	traceId := ctx.Value(constants.TRACE_KEY)
	if traceId != nil {
		if traceIdString, ok := traceId.(string); ok {
			return strings.TrimLeft(traceIdString, "0")
		}
	}
	return constants.NO_TRACE_ID
}
