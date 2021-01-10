package request

//go:generate mockgen -source=http_request.go -destination=./mocks/mock_http_request.go -package=mocks

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gola-glitch/gola-utils/http/util"
	"github.com/jtacoma/uritemplates"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gola-glitch/gola-utils/http/model"
	"github.com/gola-glitch/gola-utils/logging"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/constants"
	utilError "github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/http/client"
	"github.com/gola-glitch/gola-utils/trace"
	"go.opencensus.io/plugin/ochttp"
	openTrace "go.opencensus.io/trace"
)

const (
	ErrorTraceAttribute = "error"
)

type TraceHookFunc func([]byte) (string, error)

type HttpRequest interface {
	WithJSONBody(interface{}) HttpRequest
	WithJSONBodyNoEscapeHTML(interface{}) HttpRequest
	WithXMLBody(interface{}) HttpRequest
	WithXMLBodyTextHeader(interface{}) HttpRequest
	WithFormURLEncoded(map[string]interface{}) HttpRequest
	WithContext(context.Context) HttpRequest
	WithOauth() HttpRequest
	WithRequestBodyBytes([]byte) HttpRequest
	WithTracer(trace.Trace) HttpRequest
	WithCustomValidator(*validator.Validate) HttpRequest
	RequestTraceHook(hookFunc TraceHookFunc) HttpRequest
	ResponseTraceHook(hookFunc TraceHookFunc) HttpRequest
	ResponseAs(interface{}) HttpRequest
	ResponseStatusCodeAs(*int) HttpRequest
	ResponseHeadersAs(*map[string][]string) HttpRequest
	ResponseCookiesAs(*[]*http.Cookie) HttpRequest
	AddHeader(string, string) HttpRequest
	AddHeaders(map[string]string) HttpRequest
	AddQueryParameters(map[string]string) HttpRequest
	AddPathParameters(map[string]interface{}) HttpRequest
	AddCookie(*http.Cookie) HttpRequest
	Post(string) error
	Put(string) error
	Get(string) error
	Delete(string) error
}

type httpRequest struct {
	responseStatusCode *int
	responseModel      interface{}
	responseCookies    *[]*http.Cookie
	requestModel       interface{}
	headers            map[string]string
	cookies            []*http.Cookie
	httpClient         client.HttpClient
	ctx                context.Context
	forwardAuthHeaders bool
	validate           *validator.Validate
	requestBytes       []byte
	requestBuildError  error
	queryParameters    map[string]string
	pathParameters     map[string]interface{}
	pathTemplate       string
	responseHeaders    *map[string][]string
	enableTracing      bool
	trace              trace.Trace
	requestTraceHook   TraceHookFunc
	responseTraceHook  TraceHookFunc
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func (r httpRequest) WithCustomValidator(validator *validator.Validate) HttpRequest {
	r.validate = validator
	return r
}

func (r httpRequest) RequestTraceHook(hookFunc TraceHookFunc) HttpRequest {
	r.requestTraceHook = hookFunc
	return r
}

func (r httpRequest) ResponseTraceHook(hookFunc TraceHookFunc) HttpRequest {
	r.responseTraceHook = hookFunc
	return r
}

func (r httpRequest) ResponseAs(responseModel interface{}) HttpRequest {
	r.responseModel = responseModel
	return r
}

func (r httpRequest) ResponseStatusCodeAs(responseStatusCode *int) HttpRequest {
	r.responseStatusCode = responseStatusCode
	return r
}

func (r httpRequest) ResponseCookiesAs(responseCookies *[]*http.Cookie) HttpRequest {
	r.responseCookies = responseCookies
	return r
}

func (r httpRequest) ResponseHeadersAs(responseHeaders *map[string][]string) HttpRequest {
	r.responseHeaders = responseHeaders
	return r
}

func (r httpRequest) WithJSONBody(requestModel interface{}) HttpRequest {
	r.requestModel = requestModel
	r.headers["Content-Type"] = "application/json"
	requestBytes, err := json.Marshal(r.requestModel)

	if err != nil {
		r.requestBuildError = err
	}

	r.requestBytes = requestBytes
	return r
}

func (r httpRequest) WithTracer(t trace.Trace) HttpRequest {
	r.trace = t
	return r
}

func (r httpRequest) WithJSONBodyNoEscapeHTML(requestModel interface{}) HttpRequest {
	r.requestModel = requestModel
	r.headers["Content-Type"] = "application/json"

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(r.requestModel)

	r.requestBytes = buffer.Bytes()
	return r
}

func (r httpRequest) WithXMLBody(requestModel interface{}) HttpRequest {
	r.requestModel = requestModel
	r.headers["Content-Type"] = "application/xml"
	requestBytes, err := xml.Marshal(r.requestModel)

	if err != nil {
		r.requestBuildError = err
	}

	r.requestBytes = requestBytes
	return r
}

func (r httpRequest) WithXMLBodyTextHeader(requestModel interface{}) HttpRequest {
	r.requestModel = requestModel
	r.headers["Content-Type"] = "text/xml"
	requestBytes, err := xml.Marshal(r.requestModel)

	if err != nil {
		r.requestBuildError = err
	}

	r.requestBytes = requestBytes
	return r
}

func (r httpRequest) WithFormURLEncoded(formData map[string]interface{}) HttpRequest {
	r.requestModel = formData
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for key, value := range formData {
		switch value.(type) {
		case multipart.FileHeader:
			fileHeader := value.(multipart.FileHeader)
			file, err := fileHeader.Open()
			if err != nil {
				r.requestBuildError = err
			}

			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(key), escapeQuotes(fileHeader.Filename)))
			h.Set("Content-Type", fileHeader.Header.Get(constants.HeaderContentType))
			formFileWriter, err := bodyWriter.CreatePart(h)
			if err != nil {
				r.requestBuildError = err
			}
			if err != nil {
				r.requestBuildError = err
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				r.requestBuildError = err
			}
			if _, fileWriteError := formFileWriter.Write(content); fileWriteError != nil {
				r.requestBuildError = fileWriteError
			}
		case string:
			if fieldWriteError := bodyWriter.WriteField(key, value.(string)); fieldWriteError != nil {
				r.requestBuildError = fieldWriteError
			}
		case []byte:
			data := value.([]byte)
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(key)))
			h.Set("Content-Type", "application/json")
			byteWriter, err := bodyWriter.CreatePart(h)
			if err != nil {
				r.requestBuildError = err
			}
			if _, err = byteWriter.Write(data); err != nil {
				r.requestBuildError = err
			}
		case model.FileUploadContent:
			fileContent := value.(model.FileUploadContent)
			file, err := os.Open(fileContent.FilePath)
			if err != nil {
				logging.GetLogger(r.ctx).Error(err)
			}
			formFileWriter, err := bodyWriter.CreateFormFile(key, fileContent.FileName)
			if err != nil {
				r.requestBuildError = err
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				r.requestBuildError = err
			}
			if _, fileWriteError := formFileWriter.Write(content); fileWriteError != nil {
				r.requestBuildError = fileWriteError
			}
			err = file.Close()
			if err != nil {
				r.requestBuildError = err
			}
		default:
			r.requestBuildError = utilError.Error{
				ErrorCode:      "ERR_INVALID_REQUEST_TYPE",
				ErrorMessage:   "only multipart files and strings are supported",
				AdditionalData: nil,
			}
		}
	}

	if multipartWriterCloseError := bodyWriter.Close(); multipartWriterCloseError != nil {
		r.requestBuildError = multipartWriterCloseError
	}

	r.headers["Content-Type"] = bodyWriter.FormDataContentType()

	r.requestBytes = bodyBuf.Bytes()

	return r
}

func (r httpRequest) WithContext(ctx context.Context) HttpRequest {
	r.ctx = ctx
	return r
}

func (r httpRequest) WithOauth() HttpRequest {
	r.forwardAuthHeaders = true
	return r
}

func (r httpRequest) AddHeader(key, value string) HttpRequest {
	r.headers[key] = value
	return r
}

func (r httpRequest) AddHeaders(headers map[string]string) HttpRequest {
	for key, value := range headers {
		r.headers[key] = value
	}
	return r
}

func (r httpRequest) AddCookie(cookie *http.Cookie) HttpRequest {
	r.cookies = append(r.cookies, cookie)
	return r
}

func (r httpRequest) AddQueryParameters(queryParameters map[string]string) HttpRequest {
	r.queryParameters = queryParameters
	return r
}

func (r httpRequest) AddPathParameters(pathParameters map[string]interface{}) HttpRequest {
	r.pathParameters = pathParameters
	return r
}

func (r httpRequest) Get(url string) error {
	r.pathTemplate = url
	return r.makeRequest("GET")
}

func (r httpRequest) Post(url string) error {
	r.pathTemplate = url
	return r.makeRequest("POST")
}

func (r httpRequest) Put(url string) error {
	r.pathTemplate = url
	return r.makeRequest("PUT")
}

func (r httpRequest) Delete(url string) error {
	r.pathTemplate = url
	return r.makeRequest("DELETE")
}

func (r httpRequest) makeRequest(method string) error {
	if r.requestBuildError != nil {
		return r.requestBuildError
	}

	cleanUrl, urlFormatErr := r.removeDoubleSlashesInUrl(r.pathTemplate)
	if urlFormatErr != nil {
		return urlFormatErr
	}
	r.pathTemplate = cleanUrl

	urlWithPathParams, urlConstructionErr := r.getUrlWithPathParams()
	if urlConstructionErr != nil {
		return urlConstructionErr
	}

	httpRequest, requestError := http.NewRequest(method, urlWithPathParams, bytes.NewBuffer(r.requestBytes))

	if requestError != nil {
		return requestError
	}

	if r.ctx != nil {
		httpRequest = httpRequest.WithContext(r.ctx)
		addSessionTracingId(r.ctx, httpRequest)
	}
	if r.forwardAuthHeaders {
		if r.ctx == nil {
			return errors.New(fmt.Sprintf("Context not set for forwarding oauth headers"))
		}
		addAuthHeaders(r.ctx, httpRequest)
	}

	for k, v := range r.headers {
		httpRequest.Header.Add(k, v)
	}

	query := httpRequest.URL.Query()
	for paramKey, paramValue := range r.queryParameters {
		query.Add(paramKey, paramValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	for _, cookie := range r.cookies {
		httpRequest.AddCookie(cookie)
	}

	var dataSpan *openTrace.Span = nil
	currentContext := r.ctx
	if currentContext != nil {
		span, h := r.trace.Continue(currentContext, httpRequest)
		httpRequest = h
		span.Annotate([]openTrace.Attribute{openTrace.StringAttribute("isEndingSpan", "false")}, "log")
		dataSpan = r.logHttpRequest(span, httpRequest)
		// defer span.End() Fixes span adjustment
	}
	start := time.Now()

	response, httpError := r.httpClient.Do(httpRequest)

	_ = time.Since(start)

	_ = getStatusCode(response)

	if response != nil && r.responseStatusCode != nil {
		*r.responseStatusCode = response.StatusCode
	}

	if response != nil && r.responseHeaders != nil {
		*r.responseHeaders = response.Header
	}

	if httpError != nil {
		r.logHttpResponse("Some error occurred: "+httpError.Error(), httpRequest, dataSpan)
		return httpError
	}

	addResponseTags(response, dataSpan)
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		errorResponseBytes, readError := ioutil.ReadAll(response.Body)
		if readError != nil {
			r.logHttpResponse("Response body read Error: "+readError.Error(), httpRequest, dataSpan)
			return utilError.HttpError{
				StatusCode:   response.StatusCode,
				ResponseBody: []byte(readError.Error()),
			}
		}
		r.logHttpResponse(string(errorResponseBytes), httpRequest, dataSpan)
		return utilError.HttpError{
			StatusCode:   response.StatusCode,
			ResponseBody: errorResponseBytes,
		}
	}

	if r.responseModel != nil {
		err := r.processResponseModel(response, httpRequest, dataSpan)
		if err != nil {
			return err
		}
	}

	if r.responseModel == nil {
		r.logHttpResponse("", httpRequest, dataSpan)
	}

	if r.responseCookies != nil {
		*r.responseCookies = response.Cookies()
	}

	closeError := response.Body.Close()
	if closeError != nil {
		return closeError
	}

	return nil
}

func (r httpRequest) getUrlWithPathParams() (string, error) {
	urlTemplate, parseErr := uritemplates.Parse(r.pathTemplate)
	if parseErr != nil {
		logging.NewLoggerEntry().Error("Invalid path template. Unable to parse : ", parseErr)
		return "", errors.New(fmt.Sprintf("Invalid path template"))
	}

	urlWithPathParams, expandErr := urlTemplate.Expand(r.pathParameters)
	if expandErr != nil {
		logging.NewLoggerEntry().Error("Unable to populate path params..", expandErr)
		return "", errors.New(fmt.Sprintf("URL construction failed for path params"))
	}

	return urlWithPathParams, nil
}

func (r httpRequest) removeDoubleSlashesInUrl(url string) (string, error) {
	var scheme string
	if strings.HasPrefix(url, "https://") {
		scheme = "https://"
	} else if strings.HasPrefix(url, "http://") {
		scheme = "http://"
	} else {
		return "", errors.New(fmt.Sprintf("Url scheme missing"))
	}

	urlPath := strings.TrimPrefix(url, scheme)
	cleanUrlPath := strings.ReplaceAll(urlPath, "//", "/")

	return scheme + cleanUrlPath, nil
}

func getStatusCode(response *http.Response) string {
	statusCode := ""
	if response == nil {
		statusCode = "500"
	} else {
		statusCode = strconv.Itoa(response.StatusCode)
	}
	return statusCode
}

func addResponseTags(res *http.Response, span *openTrace.Span) {
	span.AddAttributes(openTrace.Int64Attribute(ochttp.StatusCodeAttribute, int64(res.StatusCode)))
	if res.StatusCode >= 400 {
		span.AddAttributes(openTrace.BoolAttribute(ErrorTraceAttribute, true))
	}
}

func (r httpRequest) logHttpRequest(span *openTrace.Span, httpRequest *http.Request) *openTrace.Span {
	logKeyName := "request"
	logDescription := httpRequest.Host + httpRequest.URL.RequestURI()
	var logRequest string

	_, dataSpan := openTrace.StartSpanWithRemoteParent(r.ctx, r.extractPath()+" | request/response", span.SpanContext())

	requestBodyBytes, readErr := ioutil.ReadAll(httpRequest.Body)
	defer func() {
		httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
	}()

	if readErr != nil {
		logRequest = "Error reading request body: " + readErr.Error()
	} else {
		if r.requestTraceHook != nil {
			if requestBody, err := r.requestTraceHook(requestBodyBytes); err == nil {
				logRequest = requestBody
			} else {
				logRequest = string(requestBodyBytes)
			}
		} else {
			logRequest = string(requestBodyBytes)
		}
	}

	dataSpan.Annotate([]openTrace.Attribute{
		openTrace.StringAttribute(logKeyName, logRequest),
	}, logDescription)

	return dataSpan
}

func (r httpRequest) extractPath() string {
	parsedURL, err := url.Parse(r.pathTemplate)
	if err != nil {
		logging.GetLogger(r.ctx).Warn("Error while parsing path template")
		return r.pathTemplate
	}
	return parsedURL.Path
}

func isLoggingDisabled(requestedPath string) bool {
	return strings.Contains(requestedPath, "crypto")
}

func (r httpRequest) logHttpResponse(respBody string, httpRequest *http.Request, dataSpan *openTrace.Span) {
	//TODO: Remove this if condition (don't remove the lines inside if) once deprecated methods from cryptoUtil are removed
	if dataSpan != nil {
		logKeyName := "response"
		logDescription := httpRequest.Host + httpRequest.URL.RequestURI()
		responseBodyBytes := []byte(respBody)

		if r.responseTraceHook != nil {
			if responseBody, err := r.responseTraceHook(responseBodyBytes); err == nil {
				responseBodyBytes = []byte(responseBody)
			}
		}

		dataSpan.Annotate([]openTrace.Attribute{
			openTrace.StringAttribute(logKeyName, trace.EscapeSpecialChar(responseBodyBytes)),
		}, logDescription)
		dataSpan.End()
	}
}

func (r httpRequest) processResponseModel(response *http.Response, httpRequest *http.Request, dataSpan *openTrace.Span) error {
	if twoDimByteArrPtr, is2dByteArray := r.responseModel.(*[][]byte); is2dByteArray {
		data, err := r.processMultiFormResponse(response, httpRequest, dataSpan)
		if err != nil {
			r.logHttpResponse("Response body -MIME- read Error: "+err.Error(), httpRequest, dataSpan)
			return err
		}
		*twoDimByteArrPtr = data
	} else {
		responseBytes, readError := ioutil.ReadAll(response.Body)
		if readError != nil {
			r.logHttpResponse("Response body read Error: "+readError.Error(), httpRequest, dataSpan)
			return readError
		}

		if isLoggingDisabled(httpRequest.URL.Path) {
			r.logHttpResponse("Response body not logged for security reasons", httpRequest, dataSpan)
		} else {
			r.logHttpResponse(string(responseBytes), httpRequest, dataSpan)
		}

		if strPointer, isString := r.responseModel.(*string); isString {
			*strPointer = string(responseBytes)
		} else if byteArrayPointer, isByteArray := r.responseModel.(*[]byte); isByteArray {
			*byteArrayPointer = responseBytes
		} else {
			var unmarshalError error
			contentType := response.Header.Get(constants.HeaderContentType)
			if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
				unmarshalError = xml.Unmarshal(responseBytes, r.responseModel)
			} else {
				unmarshalError = json.Unmarshal(responseBytes, r.responseModel)
			}
			if unmarshalError != nil {
				return unmarshalError
			}

			validationError := r.validateResponse(r.responseModel)
			if validationError != nil {
				return validationError
			}
		}
	}
	return nil
}

func (r httpRequest) processMultiFormResponse(response *http.Response, httpRequest *http.Request, dataSpan *openTrace.Span) ([][]byte, error) {
	_, params, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		r.logHttpResponse("Error occurred in parsing response media type: "+err.Error(), httpRequest, dataSpan)
		err := utilError.HttpError{
			StatusCode: http.StatusInternalServerError,
		}
		return nil, err
	}

	data := make([][]byte, 0)
	reader := multipart.NewReader(response.Body, params["boundary"])
	for part, err := reader.NextPart(); err == nil; part, err = reader.NextPart() {
		buf, err := ioutil.ReadAll(part)
		if err != nil {
			r.logHttpResponse("unable to read part due to "+err.Error(), httpRequest, dataSpan)
			fmt.Println("unable to read part due to", err)
			return nil, err
		}
		data = append(data, buf)
	}
	r.logHttpResponse("MIME response received", httpRequest, dataSpan)
	return data, nil
}

func addAuthHeaders(ctx context.Context, httpRequest *http.Request) {
	if ginContext, ok := ctx.(*gin.Context); ok {
		// context is Gin context, adding Gin request oauth headers
		if accessToken, accessTokenErr := util.GetAccessToken(ginContext); accessTokenErr == nil {
			httpRequest.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer "+accessToken)
		}
		if encIDToken, encIDTokenErr := util.GetEncryptedIDToken(ginContext); encIDTokenErr == nil {
			httpRequest.Header.Add(constants.ENC_ID_TOKEN_HEADER_KEY, encIDToken)
		}
		return
	}

	accessToken := ctx.Value(constants.CONTEXT_ACCESS_TOKEN)
	if accessToken != nil {
		httpRequest.Header.Add(constants.AUTHORIZATION_HEADER_KEY, "Bearer "+accessToken.(string))
	}
	encIDToken := ctx.Value(constants.CONTEXT_ENC_ID_TOKEN)
	if encIDToken != nil {
		httpRequest.Header.Add(constants.ENC_ID_TOKEN_HEADER_KEY, encIDToken.(string))
	}
}

func addSessionTracingId(ctx context.Context, httpRequest *http.Request) {
	if ginContext, ok := ctx.(*gin.Context); ok {
		sessionTracingId := ginContext.Request.Header.Get(constants.TRACING_SESSION_HEADER_KEY)
		if sessionTracingId != "" {
			httpRequest.Header.Add(constants.TRACING_SESSION_HEADER_KEY, sessionTracingId)
		}
		return
	}
	sessionTracingId, isSessionTracingId := ctx.Value(constants.TRACING_SESSION_HEADER_KEY).(string)
	if isSessionTracingId {
		httpRequest.Header.Add(constants.TRACING_SESSION_HEADER_KEY, sessionTracingId)
	}
}

func (r httpRequest) validateResponse(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		return r.validate.Struct(obj)
	}
	return nil
}

func (r httpRequest) WithRequestBodyBytes(bytes []byte) HttpRequest {
	r.requestBytes = bytes
	return r
}

// Request builder
type HttpRequestBuilder interface {
	NewRequest() HttpRequest
	NewRequestWithContext(ctx context.Context) HttpRequest
}

type requestBuilder struct {
	httpClient client.HttpClient
}

func (rb requestBuilder) NewRequest() HttpRequest {
	return httpRequest{
		httpClient: rb.httpClient,
		validate:   validator.New(),
		headers: map[string]string{
			constants.X_REQUESTED_WITH_HEADER_KEY: constants.X_REQUESTED_WITH_HEADER_VALUE,
		},
		cookies:         []*http.Cookie{},
		trace:           trace.New(),
	}
}

func (rb requestBuilder) NewRequestWithContext(ctx context.Context) HttpRequest {
	return rb.NewRequest().WithContext(ctx)
}

func NewHttpRequestBuilder(client client.HttpClient) HttpRequestBuilder {
	return requestBuilder{
		httpClient: client,
	}
}
