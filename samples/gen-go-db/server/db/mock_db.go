// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package db is a generated GoMock package.
package db

import (
	context "context"
	models "github.com/Clever/wag/samples/gen-go-db/models"
	strfmt "github.com/go-openapi/strfmt"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockInterface is a mock of Interface interface
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// SaveSimpleThing mocks base method
func (m_2 *MockInterface) SaveSimpleThing(ctx context.Context, m models.SimpleThing) error {
	ret := m_2.ctrl.Call(m_2, "SaveSimpleThing", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSimpleThing indicates an expected call of SaveSimpleThing
func (mr *MockInterfaceMockRecorder) SaveSimpleThing(ctx, m interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSimpleThing", reflect.TypeOf((*MockInterface)(nil).SaveSimpleThing), ctx, m)
}

// GetSimpleThing mocks base method
func (m *MockInterface) GetSimpleThing(ctx context.Context, name string) (*models.SimpleThing, error) {
	ret := m.ctrl.Call(m, "GetSimpleThing", ctx, name)
	ret0, _ := ret[0].(*models.SimpleThing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSimpleThing indicates an expected call of GetSimpleThing
func (mr *MockInterfaceMockRecorder) GetSimpleThing(ctx, name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSimpleThing", reflect.TypeOf((*MockInterface)(nil).GetSimpleThing), ctx, name)
}

// DeleteSimpleThing mocks base method
func (m *MockInterface) DeleteSimpleThing(ctx context.Context, name string) error {
	ret := m.ctrl.Call(m, "DeleteSimpleThing", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSimpleThing indicates an expected call of DeleteSimpleThing
func (mr *MockInterfaceMockRecorder) DeleteSimpleThing(ctx, name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSimpleThing", reflect.TypeOf((*MockInterface)(nil).DeleteSimpleThing), ctx, name)
}

// SaveThing mocks base method
func (m_2 *MockInterface) SaveThing(ctx context.Context, m models.Thing) error {
	ret := m_2.ctrl.Call(m_2, "SaveThing", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveThing indicates an expected call of SaveThing
func (mr *MockInterfaceMockRecorder) SaveThing(ctx, m interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveThing", reflect.TypeOf((*MockInterface)(nil).SaveThing), ctx, m)
}

// GetThing mocks base method
func (m *MockInterface) GetThing(ctx context.Context, name string, version int64) (*models.Thing, error) {
	ret := m.ctrl.Call(m, "GetThing", ctx, name, version)
	ret0, _ := ret[0].(*models.Thing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThing indicates an expected call of GetThing
func (mr *MockInterfaceMockRecorder) GetThing(ctx, name, version interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThing", reflect.TypeOf((*MockInterface)(nil).GetThing), ctx, name, version)
}

// GetThingsByNameAndVersion mocks base method
func (m *MockInterface) GetThingsByNameAndVersion(ctx context.Context, input GetThingsByNameAndVersionInput) ([]models.Thing, error) {
	ret := m.ctrl.Call(m, "GetThingsByNameAndVersion", ctx, input)
	ret0, _ := ret[0].([]models.Thing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThingsByNameAndVersion indicates an expected call of GetThingsByNameAndVersion
func (mr *MockInterfaceMockRecorder) GetThingsByNameAndVersion(ctx, input interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThingsByNameAndVersion", reflect.TypeOf((*MockInterface)(nil).GetThingsByNameAndVersion), ctx, input)
}

// DeleteThing mocks base method
func (m *MockInterface) DeleteThing(ctx context.Context, name string, version int64) error {
	ret := m.ctrl.Call(m, "DeleteThing", ctx, name, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteThing indicates an expected call of DeleteThing
func (mr *MockInterfaceMockRecorder) DeleteThing(ctx, name, version interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteThing", reflect.TypeOf((*MockInterface)(nil).DeleteThing), ctx, name, version)
}

// GetThingByID mocks base method
func (m *MockInterface) GetThingByID(ctx context.Context, id string) (*models.Thing, error) {
	ret := m.ctrl.Call(m, "GetThingByID", ctx, id)
	ret0, _ := ret[0].(*models.Thing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThingByID indicates an expected call of GetThingByID
func (mr *MockInterfaceMockRecorder) GetThingByID(ctx, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThingByID", reflect.TypeOf((*MockInterface)(nil).GetThingByID), ctx, id)
}

// GetThingsByNameAndCreatedAt mocks base method
func (m *MockInterface) GetThingsByNameAndCreatedAt(ctx context.Context, input GetThingsByNameAndCreatedAtInput) ([]models.Thing, error) {
	ret := m.ctrl.Call(m, "GetThingsByNameAndCreatedAt", ctx, input)
	ret0, _ := ret[0].([]models.Thing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThingsByNameAndCreatedAt indicates an expected call of GetThingsByNameAndCreatedAt
func (mr *MockInterfaceMockRecorder) GetThingsByNameAndCreatedAt(ctx, input interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThingsByNameAndCreatedAt", reflect.TypeOf((*MockInterface)(nil).GetThingsByNameAndCreatedAt), ctx, input)
}

// SaveThingWithDateRange mocks base method
func (m_2 *MockInterface) SaveThingWithDateRange(ctx context.Context, m models.ThingWithDateRange) error {
	ret := m_2.ctrl.Call(m_2, "SaveThingWithDateRange", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveThingWithDateRange indicates an expected call of SaveThingWithDateRange
func (mr *MockInterfaceMockRecorder) SaveThingWithDateRange(ctx, m interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveThingWithDateRange", reflect.TypeOf((*MockInterface)(nil).SaveThingWithDateRange), ctx, m)
}

// GetThingWithDateRange mocks base method
func (m *MockInterface) GetThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) (*models.ThingWithDateRange, error) {
	ret := m.ctrl.Call(m, "GetThingWithDateRange", ctx, name, date)
	ret0, _ := ret[0].(*models.ThingWithDateRange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThingWithDateRange indicates an expected call of GetThingWithDateRange
func (mr *MockInterfaceMockRecorder) GetThingWithDateRange(ctx, name, date interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThingWithDateRange", reflect.TypeOf((*MockInterface)(nil).GetThingWithDateRange), ctx, name, date)
}

// GetThingWithDateRangesByNameAndDate mocks base method
func (m *MockInterface) GetThingWithDateRangesByNameAndDate(ctx context.Context, input GetThingWithDateRangesByNameAndDateInput) ([]models.ThingWithDateRange, error) {
	ret := m.ctrl.Call(m, "GetThingWithDateRangesByNameAndDate", ctx, input)
	ret0, _ := ret[0].([]models.ThingWithDateRange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThingWithDateRangesByNameAndDate indicates an expected call of GetThingWithDateRangesByNameAndDate
func (mr *MockInterfaceMockRecorder) GetThingWithDateRangesByNameAndDate(ctx, input interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThingWithDateRangesByNameAndDate", reflect.TypeOf((*MockInterface)(nil).GetThingWithDateRangesByNameAndDate), ctx, input)
}

// DeleteThingWithDateRange mocks base method
func (m *MockInterface) DeleteThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) error {
	ret := m.ctrl.Call(m, "DeleteThingWithDateRange", ctx, name, date)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteThingWithDateRange indicates an expected call of DeleteThingWithDateRange
func (mr *MockInterfaceMockRecorder) DeleteThingWithDateRange(ctx, name, date interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteThingWithDateRange", reflect.TypeOf((*MockInterface)(nil).DeleteThingWithDateRange), ctx, name, date)
}

// SaveThingWithUnderscores mocks base method
func (m_2 *MockInterface) SaveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error {
	ret := m_2.ctrl.Call(m_2, "SaveThingWithUnderscores", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveThingWithUnderscores indicates an expected call of SaveThingWithUnderscores
func (mr *MockInterfaceMockRecorder) SaveThingWithUnderscores(ctx, m interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveThingWithUnderscores", reflect.TypeOf((*MockInterface)(nil).SaveThingWithUnderscores), ctx, m)
}

// GetThingWithUnderscores mocks base method
func (m *MockInterface) GetThingWithUnderscores(ctx context.Context, idApp string) (*models.ThingWithUnderscores, error) {
	ret := m.ctrl.Call(m, "GetThingWithUnderscores", ctx, idApp)
	ret0, _ := ret[0].(*models.ThingWithUnderscores)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThingWithUnderscores indicates an expected call of GetThingWithUnderscores
func (mr *MockInterfaceMockRecorder) GetThingWithUnderscores(ctx, idApp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThingWithUnderscores", reflect.TypeOf((*MockInterface)(nil).GetThingWithUnderscores), ctx, idApp)
}

// DeleteThingWithUnderscores mocks base method
func (m *MockInterface) DeleteThingWithUnderscores(ctx context.Context, idApp string) error {
	ret := m.ctrl.Call(m, "DeleteThingWithUnderscores", ctx, idApp)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteThingWithUnderscores indicates an expected call of DeleteThingWithUnderscores
func (mr *MockInterfaceMockRecorder) DeleteThingWithUnderscores(ctx, idApp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteThingWithUnderscores", reflect.TypeOf((*MockInterface)(nil).DeleteThingWithUnderscores), ctx, idApp)
}