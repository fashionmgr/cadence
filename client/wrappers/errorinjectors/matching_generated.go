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

package errorinjectors

// Code generated by gowrap. DO NOT EDIT.
// template: ../../templates/errorinjectors.tmpl
// gowrap: http://github.com/hexdigest/gowrap

import (
	"context"

	"go.uber.org/yarpc"

	"github.com/uber/cadence/client/matching"
	"github.com/uber/cadence/common/errors"
	"github.com/uber/cadence/common/log"
	"github.com/uber/cadence/common/log/tag"
	"github.com/uber/cadence/common/types"
)

const (
	msgMatchingInjectedFakeErr = "Injected fake matching client error"
)

// matchingClient implements matching.Client interface instrumented with retries
type matchingClient struct {
	client        matching.Client
	errorRate     float64
	logger        log.Logger
	fakeErrFn     func(float64) error
	forwardCallFn func(error) bool
}

// NewMatchingClient creates a new instance of matchingClient that injects error into every call with a given rate.
func NewMatchingClient(client matching.Client, errorRate float64, logger log.Logger) matching.Client {
	return &matchingClient{
		client:        client,
		errorRate:     errorRate,
		logger:        logger,
		fakeErrFn:     errors.GenerateFakeError,
		forwardCallFn: errors.ShouldForwardCall,
	}
}

func (c *matchingClient) AddActivityTask(ctx context.Context, ap1 *types.AddActivityTaskRequest, p1 ...yarpc.CallOption) (ap2 *types.AddActivityTaskResponse, err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		ap2, err = c.client.AddActivityTask(ctx, ap1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationAddActivityTask,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) AddDecisionTask(ctx context.Context, ap1 *types.AddDecisionTaskRequest, p1 ...yarpc.CallOption) (ap2 *types.AddDecisionTaskResponse, err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		ap2, err = c.client.AddDecisionTask(ctx, ap1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationAddDecisionTask,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) CancelOutstandingPoll(ctx context.Context, cp1 *types.CancelOutstandingPollRequest, p1 ...yarpc.CallOption) (err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		err = c.client.CancelOutstandingPoll(ctx, cp1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationCancelOutstandingPoll,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) DescribeTaskList(ctx context.Context, mp1 *types.MatchingDescribeTaskListRequest, p1 ...yarpc.CallOption) (dp1 *types.DescribeTaskListResponse, err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		dp1, err = c.client.DescribeTaskList(ctx, mp1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationDescribeTaskList,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) GetTaskListsByDomain(ctx context.Context, gp1 *types.GetTaskListsByDomainRequest, p1 ...yarpc.CallOption) (gp2 *types.GetTaskListsByDomainResponse, err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		gp2, err = c.client.GetTaskListsByDomain(ctx, gp1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationGetTaskListsByDomain,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) ListTaskListPartitions(ctx context.Context, mp1 *types.MatchingListTaskListPartitionsRequest, p1 ...yarpc.CallOption) (lp1 *types.ListTaskListPartitionsResponse, err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		lp1, err = c.client.ListTaskListPartitions(ctx, mp1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationListTaskListPartitions,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) PollForActivityTask(ctx context.Context, mp1 *types.MatchingPollForActivityTaskRequest, p1 ...yarpc.CallOption) (mp2 *types.MatchingPollForActivityTaskResponse, err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		mp2, err = c.client.PollForActivityTask(ctx, mp1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationPollForActivityTask,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) PollForDecisionTask(ctx context.Context, mp1 *types.MatchingPollForDecisionTaskRequest, p1 ...yarpc.CallOption) (mp2 *types.MatchingPollForDecisionTaskResponse, err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		mp2, err = c.client.PollForDecisionTask(ctx, mp1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationPollForDecisionTask,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) QueryWorkflow(ctx context.Context, mp1 *types.MatchingQueryWorkflowRequest, p1 ...yarpc.CallOption) (qp1 *types.QueryWorkflowResponse, err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		qp1, err = c.client.QueryWorkflow(ctx, mp1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationQueryWorkflow,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}

func (c *matchingClient) RespondQueryTaskCompleted(ctx context.Context, mp1 *types.MatchingRespondQueryTaskCompletedRequest, p1 ...yarpc.CallOption) (err error) {
	fakeErr := c.fakeErrFn(c.errorRate)
	var forwardCall bool
	if forwardCall = c.forwardCallFn(fakeErr); forwardCall {
		err = c.client.RespondQueryTaskCompleted(ctx, mp1, p1...)
	}

	if fakeErr != nil {
		c.logger.Error(msgMatchingInjectedFakeErr,
			tag.MatchingClientOperationRespondQueryTaskCompleted,
			tag.Error(fakeErr),
			tag.Bool(forwardCall),
			tag.ClientError(err),
		)
		err = fakeErr
		return
	}
	return
}
