// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/payloadops/plato/app/dal (interfaces: APIKeyManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dal "github.com/payloadops/plato/app/dal"
)

// MockAPIKeyManager is a mock of APIKeyManager interface.
type MockAPIKeyManager struct {
	ctrl     *gomock.Controller
	recorder *MockAPIKeyManagerMockRecorder
}

// MockAPIKeyManagerMockRecorder is the mock recorder for MockAPIKeyManager.
type MockAPIKeyManagerMockRecorder struct {
	mock *MockAPIKeyManager
}

// NewMockAPIKeyManager creates a new mock instance.
func NewMockAPIKeyManager(ctrl *gomock.Controller) *MockAPIKeyManager {
	mock := &MockAPIKeyManager{ctrl: ctrl}
	mock.recorder = &MockAPIKeyManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPIKeyManager) EXPECT() *MockAPIKeyManagerMockRecorder {
	return m.recorder
}

// CreateAPIKey mocks base method.
func (m *MockAPIKeyManager) CreateAPIKey(arg0 context.Context, arg1 string, arg2 *dal.APIKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAPIKey", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAPIKey indicates an expected call of CreateAPIKey.
func (mr *MockAPIKeyManagerMockRecorder) CreateAPIKey(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAPIKey", reflect.TypeOf((*MockAPIKeyManager)(nil).CreateAPIKey), arg0, arg1, arg2)
}

// DeleteAPIKey mocks base method.
func (m *MockAPIKeyManager) DeleteAPIKey(arg0 context.Context, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAPIKey", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAPIKey indicates an expected call of DeleteAPIKey.
func (mr *MockAPIKeyManagerMockRecorder) DeleteAPIKey(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAPIKey", reflect.TypeOf((*MockAPIKeyManager)(nil).DeleteAPIKey), arg0, arg1, arg2, arg3)
}

// GetAPIKey mocks base method.
func (m *MockAPIKeyManager) GetAPIKey(arg0 context.Context, arg1, arg2, arg3 string) (*dal.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAPIKey", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*dal.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAPIKey indicates an expected call of GetAPIKey.
func (mr *MockAPIKeyManagerMockRecorder) GetAPIKey(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIKey", reflect.TypeOf((*MockAPIKeyManager)(nil).GetAPIKey), arg0, arg1, arg2, arg3)
}

// GetAPIKeyByIDAndSecret mocks base method.
func (m *MockAPIKeyManager) GetAPIKeyByIDAndSecret(arg0 context.Context, arg1, arg2 string) (*dal.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAPIKeyByIDAndSecret", arg0, arg1, arg2)
	ret0, _ := ret[0].(*dal.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAPIKeyByIDAndSecret indicates an expected call of GetAPIKeyByIDAndSecret.
func (mr *MockAPIKeyManagerMockRecorder) GetAPIKeyByIDAndSecret(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIKeyByIDAndSecret", reflect.TypeOf((*MockAPIKeyManager)(nil).GetAPIKeyByIDAndSecret), arg0, arg1, arg2)
}

// ListAPIKeysByProject mocks base method.
func (m *MockAPIKeyManager) ListAPIKeysByProject(arg0 context.Context, arg1, arg2 string) ([]dal.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAPIKeysByProject", arg0, arg1, arg2)
	ret0, _ := ret[0].([]dal.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAPIKeysByProject indicates an expected call of ListAPIKeysByProject.
func (mr *MockAPIKeyManagerMockRecorder) ListAPIKeysByProject(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAPIKeysByProject", reflect.TypeOf((*MockAPIKeyManager)(nil).ListAPIKeysByProject), arg0, arg1, arg2)
}

// UpdateAPIKey mocks base method.
func (m *MockAPIKeyManager) UpdateAPIKey(arg0 context.Context, arg1 string, arg2 *dal.APIKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAPIKey", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAPIKey indicates an expected call of UpdateAPIKey.
func (mr *MockAPIKeyManagerMockRecorder) UpdateAPIKey(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAPIKey", reflect.TypeOf((*MockAPIKeyManager)(nil).UpdateAPIKey), arg0, arg1, arg2)
}
