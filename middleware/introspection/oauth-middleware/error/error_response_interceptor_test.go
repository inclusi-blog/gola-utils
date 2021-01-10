package error

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ErrorResponseInterceptorTest struct {
	suite.Suite
	context                  *gin.Context
	mockCtrl                 *gomock.Controller
	recorder                 *httptest.ResponseRecorder
	errorResponseInterceptor ResponseInterceptor
}

func TestErrorResponseInterceptorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorResponseInterceptorTest))
}

func (suite *ErrorResponseInterceptorTest) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("", "", nil)
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.errorResponseInterceptor = NewErrorResponseInterceptor()
}

func (suite *ErrorResponseInterceptorTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *ErrorResponseInterceptorTest) TestShouldAbortWithNoJsonBodyWhenHttpStatusIsBetween500And600() {
	oauthError := &OAuthMiddlewareError{HttpStatusCode: 500, ErrorResponse: Error{ErrorCode: "errorCode", ErrorMessage: "errorMessage"}}

	suite.errorResponseInterceptor.HandleServiceError(suite.context, oauthError)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Empty(suite.recorder.Body.String())
}

func (suite ErrorResponseInterceptorTest) TestShouldAbortWithJsonBodyWhenHttpStatusIsOtherThanBetween500And600() {
	oauthError := &OAuthMiddlewareError{HttpStatusCode: 401, ErrorResponse: Error{ErrorCode: "errorCode", ErrorMessage: "errorMessage"}}
	suite.errorResponseInterceptor.HandleServiceError(suite.context, oauthError)

	suite.Equal(http.StatusUnauthorized, suite.recorder.Code)
	expectedResponseBytes, _ := json.Marshal(gin.H{"error": "errorMessage", "errorCode": "errorCode", "errorMessage": "errorMessage"})
	suite.Equal(string(expectedResponseBytes), suite.recorder.Body.String())
}
