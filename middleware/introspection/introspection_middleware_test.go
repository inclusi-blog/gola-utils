package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/inclusi-blog/gola-utils/constants"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/http/request"
	mocks3 "github.com/inclusi-blog/gola-utils/http/request/mocks"
	error2 "github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/error"
	middlewareHttp "github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/http"
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/mocks"
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/model"
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/service"
	mocks2 "github.com/inclusi-blog/gola-utils/mocks"
	model2 "github.com/inclusi-blog/gola-utils/model"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	gohttp "net/http"
	"net/http/httptest"
	"testing"
)

type IntrospectionMiddlewareTest struct {
	suite.Suite
	mockCtrl                  *gomock.Controller
	recorder                  *httptest.ResponseRecorder
	mockTokenService          *mocks.MockTokenService
	mockProtectedUrlService   *mocks.MockProtectedUrlService
	oauthUtils                *mocks2.MockUtils
	errorRessponseInterceptor *mocks.MockResponseInterceptor
	mockIntrospectionService  *mocks.MockIntrospectionService
	ginEngine                 *gin.Engine
	context                   *gin.Context
	mockHttpRequest           *mocks3.MockHttpRequest
	mockHttpRequestBuilder    *mocks3.MockHttpRequestBuilder
}

func TestIntrospectionMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(IntrospectionMiddlewareTest))
}

func (suite *IntrospectionMiddlewareTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.mockTokenService = mocks.NewMockTokenService(suite.mockCtrl)
	suite.mockProtectedUrlService = mocks.NewMockProtectedUrlService(suite.mockCtrl)
	suite.errorRessponseInterceptor = mocks.NewMockResponseInterceptor(suite.mockCtrl)
	suite.mockIntrospectionService = mocks.NewMockIntrospectionService(suite.mockCtrl)
	suite.mockHttpRequest = mocks3.NewMockHttpRequest(suite.mockCtrl)
	suite.mockHttpRequestBuilder = mocks3.NewMockHttpRequestBuilder(suite.mockCtrl)
	suite.oauthUtils = mocks2.NewMockUtils(suite.mockCtrl)
	suite.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.context.Request, _ = gohttp.NewRequest("GET", "", nil)
	suite.ginEngine = gin.Default()
}

func (suite *IntrospectionMiddlewareTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite IntrospectionMiddlewareTest) TestShouldReturnIntrospectionMiddlewareStructWithSkipTokenDecryptTrueWhenCallingOldVersion() {
	middleware := NewIntrospectionMiddleware(suite.mockProtectedUrlService, "hydra-url")

	suite.True(middleware.isSkipTokenDecrypt())
}

func (suite IntrospectionMiddlewareTest) TestShouldReturnIntrospectionMiddlewareStructWithSkipTokenDecryptFalseWhenCallingLatestVersion() {
	middleware := NewIntrospectionAndDecryptionMiddleware(suite.mockProtectedUrlService, "hydra-url", suite.oauthUtils)

	suite.False(middleware.isSkipTokenDecrypt())
}

func (suite IntrospectionMiddlewareTest) TestIntrospectionMiddlewareShouldReturnStatusOkOnUnprotectedUrl() {
	unprotectedUrl := "/sample/unprotected/url"
	req, _ := gohttp.NewRequest("GET", unprotectedUrl, nil)
	suite.mockProtectedUrlService.EXPECT().IsProtected(unprotectedUrl).Return(false)
	w := suite.setupServerAndServe(req, gohttp.Client{})
	suite.Equal(200, w.Code)
	response, _ := ioutil.ReadAll(w.Result().Body)
	suite.Equal("This is an unprotected resource", string(response))
}

func (suite IntrospectionMiddlewareTest) TestIntrospectionMiddlewareShouldReturnUnauthorizedForProtectedUrlWithInvalidAccessToken() {
	protectedUrl := "sample/protected/url"
	req, _ := gohttp.NewRequest("GET", protectedUrl, nil)
	req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer invalid-token")
	suite.mockProtectedUrlService.EXPECT().IsProtected(protectedUrl).Return(true)

	introspectionResponse := model.IntrospectionResponse{}
	httpError := golaerror.HttpError{StatusCode: gohttp.StatusUnauthorized}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer invalid-token")
	introspectionRequest["token"] = "invalid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(gomock.Any()).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("hydra_base_admin_url/admin/oauth2/introspect").Return(httpError)

	suite.errorRessponseInterceptor.EXPECT().HandleServiceError(gomock.Any(), error2.InvalidAccessTokenError)

	suite.setupServerAndServe(req, gohttp.Client{})
}

func (suite IntrospectionMiddlewareTest) TestIntrospectionMiddlewareShouldReturnStatusOkForProtectedUrlWithValidAccessTokenAndShouldNotDecryptToken() {
	protectedUrl := "/sample/protected/url"
	req, _ := gohttp.NewRequest("GET", protectedUrl, nil)
	req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	suite.mockProtectedUrlService.EXPECT().IsProtected(protectedUrl).Return(true)

	expectedIntrospectionResponse := model.IntrospectionResponse{
		Active: true,
	}
	introspectionResponse := model.IntrospectionResponse{}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(gomock.Any()).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).DoAndReturn(func(responseHolder interface{}) request.HttpRequest {
		tempResponsePointer := responseHolder.(*model.IntrospectionResponse)
		*tempResponsePointer = expectedIntrospectionResponse
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Post("hydra_base_admin_url/admin/oauth2/introspect").Return(nil)

	w := suite.setupServerAndServe(req, gohttp.Client{})
	suite.Equal(200, w.Code)
	response, _ := ioutil.ReadAll(w.Result().Body)
	suite.Equal("This is a protected resource", string(response))
}

func (suite IntrospectionMiddlewareTest) TestIntrospectionMiddlewareShouldReturnInternalErrorWhenHttpClientReturns500ForIntrospection() {
	protectedUrl := "/sample/protected/url"
	req, _ := gohttp.NewRequest("GET", protectedUrl, nil)
	req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	suite.mockProtectedUrlService.EXPECT().IsProtected(protectedUrl).Return(true)

	introspectionResponse := model.IntrospectionResponse{}
	httpError := golaerror.HttpError{StatusCode: gohttp.StatusInternalServerError}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(gomock.Any()).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("hydra_base_admin_url/admin/oauth2/introspect").Return(httpError)
	suite.errorRessponseInterceptor.EXPECT().HandleServiceError(gomock.Any(), error2.InternalServerErrorFunc("Internal error in Hydra"))

	suite.setupServerAndServe(req, gohttp.Client{})
}

func (suite IntrospectionMiddlewareTest) TestIntrospectionMiddlewareShouldReturnInternalErrorWhenHttpClientReturns5xxForIntrospection() {
	protectedUrl := "sample/protected/url"
	req, _ := gohttp.NewRequest("GET", protectedUrl, nil)
	req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	suite.mockProtectedUrlService.EXPECT().IsProtected(protectedUrl).Return(true)

	introspectionResponse := model.IntrospectionResponse{}
	httpError := golaerror.HttpError{StatusCode: gohttp.StatusInternalServerError}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(gomock.Any()).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("hydra_base_admin_url/admin/oauth2/introspect").Return(httpError)

	suite.errorRessponseInterceptor.EXPECT().HandleServiceError(gomock.Any(), error2.InternalServerErrorFunc("Internal error in Hydra"))

	suite.setupServerAndServe(req, gohttp.Client{})
}

func (suite IntrospectionMiddlewareTest) TestIntrospectionMiddlewareShouldReturnStatusOkIfTokenValidationAndTokenDecryptionDoNotReturnAnyError() {
	protectedUrl := "/sample-decryption/protected/url"
	req, _ := gohttp.NewRequest("GET", protectedUrl, nil)
	req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	req.Header.Add(constants.ENC_ID_TOKEN_HEADER_KEY, "enc-token")
	suite.mockProtectedUrlService.EXPECT().IsProtected(protectedUrl).Return(true)

	expectedIntrospectionResponse := model.IntrospectionResponse{
		Active: true,
	}
	introspectionResponse := model.IntrospectionResponse{}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(gomock.Any()).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).DoAndReturn(func(responseHolder interface{}) request.HttpRequest {
		tempResponsePointer := responseHolder.(*model.IntrospectionResponse)
		*tempResponsePointer = expectedIntrospectionResponse
		return suite.mockHttpRequest
	})

	suite.mockHttpRequest.EXPECT().Post("hydra_base_admin_url/admin/oauth2/introspect").Return(nil)

	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(gomock.Any()).Return(model2.IdToken{}, nil)

	r := suite.setupGinServerForDecryption(req, gohttp.Client{})
	r.ServeHTTP(suite.recorder, req)

	suite.Equal(gohttp.StatusNoContent, suite.recorder.Code)
}

func (suite IntrospectionMiddlewareTest) TestIntrospectionMiddlewareShouldReturnStatusOkIfTokenValidationAndTokenDecryptionNotToBeCalledIfEncTokenMissing() {
	protectedUrl := "/sample-decryption/protected/url"
	req, _ := gohttp.NewRequest("GET", protectedUrl, nil)
	req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	suite.mockProtectedUrlService.EXPECT().IsProtected(protectedUrl).Return(true)
	expectedIntrospectionResponse := model.IntrospectionResponse{
		Active: true,
	}
	introspectionResponse := model.IntrospectionResponse{}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(gomock.Any()).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).DoAndReturn(func(responseHolder interface{}) request.HttpRequest {
		tempResponsePointer := responseHolder.(*model.IntrospectionResponse)
		*tempResponsePointer = expectedIntrospectionResponse
		return suite.mockHttpRequest
	})

	suite.mockHttpRequest.EXPECT().Post("hydra_base_admin_url/admin/oauth2/introspect").Return(nil)
	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(gomock.Any()).Times(0)

	r := suite.setupGinServerForDecryption(req, gohttp.Client{})
	r.ServeHTTP(suite.recorder, req)

	suite.Equal(gohttp.StatusNoContent, suite.recorder.Code)
}

func (suite IntrospectionMiddlewareTest) TestIntrospectionMiddlewareShouldHandleServiceErrorIfOAuthUtilsReturnsError() {
	protectedUrl := "sample-decryption/protected/url"
	suite.context.Request, _ = gohttp.NewRequest("GET", protectedUrl, nil)
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	suite.context.Request.Header.Add(constants.ENC_ID_TOKEN_HEADER_KEY, "enc-token")
	suite.mockProtectedUrlService.EXPECT().IsProtected(protectedUrl).Return(true)

	expectedIntrospectionResponse := model.IntrospectionResponse{
		Active: true,
	}
	introspectionResponse := model.IntrospectionResponse{}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(gomock.Any()).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).DoAndReturn(func(responseHolder interface{}) request.HttpRequest {
		tempResponsePointer := responseHolder.(*model.IntrospectionResponse)
		*tempResponsePointer = expectedIntrospectionResponse
		return suite.mockHttpRequest
	})

	suite.mockHttpRequest.EXPECT().Post("hydra_base_admin_url/admin/oauth2/introspect").Return(nil)

	suite.oauthUtils.EXPECT().DecodeEncryptedIdToken(gomock.Any()).Return(model2.IdToken{}, errors.New("decryption error"))
	suite.errorRessponseInterceptor.EXPECT().HandleServiceError(gomock.Any(), error2.InvalidIdTokenError)

	r := suite.setupGinServerForDecryption(suite.context.Request, gohttp.Client{})
	r.ServeHTTP(suite.recorder, suite.context.Request)
}

func (suite IntrospectionMiddlewareTest) setupServerAndServe(req *gohttp.Request, client gohttp.Client) *httptest.ResponseRecorder {
	introspectionMiddleware := suite.introspectionMiddleware(client)
	r := gin.Default()
	r.Use(introspectionMiddleware.TokenValidationMiddleware())

	r.GET("sample/unprotected/url", func(c *gin.Context) {
		c.String(gohttp.StatusOK, "This is an unprotected resource")
	})
	r.GET("sample/protected/url", func(c *gin.Context) {
		c.String(gohttp.StatusOK, "This is a protected resource")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func (suite *IntrospectionMiddlewareTest) setupGinServerForDecryption(req *gohttp.Request, client gohttp.Client) *gin.Engine {
	introspectionMiddleware := suite.introspectionAndDecryptionMiddleware(client)
	r := gin.Default()
	r.Use(introspectionMiddleware.TokenValidationMiddleware())

	r.GET("sample-decryption/protected/url", func(c *gin.Context) {
		c.Status(gohttp.StatusNoContent)
	})

	return r
}

func (suite IntrospectionMiddlewareTest) introspectionMiddleware(client gohttp.Client) introspectionMiddleware {
	hydraAdminUrl := "hydra_base_admin_url"

	return introspectionMiddleware{tokenService: service.NewTokenService(
		service.NewIntrospectionService(middlewareHttp.NewIntrospectionClient(suite.mockHttpRequestBuilder), middlewareHttp.NewIntrospectionRequestBuilder(), hydraAdminUrl), nil),
		protectedUrlService:      suite.mockProtectedUrlService,
		errorResponseInterceptor: suite.errorRessponseInterceptor,
		skipTokenDecrypt:         true}
}

func (suite IntrospectionMiddlewareTest) introspectionAndDecryptionMiddleware(client gohttp.Client) introspectionMiddleware {
	hydraAdminUrl := "hydra_base_admin_url"

	return introspectionMiddleware{tokenService: service.NewTokenService(
		service.NewIntrospectionService(middlewareHttp.NewIntrospectionClient(suite.mockHttpRequestBuilder), middlewareHttp.NewIntrospectionRequestBuilder(), hydraAdminUrl), suite.oauthUtils),
		protectedUrlService:      suite.mockProtectedUrlService,
		errorResponseInterceptor: suite.errorRessponseInterceptor,
		skipTokenDecrypt:         false}
}
