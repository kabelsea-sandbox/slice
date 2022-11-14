// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kabelsea-games/chess-backend-app/internal/caching (interfaces: AerospikeStore)

// Package cachingmock is a generated GoMock package.
package cachingmock

import (
	context "context"
	reflect "reflect"
	time "time"

	store "github.com/eko/gocache/v3/store"
	gomock "github.com/golang/mock/gomock"
)

// MockAerospikeStore is a mock of AerospikeStore interface.
type MockAerospikeStore struct {
	ctrl     *gomock.Controller
	recorder *MockAerospikeStoreMockRecorder
}

// MockAerospikeStoreMockRecorder is the mock recorder for MockAerospikeStore.
type MockAerospikeStoreMockRecorder struct {
	mock *MockAerospikeStore
}

// NewMockAerospikeStore creates a new mock instance.
func NewMockAerospikeStore(ctrl *gomock.Controller) *MockAerospikeStore {
	mock := &MockAerospikeStore{ctrl: ctrl}
	mock.recorder = &MockAerospikeStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAerospikeStore) EXPECT() *MockAerospikeStoreMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockAerospikeStore) Clear(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear.
func (mr *MockAerospikeStoreMockRecorder) Clear(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockAerospikeStore)(nil).Clear), arg0)
}

// Delete mocks base method.
func (m *MockAerospikeStore) Delete(arg0 context.Context, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAerospikeStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAerospikeStore)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockAerospikeStore) Get(arg0 context.Context, arg1 interface{}) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAerospikeStoreMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAerospikeStore)(nil).Get), arg0, arg1)
}

// GetType mocks base method.
func (m *MockAerospikeStore) GetType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetType")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetType indicates an expected call of GetType.
func (mr *MockAerospikeStoreMockRecorder) GetType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetType", reflect.TypeOf((*MockAerospikeStore)(nil).GetType))
}

// GetWithTTL mocks base method.
func (m *MockAerospikeStore) GetWithTTL(arg0 context.Context, arg1 interface{}) (interface{}, time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithTTL", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(time.Duration)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetWithTTL indicates an expected call of GetWithTTL.
func (mr *MockAerospikeStoreMockRecorder) GetWithTTL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithTTL", reflect.TypeOf((*MockAerospikeStore)(nil).GetWithTTL), arg0, arg1)
}

// Invalidate mocks base method.
func (m *MockAerospikeStore) Invalidate(arg0 context.Context, arg1 ...store.InvalidateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Invalidate", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Invalidate indicates an expected call of Invalidate.
func (mr *MockAerospikeStoreMockRecorder) Invalidate(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invalidate", reflect.TypeOf((*MockAerospikeStore)(nil).Invalidate), varargs...)
}

// Set mocks base method.
func (m *MockAerospikeStore) Set(arg0 context.Context, arg1, arg2 interface{}, arg3 ...store.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Set", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockAerospikeStoreMockRecorder) Set(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockAerospikeStore)(nil).Set), varargs...)
}
