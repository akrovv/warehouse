// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	domain "github.com/akrovv/warehouse/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockProductService is a mock of ProductService interface.
type MockProductService struct {
	ctrl     *gomock.Controller
	recorder *MockProductServiceMockRecorder
}

// MockProductServiceMockRecorder is the mock recorder for MockProductService.
type MockProductServiceMockRecorder struct {
	mock *MockProductService
}

// NewMockProductService creates a new mock instance.
func NewMockProductService(ctrl *gomock.Controller) *MockProductService {
	mock := &MockProductService{ctrl: ctrl}
	mock.recorder = &MockProductServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductService) EXPECT() *MockProductServiceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockProductService) Add(ad *domain.AddProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ad)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockProductServiceMockRecorder) Add(ad interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockProductService)(nil).Add), ad)
}

// CancelReservation mocks base method.
func (m *MockProductService) CancelReservation(wp *domain.WarehouseProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelReservation", wp)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelReservation indicates an expected call of CancelReservation.
func (mr *MockProductServiceMockRecorder) CancelReservation(wp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelReservation", reflect.TypeOf((*MockProductService)(nil).CancelReservation), wp)
}

// Create mocks base method.
func (m *MockProductService) Create(product *domain.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", product)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockProductServiceMockRecorder) Create(product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductService)(nil).Create), product)
}

// Delete mocks base method.
func (m *MockProductService) Delete(dp *domain.DeleteProduct) (*domain.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", dp)
	ret0, _ := ret[0].(*domain.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockProductServiceMockRecorder) Delete(dp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProductService)(nil).Delete), dp)
}

// Reserve mocks base method.
func (m *MockProductService) Reserve(wp *domain.WarehouseProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reserve", wp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reserve indicates an expected call of Reserve.
func (mr *MockProductServiceMockRecorder) Reserve(wp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reserve", reflect.TypeOf((*MockProductService)(nil).Reserve), wp)
}

// Transfer mocks base method.
func (m *MockProductService) Transfer(td *domain.TransferProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", td)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transfer indicates an expected call of Transfer.
func (mr *MockProductServiceMockRecorder) Transfer(td interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockProductService)(nil).Transfer), td)
}
