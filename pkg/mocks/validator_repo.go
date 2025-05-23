// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/models/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/zheli/validator-key-manager-backend/pkg/models"
)

// MockValidatorRepo is a mock of ValidatorRepo interface.
type MockValidatorRepo struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorRepoMockRecorder
}

// MockValidatorRepoMockRecorder is the mock recorder for MockValidatorRepo.
type MockValidatorRepoMockRecorder struct {
	mock *MockValidatorRepo
}

// NewMockValidatorRepo creates a new mock instance.
func NewMockValidatorRepo(ctrl *gomock.Controller) *MockValidatorRepo {
	mock := &MockValidatorRepo{ctrl: ctrl}
	mock.recorder = &MockValidatorRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidatorRepo) EXPECT() *MockValidatorRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockValidatorRepo) Create(ctx context.Context, v *models.Validator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockValidatorRepoMockRecorder) Create(ctx, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockValidatorRepo)(nil).Create), ctx, v)
}

// GetByPubkey mocks base method.
func (m *MockValidatorRepo) GetByPubkey(ctx context.Context, pubkey string) (*models.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPubkey", ctx, pubkey)
	ret0, _ := ret[0].(*models.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPubkey indicates an expected call of GetByPubkey.
func (mr *MockValidatorRepoMockRecorder) GetByPubkey(ctx, pubkey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPubkey", reflect.TypeOf((*MockValidatorRepo)(nil).GetByPubkey), ctx, pubkey)
}

// List mocks base method.
func (m *MockValidatorRepo) List(ctx context.Context, filters map[string]interface{}) ([]models.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filters)
	ret0, _ := ret[0].([]models.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockValidatorRepoMockRecorder) List(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockValidatorRepo)(nil).List), ctx, filters)
}

// UpdateStatus mocks base method.
func (m *MockValidatorRepo) UpdateStatus(ctx context.Context, pubkey, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, pubkey, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockValidatorRepoMockRecorder) UpdateStatus(ctx, pubkey, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockValidatorRepo)(nil).UpdateStatus), ctx, pubkey, status)
}
