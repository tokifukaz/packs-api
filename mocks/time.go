// Code generated by MockGen. DO NOT EDIT.
// Source: internal/utils/time.go

// Package mock_utils is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockTime is a mock of Time interface.
type MockTime struct {
	ctrl     *gomock.Controller
	recorder *MockTimeMockRecorder
}

// MockTimeMockRecorder is the mock recorder for MockTime.
type MockTimeMockRecorder struct {
	mock *MockTime
}

// NewMockTime creates a new mock instance.
func NewMockTime(ctrl *gomock.Controller) *MockTime {
	mock := &MockTime{ctrl: ctrl}
	mock.recorder = &MockTimeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTime) EXPECT() *MockTimeMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *MockTime) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockTimeMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockTime)(nil).Now))
}

// Pattern mocks base method.
func (m *MockTime) Pattern() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pattern")
	ret0, _ := ret[0].(string)
	return ret0
}

// Pattern indicates an expected call of Pattern.
func (mr *MockTimeMockRecorder) Pattern() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pattern", reflect.TypeOf((*MockTime)(nil).Pattern))
}
