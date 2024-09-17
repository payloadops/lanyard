// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/payloadops/lanyard/app/dal (interfaces: DynamoDBAPI)
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=mocks/mock_dynamo_db_client.go github.com/payloadops/lanyard/app/dal DynamoDBAPI
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gomock "go.uber.org/mock/gomock"
)

// MockDynamoDBAPI is a mock of DynamoDBAPI interface.
type MockDynamoDBAPI struct {
	ctrl     *gomock.Controller
	recorder *MockDynamoDBAPIMockRecorder
}

// MockDynamoDBAPIMockRecorder is the mock recorder for MockDynamoDBAPI.
type MockDynamoDBAPIMockRecorder struct {
	mock *MockDynamoDBAPI
}

// NewMockDynamoDBAPI creates a new mock instance.
func NewMockDynamoDBAPI(ctrl *gomock.Controller) *MockDynamoDBAPI {
	mock := &MockDynamoDBAPI{ctrl: ctrl}
	mock.recorder = &MockDynamoDBAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDynamoDBAPI) EXPECT() *MockDynamoDBAPIMockRecorder {
	return m.recorder
}

// DeleteItem mocks base method.
func (m *MockDynamoDBAPI) DeleteItem(arg0 context.Context, arg1 *dynamodb.DeleteItemInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteItem", varargs...)
	ret0, _ := ret[0].(*dynamodb.DeleteItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockDynamoDBAPIMockRecorder) DeleteItem(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockDynamoDBAPI)(nil).DeleteItem), varargs...)
}

// GetItem mocks base method.
func (m *MockDynamoDBAPI) GetItem(arg0 context.Context, arg1 *dynamodb.GetItemInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetItem", varargs...)
	ret0, _ := ret[0].(*dynamodb.GetItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockDynamoDBAPIMockRecorder) GetItem(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockDynamoDBAPI)(nil).GetItem), varargs...)
}

// PutItem mocks base method.
func (m *MockDynamoDBAPI) PutItem(arg0 context.Context, arg1 *dynamodb.PutItemInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PutItem", varargs...)
	ret0, _ := ret[0].(*dynamodb.PutItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutItem indicates an expected call of PutItem.
func (mr *MockDynamoDBAPIMockRecorder) PutItem(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutItem", reflect.TypeOf((*MockDynamoDBAPI)(nil).PutItem), varargs...)
}

// Query mocks base method.
func (m *MockDynamoDBAPI) Query(arg0 context.Context, arg1 *dynamodb.QueryInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*dynamodb.QueryOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockDynamoDBAPIMockRecorder) Query(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDynamoDBAPI)(nil).Query), varargs...)
}

// UpdateItem mocks base method.
func (m *MockDynamoDBAPI) UpdateItem(arg0 context.Context, arg1 *dynamodb.UpdateItemInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateItem", varargs...)
	ret0, _ := ret[0].(*dynamodb.UpdateItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateItem indicates an expected call of UpdateItem.
func (mr *MockDynamoDBAPIMockRecorder) UpdateItem(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItem", reflect.TypeOf((*MockDynamoDBAPI)(nil).UpdateItem), varargs...)
}
