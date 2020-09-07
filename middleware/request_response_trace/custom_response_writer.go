package request_response_trace

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (writer responseBodyWriter) Write(bytesArray []byte) (int, error) {
	writer.body.Write(bytesArray)
	return writer.ResponseWriter.Write(bytesArray)
}

func NewCustomResponseWriter(responseWriter gin.ResponseWriter) *responseBodyWriter {
	return &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: responseWriter}
}
