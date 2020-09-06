// Code generated by MockGen. DO NOT EDIT.
// Source: trace/trace.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	trace "go.opencensus.io/trace"
	http "net/http"
	reflect "reflect"
)

// MockTrace is a mock of Trace interface
type MockTrace struct {
	ctrl     *gomock.Controller
	recorder *MockTraceMockRecorder
}

// MockTraceMockRecorder is the mock recorder for MockTrace
type MockTraceMockRecorder struct {
	mock *MockTrace
}

// NewMockTrace creates a new mock instance
func NewMockTrace(ctrl *gomock.Controller) *MockTrace {
	mock := &MockTrace{ctrl: ctrl}
	mock.recorder = &MockTraceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTrace) EXPECT() *MockTraceMockRecorder {
	return m.recorder
}

// Continue mocks base method
func (m *MockTrace) Continue(ctx context.Context, httpRequest *http.Request) (*trace.Span, *http.Request) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Continue", ctx, httpRequest)
	ret0, _ := ret[0].(*trace.Span)
	ret1, _ := ret[1].(*http.Request)
	return ret0, ret1
}

// Continue indicates an expected call of Continue
func (mr *MockTraceMockRecorder) Continue(ctx, httpRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Continue", reflect.TypeOf((*MockTrace)(nil).Continue), ctx, httpRequest)
}
