// The MIT License (MIT)
//
// Copyright (c) 2017-2022 Uber Technologies Inc.
//
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

package replication

import (
	"context"
	"strconv"
	"time"

	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/log"
	"github.com/uber/cadence/common/log/tag"
	"github.com/uber/cadence/common/metrics"
	"github.com/uber/cadence/common/persistence"
	"github.com/uber/cadence/common/types"
)

type (
	// TaskAckManager is the ack manager for replication tasks
	TaskAckManager struct {
		ackLevels ackLevelStore

		scope  metrics.Scope
		logger log.Logger

		reader taskReader
		store  *TaskStore
	}

	ackLevelStore interface {
		GetTransferMaxReadLevel() int64

		GetClusterReplicationLevel(cluster string) int64
		UpdateClusterReplicationLevel(cluster string, lastTaskID int64) error
	}
	taskReader interface {
		Read(ctx context.Context, readLevel int64, maxReadLevel int64) ([]*persistence.ReplicationTaskInfo, bool, error)
	}
)

// NewTaskAckManager initializes a new replication task ack manager
func NewTaskAckManager(
	shardID int,
	ackLevels ackLevelStore,
	metricsClient metrics.Client,
	logger log.Logger,
	reader taskReader,
	store *TaskStore,
) TaskAckManager {

	return TaskAckManager{
		ackLevels: ackLevels,
		scope: metricsClient.Scope(
			metrics.ReplicatorQueueProcessorScope,
			metrics.InstanceTag(strconv.Itoa(shardID)),
		),
		logger: logger.WithTags(tag.ComponentReplicationAckManager),
		reader: reader,
		store:  store,
	}
}

func (t *TaskAckManager) GetTasks(ctx context.Context, pollingCluster string, lastReadTaskID int64) (*types.ReplicationMessages, error) {
	if lastReadTaskID == common.EmptyMessageID {
		lastReadTaskID = t.ackLevels.GetClusterReplicationLevel(pollingCluster)
	}

	taskGeneratedTimer := t.scope.StartTimer(metrics.TaskLatency)

	tasks, hasMore, err := t.reader.Read(ctx, lastReadTaskID, t.ackLevels.GetTransferMaxReadLevel())
	if err != nil {
		return nil, err
	}

	var replicationTasks []*types.ReplicationTask
	readLevel := lastReadTaskID
TaskInfoLoop:
	for _, task := range tasks {
		replicationTask, err := t.store.Get(ctx, pollingCluster, *task)

		switch err.(type) {
		case nil:
			// No action
		case *types.BadRequestError, *types.InternalDataInconsistencyError, *types.EntityNotExistsError:
			t.logger.Warn("Failed to get replication task.", tag.Error(err))
		default:
			t.logger.Error("Failed to get replication task. Return what we have so far.", tag.Error(err))
			hasMore = true
			break TaskInfoLoop
		}
		readLevel = task.TaskID
		if replicationTask != nil {
			replicationTasks = append(replicationTasks, replicationTask)
		}
	}
	taskGeneratedTimer.Stop()

	t.scope.RecordTimer(metrics.ReplicationTasksLag, time.Duration(t.ackLevels.GetTransferMaxReadLevel()-readLevel))
	t.scope.RecordTimer(metrics.ReplicationTasksFetched, time.Duration(len(tasks)))
	t.scope.RecordTimer(metrics.ReplicationTasksReturned, time.Duration(len(replicationTasks)))
	t.scope.RecordTimer(metrics.ReplicationTasksReturnedDiff, time.Duration(len(tasks)-len(replicationTasks)))

	if err := t.ackLevels.UpdateClusterReplicationLevel(pollingCluster, lastReadTaskID); err != nil {
		t.logger.Error("error updating replication level for shard", tag.Error(err), tag.OperationFailed)
	}

	if err := t.store.Ack(pollingCluster, lastReadTaskID); err != nil {
		t.logger.Error("error updating replication level for hydrated task store", tag.Error(err), tag.OperationFailed)
	}

	t.logger.Debug("Get replication tasks", tag.SourceCluster(pollingCluster), tag.ShardReplicationAck(lastReadTaskID), tag.ReadLevel(readLevel))
	return &types.ReplicationMessages{
		ReplicationTasks:       replicationTasks,
		HasMore:                hasMore,
		LastRetrievedMessageID: readLevel,
	}, nil
}
