// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	domain "github.com/akrovv/warehouse/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockWarehouseService is a mock of WarehouseService interface.
type MockWarehouseService struct {
	ctrl     *gomock.Controller
	recorder *MockWarehouseServiceMockRecorder
}

// MockWarehouseServiceMockRecorder is the mock recorder for MockWarehouseService.
type MockWarehouseServiceMockRecorder struct {
	mock *MockWarehouseService
}

// NewMockWarehouseService creates a new mock instance.
func NewMockWarehouseService(ctrl *gomock.Controller) *MockWarehouseService {
	mock := &MockWarehouseService{ctrl: ctrl}
	mock.recorder = &MockWarehouseServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWarehouseService) EXPECT() *MockWarehouseServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockWarehouseService) Create(warehouse *domain.Warehouse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", warehouse)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockWarehouseServiceMockRecorder) Create(warehouse interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWarehouseService)(nil).Create), warehouse)
}

// GetLeftOvers mocks base method.
func (m *MockWarehouseService) GetLeftOvers(gw *domain.GetFromWarehouse) ([]domain.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLeftOvers", gw)
	ret0, _ := ret[0].([]domain.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLeftOvers indicates an expected call of GetLeftOvers.
func (mr *MockWarehouseServiceMockRecorder) GetLeftOvers(gw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLeftOvers", reflect.TypeOf((*MockWarehouseService)(nil).GetLeftOvers), gw)
}
