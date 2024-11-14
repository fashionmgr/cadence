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

package matching

import (
	"math/rand"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/cache"
	"github.com/uber/cadence/common/log"
	"github.com/uber/cadence/common/log/tag"
	"github.com/uber/cadence/common/types"
)

const _backlogThreshold = 100

type (
	weightSelector struct {
		sync.RWMutex
		weights     []int64
		initialized bool
		threshold   int64
	}

	weightedLoadBalancer struct {
		fallbackLoadBalancer LoadBalancer
		provider             PartitionConfigProvider
		weightCache          cache.Cache
		logger               log.Logger
	}
)

func newWeightSelector(n int, threshold int64) *weightSelector {
	pw := &weightSelector{weights: make([]int64, n), threshold: threshold}
	for i := range pw.weights {
		pw.weights[i] = -1
	}
	return pw
}

func (pw *weightSelector) pick() int {
	cumulativeWeights := make([]int64, len(pw.weights))
	totalWeight := int64(0)
	pw.RLock()
	defer pw.RUnlock()
	if !pw.initialized {
		return -1
	}
	shouldDrain := false
	for i, w := range pw.weights {
		totalWeight += w
		cumulativeWeights[i] = totalWeight
		if w > pw.threshold {
			shouldDrain = true // only enable weight selection if backlog size is larger than the threshold
		}
	}
	if totalWeight <= 0 || !shouldDrain {
		return -1
	}
	r := rand.Int63n(totalWeight)
	index := sort.Search(len(cumulativeWeights), func(i int) bool {
		return cumulativeWeights[i] > r
	})
	return index
}

func (pw *weightSelector) update(n, p int, weight int64) {
	pw.Lock()
	defer pw.Unlock()
	if n > len(pw.weights) {
		newWeights := make([]int64, n)
		copy(newWeights, pw.weights)
		for i := len(pw.weights); i < n; i++ {
			newWeights[i] = -1
		}
		pw.weights = newWeights
		pw.initialized = false
	} else if n < len(pw.weights) {
		pw.weights = pw.weights[:n]
	}
	pw.weights[p] = weight
	for _, w := range pw.weights {
		if w == -1 {
			return
		}
	}
	pw.initialized = true
}

func NewWeightedLoadBalancer(
	lb LoadBalancer,
	provider PartitionConfigProvider,
	logger log.Logger,
) LoadBalancer {
	return &weightedLoadBalancer{
		fallbackLoadBalancer: lb,
		provider:             provider,
		weightCache: cache.New(&cache.Options{
			TTL:             0,
			InitialCapacity: 100,
			Pin:             false,
			MaxCount:        3000,
			ActivelyEvict:   false,
		}),
		logger: logger,
	}
}

func (lb *weightedLoadBalancer) PickWritePartition(
	domainID string,
	taskList types.TaskList,
	taskListType int,
	forwardedFrom string,
) string {
	return lb.fallbackLoadBalancer.PickWritePartition(domainID, taskList, taskListType, forwardedFrom)
}

func (lb *weightedLoadBalancer) PickReadPartition(
	domainID string,
	taskList types.TaskList,
	taskListType int,
	forwardedFrom string,
) string {
	if forwardedFrom != "" || taskList.GetKind() == types.TaskListKindSticky {
		return taskList.GetName()
	}
	if strings.HasPrefix(taskList.GetName(), common.ReservedTaskListPrefix) {
		return taskList.GetName()
	}
	taskListKey := key{
		domainID:     domainID,
		taskListName: taskList.GetName(),
		taskListType: taskListType,
	}
	wI := lb.weightCache.Get(taskListKey)
	if wI == nil {
		return lb.fallbackLoadBalancer.PickReadPartition(domainID, taskList, taskListType, forwardedFrom)
	}
	w, ok := wI.(*weightSelector)
	if !ok {
		return lb.fallbackLoadBalancer.PickReadPartition(domainID, taskList, taskListType, forwardedFrom)
	}
	p := w.pick()
	lb.logger.Debug("pick read partition", tag.WorkflowDomainID(domainID), tag.WorkflowTaskListName(taskList.GetName()), tag.WorkflowTaskListType(taskListType), tag.Dynamic("weights", w.weights), tag.Dynamic("tasklist-partition", p))
	if p < 0 {
		return lb.fallbackLoadBalancer.PickReadPartition(domainID, taskList, taskListType, forwardedFrom)
	}
	return getPartitionTaskListName(taskList.GetName(), p)
}

func (lb *weightedLoadBalancer) UpdateWeight(
	domainID string,
	taskList types.TaskList,
	taskListType int,
	forwardedFrom string,
	partition string,
	weight int64,
) {
	if forwardedFrom != "" || taskList.GetKind() == types.TaskListKindSticky {
		return
	}
	if strings.HasPrefix(taskList.GetName(), common.ReservedTaskListPrefix) {
		return
	}
	p := 0
	if partition != taskList.GetName() {
		var err error
		p, err = strconv.Atoi(path.Base(partition))
		if err != nil {
			return
		}
	}
	taskListKey := key{
		domainID:     domainID,
		taskListName: taskList.GetName(),
		taskListType: taskListType,
	}
	n := lb.provider.GetNumberOfReadPartitions(domainID, taskList, taskListType)
	if n <= 1 {
		lb.weightCache.Delete(taskListKey)
		return
	}
	wI := lb.weightCache.Get(taskListKey)
	if wI == nil {
		var err error
		w := newWeightSelector(n, _backlogThreshold)
		wI, err = lb.weightCache.PutIfNotExist(taskListKey, w)
		if err != nil {
			return
		}
	}
	w, ok := wI.(*weightSelector)
	if !ok {
		return
	}
	lb.logger.Debug("update tasklist partition weight", tag.WorkflowDomainID(domainID), tag.WorkflowTaskListName(taskList.GetName()), tag.WorkflowTaskListType(taskListType), tag.Dynamic("weights", w.weights), tag.Dynamic("tasklist-partition", p), tag.Dynamic("weight", weight))
	w.update(n, p, weight)
}
