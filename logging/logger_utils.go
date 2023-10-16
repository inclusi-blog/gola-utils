package logging

import (
	"context"
	"github.com/inclusi-blog/gola-utils/constants"

	"github.com/gin-gonic/gin"
)

func GetLogger(c context.Context) *golaLoggerEntry {
	if c == nil {
		return NewLoggerEntry().WithContext(c).WithField(constants.TRACE_KEY, constants.NO_TRACE_ID)
	}
	if ginContext, ok := c.(*gin.Context); ok {
		logger, ok := ginContext.Get(constants.LOGGER_KEY)
		if ok {
			return logger.(*golaLoggerEntry)
		}
	}
	contextLogger := c.Value(constants.LOGGER_KEY)
	if contextLogger != nil {
		if logger, ok := contextLogger.(*golaLoggerEntry); ok {
			if ok {
				return logger
			}
		}
	}
	return NewLoggerEntry().WithContext(c).WithField(constants.TRACE_KEY, constants.NO_TRACE_ID)
}
