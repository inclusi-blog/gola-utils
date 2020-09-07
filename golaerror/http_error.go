package golaerror

import (
	"strconv"
)

type HttpError struct {
	StatusCode   int
	ResponseBody []byte
}

func (httpError HttpError) Error() string {
	return "StatusCode : " + strconv.Itoa(httpError.StatusCode) +
		", ResponseBody : " + string(httpError.ResponseBody)
}
