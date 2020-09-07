package golaerror

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type httpError struct {
	suite.Suite
}

func TestHttpErrorTestSuite(t *testing.T) {
	suite.Run(t, new(httpError))
}

func (suite httpError) TestShouldCreateNewError() {
	httpError := HttpError{
		StatusCode:   500,
		ResponseBody: []byte("body"),
	}

	suite.Equal("StatusCode : 500, ResponseBody : body", httpError.Error())
}
