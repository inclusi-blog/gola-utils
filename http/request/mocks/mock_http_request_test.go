package mocks

import (
	"github.com/inclusi-blog/gola-utils/http/request"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MockHttpRequestTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
}

func TestMockHttpRequestTestSuite(t *testing.T) {
	suite.Run(t, new(MockHttpRequestTestSuite))
}

func (suite *MockHttpRequestTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
}

func (suite MockHttpRequestTestSuite) TestMockGenerated() {
	var httpRequest request.HttpRequest
	httpRequest = NewMockHttpRequest(suite.mockCtrl)
	println(httpRequest)
}
