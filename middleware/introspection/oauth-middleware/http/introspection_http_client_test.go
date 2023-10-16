package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/constants"
	"github.com/inclusi-blog/gola-utils/golaerror"
	request2 "github.com/inclusi-blog/gola-utils/http/request"
	"github.com/inclusi-blog/gola-utils/http/request/mocks"
	error2 "github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/error"
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	gohttp "net/http"
	"net/http/httptest"
	"testing"
)

type IntrospectionHttpClientTest struct {
	suite.Suite
	context                 *gin.Context
	mockCtrl                *gomock.Controller
	introspectionHttpClient IntrospectionHttpClient
	mockHttpRequest         *mocks.MockHttpRequest
	mockHttpRequestBuilder  *mocks.MockHttpRequestBuilder
	recorder                *httptest.ResponseRecorder
	ginEngine               *gin.Engine
}

func TestIntrospectionMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(IntrospectionHttpClientTest))
}

func (suite *IntrospectionHttpClientTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.mockHttpRequest = mocks.NewMockHttpRequest(suite.mockCtrl)
	suite.mockHttpRequestBuilder = mocks.NewMockHttpRequestBuilder(suite.mockCtrl)
}

func (suite IntrospectionHttpClientTest) TestClientShouldReturnStatusOkOnActiveResponse() {
	suite.context.Request, _ = gohttp.NewRequest("GET", "random/url", nil)
	suite.introspectionHttpClient = NewIntrospectionClient(suite.mockHttpRequestBuilder)

	expectedIntrospectionResponse := model.IntrospectionResponse{
		Active: true,
	}
	introspectionResponse := model.IntrospectionResponse{}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(suite.context).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).DoAndReturn(func(responseHolder interface{}) request2.HttpRequest {
		tempResponsePointer := responseHolder.(*model.IntrospectionResponse)
		*tempResponsePointer = expectedIntrospectionResponse
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Post("random/url").Return(nil)

	expectedResponse, _ := json.Marshal(model.IntrospectionResponse{Active: true})

	actualResponse, actualError := suite.introspectionHttpClient.Introspect(suite.context, suite.context.Request)
	suite.Nil(actualError)
	suite.Equal(string(expectedResponse), string(actualResponse))
}

func (suite IntrospectionHttpClientTest) TestClientShouldReturnAuthenticationErrorOnInActiveResponse() {
	suite.context.Request, _ = gohttp.NewRequest("GET", "random/url", nil)
	suite.introspectionHttpClient = NewIntrospectionClient(suite.mockHttpRequestBuilder)

	expectedIntrospectionResponse := model.IntrospectionResponse{
		Active: false,
	}
	introspectionResponse := model.IntrospectionResponse{}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer invalid-token")
	introspectionRequest["token"] = "invalid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(suite.context).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).DoAndReturn(func(responseHolder interface{}) request2.HttpRequest {
		tempResponsePointer := responseHolder.(*model.IntrospectionResponse)
		*tempResponsePointer = expectedIntrospectionResponse
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Post("random/url").Return(nil)

	actualResponse, actualError := suite.introspectionHttpClient.Introspect(suite.context, suite.context.Request)
	expectedResponse, _ := json.Marshal(model.IntrospectionResponse{Active: false})
	suite.Equal(string(expectedResponse), string(actualResponse))
	suite.Equal(error2.AuthenticationError{}, actualError)
}

func (suite IntrospectionHttpClientTest) TestClientShouldReturnHydraInternalServerErrorOn500Response() {
	suite.context.Request, _ = gohttp.NewRequest("GET", "random/url", nil)
	suite.introspectionHttpClient = NewIntrospectionClient(suite.mockHttpRequestBuilder)

	introspectionResponse := model.IntrospectionResponse{}
	httpError := golaerror.HttpError{StatusCode: gohttp.StatusInternalServerError}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(suite.context).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("random/url").Return(httpError)

	_, actualError := suite.introspectionHttpClient.Introspect(suite.context, suite.context.Request)
	suite.Equal(error2.HydraInternalServerError{}, actualError)
}

func (suite IntrospectionHttpClientTest) TestClientShouldReturnAuthenticationErrorForServerErrorOtherThan500() {
	suite.context.Request, _ = gohttp.NewRequest("GET", "random/url", nil)
	suite.introspectionHttpClient = NewIntrospectionClient(suite.mockHttpRequestBuilder)

	introspectionResponse := model.IntrospectionResponse{}
	httpError := golaerror.HttpError{StatusCode: gohttp.StatusUnauthorized}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(suite.context).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("random/url").Return(httpError)

	_, actualError := suite.introspectionHttpClient.Introspect(suite.context, suite.context.Request)
	suite.Equal(error2.AuthenticationError{}, actualError)
}

func (suite IntrospectionHttpClientTest) TestClientShouldReturnHydraInternalServerErrorOn5xxResponse() {
	suite.context.Request, _ = gohttp.NewRequest("GET", "random/url", nil)
	suite.introspectionHttpClient = NewIntrospectionClient(suite.mockHttpRequestBuilder)

	introspectionResponse := model.IntrospectionResponse{}
	httpError := golaerror.HttpError{StatusCode: gohttp.StatusHTTPVersionNotSupported}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(suite.context).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("random/url").Return(httpError)

	_, actualError := suite.introspectionHttpClient.Introspect(suite.context, suite.context.Request)
	suite.Equal(error2.HydraInternalServerError{}, actualError)
}

func (suite IntrospectionHttpClientTest) TestClientShouldReturnHttpErrorOnEmptyUrl() {
	suite.context.Request, _ = gohttp.NewRequest("GET", "", nil)
	suite.introspectionHttpClient = NewIntrospectionClient(suite.mockHttpRequestBuilder)

	introspectionResponse := model.IntrospectionResponse{}

	httpError := golaerror.HttpError{StatusCode: gohttp.StatusInternalServerError}

	introspectionRequest := make(map[string]interface{})
	suite.context.Request.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer valid-token")
	introspectionRequest["token"] = "valid-token"

	suite.mockHttpRequestBuilder.EXPECT().NewRequestWithContext(suite.context).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Content-Type", "application/x-www-form-urlencoded").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddHeader("Accept", "application/json").Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(introspectionRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&introspectionResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("").Return(httpError)
	_, actualError := suite.introspectionHttpClient.Introspect(suite.context, suite.context.Request)
	suite.Equal(error2.HydraInternalServerError{}, actualError)
}

func httpClientForTokenIntrospection(isTokenActive bool) gohttp.Client {
	return NewHttpTestClient(func(req *gohttp.Request) *gohttp.Response {
		response, _ := json.Marshal(model.IntrospectionResponse{Active: isTokenActive})
		return &gohttp.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(string(response))),
			Header:     make(gohttp.Header),
		}
	})
}

func httpClientWithResponseServerError(statusCode int) gohttp.Client {
	return NewHttpTestClient(func(req *gohttp.Request) *gohttp.Response {
		return &gohttp.Response{
			StatusCode: statusCode,
			Header:     make(gohttp.Header),
		}
	})
}

func httpClientWithInvalidResponse() gohttp.Client {
	return NewHttpTestClient(func(req *gohttp.Request) *gohttp.Response {
		return &gohttp.Response{
			StatusCode: 200,
			Header:     make(gohttp.Header),
			Body:       errReader(0),
		}
	})
}

type errReader int

func (errReader) Close() (err error) {
	return nil
}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("could not ready body")
}

func httpClientWithResponseHttpError() gohttp.Client {
	client, _ := NewHttpTestClientWithError(func(req *gohttp.Request) (*gohttp.Response, error) {
		return nil, error2.AuthenticationError{}
	})
	return client
}

type RoundTripFunc func(req *gohttp.Request) *gohttp.Response

type RoundTripFuncWithError func(req *gohttp.Request) (*gohttp.Response, error)

func (f RoundTripFuncWithError) RoundTrip(req *gohttp.Request) (*gohttp.Response, error) {
	return nil, error2.AuthenticationError{}
}

func (f RoundTripFunc) RoundTrip(req *gohttp.Request) (*gohttp.Response, error) {
	return f(req), nil
}

func NewHttpTestClient(fn RoundTripFunc) gohttp.Client {
	return gohttp.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewHttpTestClientWithError(fn RoundTripFuncWithError) (gohttp.Client, error) {
	return gohttp.Client{
		Transport: RoundTripFuncWithError(fn),
	}, error2.AuthenticationError{}
}
