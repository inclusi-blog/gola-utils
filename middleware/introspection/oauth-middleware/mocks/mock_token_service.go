// Code generated by MockGen. DO NOT EDIT.
// Source: token_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gin "github.com/gin-gonic/gin"
	error "github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/error"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTokenService is a mock of TokenService interface
type MockTokenService struct {
	ctrl     *gomock.Controller
	recorder *MockTokenServiceMockRecorder
}

// MockTokenServiceMockRecorder is the mock recorder for MockTokenService
type MockTokenServiceMockRecorder struct {
	mock *MockTokenService
}

// NewMockTokenService creates a new mock instance
func NewMockTokenService(ctrl *gomock.Controller) *MockTokenService {
	mock := &MockTokenService{ctrl: ctrl}
	mock.recorder = &MockTokenServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTokenService) EXPECT() *MockTokenServiceMockRecorder {
	return m.recorder
}

// Validate mocks base method
func (m *MockTokenService) Validate(ctx *gin.Context) *error.OAuthMiddlewareError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", ctx)
	ret0, _ := ret[0].(*error.OAuthMiddlewareError)
	return ret0
}

// Validate indicates an expected call of Validate
func (mr *MockTokenServiceMockRecorder) Validate(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockTokenService)(nil).Validate), ctx)
}

// Decrypt mocks base method
func (m *MockTokenService) Decrypt(ctx *gin.Context) *error.OAuthMiddlewareError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decrypt", ctx)
	ret0, _ := ret[0].(*error.OAuthMiddlewareError)
	return ret0
}

// Decrypt indicates an expected call of Decrypt
func (mr *MockTokenServiceMockRecorder) Decrypt(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrypt", reflect.TypeOf((*MockTokenService)(nil).Decrypt), ctx)
}
