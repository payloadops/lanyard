// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/payloadops/plato/app/dal (interfaces: ProjectManager)
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=mocks/mock_project_db_client.go github.com/payloadops/plato/app/dal ProjectManager
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dal "github.com/payloadops/plato/app/dal"
	gomock "go.uber.org/mock/gomock"
)

// MockProjectManager is a mock of ProjectManager interface.
type MockProjectManager struct {
	ctrl     *gomock.Controller
	recorder *MockProjectManagerMockRecorder
}

// MockProjectManagerMockRecorder is the mock recorder for MockProjectManager.
type MockProjectManagerMockRecorder struct {
	mock *MockProjectManager
}

// NewMockProjectManager creates a new mock instance.
func NewMockProjectManager(ctrl *gomock.Controller) *MockProjectManager {
	mock := &MockProjectManager{ctrl: ctrl}
	mock.recorder = &MockProjectManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectManager) EXPECT() *MockProjectManagerMockRecorder {
	return m.recorder
}

// CreateProject mocks base method.
func (m *MockProjectManager) CreateProject(arg0 context.Context, arg1 *dal.Project) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockProjectManagerMockRecorder) CreateProject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProjectManager)(nil).CreateProject), arg0, arg1)
}

// DeleteProject mocks base method.
func (m *MockProjectManager) DeleteProject(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject.
func (mr *MockProjectManagerMockRecorder) DeleteProject(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockProjectManager)(nil).DeleteProject), arg0, arg1, arg2)
}

// GetProject mocks base method.
func (m *MockProjectManager) GetProject(arg0 context.Context, arg1, arg2 string) (*dal.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", arg0, arg1, arg2)
	ret0, _ := ret[0].(*dal.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject.
func (mr *MockProjectManagerMockRecorder) GetProject(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockProjectManager)(nil).GetProject), arg0, arg1, arg2)
}

// ListProjectsByOrganization mocks base method.
func (m *MockProjectManager) ListProjectsByOrganization(arg0 context.Context, arg1 string) ([]dal.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProjectsByOrganization", arg0, arg1)
	ret0, _ := ret[0].([]dal.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProjectsByOrganization indicates an expected call of ListProjectsByOrganization.
func (mr *MockProjectManagerMockRecorder) ListProjectsByOrganization(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProjectsByOrganization", reflect.TypeOf((*MockProjectManager)(nil).ListProjectsByOrganization), arg0, arg1)
}

// UpdateProject mocks base method.
func (m *MockProjectManager) UpdateProject(arg0 context.Context, arg1 *dal.Project) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProject", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProject indicates an expected call of UpdateProject.
func (mr *MockProjectManagerMockRecorder) UpdateProject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProject", reflect.TypeOf((*MockProjectManager)(nil).UpdateProject), arg0, arg1)
}
