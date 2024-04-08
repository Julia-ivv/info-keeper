// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage (interfaces: Repositorier)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	storage "github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockRepositorier is a mock of Repositorier interface.
type MockRepositorier struct {
	ctrl     *gomock.Controller
	recorder *MockRepositorierMockRecorder
}

// MockRepositorierMockRecorder is the mock recorder for MockRepositorier.
type MockRepositorierMockRecorder struct {
	mock *MockRepositorier
}

// NewMockRepositorier creates a new mock instance.
func NewMockRepositorier(ctrl *gomock.Controller) *MockRepositorier {
	mock := &MockRepositorier{ctrl: ctrl}
	mock.recorder = &MockRepositorierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositorier) EXPECT() *MockRepositorierMockRecorder {
	return m.recorder
}

// AddBinaryRecord mocks base method.
func (m *MockRepositorier) AddBinaryRecord(arg0 context.Context, arg1 string, arg2, arg3, arg4 []byte, arg5 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBinaryRecord", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBinaryRecord indicates an expected call of AddBinaryRecord.
func (mr *MockRepositorierMockRecorder) AddBinaryRecord(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBinaryRecord", reflect.TypeOf((*MockRepositorier)(nil).AddBinaryRecord), arg0, arg1, arg2, arg3, arg4, arg5)
}

// AddCard mocks base method.
func (m *MockRepositorier) AddCard(arg0 context.Context, arg1 string, arg2, arg3, arg4, arg5, arg6 []byte, arg7 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCard", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCard indicates an expected call of AddCard.
func (mr *MockRepositorierMockRecorder) AddCard(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCard", reflect.TypeOf((*MockRepositorier)(nil).AddCard), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}

// AddLoginPwd mocks base method.
func (m *MockRepositorier) AddLoginPwd(arg0 context.Context, arg1 string, arg2, arg3, arg4, arg5 []byte, arg6 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLoginPwd", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLoginPwd indicates an expected call of AddLoginPwd.
func (mr *MockRepositorierMockRecorder) AddLoginPwd(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLoginPwd", reflect.TypeOf((*MockRepositorier)(nil).AddLoginPwd), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// AddSyncData mocks base method.
func (m *MockRepositorier) AddSyncData(arg0 context.Context, arg1 string, arg2 []storage.Card, arg3 []storage.LoginPwd, arg4 []storage.TextRecord, arg5 []storage.BinaryRecord) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSyncData", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSyncData indicates an expected call of AddSyncData.
func (mr *MockRepositorierMockRecorder) AddSyncData(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSyncData", reflect.TypeOf((*MockRepositorier)(nil).AddSyncData), arg0, arg1, arg2, arg3, arg4, arg5)
}

// AddTextRecord mocks base method.
func (m *MockRepositorier) AddTextRecord(arg0 context.Context, arg1 string, arg2, arg3, arg4 []byte, arg5 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTextRecord", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTextRecord indicates an expected call of AddTextRecord.
func (mr *MockRepositorierMockRecorder) AddTextRecord(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTextRecord", reflect.TypeOf((*MockRepositorier)(nil).AddTextRecord), arg0, arg1, arg2, arg3, arg4, arg5)
}

// AuthUser mocks base method.
func (m *MockRepositorier) AuthUser(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AuthUser indicates an expected call of AuthUser.
func (mr *MockRepositorierMockRecorder) AuthUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthUser", reflect.TypeOf((*MockRepositorier)(nil).AuthUser), arg0, arg1, arg2)
}

// Close mocks base method.
func (m *MockRepositorier) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRepositorierMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepositorier)(nil).Close))
}

// GetBinaryRecord mocks base method.
func (m *MockRepositorier) GetBinaryRecord(arg0 context.Context, arg1 string, arg2 []byte) (storage.BinaryRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBinaryRecord", arg0, arg1, arg2)
	ret0, _ := ret[0].(storage.BinaryRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBinaryRecord indicates an expected call of GetBinaryRecord.
func (mr *MockRepositorierMockRecorder) GetBinaryRecord(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBinaryRecord", reflect.TypeOf((*MockRepositorier)(nil).GetBinaryRecord), arg0, arg1, arg2)
}

// GetCard mocks base method.
func (m *MockRepositorier) GetCard(arg0 context.Context, arg1 string, arg2 []byte) (storage.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCard", arg0, arg1, arg2)
	ret0, _ := ret[0].(storage.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCard indicates an expected call of GetCard.
func (mr *MockRepositorierMockRecorder) GetCard(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCard", reflect.TypeOf((*MockRepositorier)(nil).GetCard), arg0, arg1, arg2)
}

// GetLastSyncTime mocks base method.
func (m *MockRepositorier) GetLastSyncTime(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastSyncTime", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastSyncTime indicates an expected call of GetLastSyncTime.
func (mr *MockRepositorierMockRecorder) GetLastSyncTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastSyncTime", reflect.TypeOf((*MockRepositorier)(nil).GetLastSyncTime), arg0, arg1)
}

// GetLoginPwd mocks base method.
func (m *MockRepositorier) GetLoginPwd(arg0 context.Context, arg1 string, arg2, arg3 []byte) (storage.LoginPwd, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoginPwd", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(storage.LoginPwd)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoginPwd indicates an expected call of GetLoginPwd.
func (mr *MockRepositorierMockRecorder) GetLoginPwd(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoginPwd", reflect.TypeOf((*MockRepositorier)(nil).GetLoginPwd), arg0, arg1, arg2, arg3)
}

// GetTextRecord mocks base method.
func (m *MockRepositorier) GetTextRecord(arg0 context.Context, arg1 string, arg2 []byte) (storage.TextRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTextRecord", arg0, arg1, arg2)
	ret0, _ := ret[0].(storage.TextRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTextRecord indicates an expected call of GetTextRecord.
func (mr *MockRepositorierMockRecorder) GetTextRecord(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTextRecord", reflect.TypeOf((*MockRepositorier)(nil).GetTextRecord), arg0, arg1, arg2)
}

// GetUserBinaryRecordsAfterTime mocks base method.
func (m *MockRepositorier) GetUserBinaryRecordsAfterTime(arg0 context.Context, arg1, arg2 string) ([]storage.BinaryRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBinaryRecordsAfterTime", arg0, arg1, arg2)
	ret0, _ := ret[0].([]storage.BinaryRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBinaryRecordsAfterTime indicates an expected call of GetUserBinaryRecordsAfterTime.
func (mr *MockRepositorierMockRecorder) GetUserBinaryRecordsAfterTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBinaryRecordsAfterTime", reflect.TypeOf((*MockRepositorier)(nil).GetUserBinaryRecordsAfterTime), arg0, arg1, arg2)
}

// GetUserCardsAfterTime mocks base method.
func (m *MockRepositorier) GetUserCardsAfterTime(arg0 context.Context, arg1, arg2 string) ([]storage.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCardsAfterTime", arg0, arg1, arg2)
	ret0, _ := ret[0].([]storage.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCardsAfterTime indicates an expected call of GetUserCardsAfterTime.
func (mr *MockRepositorierMockRecorder) GetUserCardsAfterTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCardsAfterTime", reflect.TypeOf((*MockRepositorier)(nil).GetUserCardsAfterTime), arg0, arg1, arg2)
}

// GetUserLoginsPwdsAfterTime mocks base method.
func (m *MockRepositorier) GetUserLoginsPwdsAfterTime(arg0 context.Context, arg1, arg2 string) ([]storage.LoginPwd, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserLoginsPwdsAfterTime", arg0, arg1, arg2)
	ret0, _ := ret[0].([]storage.LoginPwd)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserLoginsPwdsAfterTime indicates an expected call of GetUserLoginsPwdsAfterTime.
func (mr *MockRepositorierMockRecorder) GetUserLoginsPwdsAfterTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserLoginsPwdsAfterTime", reflect.TypeOf((*MockRepositorier)(nil).GetUserLoginsPwdsAfterTime), arg0, arg1, arg2)
}

// GetUserTextRecordsAfterTime mocks base method.
func (m *MockRepositorier) GetUserTextRecordsAfterTime(arg0 context.Context, arg1, arg2 string) ([]storage.TextRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserTextRecordsAfterTime", arg0, arg1, arg2)
	ret0, _ := ret[0].([]storage.TextRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserTextRecordsAfterTime indicates an expected call of GetUserTextRecordsAfterTime.
func (mr *MockRepositorierMockRecorder) GetUserTextRecordsAfterTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserTextRecordsAfterTime", reflect.TypeOf((*MockRepositorier)(nil).GetUserTextRecordsAfterTime), arg0, arg1, arg2)
}

// RegUser mocks base method.
func (m *MockRepositorier) RegUser(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegUser indicates an expected call of RegUser.
func (mr *MockRepositorierMockRecorder) RegUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegUser", reflect.TypeOf((*MockRepositorier)(nil).RegUser), arg0, arg1, arg2)
}

// UpdateBinaryRecord mocks base method.
func (m *MockRepositorier) UpdateBinaryRecord(arg0 context.Context, arg1 string, arg2, arg3, arg4 []byte, arg5 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBinaryRecord", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBinaryRecord indicates an expected call of UpdateBinaryRecord.
func (mr *MockRepositorierMockRecorder) UpdateBinaryRecord(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBinaryRecord", reflect.TypeOf((*MockRepositorier)(nil).UpdateBinaryRecord), arg0, arg1, arg2, arg3, arg4, arg5)
}

// UpdateCard mocks base method.
func (m *MockRepositorier) UpdateCard(arg0 context.Context, arg1 string, arg2, arg3, arg4, arg5, arg6 []byte, arg7 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCard", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCard indicates an expected call of UpdateCard.
func (mr *MockRepositorierMockRecorder) UpdateCard(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCard", reflect.TypeOf((*MockRepositorier)(nil).UpdateCard), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}

// UpdateLastSyncTime mocks base method.
func (m *MockRepositorier) UpdateLastSyncTime(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLastSyncTime", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLastSyncTime indicates an expected call of UpdateLastSyncTime.
func (mr *MockRepositorierMockRecorder) UpdateLastSyncTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLastSyncTime", reflect.TypeOf((*MockRepositorier)(nil).UpdateLastSyncTime), arg0, arg1, arg2)
}

// UpdateLoginPwd mocks base method.
func (m *MockRepositorier) UpdateLoginPwd(arg0 context.Context, arg1 string, arg2, arg3, arg4, arg5 []byte, arg6 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLoginPwd", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLoginPwd indicates an expected call of UpdateLoginPwd.
func (mr *MockRepositorierMockRecorder) UpdateLoginPwd(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLoginPwd", reflect.TypeOf((*MockRepositorier)(nil).UpdateLoginPwd), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// UpdateTextRecord mocks base method.
func (m *MockRepositorier) UpdateTextRecord(arg0 context.Context, arg1 string, arg2, arg3, arg4 []byte, arg5 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTextRecord", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTextRecord indicates an expected call of UpdateTextRecord.
func (mr *MockRepositorierMockRecorder) UpdateTextRecord(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTextRecord", reflect.TypeOf((*MockRepositorier)(nil).UpdateTextRecord), arg0, arg1, arg2, arg3, arg4, arg5)
}
