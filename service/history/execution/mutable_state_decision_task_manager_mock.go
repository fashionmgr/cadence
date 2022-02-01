// The MIT License (MIT)

// Copyright (c) 2017-2020 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Code generated by MockGen. DO NOT EDIT.
// Source: mutable_state_decision_task_manager.go

// Package execution is a generated GoMock package.
package execution

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	types "github.com/uber/cadence/common/types"
)

// MockmutableStateDecisionTaskManager is a mock of mutableStateDecisionTaskManager interface.
type MockmutableStateDecisionTaskManager struct {
	ctrl     *gomock.Controller
	recorder *MockmutableStateDecisionTaskManagerMockRecorder
}

// MockmutableStateDecisionTaskManagerMockRecorder is the mock recorder for MockmutableStateDecisionTaskManager.
type MockmutableStateDecisionTaskManagerMockRecorder struct {
	mock *MockmutableStateDecisionTaskManager
}

// NewMockmutableStateDecisionTaskManager creates a new mock instance.
func NewMockmutableStateDecisionTaskManager(ctrl *gomock.Controller) *MockmutableStateDecisionTaskManager {
	mock := &MockmutableStateDecisionTaskManager{ctrl: ctrl}
	mock.recorder = &MockmutableStateDecisionTaskManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmutableStateDecisionTaskManager) EXPECT() *MockmutableStateDecisionTaskManagerMockRecorder {
	return m.recorder
}

// AddDecisionTaskCompletedEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddDecisionTaskCompletedEvent(scheduleEventID, startedEventID int64, request *types.RespondDecisionTaskCompletedRequest, maxResetPoints int) (*types.HistoryEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecisionTaskCompletedEvent", scheduleEventID, startedEventID, request, maxResetPoints)
	ret0, _ := ret[0].(*types.HistoryEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDecisionTaskCompletedEvent indicates an expected call of AddDecisionTaskCompletedEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddDecisionTaskCompletedEvent(scheduleEventID, startedEventID, request, maxResetPoints interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTaskCompletedEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddDecisionTaskCompletedEvent), scheduleEventID, startedEventID, request, maxResetPoints)
}

// AddDecisionTaskFailedEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddDecisionTaskFailedEvent(scheduleEventID, startedEventID int64, cause types.DecisionTaskFailedCause, details []byte, identity, reason, binChecksum, baseRunID, newRunID string, forkEventVersion int64) (*types.HistoryEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecisionTaskFailedEvent", scheduleEventID, startedEventID, cause, details, identity, reason, binChecksum, baseRunID, newRunID, forkEventVersion)
	ret0, _ := ret[0].(*types.HistoryEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDecisionTaskFailedEvent indicates an expected call of AddDecisionTaskFailedEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddDecisionTaskFailedEvent(scheduleEventID, startedEventID, cause, details, identity, reason, binChecksum, baseRunID, newRunID, forkEventVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTaskFailedEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddDecisionTaskFailedEvent), scheduleEventID, startedEventID, cause, details, identity, reason, binChecksum, baseRunID, newRunID, forkEventVersion)
}

// AddDecisionTaskResetTimeoutEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddDecisionTaskResetTimeoutEvent(scheduleEventID int64, baseRunID, newRunID string, forkEventVersion int64, reason string) (*types.HistoryEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecisionTaskResetTimeoutEvent", scheduleEventID, baseRunID, newRunID, forkEventVersion, reason)
	ret0, _ := ret[0].(*types.HistoryEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDecisionTaskResetTimeoutEvent indicates an expected call of AddDecisionTaskResetTimeoutEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddDecisionTaskResetTimeoutEvent(scheduleEventID, baseRunID, newRunID, forkEventVersion, reason interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTaskResetTimeoutEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddDecisionTaskResetTimeoutEvent), scheduleEventID, baseRunID, newRunID, forkEventVersion, reason)
}

// AddDecisionTaskScheduleToStartTimeoutEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddDecisionTaskScheduleToStartTimeoutEvent(scheduleEventID int64) (*types.HistoryEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecisionTaskScheduleToStartTimeoutEvent", scheduleEventID)
	ret0, _ := ret[0].(*types.HistoryEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDecisionTaskScheduleToStartTimeoutEvent indicates an expected call of AddDecisionTaskScheduleToStartTimeoutEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddDecisionTaskScheduleToStartTimeoutEvent(scheduleEventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTaskScheduleToStartTimeoutEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddDecisionTaskScheduleToStartTimeoutEvent), scheduleEventID)
}

// AddDecisionTaskScheduledEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddDecisionTaskScheduledEvent(bypassTaskGeneration bool) (*DecisionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecisionTaskScheduledEvent", bypassTaskGeneration)
	ret0, _ := ret[0].(*DecisionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDecisionTaskScheduledEvent indicates an expected call of AddDecisionTaskScheduledEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddDecisionTaskScheduledEvent(bypassTaskGeneration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTaskScheduledEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddDecisionTaskScheduledEvent), bypassTaskGeneration)
}

// AddDecisionTaskScheduledEventAsHeartbeat mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddDecisionTaskScheduledEventAsHeartbeat(bypassTaskGeneration bool, originalScheduledTimestamp int64) (*DecisionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecisionTaskScheduledEventAsHeartbeat", bypassTaskGeneration, originalScheduledTimestamp)
	ret0, _ := ret[0].(*DecisionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDecisionTaskScheduledEventAsHeartbeat indicates an expected call of AddDecisionTaskScheduledEventAsHeartbeat.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddDecisionTaskScheduledEventAsHeartbeat(bypassTaskGeneration, originalScheduledTimestamp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTaskScheduledEventAsHeartbeat", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddDecisionTaskScheduledEventAsHeartbeat), bypassTaskGeneration, originalScheduledTimestamp)
}

// AddDecisionTaskStartedEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddDecisionTaskStartedEvent(scheduleEventID int64, requestID string, request *types.PollForDecisionTaskRequest) (*types.HistoryEvent, *DecisionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecisionTaskStartedEvent", scheduleEventID, requestID, request)
	ret0, _ := ret[0].(*types.HistoryEvent)
	ret1, _ := ret[1].(*DecisionInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddDecisionTaskStartedEvent indicates an expected call of AddDecisionTaskStartedEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddDecisionTaskStartedEvent(scheduleEventID, requestID, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTaskStartedEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddDecisionTaskStartedEvent), scheduleEventID, requestID, request)
}

// AddDecisionTaskTimedOutEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddDecisionTaskTimedOutEvent(scheduleEventID, startedEventID int64) (*types.HistoryEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecisionTaskTimedOutEvent", scheduleEventID, startedEventID)
	ret0, _ := ret[0].(*types.HistoryEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDecisionTaskTimedOutEvent indicates an expected call of AddDecisionTaskTimedOutEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddDecisionTaskTimedOutEvent(scheduleEventID, startedEventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTaskTimedOutEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddDecisionTaskTimedOutEvent), scheduleEventID, startedEventID)
}

// AddFirstDecisionTaskScheduled mocks base method.
func (m *MockmutableStateDecisionTaskManager) AddFirstDecisionTaskScheduled(startEvent *types.HistoryEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFirstDecisionTaskScheduled", startEvent)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFirstDecisionTaskScheduled indicates an expected call of AddFirstDecisionTaskScheduled.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) AddFirstDecisionTaskScheduled(startEvent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFirstDecisionTaskScheduled", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).AddFirstDecisionTaskScheduled), startEvent)
}

// CreateTransientDecisionEvents mocks base method.
func (m *MockmutableStateDecisionTaskManager) CreateTransientDecisionEvents(decision *DecisionInfo, identity string) (*types.HistoryEvent, *types.HistoryEvent) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransientDecisionEvents", decision, identity)
	ret0, _ := ret[0].(*types.HistoryEvent)
	ret1, _ := ret[1].(*types.HistoryEvent)
	return ret0, ret1
}

// CreateTransientDecisionEvents indicates an expected call of CreateTransientDecisionEvents.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) CreateTransientDecisionEvents(decision, identity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransientDecisionEvents", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).CreateTransientDecisionEvents), decision, identity)
}

// DeleteDecision mocks base method.
func (m *MockmutableStateDecisionTaskManager) DeleteDecision() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteDecision")
}

// DeleteDecision indicates an expected call of DeleteDecision.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) DeleteDecision() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDecision", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).DeleteDecision))
}

// FailDecision mocks base method.
func (m *MockmutableStateDecisionTaskManager) FailDecision(incrementAttempt bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FailDecision", incrementAttempt)
}

// FailDecision indicates an expected call of FailDecision.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) FailDecision(incrementAttempt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FailDecision", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).FailDecision), incrementAttempt)
}

// GetDecisionInfo mocks base method.
func (m *MockmutableStateDecisionTaskManager) GetDecisionInfo(scheduleEventID int64) (*DecisionInfo, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDecisionInfo", scheduleEventID)
	ret0, _ := ret[0].(*DecisionInfo)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetDecisionInfo indicates an expected call of GetDecisionInfo.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) GetDecisionInfo(scheduleEventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDecisionInfo", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).GetDecisionInfo), scheduleEventID)
}

// GetDecisionScheduleToStartTimeout mocks base method.
func (m *MockmutableStateDecisionTaskManager) GetDecisionScheduleToStartTimeout() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDecisionScheduleToStartTimeout")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// GetDecisionScheduleToStartTimeout indicates an expected call of GetDecisionScheduleToStartTimeout.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) GetDecisionScheduleToStartTimeout() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDecisionScheduleToStartTimeout", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).GetDecisionScheduleToStartTimeout))
}

// GetInFlightDecision mocks base method.
func (m *MockmutableStateDecisionTaskManager) GetInFlightDecision() (*DecisionInfo, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInFlightDecision")
	ret0, _ := ret[0].(*DecisionInfo)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetInFlightDecision indicates an expected call of GetInFlightDecision.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) GetInFlightDecision() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInFlightDecision", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).GetInFlightDecision))
}

// GetPendingDecision mocks base method.
func (m *MockmutableStateDecisionTaskManager) GetPendingDecision() (*DecisionInfo, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingDecision")
	ret0, _ := ret[0].(*DecisionInfo)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetPendingDecision indicates an expected call of GetPendingDecision.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) GetPendingDecision() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingDecision", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).GetPendingDecision))
}

// HasInFlightDecision mocks base method.
func (m *MockmutableStateDecisionTaskManager) HasInFlightDecision() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasInFlightDecision")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasInFlightDecision indicates an expected call of HasInFlightDecision.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) HasInFlightDecision() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasInFlightDecision", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).HasInFlightDecision))
}

// HasPendingDecision mocks base method.
func (m *MockmutableStateDecisionTaskManager) HasPendingDecision() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasPendingDecision")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasPendingDecision indicates an expected call of HasPendingDecision.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) HasPendingDecision() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasPendingDecision", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).HasPendingDecision))
}

// HasProcessedOrPendingDecision mocks base method.
func (m *MockmutableStateDecisionTaskManager) HasProcessedOrPendingDecision() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasProcessedOrPendingDecision")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasProcessedOrPendingDecision indicates an expected call of HasProcessedOrPendingDecision.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) HasProcessedOrPendingDecision() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasProcessedOrPendingDecision", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).HasProcessedOrPendingDecision))
}

// ReplicateDecisionTaskCompletedEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) ReplicateDecisionTaskCompletedEvent(event *types.HistoryEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplicateDecisionTaskCompletedEvent", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplicateDecisionTaskCompletedEvent indicates an expected call of ReplicateDecisionTaskCompletedEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) ReplicateDecisionTaskCompletedEvent(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplicateDecisionTaskCompletedEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).ReplicateDecisionTaskCompletedEvent), event)
}

// ReplicateDecisionTaskFailedEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) ReplicateDecisionTaskFailedEvent() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplicateDecisionTaskFailedEvent")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplicateDecisionTaskFailedEvent indicates an expected call of ReplicateDecisionTaskFailedEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) ReplicateDecisionTaskFailedEvent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplicateDecisionTaskFailedEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).ReplicateDecisionTaskFailedEvent))
}

// ReplicateDecisionTaskScheduledEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) ReplicateDecisionTaskScheduledEvent(version, scheduleID int64, taskList string, startToCloseTimeoutSeconds int32, attempt, scheduleTimestamp, originalScheduledTimestamp int64) (*DecisionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplicateDecisionTaskScheduledEvent", version, scheduleID, taskList, startToCloseTimeoutSeconds, attempt, scheduleTimestamp, originalScheduledTimestamp)
	ret0, _ := ret[0].(*DecisionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplicateDecisionTaskScheduledEvent indicates an expected call of ReplicateDecisionTaskScheduledEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) ReplicateDecisionTaskScheduledEvent(version, scheduleID, taskList, startToCloseTimeoutSeconds, attempt, scheduleTimestamp, originalScheduledTimestamp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplicateDecisionTaskScheduledEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).ReplicateDecisionTaskScheduledEvent), version, scheduleID, taskList, startToCloseTimeoutSeconds, attempt, scheduleTimestamp, originalScheduledTimestamp)
}

// ReplicateDecisionTaskStartedEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) ReplicateDecisionTaskStartedEvent(decision *DecisionInfo, version, scheduleID, startedID int64, requestID string, timestamp int64) (*DecisionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplicateDecisionTaskStartedEvent", decision, version, scheduleID, startedID, requestID, timestamp)
	ret0, _ := ret[0].(*DecisionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplicateDecisionTaskStartedEvent indicates an expected call of ReplicateDecisionTaskStartedEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) ReplicateDecisionTaskStartedEvent(decision, version, scheduleID, startedID, requestID, timestamp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplicateDecisionTaskStartedEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).ReplicateDecisionTaskStartedEvent), decision, version, scheduleID, startedID, requestID, timestamp)
}

// ReplicateDecisionTaskTimedOutEvent mocks base method.
func (m *MockmutableStateDecisionTaskManager) ReplicateDecisionTaskTimedOutEvent(timeoutType types.TimeoutType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplicateDecisionTaskTimedOutEvent", timeoutType)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplicateDecisionTaskTimedOutEvent indicates an expected call of ReplicateDecisionTaskTimedOutEvent.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) ReplicateDecisionTaskTimedOutEvent(timeoutType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplicateDecisionTaskTimedOutEvent", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).ReplicateDecisionTaskTimedOutEvent), timeoutType)
}

// ReplicateTransientDecisionTaskScheduled mocks base method.
func (m *MockmutableStateDecisionTaskManager) ReplicateTransientDecisionTaskScheduled() (*DecisionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplicateTransientDecisionTaskScheduled")
	ret0, _ := ret[0].(*DecisionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplicateTransientDecisionTaskScheduled indicates an expected call of ReplicateTransientDecisionTaskScheduled.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) ReplicateTransientDecisionTaskScheduled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplicateTransientDecisionTaskScheduled", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).ReplicateTransientDecisionTaskScheduled))
}

// UpdateDecision mocks base method.
func (m *MockmutableStateDecisionTaskManager) UpdateDecision(decision *DecisionInfo) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateDecision", decision)
}

// UpdateDecision indicates an expected call of UpdateDecision.
func (mr *MockmutableStateDecisionTaskManagerMockRecorder) UpdateDecision(decision interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDecision", reflect.TypeOf((*MockmutableStateDecisionTaskManager)(nil).UpdateDecision), decision)
}
