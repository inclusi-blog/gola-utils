package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/constants"
	"github.com/gola-glitch/gola-utils/mocks"
	"github.com/gola-glitch/gola-utils/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TokenMiddlewareTest struct {
	suite.Suite
	ginEngine      *gin.Engine
	mockController *gomock.Controller
	oauthUtils     *mocks.MockUtils
}

func TestTokenMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(TokenMiddlewareTest))
}

func (suite *TokenMiddlewareTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.oauthUtils = mocks.NewMockUtils(suite.mockController)
}

func (suite TokenMiddlewareTest) TestShouldAddIdTokenToContext() {

	router := gin.New()
	middleware := NewTokenMiddleware(suite.oauthUtils)
	expectedIdTokenVal := model.IdToken{
		UserId:    "dsf@as.com",
	}

	router.GET("/test", middleware.DecryptIdToken(), func(context *gin.Context) {
		idTokenVal, _ := context.Get("encIdToken")
		idTokenVal = idTokenVal.(model.IdToken)
		assert.Equal(suite.T(), expectedIdTokenVal, idTokenVal)
		context.JSON(http.StatusOK, gin.H{})
	})

	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(gomock.Any()).Return(expectedIdTokenVal, nil).Times(1)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer accessToken")
	req.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, "idToken")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Result().StatusCode)
}

func (suite TokenMiddlewareTest) TestShouldReturnUnAuthorized_WhenAccessTokenNotExists() {

	router := gin.New()
	middleware := NewTokenMiddleware(suite.oauthUtils)

	router.GET("/test", middleware.DecryptIdToken(), nil)
	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(gomock.Any()).Return(model.IdToken{
		UserId:    "dsf@as.com",
	}, nil).Times(1)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, "idToken")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Result().StatusCode)
}

func (suite TokenMiddlewareTest) TestShouldReturnUnAuthorized_WhenIdTokenNotExists() {

	router := gin.New()
	middleware := NewTokenMiddleware(suite.oauthUtils)

	router.GET("/test", middleware.DecryptIdToken(), nil)
	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(gomock.Any()).Return(model.IdToken{
		UserId:    "dsf@as.com",
	}, nil).Times(1)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer accessToken")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Result().StatusCode)
}

func (suite TokenMiddlewareTest) TestShouldReturnInternalServerError_WhenDecodeTokenFails() {

	router := gin.New()
	middleware := NewTokenMiddleware(suite.oauthUtils)

	router.GET("/test", middleware.DecryptIdToken(), nil)
	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(gomock.Any()).
		Return(model.IdToken{}, errors.New("some error")).Times(1)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer accessToken")
	req.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, "idToken")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Result().StatusCode)
}
