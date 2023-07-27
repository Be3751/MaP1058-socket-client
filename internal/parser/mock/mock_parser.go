// Code generated by MockGen. DO NOT EDIT.
// Source: parser.go

// Package mock_parser is a generated GoMock package.
package mock_parser

import (
	reflect "reflect"

	model "github.com/Be3751/MaP1058-socket-client/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockParser is a mock of Parser interface.
type MockParser struct {
	ctrl     *gomock.Controller
	recorder *MockParserMockRecorder
}

// MockParserMockRecorder is the mock recorder for MockParser.
type MockParserMockRecorder struct {
	mock *MockParser
}

// NewMockParser creates a new mock instance.
func NewMockParser(ctrl *gomock.Controller) *MockParser {
	mock := &MockParser{ctrl: ctrl}
	mock.recorder = &MockParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParser) EXPECT() *MockParserMockRecorder {
	return m.recorder
}

// ToSignals mocks base method.
func (m *MockParser) ToSignals(adSignals []byte) (*model.Signals, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToSignals", adSignals)
	ret0, _ := ret[0].(*model.Signals)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToSignals indicates an expected call of ToSignals.
func (mr *MockParserMockRecorder) ToSignals(adSignals interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToSignals", reflect.TypeOf((*MockParser)(nil).ToSignals), adSignals)
}
