// Code generated by MockGen. DO NOT EDIT.
// Source: error_response_interceptor.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gin "github.com/gin-gonic/gin"
	error "github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/error"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockResponseInterceptor is a mock of ResponseInterceptor interface
type MockResponseInterceptor struct {
	ctrl     *gomock.Controller
	recorder *MockResponseInterceptorMockRecorder
}

// MockResponseInterceptorMockRecorder is the mock recorder for MockResponseInterceptor
type MockResponseInterceptorMockRecorder struct {
	mock *MockResponseInterceptor
}

// NewMockResponseInterceptor creates a new mock instance
func NewMockResponseInterceptor(ctrl *gomock.Controller) *MockResponseInterceptor {
	mock := &MockResponseInterceptor{ctrl: ctrl}
	mock.recorder = &MockResponseInterceptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResponseInterceptor) EXPECT() *MockResponseInterceptorMockRecorder {
	return m.recorder
}

// HandleServiceError mocks base method
func (m *MockResponseInterceptor) HandleServiceError(ctx *gin.Context, serviceError *error.OAuthMiddlewareError) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleServiceError", ctx, serviceError)
}

// HandleServiceError indicates an expected call of HandleServiceError
func (mr *MockResponseInterceptorMockRecorder) HandleServiceError(ctx, serviceError interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleServiceError", reflect.TypeOf((*MockResponseInterceptor)(nil).HandleServiceError), ctx, serviceError)
}