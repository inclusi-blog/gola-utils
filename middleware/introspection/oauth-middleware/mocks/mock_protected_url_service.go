// Code generated by MockGen. DO NOT EDIT.
// Source: protected_url_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockProtectedUrlService is a mock of ProtectedUrlService interface
type MockProtectedUrlService struct {
	ctrl     *gomock.Controller
	recorder *MockProtectedUrlServiceMockRecorder
}

// MockProtectedUrlServiceMockRecorder is the mock recorder for MockProtectedUrlService
type MockProtectedUrlServiceMockRecorder struct {
	mock *MockProtectedUrlService
}

// NewMockProtectedUrlService creates a new mock instance
func NewMockProtectedUrlService(ctrl *gomock.Controller) *MockProtectedUrlService {
	mock := &MockProtectedUrlService{ctrl: ctrl}
	mock.recorder = &MockProtectedUrlServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProtectedUrlService) EXPECT() *MockProtectedUrlServiceMockRecorder {
	return m.recorder
}

// IsProtected mocks base method
func (m *MockProtectedUrlService) IsProtected(url string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsProtected", url)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsProtected indicates an expected call of IsProtected
func (mr *MockProtectedUrlServiceMockRecorder) IsProtected(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsProtected", reflect.TypeOf((*MockProtectedUrlService)(nil).IsProtected), url)
}
