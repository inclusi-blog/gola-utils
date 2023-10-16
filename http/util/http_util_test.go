package util

import (
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/constants"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HttpUtilTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	recorder    *httptest.ResponseRecorder
	ginContext  *gin.Context
	accessToken string
	encIDToken  string
}

func TestHttpUtilTestSuite(t *testing.T) {
	suite.Run(t, new(HttpUtilTestSuite))
}

func (suite *HttpUtilTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.ginContext, _ = gin.CreateTestContext(suite.recorder)
	suite.ginContext.Request, _ = http.NewRequest("GET", "dummyUrl", nil)
	suite.accessToken = "accessToken"
	suite.encIDToken = "encIDToken"
}

func (suite HttpUtilTestSuite) TestGetAccessTokenFromRequestHeader() {
	suite.ginContext.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer "+suite.accessToken)
	cookieWithAccessToken := http.Cookie{Name: constants.COOKIE_ACCESS_TOKEN, Value: suite.accessToken}
	suite.ginContext.Request.AddCookie(&cookieWithAccessToken)

	accessToken, err := GetAccessToken(suite.ginContext)

	suite.Equal(suite.accessToken, accessToken)
	suite.Nil(err)
}

func (suite HttpUtilTestSuite) TestGetAccessTokenShouldReturnErrorIfAuthorizationHeaderIsEmpty() {
	cookieWithoutAccessToken := http.Cookie{Name: "invalid", Value: suite.accessToken}
	suite.ginContext.Request.AddCookie(&cookieWithoutAccessToken)
	accessToken, err := GetAccessToken(suite.ginContext)

	suite.Equal("", accessToken)
	suite.NotNil(err)
	suite.Equal("invalid bearer/cookie header", err.Error())
}

func (suite HttpUtilTestSuite) TestGetAccessTokenShouldReturnErrorIfAuthorizationHeaderIsInvalidWithoutBearer() {
	suite.ginContext.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, suite.accessToken)

	accessToken, err := GetAccessToken(suite.ginContext)

	suite.Equal("", accessToken)
	suite.NotNil(err)
	suite.Equal("invalid bearer/cookie header", err.Error())
}

func (suite HttpUtilTestSuite) TestGetAccessTokenShouldReturnErrorIfAuthorizationHeaderIsInvalid() {
	suite.ginContext.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, suite.accessToken)
	cookieWithoutAccessToken := http.Cookie{Name: "invalid", Value: suite.accessToken}
	suite.ginContext.Request.AddCookie(&cookieWithoutAccessToken)

	accessToken, err := GetAccessToken(suite.ginContext)

	suite.Equal("", accessToken)
	suite.NotNil(err)
	suite.Equal("invalid bearer/cookie header", err.Error())
}

func (suite HttpUtilTestSuite) TestGetAccessTokenFromCookieIfAuthorizationHeaderIsNotAvailable() {
	cookieWithAccessToken := http.Cookie{Name: constants.COOKIE_ACCESS_TOKEN, Value: suite.accessToken}
	suite.ginContext.Request.AddCookie(&cookieWithAccessToken)

	accessToken, err := GetAccessToken(suite.ginContext)

	suite.Equal(suite.accessToken, accessToken)
	suite.Nil(err)
}

func (suite HttpUtilTestSuite) TestGetAccessTokenShouldReturnErrorIfHeaderAndCookieAreNotAvailable() {
	accessToken, err := GetAccessToken(suite.ginContext)

	suite.Equal("", accessToken)
	suite.NotNil(err)
	suite.Equal("invalid bearer/cookie header", err.Error())
}

func (suite HttpUtilTestSuite) TestGetEncIDTokenFromRequestHeader() {
	suite.ginContext.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, suite.encIDToken)

	encIDToken, err := GetEncryptedIDToken(suite.ginContext)

	suite.Equal(suite.encIDToken, encIDToken)
	suite.Nil(err)
}

func (suite HttpUtilTestSuite) TestGetEncIDTokenShouldReturnErrorIfRequestHeaderIsEmpty() {
	cookieWithoutEncIDToken := http.Cookie{Name: "invalid", Value: suite.encIDToken}
	suite.ginContext.Request.AddCookie(&cookieWithoutEncIDToken)

	encIDToken, err := GetEncryptedIDToken(suite.ginContext)

	suite.Equal("", encIDToken)
	suite.NotNil(err)
	suite.Equal("invalid enc-id-token header/cookie", err.Error())
}

func (suite HttpUtilTestSuite) TestGetEncIDTokenFromCookieIfRequestHeaderIsEmpty() {
	cookieWithEncIDToken := http.Cookie{Name: constants.COOKIE_ENC_ID_TOKEN, Value: suite.encIDToken}
	suite.ginContext.Request.AddCookie(&cookieWithEncIDToken)

	encIDToken, err := GetEncryptedIDToken(suite.ginContext)

	suite.Equal(suite.encIDToken, encIDToken)
	suite.Nil(err)
}

func (suite HttpUtilTestSuite) TestGetEncIDTokenShouldReturnErrorIfHeaderAndCookieIsInvalid() {
	encIDToken, err := GetEncryptedIDToken(suite.ginContext)

	suite.Equal("", encIDToken)
	suite.NotNil(err)
	suite.Equal("invalid enc-id-token header/cookie", err.Error())
}

func (suite HttpUtilTestSuite) TestFormBearerAuthorizationHeaderShouldReturnAuthHeader() {
	authHeader := FormBearerAuthorizationHeader(suite.accessToken)

	suite.Equal("Bearer accessToken", authHeader)
}
