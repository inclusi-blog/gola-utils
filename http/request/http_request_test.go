package request

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gola-glitch/gola-utils/constants"
	utilError "github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/http/model"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/http/client/mocks"
	traceMocks "github.com/gola-glitch/gola-utils/trace/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
)

type dummyRequest struct {
	FieldA string  `json:"fieldA"`
	FieldB float64 `json:"fieldB"`
}

type dummyResponse struct {
	ResponseFieldA string `json:"responseFieldA" validate:"required"`
}

type dummyXMLRequest struct {
	XMLName xml.Name `xml:"payload"`
	FieldA  string   `xml:"fieldA"`
	FieldB  float64  `xml:"fieldB"`
}

type dummyXMLResponse struct {
	ResponseFieldA string `xml:"responseFieldA" validate:"required"`
}

type dummyErrorResponse struct {
	ErrorCode    int
	ErrorMessage string
}

type dummyResponseWithCustomValidation struct {
	ResponseFieldA string `json:"responseFieldA" validate:"isBool"`
}

type SpanExporter struct {
	spans []*trace.SpanData
}

func (s *SpanExporter) ExportSpan(d *trace.SpanData) {
	s.spans = append(s.spans, d)
}
func (s *SpanExporter) Reset() {
	s.spans = []*trace.SpanData{}
}

func dummyValidationIsBool(fl validator.FieldLevel) bool {
	return strings.ToLower(fl.Field().String()) == "true" ||
		strings.ToLower(fl.Field().String()) == "false"
}

type HttpRequestTestSuite struct {
	suite.Suite
	mockCtrl           *gomock.Controller
	mockHttpClient     *mocks.MockHttpClient
	recorder           *httptest.ResponseRecorder
	context            *gin.Context
	httpRequestBuilder HttpRequestBuilder
	requestBody        dummyRequest
	requestBodyAsBytes []byte
	XMLRequestBody     dummyXMLRequest
	url                string
	trace              *traceMocks.MockTrace
}

func TestHttpRequestTestSuite(t *testing.T) {
	suite.Run(t, new(HttpRequestTestSuite))
}

func (suite *HttpRequestTestSuite) SetupTest() {
	suite.url = "http://dummyurl.com"
	suite.requestBody = dummyRequest{FieldA: "sample value"}
	suite.requestBodyAsBytes, _ = json.Marshal(suite.requestBody)
	suite.XMLRequestBody = dummyXMLRequest{
		XMLName: xml.Name{},
		FieldA:  "A",
		FieldB:  1,
	}

	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockHttpClient = mocks.NewMockHttpClient(suite.mockCtrl)
	suite.httpRequestBuilder = NewHttpRequestBuilder(suite.mockHttpClient)
	suite.trace = traceMocks.NewMockTrace(suite.mockCtrl)
}

func (suite HttpRequestTestSuite) TestShouldMakeRequest() {
	defer suite.mockCtrl.Finish()
	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyResponse{ResponseFieldA: "sample response value"}
			responseJSONString, _ := json.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: constants.CONTEXT_ENC_ID_TOKEN, Value: "dummyencIdToken"}
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()
				requestByte, _ := ioutil.ReadAll(actualRequest.Body)
				ioutil.NopCloser(bytes.NewBuffer(requestByte))
				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "application/json") &&
					actualRequest.Header.Get("UserId") == "ABC123" &&
					actualRequest.Header.Get("Password") == "PWD456" &&
					actualRequest.Header.Get("Session-Tracing-ID") == "" &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value &&
					(actualRequest.Method == "GET" || string(requestByte) == string(suite.requestBodyAsBytes)) {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})
			var actualResponse dummyResponse
			var responseStatusCode int
			headers := map[string]string{
				"UserId":   "ABC123",
				"Password": "PWD456",
			}
			builder := suite.httpRequestBuilder.
				NewRequest().
				WithContext(context.Background()).
				ResponseStatusCodeAs(&responseStatusCode).
				ResponseAs(&actualResponse).AddHeaders(headers).
				AddCookie(cookie)

			var err error
			if test.method == "POST" {
				err = builder.
					WithJSONBody(suite.requestBody).
					Post(suite.url)

			} else if test.method == "GET" {
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithJSONBody(suite.requestBody).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(http.StatusOK, responseStatusCode)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldMakeRequestWithRequestBodyAndEscapeSpecialCharacters() {
	defer suite.mockCtrl.Finish()

	expectedResponse := dummyResponse{ResponseFieldA: "sample response value"}
	responseJSONString, _ := json.Marshal(expectedResponse)

	requestWithSpecialChar := dummyRequest{FieldA: "sample & value"}
	requestBodyAsBytes, _ := json.Marshal(requestWithSpecialChar)
	cookie := &http.Cookie{Name: constants.CONTEXT_ENC_ID_TOKEN, Value: "dummyencIdToken"}
	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		actualCookies := actualRequest.Cookies()
		requestByte, _ := ioutil.ReadAll(actualRequest.Body)
		ioutil.NopCloser(bytes.NewBuffer(requestByte))
		if actualRequest.URL.String() == suite.url &&
			actualRequest.Method == "POST" &&
			actualRequest.Header.Get("Content-Type") == "application/json" &&
			actualRequest.Header.Get("UserId") == "ABC123" &&
			actualRequest.Header.Get("Password") == "PWD456" &&
			actualRequest.Header.Get("Session-Tracing-ID") == "" &&
			len(actualCookies) == 1 &&
			actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value &&
			string(requestByte) == string(requestBodyAsBytes) {
			return &response, nil
		}
		return &http.Response{}, errors.New("request not matching")
	})
	var actualResponse dummyResponse
	var responseStatusCode int
	headers := map[string]string{
		"UserId":   "ABC123",
		"Password": "PWD456",
	}
	builder := suite.httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		ResponseStatusCodeAs(&responseStatusCode).
		ResponseAs(&actualResponse).AddHeaders(headers).
		AddCookie(cookie)

	err := builder.
		WithJSONBody(requestWithSpecialChar).
		Post(suite.url)

	suite.Nil(err)
	suite.Equal(http.StatusOK, responseStatusCode)
	suite.Equal(expectedResponse, actualResponse)

}

func (suite HttpRequestTestSuite) TestShouldMakeRequestWithJSONBodyNoEscapeHTML() {

	requestBodyWithSpecialChar := dummyRequest{FieldA: "sample & value"}
	var requestBodyWithSpecialCharBuffer bytes.Buffer
	encoder := json.NewEncoder(&requestBodyWithSpecialCharBuffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(requestBodyWithSpecialChar)
	requestBodyWithSpecialCharAsBytes := requestBodyWithSpecialCharBuffer.Bytes()

	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyResponse{ResponseFieldA: "sample response value"}
			responseJSONString, _ := json.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: constants.CONTEXT_ENC_ID_TOKEN, Value: "dummyencIdToken"}
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()
				requestByte, _ := ioutil.ReadAll(actualRequest.Body)
				println("expected:", string(requestByte) == string(requestBodyWithSpecialCharAsBytes))
				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "application/json") &&
					actualRequest.Header.Get("UserId") == "ABC123" &&
					actualRequest.Header.Get("Password") == "PWD456" &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value &&
					(actualRequest.Method == "GET" || string(requestByte) == string(requestBodyWithSpecialCharAsBytes)) {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			var actualResponse dummyResponse
			builder := suite.
				httpRequestBuilder.
				NewRequest().
				WithContext(context.Background()).
				ResponseAs(&actualResponse).
				AddHeader("UserId", "ABC123").
				AddHeader("Password", "PWD456").
				AddCookie(cookie)

			var err error
			if test.method == "POST" {
				err = builder.
					WithJSONBodyNoEscapeHTML(requestBodyWithSpecialChar).
					Post(suite.url)

			} else if test.method == "GET" {
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithJSONBodyNoEscapeHTML(requestBodyWithSpecialChar).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldMakeRequestWithXMLBody() {

	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyXMLResponse{ResponseFieldA: "Success"}
			responseXMLString, _ := xml.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: constants.CONTEXT_ENC_ID_TOKEN, Value: "dummyencIdToken"}
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseXMLString))}
			response.Header = http.Header{}
			response.Header.Set(constants.HeaderContentType, "application/xml;charset=windows-1252")
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()
				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "application/xml") &&
					actualRequest.Header.Get("UserId") == "ABC123" &&
					actualRequest.Header.Get("Password") == "PWD456" &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			actualResponse := dummyXMLResponse{}
			builder := suite.
				httpRequestBuilder.
				NewRequest().
				WithContext(context.Background()).
				ResponseAs(&actualResponse).
				AddHeader("UserId", "ABC123").
				AddHeader("Password", "PWD456").
				AddCookie(cookie)

			var err error
			if test.method == "POST" {
				err = builder.
					WithXMLBody(suite.XMLRequestBody).
					Post(suite.url)
			} else if test.method == "GET" {
				builder.AddHeader("Content-Type", "application/xml")
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithXMLBody(suite.XMLRequestBody).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldMakeRequestWithXMLBodyTextHeader() {

	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyXMLResponse{ResponseFieldA: "Success"}
			responseXMLString, _ := xml.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: constants.CONTEXT_ENC_ID_TOKEN, Value: "dummyencIdToken"}
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseXMLString))}
			response.Header = http.Header{}
			response.Header.Set(constants.HeaderContentType, "text/xml;charset=windows-1252")
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()
				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "text/xml") &&
					actualRequest.Header.Get("UserId") == "ABC123" &&
					actualRequest.Header.Get("Password") == "PWD456" &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			actualResponse := dummyXMLResponse{}
			builder := suite.
				httpRequestBuilder.
				NewRequest().
				WithContext(context.Background()).
				ResponseAs(&actualResponse).
				AddHeader("UserId", "ABC123").
				AddHeader("Password", "PWD456").
				AddCookie(cookie)

			var err error
			if test.method == "POST" {
				err = builder.
					WithXMLBodyTextHeader(suite.XMLRequestBody).
					Post(suite.url)
			} else if test.method == "GET" {
				builder.AddHeader("Content-Type", "text/xml")
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithXMLBodyTextHeader(suite.XMLRequestBody).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldMakeRequestWithXMLBodyTextHeaderWhenRequestInvalid() {

	var testCases = []struct {
		method string
	}{
		{method: "POST"},
	}

	suite.requestBody.FieldB = math.Inf(1)

	requestChan := make(chan string)

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyXMLResponse{ResponseFieldA: "Success"}
			responseXMLString, _ := xml.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: constants.CONTEXT_ENC_ID_TOKEN, Value: "dummyencIdToken"}
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseXMLString))}
			response.Header = http.Header{}
			response.Header.Set(constants.HeaderContentType, "text/xml;charset=windows-1252")
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()
				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "text/xml") &&
					actualRequest.Header.Get("UserId") == "ABC123" &&
					actualRequest.Header.Get("Password") == "PWD456" &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			builder := suite.
				httpRequestBuilder.
				NewRequest()
			var err error
			err = builder.
				WithXMLBodyTextHeader(requestChan).
				Post(suite.url)
			suite.NotNil(err)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldMakeRequestWithXMLBodyAndContentIsTextXML() {

	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyXMLResponse{ResponseFieldA: "Success"}
			responseXMLString, _ := xml.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: constants.CONTEXT_ENC_ID_TOKEN, Value: "dummyencIdToken"}
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseXMLString))}
			response.Header = http.Header{}
			response.Header.Set(constants.HeaderContentType, "text/xml;charset=windows-1252")
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()
				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "application/xml") &&
					actualRequest.Header.Get("UserId") == "ABC123" &&
					actualRequest.Header.Get("Password") == "PWD456" &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			actualResponse := dummyXMLResponse{}
			builder := suite.
				httpRequestBuilder.
				NewRequest().
				WithContext(context.Background()).
				ResponseAs(&actualResponse).
				AddHeader("UserId", "ABC123").
				AddHeader("Password", "PWD456").
				AddCookie(cookie)

			var err error
			if test.method == "POST" {
				err = builder.
					WithXMLBody(suite.XMLRequestBody).
					Post(suite.url)
			} else if test.method == "GET" {
				builder.AddHeader("Content-Type", "application/xml")
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithXMLBody(suite.XMLRequestBody).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldMakeRequestWithFormURLEncodedBody() {
	dummyResponse := dummyXMLResponse{ResponseFieldA: "Success"}
	responseXMLString, _ := xml.Marshal(dummyResponse)
	expectedResponse := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseXMLString))}
	expectedResponse.Header = http.Header{}
	expectedResponse.Header.Set(constants.HeaderContentType, "application/xml")

	cookie := &http.Cookie{Name: constants.CONTEXT_ENC_ID_TOKEN, Value: "dummyencIdToken"}

	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		actualCookies := actualRequest.Cookies()
		if actualRequest.URL.String() == suite.url &&
			actualRequest.Method == "POST" &&
			strings.Contains(actualRequest.Header.Get("Content-Type"), "multipart/form-data") &&
			actualRequest.Header.Get("UserId") == "ABC123" &&
			actualRequest.Header.Get("Password") == "PWD456" &&
			len(actualCookies) == 1 &&
			actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value {
			return &expectedResponse, nil
		}
		return &http.Response{}, errors.New("request not matching")
	})

	actualResponse := dummyXMLResponse{}
	builder := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		ResponseAs(&actualResponse).
		AddHeader("UserId", "ABC123").
		AddHeader("Password", "PWD456").
		AddCookie(cookie)

	err := builder.
		WithFormURLEncoded(map[string]interface{}{}).
		Post(suite.url)

	suite.Nil(err)
	suite.Equal(dummyResponse, actualResponse)
}

func (suite HttpRequestTestSuite) TestShouldProduceErrorOnBuildRequest() {
	formData := make(map[string]interface{}, 0)
	formData["file"] = multipart.FileHeader{
		Filename: "name",
		Header:   nil,
		Size:     5,
	}

	err := suite.httpRequestBuilder.
		NewRequest().
		WithFormURLEncoded(formData).
		Post("url")

	suite.NotNil(err)
}

func (suite HttpRequestTestSuite) TestShouldNotProduceErrorForFormUrlWithModel() {
	dummyResponse := dummyXMLResponse{ResponseFieldA: "Success"}
	responseXMLString, _ := xml.Marshal(dummyResponse)
	expectedResponse := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseXMLString))}
	expectedResponse.Header = http.Header{}
	expectedResponse.Header.Set(constants.HeaderContentType, "application/xml")
	formData := make(map[string]interface{}, 0)
	formData["fileContent"] = model.FileUploadContent{
		FilePath: "testdata/AC_ENTRY_POSTING_FEED_20201010150405",
		FileName: "AC_ENTRY_POSTING_FEED_20201010150405",
	}

	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		if actualRequest.URL.String() == suite.url &&
			actualRequest.Method == "POST" &&
			strings.Contains(actualRequest.Header.Get("Content-Type"), "multipart/form-data") {
			return &expectedResponse, nil
		}
		return &http.Response{}, errors.New("request not matching")
	})
	actualResponse := dummyXMLResponse{}

	err := suite.httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		ResponseAs(&actualResponse).
		WithFormURLEncoded(formData).
		Post("http://dummyurl.com")

	suite.Nil(err)
	suite.Equal(dummyResponse, actualResponse)

}

func (suite HttpRequestTestSuite) TestShouldNotProduceErrorForFormUrlWithByteArray() {
	dummyResponse := dummyXMLResponse{ResponseFieldA: "Success"}
	responseXMLString, _ := xml.Marshal(dummyResponse)
	expectedResponse := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseXMLString))}
	expectedResponse.Header = http.Header{}
	expectedResponse.Header.Set(constants.HeaderContentType, "application/xml")
	formData := make(map[string]interface{}, 0)
	formData["fileContent"] = []byte("test string")

	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		if actualRequest.URL.String() == suite.url &&
			actualRequest.Method == "POST" &&
			strings.Contains(actualRequest.Header.Get("Content-Type"), "multipart/form-data") {
			return &expectedResponse, nil
		}
		return &http.Response{}, errors.New("request not matching")
	})
	actualResponse := dummyXMLResponse{}

	err := suite.httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		ResponseAs(&actualResponse).
		WithFormURLEncoded(formData).
		Post("http://dummyurl.com")

	suite.Nil(err)
	suite.Equal(dummyResponse, actualResponse)
}

func (suite HttpRequestTestSuite) TestShouldProduceErrorOnBuildRequestForInvalidTypeOfData() {
	formData := make(map[string]interface{}, 0)
	formData["file"] = 0

	err := suite.httpRequestBuilder.
		NewRequest().
		WithFormURLEncoded(formData).
		Post("url")

	suite.NotNil(err)
	suite.EqualError(err, "ErrorCode: ERR_INVALID_REQUEST_TYPE ErrorMessage: only multipart files and strings are supported")
}

func (suite HttpRequestTestSuite) TestShouldMakeRequestWithQueryParameters() {

	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyXMLResponse{ResponseFieldA: "Success"}
			responseXMLString, _ := xml.Marshal(expectedResponse)

			urlWithParams := suite.url + "?param1=value1&param2=value2"
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseXMLString))}
			response.Header = http.Header{}
			response.Header.Set(constants.HeaderContentType, "application/xml")
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				if actualRequest.URL.String() == urlWithParams &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "application/xml") {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			queryParams := map[string]string{"param1": "value1", "param2": "value2"}
			actualResponse := dummyXMLResponse{}
			builder := suite.
				httpRequestBuilder.
				NewRequest().
				WithContext(context.Background()).
				AddQueryParameters(queryParams).
				ResponseAs(&actualResponse)

			var err error
			if test.method == "POST" {
				err = builder.
					WithXMLBody(suite.XMLRequestBody).
					Post(suite.url)
			} else if test.method == "GET" {
				builder.AddHeader("Content-Type", "application/xml")
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithXMLBody(suite.XMLRequestBody).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldPassOauthCookiesWhenContextSetAndOauthEnabled() {
	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyResponse{ResponseFieldA: "sample response value"}
			responseJSONString, _ := json.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: "sampleCookie", Value: "sampleCookieValue"}
			ctx := context.WithValue(context.Background(), constants.CONTEXT_ACCESS_TOKEN, "accessTokenValue")
			ctx = context.WithValue(ctx, constants.CONTEXT_ENC_ID_TOKEN, "encIdTokenValue")
			ctx = context.WithValue(ctx, "Session-Tracing-ID", "QWERTY-1234")
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()

				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "application/json") &&
					actualRequest.Header.Get("Session-Tracing-ID") == "QWERTY-1234" &&
					actualRequest.Header.Get(constants.AUTHORIZATION_HEADER_KEY) == "Bearer accessTokenValue" &&
					actualRequest.Header.Get(constants.ENC_ID_TOKEN_HEADER_KEY) == "encIdTokenValue" &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			actualResponse := dummyResponse{}
			builder := suite.httpRequestBuilder.
				NewRequest().
				ResponseAs(&actualResponse).
				AddCookie(cookie).
				WithContext(ctx).
				WithOauth()

			var err error
			if test.method == "POST" {
				err = builder.
					WithJSONBody(suite.requestBody).
					Post(suite.url)
			} else if test.method == "GET" {
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithJSONBody(suite.requestBody).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldPassOauthCookiesFromGinContextWhenContextSetAndOauthEnabled() {
	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
			ginContext.Request = httptest.NewRequest("", "/", strings.NewReader(""))
			ginContext.Request.Header.Add("Session-Tracing-ID", "QWERTY-1234")
			expectedResponse := dummyResponse{ResponseFieldA: "sample response value"}
			responseJSONString, _ := json.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: "sampleCookie", Value: "sampleCookieValue"}
			ginContext.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer accessTokenValue")
			ginContext.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, "encIdTokenValue")
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()

				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "application/json") &&
					actualRequest.Header.Get("Session-Tracing-ID") == "QWERTY-1234" &&
					actualRequest.Header.Get(constants.AUTHORIZATION_HEADER_KEY) == "Bearer accessTokenValue" &&
					actualRequest.Header.Get(constants.ENC_ID_TOKEN_HEADER_KEY) == "encIdTokenValue" &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			actualResponse := dummyResponse{}
			builder := suite.httpRequestBuilder.
				NewRequest().
				ResponseAs(&actualResponse).
				AddCookie(cookie).
				WithContext(ginContext).
				WithOauth()

			var err error
			if test.method == "POST" {
				err = builder.
					WithJSONBody(suite.requestBody).
					Post(suite.url)
			} else if test.method == "GET" {
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithJSONBody(suite.requestBody).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldNotPassOauthCookiesWhenContextSetAndOauthDisabled() {

	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			expectedResponse := dummyResponse{ResponseFieldA: "sample response value"}
			responseJSONString, _ := json.Marshal(expectedResponse)

			cookie := &http.Cookie{Name: "sampleCookie", Value: "sampleCookieValue"}
			ctx := context.WithValue(context.Background(), constants.CONTEXT_ACCESS_TOKEN, "access token")
			ctx = context.WithValue(ctx, constants.CONTEXT_ENC_ID_TOKEN, "enc id-token")
			response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				actualCookies := actualRequest.Cookies()

				if actualRequest.URL.String() == suite.url &&
					actualRequest.Method == test.method &&
					(actualRequest.Method == "GET" || actualRequest.Header.Get("Content-Type") == "application/json") &&
					len(actualCookies) == 1 &&
					actualCookies[0].Name == cookie.Name && actualCookies[0].Value == cookie.Value {
					return &response, nil
				}
				return &http.Response{}, errors.New("request not matching")
			})

			actualResponse := dummyResponse{}
			builder := suite.
				httpRequestBuilder.
				NewRequest().
				ResponseAs(&actualResponse).
				AddCookie(cookie).
				WithContext(ctx)

			var err error
			if test.method == "POST" {
				err = builder.
					WithJSONBody(suite.requestBody).
					Post(suite.url)
			} else if test.method == "GET" {
				err = builder.Get(suite.url)
			} else if test.method == "PUT" {
				err = builder.
					WithJSONBody(suite.requestBody).
					Put(suite.url)
			}

			suite.Nil(err)
			suite.Equal(expectedResponse, actualResponse)
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfContextNoteSetButOauthEnabled() {
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).Times(0)

	actualResponse := dummyResponse{}
	err := suite.
		httpRequestBuilder.
		NewRequest().
		ResponseAs(&actualResponse).
		WithOauth().
		Get(suite.url)

	suite.NotNil(err)
	suite.Equal(err, errors.New("Context not set for forwarding oauth headers"))
}

func (suite HttpRequestTestSuite) TestShouldMakePostRequestAndGetResponseHeaders() {
	response := http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(""))),
		Header: map[string][]string{
			"testHeader1": {"testHeader1Val1"},
			"testHeader2": {"testHeader2Val1", "testHeader2Val1"},
		},
	}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	var actualResponseHeaders map[string][]string
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseHeadersAs(&actualResponseHeaders).
		Post(suite.url)

	suite.NotNil(actualResponseHeaders)
	suite.Equal(2, len(actualResponseHeaders))
	suite.Equal([]string{"testHeader1Val1"}, actualResponseHeaders["testHeader1"])
	suite.Equal([]string{"testHeader2Val1", "testHeader2Val1"}, actualResponseHeaders["testHeader2"])
	suite.Nil(err)
}

func (suite HttpRequestTestSuite) TestShouldMakePostRequestAndUnmarshalToString() {
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := ""
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.Equal(string(expectedJSONResponse), actualResponse)
	suite.Nil(err)
}

func (suite HttpRequestTestSuite) TestShouldMakePostRequestAndUnmarshalToByteArray() {
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	var actualResponse []byte
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.Equal(expectedJSONResponse, actualResponse)
	suite.Nil(err)
}

func (suite HttpRequestTestSuite) TestShouldMakePostRequestAndUnmarshalTo2dByteArray() {
	responseModel := `--34b21
	Content-Type: application/json
	Content-Disposition: form-data; name="text"

	{
		"result": "successs"
	}
	--34b21
	Content-Type: application/octet-stream
	Content-Disposition: form-data; filename="sample data"

	Sample response

	--34b21--`
	serverResponse, _ := json.Marshal(responseModel)

	response := http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"multipart/form-data; charset=UTF-8; boundary=\"34b21\""},
		},
		Body:          ioutil.NopCloser(bytes.NewBuffer(serverResponse)),
		ContentLength: int64(len(serverResponse)),
	}

	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	var actualResponse [][]byte
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.Nil(err)
}

func (suite HttpRequestTestSuite) TestShouldMakeDeleteRequestAndUnmarshalToByteArray() {
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	var actualResponse []byte
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		ResponseAs(&actualResponse).
		Delete(suite.url)

	suite.Equal(expectedJSONResponse, actualResponse)
	suite.Nil(err)
}

func (suite HttpRequestTestSuite) TestShouldMakePutRequestAndUnmarshalToString() {
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := ""
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Put(suite.url)

	suite.Equal(string(expectedJSONResponse), actualResponse)
	suite.Nil(err)
}

func (suite HttpRequestTestSuite) TestShouldMakePutRequestAndUnmarshalToByteArray() {
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	var actualResponse []byte
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Put(suite.url)

	suite.Equal(expectedJSONResponse, actualResponse)
	suite.Nil(err)
}

func (suite HttpRequestTestSuite) TestShouldMakeGetRequestAndRetrieveCookies() {
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	responseHeaders := http.Header{}
	responseHeaders.Add("Set-Cookie", "CookieA=CookieA_Value")
	response := http.Response{StatusCode: http.StatusOK, Header: responseHeaders, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	var actualCookies = &[]*http.Cookie{}
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		ResponseCookiesAs(actualCookies).
		Get(suite.url)

	suite.Equal("CookieA", (*actualCookies)[0].Name)
	suite.Equal("CookieA_Value", (*actualCookies)[0].Value)
	suite.Nil(err)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfMarshalingRequestBodyFails() {
	suite.requestBody.FieldB = math.Inf(1)
	actualResponse := dummyResponse{}
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.NotNil(err)
	suite.Equal(dummyResponse{}, actualResponse)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfUnmarshalResponseFails() {
	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer([]byte("invalid JSON")))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := dummyResponse{}
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.NotNil(err)
	suite.Equal(dummyResponse{}, actualResponse)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorForNon200Response() {
	expectedErrorResponse := dummyErrorResponse{
		ErrorCode:    619,
		ErrorMessage: "SOME ERROR MESSAGE",
	}

	errorResponseJSONString, _ := json.Marshal(expectedErrorResponse)

	response := http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer(errorResponseJSONString))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := dummyResponse{}
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.NotNil(err)
	_, ok := err.(utilError.HttpError)
	suite.Equal(true, ok)
	suite.Equal(dummyResponse{}, actualResponse)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorForNon200UnParsableResponse() {
	errorResponseJSONString, _ := json.Marshal("{")

	response := http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer(errorResponseJSONString))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := dummyResponse{}
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		Post(suite.url)

	suite.NotNil(err)
	_, ok := err.(utilError.HttpError)
	suite.Equal(true, ok)
	suite.Equal(`StatusCode : 500, ResponseBody : "{"`, err.Error())
	suite.Equal(dummyResponse{}, actualResponse)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfHttpRequestFails() {
	expectedError := errors.New("error connecting")
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusUnauthorized}, expectedError
	})

	var responseStatusCode int
	actualResponse := dummyResponse{}
	actualError := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseStatusCodeAs(&responseStatusCode).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.Equal(expectedError, actualError)
	suite.Equal(http.StatusUnauthorized, responseStatusCode)
	suite.Equal(dummyResponse{}, actualResponse)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfCreatingHttpRequestFails() {
	actualResponse := dummyResponse{}
	actualError := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(string([]byte{0x7f}))

	suite.NotNil(actualError)
	suite.Equal(dummyResponse{}, actualResponse)
}

type errReader int

func (errReader) Close() (err error) {
	return nil
}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("could not ready body")
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfReadingResponseBodyFails() {
	response := http.Response{StatusCode: http.StatusOK, Body: errReader(0)}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := dummyResponse{}
	actualError := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.NotNil(actualError)
	suite.Equal(dummyResponse{}, actualResponse)
}

type errCloser struct {
	responseBytes []byte
	reader        io.Reader
}

func (errCloser) Close() (err error) {
	return errors.New("could not close body")
}

func (ec errCloser) Read(p []byte) (n int, err error) {
	return ec.reader.Read(p)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfClosingResponseBodyFails() {
	expectedResponse := dummyResponse{ResponseFieldA: "sample response value"}
	responseJSONString, _ := json.Marshal(expectedResponse)

	response := http.Response{StatusCode: http.StatusOK, Body: errCloser{reader: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := dummyResponse{}
	actualError := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.NotNil(actualError)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfValidatingResponseModelFails() {
	expectedResponse := dummyResponse{ResponseFieldA: ""}
	responseJSONString, _ := json.Marshal(expectedResponse)

	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := dummyResponse{}
	actualError := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.NotNil(actualError)
}

func (suite HttpRequestTestSuite) TestShouldReturnErrorIfCustomValidatingResponseModelFails() {
	expectedResponse := dummyResponseWithCustomValidation{ResponseFieldA: "Non Boolean value"}
	responseJSONString, _ := json.Marshal(expectedResponse)

	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(responseJSONString))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := dummyResponseWithCustomValidation{}
	validate := validator.New()
	_ = validate.RegisterValidation("isBool", dummyValidationIsBool)

	actualError := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		WithCustomValidator(validate).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.NotNil(actualError)
}

//Marshall/unmarshall XML
//add get

func (suite HttpRequestTestSuite) TestShouldStartTracingWhenContextSet() {
	defer suite.mockCtrl.Finish()
	var testCases = []struct {
		method string
	}{
		{method: "POST"},
		{method: "GET"},
	}

	for _, test := range testCases {
		suite.Run(test.method, func() {
			type contextKey struct{}
			ctx := context.WithValue(context.Background(), contextKey{}, &trace.Span{})
			suite.trace.EXPECT().Continue(ctx, gomock.Any()).MaxTimes(2).DoAndReturn(func(ctx context.Context, h *http.Request) (*trace.Span, *http.Request) {
				return &trace.Span{}, h
			})
			responseModel := dummyResponse{ResponseFieldA: "sample response value"}
			expectedJSONResponse, _ := json.Marshal(responseModel)

			responseHeaders := http.Header{}
			responseHeaders.Add("Set-Cookie", "CookieA=CookieA_Value")
			response := http.Response{StatusCode: http.StatusOK, Header: responseHeaders, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
			suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
				return &response, nil
			})

			_ = suite.
				httpRequestBuilder.
				NewRequest().
				WithContext(ctx).
				WithTracer(suite.trace).
				Get("http://another-service.com/")
		})
	}
}

func (suite HttpRequestTestSuite) TestShouldNotUpdateRequestWithMessageIdWhenTraceIdIsNotPresent() {
	defer suite.mockCtrl.Finish()

	type contextKey struct{}
	ctx := context.WithValue(context.Background(), contextKey{}, &trace.Span{})
	ctx = context.WithValue(ctx, constants.TRACE_KEY, constants.NO_TRACE_ID)

	suite.trace.EXPECT().Continue(ctx, gomock.Any()).MaxTimes(2).DoAndReturn(func(ctx context.Context, h *http.Request) (*trace.Span, *http.Request) {
		return &trace.Span{}, h
	})
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	responseHeaders := http.Header{}
	responseHeaders.Add("Set-Cookie", "CookieA=CookieA_Value")
	responseHeaders.Add("Content-Type", "CookieA=CookieA_Value")
	response := http.Response{StatusCode: http.StatusOK, Header: responseHeaders, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		requestBytes, _ := ioutil.ReadAll(actualRequest.Body)
		request := cotsRequest{}
		json.Unmarshal(requestBytes, &request)
		expectedRequest := cotsRequest{
			Wrapper{
				MessageHeader: MessageHeader{
					HeaderField: HeaderField{
						ApplicationId: "gola",
						MessageId:     "message-id",
					},
				},
			},
		}

		suite.Equal(expectedRequest, request)
		return &response, nil
	})

	request := cotsRequest{
		Wrapper{
			MessageHeader: MessageHeader{
				HeaderField: HeaderField{
					ApplicationId: "gola",
					MessageId:     "message-id",
				},
			},
		},
	}

	_ = suite.
		httpRequestBuilder.
		NewRequest().
		WithJSONBody(request).
		WithContext(ctx).
		WithTracer(suite.trace).
		Get("http://another-service.com/")

}

func (suite HttpRequestTestSuite) TestShouldAddHTTPStatusCodeToDataSpan() {
	e := SpanExporter{}
	trace.RegisterExporter(&e)
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	ctx, s := trace.StartSpan(context.Background(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := ""
	_ = suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(ctx).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	s.End()
	suite.Equal(e.spans[0].Attributes[ochttp.StatusCodeAttribute], int64(200))
	suite.Equal(e.spans[0].Attributes[ErrorTraceAttribute], nil)
}

func (suite HttpRequestTestSuite) TestShouldAddHTTPStatusCodeAndErrorToDataSpan() {
	e := SpanExporter{}
	trace.RegisterExporter(&e)
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	ctx, s := trace.StartSpan(context.Background(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	response := http.Response{StatusCode: http.StatusNotFound, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := ""
	_ = suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(ctx).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	s.End()
	suite.Equal(e.spans[0].Attributes[ochttp.StatusCodeAttribute], int64(404))
	suite.Equal(e.spans[0].Attributes[ErrorTraceAttribute], true)
}

func (suite HttpRequestTestSuite) TestShouldGiveSuccessfulResponseFor3XXStatusCode() {
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	actualResponse := ""
	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(context.Background()).
		WithJSONBody(suite.requestBody).
		ResponseAs(&actualResponse).
		Post(suite.url)

	suite.Nil(err)
	suite.Equal(string(expectedJSONResponse), actualResponse)
}

func (suite HttpRequestTestSuite) TestShouldCloseSpanWhenResponseModelIsNotPassed() {
	e := SpanExporter{}
	trace.RegisterExporter(&e)
	responseModel := dummyResponse{ResponseFieldA: "sample response value"}
	expectedJSONResponse, _ := json.Marshal(responseModel)

	ctx, s := trace.StartSpan(context.Background(), "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	response := http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(expectedJSONResponse))}
	suite.mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(actualRequest *http.Request) (*http.Response, error) {
		return &response, nil
	})

	err := suite.
		httpRequestBuilder.
		NewRequest().
		WithContext(ctx).
		WithJSONBody(suite.requestBody).
		Post(suite.url)

	s.End()
	suite.Nil(err)
	suite.Equal(e.spans[0].Attributes[ochttp.StatusCodeAttribute], int64(200))
	suite.Equal(len(e.spans[0].Annotations), 2)
}

type cotsRequest struct {
	Wrapper `json:"random_key"`
}

type Wrapper struct {
	MessageHeader MessageHeader `json:"msgHdr"`
}

type MessageHeader struct {
	HeaderField HeaderField `json:"hdrFlds"`
}

type HeaderField struct {
	ConversationID string `json:"cnvId"`
	MessageId      string `json:"msgId"`
	ApplicationId  string `json:"appId"`
	TimeStamp      string `json:"timestamp"`
}
