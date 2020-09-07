package golaerror

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorTestSuite struct {
	suite.Suite
}

func TestErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

func (suite ErrorTestSuite) TestShouldCreateNewError() {
	golaError := New("ERR_TEST_ERROR_CODE", "Some error message", nil)

	suite.Equal("ERR_TEST_ERROR_CODE", golaError.ErrorCode)
	suite.Equal("Some error message", golaError.ErrorMessage)
}

func (suite ErrorTestSuite) TestShouldMarshalErrorWithoutAdditionalData() {
	golaError := New("ERR_TEST_ERROR_CODE", "Some error message", nil)

	bytes, _ := json.Marshal(golaError)
	suite.Equal(`{"errorCode":"ERR_TEST_ERROR_CODE","errorMessage":"Some error message"}`, string(bytes))
}

func (suite ErrorTestSuite) TestShouldMarshalErrorWithAdditionalData() {
	golaError := New("ERR_TEST_ERROR_CODE", "Some error message", "Some additional data")

	bytes, _ := json.Marshal(golaError)
	suite.Equal(`{"errorCode":"ERR_TEST_ERROR_CODE","errorMessage":"Some error message","additionalData":"Some additional data"}`, string(bytes))
}
