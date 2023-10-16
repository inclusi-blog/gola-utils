package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/constants"
	"github.com/inclusi-blog/gola-utils/golang_error"
	error2 "github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/error"
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/mocks"
	oauth "github.com/inclusi-blog/gola-utils/mocks"
	"github.com/inclusi-blog/gola-utils/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TokenServiceTest struct {
	suite.Suite
	recorder                 *httptest.ResponseRecorder
	context                  *gin.Context
	mockCtrl                 *gomock.Controller
	oauthUtils               *oauth.MockUtils
	mockIntrospectionService *mocks.MockIntrospectionService
	tokenService             TokenService
}

func TestTokenServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TokenServiceTest))
}

func (suite *TokenServiceTest) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "dummyUrl", nil)
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.oauthUtils = oauth.NewMockUtils(suite.mockCtrl)
	suite.mockIntrospectionService = mocks.NewMockIntrospectionService(suite.mockCtrl)
	suite.tokenService = NewTokenService(suite.mockIntrospectionService, suite.oauthUtils)
}

func (suite *TokenServiceTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite TokenServiceTest) TestTokenServiceShouldReturnUnAuthorizedOnEmptyAccessToken() {
	suite.setupHeader("")
	err := suite.tokenService.Validate(suite.context)

	suite.Equal(error2.InvalidAccessTokenError, err)
}

func (suite TokenServiceTest) TestTokenServiceShouldReturnUnAuthorizedWhenThereIsAnErrorFetchingHeader() {
	err := suite.tokenService.Validate(suite.context)

	suite.Equal(error2.InvalidAccessTokenError, err)
}

func (suite TokenServiceTest) TestTokenServiceShouldReturnValidLoginOnValidAccessToken() {
	accessToken := "valid-access-token"
	suite.setupHeader(accessToken)
	suite.mockIntrospectionService.EXPECT().Introspect(suite.context, accessToken).Return([]byte("valid response"), nil)
	err := suite.tokenService.Validate(suite.context)

	suite.Nil(err)
}

func (suite TokenServiceTest) TestTokenServiceShouldReturnInvalidLoginOnInvalidAccessToken() {
	accessToken := "invalid-access-token"
	suite.setupHeader(accessToken)
	suite.mockIntrospectionService.EXPECT().Introspect(suite.context, accessToken).Return([]byte("invalid response"), error2.AuthenticationError{})

	err := suite.tokenService.Validate(suite.context)

	suite.Equal(error2.InvalidAccessTokenError, err)
}

func (suite TokenServiceTest) TestTokenServiceShouldReturnServerErrorWhenTokenIntrospectionReturnsServerError() {
	accessToken := "invalid-access-token"
	suite.setupHeader(accessToken)
	hydraError := error2.HydraInternalServerError{}
	suite.mockIntrospectionService.EXPECT().Introspect(suite.context, accessToken).Return([]byte("invalid response"), hydraError)

	err := suite.tokenService.Validate(suite.context)

	suite.Equal(error2.InternalServerErrorFunc(hydraError.Error()), err)
}

func (suite TokenServiceTest) setupHeader(accessToken string) {
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer "+accessToken)
}

func (suite TokenServiceTest) TestDecryptShouldReturnUnAuthorizedIfOAuthUtilsReturnsNonServerError() {
	utilError := errors.New("decryption error")
	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(suite.context).Return(model.IdToken{}, utilError)

	err := suite.tokenService.Decrypt(suite.context)

	suite.Equal(error2.InvalidIdTokenError, err)
}

func (suite TokenServiceTest) TestDecryptShouldReturnServerErrorIfOAuthUtilsReturnsServerError() {
	utilError := golang_error.InternalServerError{}
	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(suite.context).Return(model.IdToken{}, utilError)

	err := suite.tokenService.Decrypt(suite.context)

	suite.Equal(error2.InternalServerErrorFunc("error in token decryption"), err)
}

func (suite TokenServiceTest) TestDecryptShouldNotReturnAnyErrorAfterSuccessfulDecryption() {
	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(suite.context).Return(model.IdToken{UserId: "user-id"}, nil)
	err := suite.tokenService.Decrypt(suite.context)
	suite.Nil(err)
}
