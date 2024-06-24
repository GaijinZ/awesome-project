// Code generated by MockGen. DO NOT EDIT.
// Source: ../repositories/category_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "awesomeProject/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCategorer is a mock of Categorer interface.
type MockCategorer struct {
	ctrl     *gomock.Controller
	recorder *MockCategorerMockRecorder
}

// MockCategorerMockRecorder is the mock recorder for MockCategorer.
type MockCategorerMockRecorder struct {
	mock *MockCategorer
}

// NewMockCategorer creates a new mock instance.
func NewMockCategorer(ctrl *gomock.Controller) *MockCategorer {
	mock := &MockCategorer{ctrl: ctrl}
	mock.recorder = &MockCategorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategorer) EXPECT() *MockCategorerMockRecorder {
	return m.recorder
}

// CreateCategory mocks base method.
func (m *MockCategorer) CreateCategory(category models.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", category)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCategory indicates an expected call of CreateCategory.
func (mr *MockCategorerMockRecorder) CreateCategory(category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockCategorer)(nil).CreateCategory), category)
}

// DeleteCategory mocks base method.
func (m *MockCategorer) DeleteCategory(categoryID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategory", categoryID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategory indicates an expected call of DeleteCategory.
func (mr *MockCategorerMockRecorder) DeleteCategory(categoryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategory", reflect.TypeOf((*MockCategorer)(nil).DeleteCategory), categoryID)
}

// GetCategory mocks base method.
func (m *MockCategorer) GetCategory(categoryID string) (*models.CategoryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategory", categoryID)
	ret0, _ := ret[0].(*models.CategoryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategory indicates an expected call of GetCategory.
func (mr *MockCategorerMockRecorder) GetCategory(categoryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategory", reflect.TypeOf((*MockCategorer)(nil).GetCategory), categoryID)
}

// UpdateCategory mocks base method.
func (m *MockCategorer) UpdateCategory(category models.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", category)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCategory indicates an expected call of UpdateCategory.
func (mr *MockCategorerMockRecorder) UpdateCategory(category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockCategorer)(nil).UpdateCategory), category)
}
