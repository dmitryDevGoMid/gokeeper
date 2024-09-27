// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_user is a generated GoMock package.
package mock_user

import (
	context "context"
	user "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserRepository) Create(arg0 context.Context, arg1 *user.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), arg0, arg1)
}

// CreateCardByUser mocks base method.
func (m *MockUserRepository) CreateCardByUser(ctx context.Context, user *user.SaveData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCardByUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCardByUser indicates an expected call of CreateCardByUser.
func (mr *MockUserRepositoryMockRecorder) CreateCardByUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCardByUser", reflect.TypeOf((*MockUserRepository)(nil).CreateCardByUser), ctx, user)
}

// CreatePasswordByUser mocks base method.
func (m *MockUserRepository) CreatePasswordByUser(ctx context.Context, user *user.SaveData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePasswordByUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePasswordByUser indicates an expected call of CreatePasswordByUser.
func (mr *MockUserRepositoryMockRecorder) CreatePasswordByUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePasswordByUser", reflect.TypeOf((*MockUserRepository)(nil).CreatePasswordByUser), ctx, user)
}

// DelereCardById mocks base method.
func (m *MockUserRepository) DelereCardById(ctx context.Context, id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelereCardById", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelereCardById indicates an expected call of DelereCardById.
func (mr *MockUserRepositoryMockRecorder) DelereCardById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelereCardById", reflect.TypeOf((*MockUserRepository)(nil).DelereCardById), ctx, id)
}

// DelerePasswordById mocks base method.
func (m *MockUserRepository) DelerePasswordById(ctx context.Context, id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelerePasswordById", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelerePasswordById indicates an expected call of DelerePasswordById.
func (mr *MockUserRepositoryMockRecorder) DelerePasswordById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelerePasswordById", reflect.TypeOf((*MockUserRepository)(nil).DelerePasswordById), ctx, id)
}

// DeleteByUsername mocks base method.
func (m *MockUserRepository) DeleteByUsername(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUsername", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUsername indicates an expected call of DeleteByUsername.
func (mr *MockUserRepositoryMockRecorder) DeleteByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUsername", reflect.TypeOf((*MockUserRepository)(nil).DeleteByUsername), arg0, arg1)
}

// GetByUsername mocks base method.
func (m *MockUserRepository) GetByUsername(arg0 context.Context, arg1 string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", arg0, arg1)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsername indicates an expected call of GetByUsername.
func (mr *MockUserRepositoryMockRecorder) GetByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockUserRepository)(nil).GetByUsername), arg0, arg1)
}

// GetCardsByUser mocks base method.
func (m *MockUserRepository) GetCardsByUser(ctx context.Context, userData *user.User) (*[]user.ResponseSaveData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCardsByUser", ctx, userData)
	ret0, _ := ret[0].(*[]user.ResponseSaveData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCardsByUser indicates an expected call of GetCardsByUser.
func (mr *MockUserRepositoryMockRecorder) GetCardsByUser(ctx, userData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCardsByUser", reflect.TypeOf((*MockUserRepository)(nil).GetCardsByUser), ctx, userData)
}

// GetPasswordByUser mocks base method.
func (m *MockUserRepository) GetPasswordByUser(ctx context.Context, userData *user.User) (*[]user.ResponseSaveData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPasswordByUser", ctx, userData)
	ret0, _ := ret[0].(*[]user.ResponseSaveData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPasswordByUser indicates an expected call of GetPasswordByUser.
func (mr *MockUserRepositoryMockRecorder) GetPasswordByUser(ctx, userData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswordByUser", reflect.TypeOf((*MockUserRepository)(nil).GetPasswordByUser), ctx, userData)
}

// UpdateCardByKey mocks base method.
func (m *MockUserRepository) UpdateCardByKey(ctx context.Context, user *user.SaveData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCardByKey", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCardByKey indicates an expected call of UpdateCardByKey.
func (mr *MockUserRepositoryMockRecorder) UpdateCardByKey(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCardByKey", reflect.TypeOf((*MockUserRepository)(nil).UpdateCardByKey), ctx, user)
}

// UpdatePassword mocks base method.
func (m *MockUserRepository) UpdatePassword(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockUserRepositoryMockRecorder) UpdatePassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUserRepository)(nil).UpdatePassword), arg0, arg1, arg2)
}

// UpdatePasswordByKey mocks base method.
func (m *MockUserRepository) UpdatePasswordByKey(ctx context.Context, user *user.SaveData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePasswordByKey", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePasswordByKey indicates an expected call of UpdatePasswordByKey.
func (mr *MockUserRepositoryMockRecorder) UpdatePasswordByKey(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePasswordByKey", reflect.TypeOf((*MockUserRepository)(nil).UpdatePasswordByKey), ctx, user)
}
