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

package failure

import (
	"context"
	"strings"

	"github.com/uber/cadence/common/types"
	"github.com/uber/cadence/service/worker/diagnostics/invariant"
)

// Failure is an invariant that will be used to identify the different failures in the workflow execution history
type Failure invariant.Invariant

type failure struct {
	workflowExecutionHistory *types.GetWorkflowExecutionHistoryResponse
	domain                   string
}

type Params struct {
	WorkflowExecutionHistory *types.GetWorkflowExecutionHistoryResponse
	Domain                   string
}

func NewInvariant(p Params) Failure {
	return &failure{
		workflowExecutionHistory: p.WorkflowExecutionHistory,
		domain:                   p.Domain,
	}
}

func (f *failure) Check(context.Context) ([]invariant.InvariantCheckResult, error) {
	result := make([]invariant.InvariantCheckResult, 0)
	events := f.workflowExecutionHistory.GetHistory().GetEvents()
	for _, event := range events {
		if event.GetWorkflowExecutionFailedEventAttributes() != nil && event.WorkflowExecutionFailedEventAttributes.Reason != nil {
			attr := event.WorkflowExecutionFailedEventAttributes
			reason := attr.Reason
			identity := fetchIdentity(attr, events)
			result = append(result, invariant.InvariantCheckResult{
				InvariantType: WorkflowFailed.String(),
				Reason:        errorTypeFromReason(*reason).String(),
				Metadata:      invariant.MarshalData(FailureMetadata{Identity: identity}),
			})
		}
		if event.GetActivityTaskFailedEventAttributes() != nil && event.ActivityTaskFailedEventAttributes.Reason != nil {
			attr := event.ActivityTaskFailedEventAttributes
			reason := attr.Reason
			scheduled := fetchScheduledEvent(attr, events)
			started := fetchStartedEvent(attr, events)
			result = append(result, invariant.InvariantCheckResult{
				InvariantType: ActivityFailed.String(),
				Reason:        errorTypeFromReason(*reason).String(),
				Metadata: invariant.MarshalData(FailureMetadata{
					Identity:          attr.Identity,
					ActivityScheduled: scheduled,
					ActivityStarted:   started,
				}),
			})
		}
	}
	return result, nil
}

func errorTypeFromReason(reason string) ErrorType {
	if strings.Contains(reason, "Generic") {
		return GenericError
	}
	if strings.Contains(reason, "Panic") {
		return PanicError
	}
	if strings.Contains(reason, "Timeout") {
		return TimeoutError
	}
	return CustomError
}

func fetchIdentity(attr *types.WorkflowExecutionFailedEventAttributes, events []*types.HistoryEvent) string {
	for _, event := range events {
		if event.ID == attr.DecisionTaskCompletedEventID {
			return event.GetDecisionTaskCompletedEventAttributes().Identity
		}
	}
	return ""
}

func fetchScheduledEvent(attr *types.ActivityTaskFailedEventAttributes, events []*types.HistoryEvent) *types.ActivityTaskScheduledEventAttributes {
	for _, event := range events {
		if event.ID == attr.GetScheduledEventID() {
			return event.GetActivityTaskScheduledEventAttributes()
		}
	}
	return nil
}

func fetchStartedEvent(attr *types.ActivityTaskFailedEventAttributes, events []*types.HistoryEvent) *types.ActivityTaskStartedEventAttributes {
	for _, event := range events {
		if event.ID == attr.GetStartedEventID() {
			return event.GetActivityTaskStartedEventAttributes()
		}
	}
	return nil
}

func (f *failure) RootCause(ctx context.Context, issues []invariant.InvariantCheckResult) ([]invariant.InvariantRootCauseResult, error) {
	result := make([]invariant.InvariantRootCauseResult, 0)
	for _, issue := range issues {
		switch issue.Reason {
		case GenericError.String():
			result = append(result, invariant.InvariantRootCauseResult{
				RootCause: invariant.RootCauseTypeServiceSideIssue,
				Metadata:  issue.Metadata,
			})
		case PanicError.String():
			result = append(result, invariant.InvariantRootCauseResult{
				RootCause: invariant.RootCauseTypeServiceSidePanic,
				Metadata:  issue.Metadata,
			})
		case CustomError.String():
			result = append(result, invariant.InvariantRootCauseResult{
				RootCause: invariant.RootCauseTypeServiceSideCustomError,
				Metadata:  issue.Metadata,
			})
		}
	}
	return result, nil
}
