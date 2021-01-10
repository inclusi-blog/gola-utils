package service

import (
	"github.com/gin-gonic/gin"
	error2 "github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/error"
	"github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type IntrospectionServiceTest struct {
	suite.Suite
	mockCtrl                        *gomock.Controller
	context                         *gin.Context
	mockIntrospectionHttpClient     *mocks.MockIntrospectionHttpClient
	mockIntrospectionRequestBuilder *mocks.MockIntrospectionRequestBuilder
	introspectionService            IntrospectionService
	hydraAdminUrl                   string
}

func TestIntrospectionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(IntrospectionServiceTest))
}

func (suite *IntrospectionServiceTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.mockIntrospectionHttpClient = mocks.NewMockIntrospectionHttpClient(suite.mockCtrl)
	suite.mockIntrospectionRequestBuilder = mocks.NewMockIntrospectionRequestBuilder(suite.mockCtrl)
	suite.hydraAdminUrl = "admin-url"
	suite.introspectionService = NewIntrospectionService(suite.mockIntrospectionHttpClient, suite.mockIntrospectionRequestBuilder, suite.hydraAdminUrl)
}

func (suite *IntrospectionServiceTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite IntrospectionServiceTest) TestIntrospectionServiceShouldSendErrorIfRequestBuilderSendsError() {
	accessToken := "access_token"
	suite.context.Request, _ = http.NewRequest("", "", nil)
	suite.mockIntrospectionRequestBuilder.EXPECT().Build(suite.hydraAdminUrl, accessToken).Return(suite.context.Request, error2.AuthenticationError{})
	response, err := suite.introspectionService.Introspect(suite.context, accessToken)
	suite.Nil(response)
	suite.Equal(error2.AuthenticationError{}, err)
}

func (suite IntrospectionServiceTest) TestIntrospectionServiceShouldReturnResponseOfIntrospectHttpClientOnNoErrorWhileBuildingRequest() {
	accessToken := "access_token"
	suite.context.Request, _ = http.NewRequest("", "", nil)
	suite.mockIntrospectionRequestBuilder.EXPECT().Build(suite.hydraAdminUrl, accessToken).Return(suite.context.Request, nil)
	expectedResponse := []byte("abc")
	suite.mockIntrospectionHttpClient.EXPECT().Introspect(suite.context, suite.context.Request).Return(expectedResponse, error2.AuthenticationError{})
	response, err := suite.introspectionService.Introspect(suite.context, accessToken)
	suite.Equal(expectedResponse, response)
	suite.Equal(error2.AuthenticationError{}, err)
}
