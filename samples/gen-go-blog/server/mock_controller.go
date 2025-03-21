// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package server is a generated GoMock package.
package server

import (
	context "context"
	reflect "reflect"

	models "github.com/Clever/wag/samples/gen-go-blog/models/v9"
	gomock "github.com/golang/mock/gomock"
)

// MockController is a mock of Controller interface.
type MockController struct {
	ctrl     *gomock.Controller
	recorder *MockControllerMockRecorder
}

// MockControllerMockRecorder is the mock recorder for MockController.
type MockControllerMockRecorder struct {
	mock *MockController
}

// NewMockController creates a new mock instance.
func NewMockController(ctrl *gomock.Controller) *MockController {
	mock := &MockController{ctrl: ctrl}
	mock.recorder = &MockControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockController) EXPECT() *MockControllerMockRecorder {
	return m.recorder
}

// GetSectionsForStudent mocks base method.
func (m *MockController) GetSectionsForStudent(ctx context.Context, studentID string) ([]models.Section, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSectionsForStudent", ctx, studentID)
	ret0, _ := ret[0].([]models.Section)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSectionsForStudent indicates an expected call of GetSectionsForStudent.
func (mr *MockControllerMockRecorder) GetSectionsForStudent(ctx, studentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSectionsForStudent", reflect.TypeOf((*MockController)(nil).GetSectionsForStudent), ctx, studentID)
}

// PostGradeFileForStudent mocks base method.
func (m *MockController) PostGradeFileForStudent(ctx context.Context, i *models.PostGradeFileForStudentInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostGradeFileForStudent", ctx, i)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostGradeFileForStudent indicates an expected call of PostGradeFileForStudent.
func (mr *MockControllerMockRecorder) PostGradeFileForStudent(ctx, i interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostGradeFileForStudent", reflect.TypeOf((*MockController)(nil).PostGradeFileForStudent), ctx, i)
}

// PostSectionsForStudent mocks base method.
func (m *MockController) PostSectionsForStudent(ctx context.Context, i *models.PostSectionsForStudentInput) ([]models.Section, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostSectionsForStudent", ctx, i)
	ret0, _ := ret[0].([]models.Section)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostSectionsForStudent indicates an expected call of PostSectionsForStudent.
func (mr *MockControllerMockRecorder) PostSectionsForStudent(ctx, i interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostSectionsForStudent", reflect.TypeOf((*MockController)(nil).PostSectionsForStudent), ctx, i)
}
