// Code generated by MockGen. DO NOT EDIT.
// Source: http/request/http_request.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	request "github.com/gola-glitch/gola-utils/http/request"
	trace "github.com/gola-glitch/gola-utils/trace"
	gomock "github.com/golang/mock/gomock"
	validator_v9 "gopkg.in/go-playground/validator.v9"
	http "net/http"
	reflect "reflect"
)

// MockHttpRequest is a mock of HttpRequest interface
type MockHttpRequest struct {
	ctrl     *gomock.Controller
	recorder *MockHttpRequestMockRecorder
}

// MockHttpRequestMockRecorder is the mock recorder for MockHttpRequest
type MockHttpRequestMockRecorder struct {
	mock *MockHttpRequest
}

// NewMockHttpRequest creates a new mock instance
func NewMockHttpRequest(ctrl *gomock.Controller) *MockHttpRequest {
	mock := &MockHttpRequest{ctrl: ctrl}
	mock.recorder = &MockHttpRequestMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHttpRequest) EXPECT() *MockHttpRequestMockRecorder {
	return m.recorder
}

// WithJSONBody mocks base method
func (m *MockHttpRequest) WithJSONBody(arg0 interface{}) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithJSONBody", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithJSONBody indicates an expected call of WithJSONBody
func (mr *MockHttpRequestMockRecorder) WithJSONBody(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithJSONBody", reflect.TypeOf((*MockHttpRequest)(nil).WithJSONBody), arg0)
}

// WithJSONBodyNoEscapeHTML mocks base method
func (m *MockHttpRequest) WithJSONBodyNoEscapeHTML(arg0 interface{}) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithJSONBodyNoEscapeHTML", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithJSONBodyNoEscapeHTML indicates an expected call of WithJSONBodyNoEscapeHTML
func (mr *MockHttpRequestMockRecorder) WithJSONBodyNoEscapeHTML(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithJSONBodyNoEscapeHTML", reflect.TypeOf((*MockHttpRequest)(nil).WithJSONBodyNoEscapeHTML), arg0)
}

// WithXMLBody mocks base method
func (m *MockHttpRequest) WithXMLBody(arg0 interface{}) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithXMLBody", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithXMLBody indicates an expected call of WithXMLBody
func (mr *MockHttpRequestMockRecorder) WithXMLBody(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithXMLBody", reflect.TypeOf((*MockHttpRequest)(nil).WithXMLBody), arg0)
}

// WithXMLBodyTextHeader mocks base method
func (m *MockHttpRequest) WithXMLBodyTextHeader(arg0 interface{}) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithXMLBodyTextHeader", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithXMLBodyTextHeader indicates an expected call of WithXMLBodyTextHeader
func (mr *MockHttpRequestMockRecorder) WithXMLBodyTextHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithXMLBodyTextHeader", reflect.TypeOf((*MockHttpRequest)(nil).WithXMLBodyTextHeader), arg0)
}

// WithFormURLEncoded mocks base method
func (m *MockHttpRequest) WithFormURLEncoded(arg0 map[string]interface{}) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithFormURLEncoded", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithFormURLEncoded indicates an expected call of WithFormURLEncoded
func (mr *MockHttpRequestMockRecorder) WithFormURLEncoded(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithFormURLEncoded", reflect.TypeOf((*MockHttpRequest)(nil).WithFormURLEncoded), arg0)
}

// WithContext mocks base method
func (m *MockHttpRequest) WithContext(arg0 context.Context) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithContext", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithContext indicates an expected call of WithContext
func (mr *MockHttpRequestMockRecorder) WithContext(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithContext", reflect.TypeOf((*MockHttpRequest)(nil).WithContext), arg0)
}

// WithOauth mocks base method
func (m *MockHttpRequest) WithOauth() request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithOauth")
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithOauth indicates an expected call of WithOauth
func (mr *MockHttpRequestMockRecorder) WithOauth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithOauth", reflect.TypeOf((*MockHttpRequest)(nil).WithOauth))
}

// WithTracer mocks base method
func (m *MockHttpRequest) WithTracer(arg0 trace.Trace) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTracer", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithTracer indicates an expected call of WithTracer
func (mr *MockHttpRequestMockRecorder) WithTracer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTracer", reflect.TypeOf((*MockHttpRequest)(nil).WithTracer), arg0)
}

// WithCustomValidator mocks base method
func (m *MockHttpRequest) WithCustomValidator(arg0 *validator_v9.Validate) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithCustomValidator", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// WithCustomValidator indicates an expected call of WithCustomValidator
func (mr *MockHttpRequestMockRecorder) WithCustomValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithCustomValidator", reflect.TypeOf((*MockHttpRequest)(nil).WithCustomValidator), arg0)
}

// ResponseAs mocks base method
func (m *MockHttpRequest) ResponseAs(arg0 interface{}) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResponseAs", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// ResponseAs indicates an expected call of ResponseAs
func (mr *MockHttpRequestMockRecorder) ResponseAs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseAs", reflect.TypeOf((*MockHttpRequest)(nil).ResponseAs), arg0)
}

// ResponseStatusCodeAs mocks base method
func (m *MockHttpRequest) ResponseStatusCodeAs(arg0 *int) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResponseStatusCodeAs", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// ResponseStatusCodeAs indicates an expected call of ResponseStatusCodeAs
func (mr *MockHttpRequestMockRecorder) ResponseStatusCodeAs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseStatusCodeAs", reflect.TypeOf((*MockHttpRequest)(nil).ResponseStatusCodeAs), arg0)
}

// ResponseHeadersAs mocks base method
func (m *MockHttpRequest) ResponseHeadersAs(arg0 *map[string][]string) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResponseHeadersAs", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// ResponseHeadersAs indicates an expected call of ResponseHeadersAs
func (mr *MockHttpRequestMockRecorder) ResponseHeadersAs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseHeadersAs", reflect.TypeOf((*MockHttpRequest)(nil).ResponseHeadersAs), arg0)
}

// ResponseCookiesAs mocks base method
func (m *MockHttpRequest) ResponseCookiesAs(arg0 *[]*http.Cookie) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResponseCookiesAs", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// ResponseCookiesAs indicates an expected call of ResponseCookiesAs
func (mr *MockHttpRequestMockRecorder) ResponseCookiesAs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseCookiesAs", reflect.TypeOf((*MockHttpRequest)(nil).ResponseCookiesAs), arg0)
}

// AddHeader mocks base method
func (m *MockHttpRequest) AddHeader(arg0, arg1 string) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddHeader", arg0, arg1)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// AddHeader indicates an expected call of AddHeader
func (mr *MockHttpRequestMockRecorder) AddHeader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHeader", reflect.TypeOf((*MockHttpRequest)(nil).AddHeader), arg0, arg1)
}

// AddHeaders mocks base method
func (m *MockHttpRequest) AddHeaders(arg0 map[string]string) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddHeaders", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// AddHeaders indicates an expected call of AddHeaders
func (mr *MockHttpRequestMockRecorder) AddHeaders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHeaders", reflect.TypeOf((*MockHttpRequest)(nil).AddHeaders), arg0)
}

// AddQueryParameters mocks base method
func (m *MockHttpRequest) AddQueryParameters(arg0 map[string]string) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddQueryParameters", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// AddQueryParameters indicates an expected call of AddQueryParameters
func (mr *MockHttpRequestMockRecorder) AddQueryParameters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddQueryParameters", reflect.TypeOf((*MockHttpRequest)(nil).AddQueryParameters), arg0)
}

// AddCookie mocks base method
func (m *MockHttpRequest) AddCookie(arg0 *http.Cookie) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCookie", arg0)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// AddCookie indicates an expected call of AddCookie
func (mr *MockHttpRequestMockRecorder) AddCookie(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCookie", reflect.TypeOf((*MockHttpRequest)(nil).AddCookie), arg0)
}

// Post mocks base method
func (m *MockHttpRequest) Post(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Post indicates an expected call of Post
func (mr *MockHttpRequestMockRecorder) Post(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockHttpRequest)(nil).Post), arg0)
}

// Put mocks base method
func (m *MockHttpRequest) Put(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockHttpRequestMockRecorder) Put(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockHttpRequest)(nil).Put), arg0)
}

// Get mocks base method
func (m *MockHttpRequest) Get(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockHttpRequestMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockHttpRequest)(nil).Get), arg0)
}

// Delete mocks base method
func (m *MockHttpRequest) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockHttpRequestMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockHttpRequest)(nil).Delete), arg0)
}

// MockHttpRequestBuilder is a mock of HttpRequestBuilder interface
type MockHttpRequestBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockHttpRequestBuilderMockRecorder
}

// MockHttpRequestBuilderMockRecorder is the mock recorder for MockHttpRequestBuilder
type MockHttpRequestBuilderMockRecorder struct {
	mock *MockHttpRequestBuilder
}

// NewMockHttpRequestBuilder creates a new mock instance
func NewMockHttpRequestBuilder(ctrl *gomock.Controller) *MockHttpRequestBuilder {
	mock := &MockHttpRequestBuilder{ctrl: ctrl}
	mock.recorder = &MockHttpRequestBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHttpRequestBuilder) EXPECT() *MockHttpRequestBuilderMockRecorder {
	return m.recorder
}

// NewRequest mocks base method
func (m *MockHttpRequestBuilder) NewRequest() request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRequest")
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// NewRequest indicates an expected call of NewRequest
func (mr *MockHttpRequestBuilderMockRecorder) NewRequest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRequest", reflect.TypeOf((*MockHttpRequestBuilder)(nil).NewRequest))
}

// NewRequestWithContext mocks base method
func (m *MockHttpRequestBuilder) NewRequestWithContext(ctx context.Context) request.HttpRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRequestWithContext", ctx)
	ret0, _ := ret[0].(request.HttpRequest)
	return ret0
}

// NewRequestWithContext indicates an expected call of NewRequestWithContext
func (mr *MockHttpRequestBuilderMockRecorder) NewRequestWithContext(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRequestWithContext", reflect.TypeOf((*MockHttpRequestBuilder)(nil).NewRequestWithContext), ctx)
}