// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/payloadops/lanyard/app/dal (interfaces: ActorManager)
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=mocks/mock_actor_db_client.go github.com/payloadops/lanyard/app/dal ActorManager
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dal "github.com/payloadops/lanyard/app/dal"
	gomock "go.uber.org/mock/gomock"
)

// MockActorManager is a mock of ActorManager interface.
type MockActorManager struct {
	ctrl     *gomock.Controller
	recorder *MockActorManagerMockRecorder
}

// MockActorManagerMockRecorder is the mock recorder for MockActorManager.
type MockActorManagerMockRecorder struct {
	mock *MockActorManager
}

// NewMockActorManager creates a new mock instance.
func NewMockActorManager(ctrl *gomock.Controller) *MockActorManager {
	mock := &MockActorManager{ctrl: ctrl}
	mock.recorder = &MockActorManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActorManager) EXPECT() *MockActorManagerMockRecorder {
	return m.recorder
}

// CreateActor mocks base method.
func (m *MockActorManager) CreateActor(arg0 context.Context, arg1 string, arg2 *dal.Actor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActor", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateActor indicates an expected call of CreateActor.
func (mr *MockActorManagerMockRecorder) CreateActor(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActor", reflect.TypeOf((*MockActorManager)(nil).CreateActor), arg0, arg1, arg2)
}

// DeleteActor mocks base method.
func (m *MockActorManager) DeleteActor(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActor", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActor indicates an expected call of DeleteActor.
func (mr *MockActorManagerMockRecorder) DeleteActor(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActor", reflect.TypeOf((*MockActorManager)(nil).DeleteActor), arg0, arg1, arg2)
}

// GetActor mocks base method.
func (m *MockActorManager) GetActor(arg0 context.Context, arg1, arg2 string) (*dal.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActor", arg0, arg1, arg2)
	ret0, _ := ret[0].(*dal.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActor indicates an expected call of GetActor.
func (mr *MockActorManagerMockRecorder) GetActor(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActor", reflect.TypeOf((*MockActorManager)(nil).GetActor), arg0, arg1, arg2)
}

// ListActors mocks base method.
func (m *MockActorManager) ListActors(arg0 context.Context, arg1 string) ([]dal.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListActors", arg0, arg1)
	ret0, _ := ret[0].([]dal.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListActors indicates an expected call of ListActors.
func (mr *MockActorManagerMockRecorder) ListActors(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListActors", reflect.TypeOf((*MockActorManager)(nil).ListActors), arg0, arg1)
}

// UpdateActor mocks base method.
func (m *MockActorManager) UpdateActor(arg0 context.Context, arg1 string, arg2 *dal.Actor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateActor", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateActor indicates an expected call of UpdateActor.
func (mr *MockActorManagerMockRecorder) UpdateActor(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateActor", reflect.TypeOf((*MockActorManager)(nil).UpdateActor), arg0, arg1, arg2)
}
