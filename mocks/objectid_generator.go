// Code generated by MockGen. DO NOT EDIT.
// Source: internal/utils/objectid_generator.go

// Package mock_utils is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockObjectIDGenerator is a mock of ObjectIDGenerator interface.
type MockObjectIDGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockObjectIDGeneratorMockRecorder
}

// MockObjectIDGeneratorMockRecorder is the mock recorder for MockObjectIDGenerator.
type MockObjectIDGeneratorMockRecorder struct {
	mock *MockObjectIDGenerator
}

// NewMockObjectIDGenerator creates a new mock instance.
func NewMockObjectIDGenerator(ctrl *gomock.Controller) *MockObjectIDGenerator {
	mock := &MockObjectIDGenerator{ctrl: ctrl}
	mock.recorder = &MockObjectIDGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObjectIDGenerator) EXPECT() *MockObjectIDGeneratorMockRecorder {
	return m.recorder
}

// GenerateRandomObjectID mocks base method.
func (m *MockObjectIDGenerator) GenerateRandomObjectID() primitive.ObjectID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRandomObjectID")
	ret0, _ := ret[0].(primitive.ObjectID)
	return ret0
}

// GenerateRandomObjectID indicates an expected call of GenerateRandomObjectID.
func (mr *MockObjectIDGeneratorMockRecorder) GenerateRandomObjectID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRandomObjectID", reflect.TypeOf((*MockObjectIDGenerator)(nil).GenerateRandomObjectID))
}

// ParseObjectID mocks base method.
func (m *MockObjectIDGenerator) ParseObjectID(id string) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseObjectID", id)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseObjectID indicates an expected call of ParseObjectID.
func (mr *MockObjectIDGeneratorMockRecorder) ParseObjectID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseObjectID", reflect.TypeOf((*MockObjectIDGenerator)(nil).ParseObjectID), id)
}
