package trace

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/constants"
	"strings"
)

const (
	validValueMinChar   = 0
	validValueMaxChar   = 126

	NO_BODY_CONTENT             = "NO BODY CONTENT"
	BASE64_ENCODED_CONTENT      = "BASE64_ENCODED_CONTENT: "
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

func EscapeSpecialChar(value []byte) string {
	if value == nil || len(value) == 0 {
		return NO_BODY_CONTENT
	}
	containsSpecialChar := false
	for _, c := range value {
		if (c < validValueMinChar) || (c > validValueMaxChar) {
			containsSpecialChar = true
		}
	}
	if containsSpecialChar {
		return BASE64_ENCODED_CONTENT + base64.StdEncoding.EncodeToString(value)
	}
	return string(value)
}
