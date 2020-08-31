package logging

import (
	"github.com/gin-gonic/gin"
	"gola-utils/constants"
)

func LoggingMiddleware(entry *golaLoggerEntry) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get(constants.TRACE_ID_HTTP_HEADER)
		if traceID == "" {
			traceID = constants.NO_TRACE_ID
		}
		//standardLogger := logrus.StandardLogger()
		jsonFormat := c.Request.Header.Get(constants.JSON_FORMATTER_HTTP_HEADER)
		if jsonFormat != "" {
			entry.SetFormatter(constants.JSON)
		}
		//for _, h := range hooks {
		// standardLogger.AddHook(h)
		//}
		logger := entry.WithField(constants.TRACE_KEY, traceID).WithContext(c.Request.Context())
		c.Set(constants.LOGGER_KEY, logger)
		c.Set(constants.TRACE_KEY, traceID)
		c.Next()
	}
}
