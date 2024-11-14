// Copyright (c) 2021 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package persistence

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/dynamicconfig"
	"github.com/uber/cadence/common/log"
	"github.com/uber/cadence/common/log/tag"
	"github.com/uber/cadence/common/types"
)

type (
	visibilityTripleManager struct {
		logger                       log.Logger
		dbVisibilityManager          VisibilityManager
		destinationVisibilityManager VisibilityManager
		sourceVisibilityManager      VisibilityManager
		readModeIsFromDestination    dynamicconfig.BoolPropertyFnWithDomainFilter
		readModeIsFromSource         dynamicconfig.BoolPropertyFnWithDomainFilter
		writeMode                    dynamicconfig.StringPropertyFn
		logCustomerQueryParameter    dynamicconfig.BoolPropertyFnWithDomainFilter
		readModeIsDouble             dynamicconfig.BoolPropertyFnWithDomainFilter
	}
)

const (
	VisibilityOverridePrimary   = "Primary"
	VisibilityOverrideSecondary = "Secondary"
	ContextKey                  = ResponseComparatorContextKey("visibility-override")
)

// ResponseComparatorContextKey is for Pinot/ES response comparator. This struct will be passed into ctx as a key.
type ResponseComparatorContextKey string

type OperationType string

var Operation = struct {
	LIST  OperationType
	COUNT OperationType
}{
	LIST:  "list",
	COUNT: "count",
}

var _ VisibilityManager = (*visibilityTripleManager)(nil)

// NewVisibilityTripleManager create a visibility manager that operate on DB or advanced visibility based on dynamic config.
// For Pinot migration, Pinot is the destination visibility manager, ES is the source visibility manager, and DB is the fallback.
// For OpenSearch migration, OS is the destination visibility manager, ES is the source visibility manager, and DB is the fallback.
func NewVisibilityTripleManager(
	dbVisibilityManager VisibilityManager, // one of the VisibilityManager can be nil
	destinationVisibilityManager VisibilityManager,
	sourceVisibilityManager VisibilityManager,
	readModeIsFromDestination dynamicconfig.BoolPropertyFnWithDomainFilter,
	readModeIsFromSource dynamicconfig.BoolPropertyFnWithDomainFilter,
	visWritingMode dynamicconfig.StringPropertyFn,
	logCustomerQueryParameter dynamicconfig.BoolPropertyFnWithDomainFilter,
	readModeIsDouble dynamicconfig.BoolPropertyFnWithDomainFilter,
	logger log.Logger,
) VisibilityManager {
	if dbVisibilityManager == nil && destinationVisibilityManager == nil && sourceVisibilityManager == nil {
		logger.Fatal("require one of dbVisibilityManager or pinotVisibilityManager or esVisibilityManager")
		return nil
	}
	return &visibilityTripleManager{
		dbVisibilityManager:          dbVisibilityManager,
		destinationVisibilityManager: destinationVisibilityManager,
		sourceVisibilityManager:      sourceVisibilityManager,
		readModeIsFromDestination:    readModeIsFromDestination,
		readModeIsFromSource:         readModeIsFromSource,
		writeMode:                    visWritingMode,
		logger:                       logger,
		logCustomerQueryParameter:    logCustomerQueryParameter,
		readModeIsDouble:             readModeIsDouble,
	}
}

func (v *visibilityTripleManager) Close() {
	if v.dbVisibilityManager != nil {
		v.dbVisibilityManager.Close()
	}
	if v.destinationVisibilityManager != nil {
		v.destinationVisibilityManager.Close()
	}
	if v.sourceVisibilityManager != nil {
		v.sourceVisibilityManager.Close()
	}
}

func (v *visibilityTripleManager) GetName() string {
	if v.destinationVisibilityManager != nil {
		return v.destinationVisibilityManager.GetName()
	} else if v.sourceVisibilityManager != nil {
		return v.sourceVisibilityManager.GetName()
	}
	return v.dbVisibilityManager.GetName()
}

func (v *visibilityTripleManager) RecordWorkflowExecutionStarted(
	ctx context.Context,
	request *RecordWorkflowExecutionStartedRequest,
) error {
	return v.chooseVisibilityManagerForWrite(
		ctx,
		func() error {
			return v.dbVisibilityManager.RecordWorkflowExecutionStarted(ctx, request)
		},
		func() error {
			return v.sourceVisibilityManager.RecordWorkflowExecutionStarted(ctx, request)
		},
		func() error {
			return v.destinationVisibilityManager.RecordWorkflowExecutionStarted(ctx, request)
		},
	)
}

func (v *visibilityTripleManager) RecordWorkflowExecutionClosed(
	ctx context.Context,
	request *RecordWorkflowExecutionClosedRequest,
) error {
	return v.chooseVisibilityManagerForWrite(
		ctx,
		func() error {
			return v.dbVisibilityManager.RecordWorkflowExecutionClosed(ctx, request)
		},
		func() error {
			return v.sourceVisibilityManager.RecordWorkflowExecutionClosed(ctx, request)
		},
		func() error {
			return v.destinationVisibilityManager.RecordWorkflowExecutionClosed(ctx, request)
		},
	)
}

func (v *visibilityTripleManager) RecordWorkflowExecutionUninitialized(
	ctx context.Context,
	request *RecordWorkflowExecutionUninitializedRequest,
) error {
	return v.chooseVisibilityManagerForWrite(
		ctx,
		func() error {
			return v.dbVisibilityManager.RecordWorkflowExecutionUninitialized(ctx, request)
		},
		func() error {
			return v.sourceVisibilityManager.RecordWorkflowExecutionUninitialized(ctx, request)
		},
		func() error {
			return v.destinationVisibilityManager.RecordWorkflowExecutionUninitialized(ctx, request)
		},
	)
}

func (v *visibilityTripleManager) DeleteWorkflowExecution(
	ctx context.Context,
	request *VisibilityDeleteWorkflowExecutionRequest,
) error {
	return v.chooseVisibilityManagerForWrite(
		ctx,
		func() error {
			return v.dbVisibilityManager.DeleteWorkflowExecution(ctx, request)
		},
		func() error {
			return v.sourceVisibilityManager.DeleteWorkflowExecution(ctx, request)
		},
		func() error {
			return v.destinationVisibilityManager.DeleteWorkflowExecution(ctx, request)
		},
	)
}

func (v *visibilityTripleManager) DeleteUninitializedWorkflowExecution(
	ctx context.Context,
	request *VisibilityDeleteWorkflowExecutionRequest,
) error {
	return v.chooseVisibilityManagerForWrite(
		ctx,
		func() error {
			return v.dbVisibilityManager.DeleteUninitializedWorkflowExecution(ctx, request)
		},
		func() error {
			return v.sourceVisibilityManager.DeleteUninitializedWorkflowExecution(ctx, request)
		},
		func() error {
			return v.destinationVisibilityManager.DeleteUninitializedWorkflowExecution(ctx, request)
		},
	)
}

func (v *visibilityTripleManager) UpsertWorkflowExecution(
	ctx context.Context,
	request *UpsertWorkflowExecutionRequest,
) error {
	return v.chooseVisibilityManagerForWrite(
		ctx,
		func() error {
			return v.dbVisibilityManager.UpsertWorkflowExecution(ctx, request)
		},
		func() error {
			return v.sourceVisibilityManager.UpsertWorkflowExecution(ctx, request)
		},
		func() error {
			return v.destinationVisibilityManager.UpsertWorkflowExecution(ctx, request)
		},
	)
}

func (v *visibilityTripleManager) chooseVisibilityModeForAdmin() string {
	switch {
	case v.dbVisibilityManager != nil && v.sourceVisibilityManager != nil && v.destinationVisibilityManager != nil:
		return common.AdvancedVisibilityMigrationWritingModeTriple
	case v.sourceVisibilityManager != nil && v.destinationVisibilityManager != nil:
		return common.AdvancedVisibilityMigrationWritingModeDual
	case v.destinationVisibilityManager != nil:
		return common.AdvancedVisibilityMigrationWritingModeDestination
	case v.dbVisibilityManager != nil:
		return common.AdvancedVisibilityMigrationWritingModeOff
	default:
		return "INVALID_ADMIN_MODE"
	}
}

func (v *visibilityTripleManager) chooseVisibilityManagerForWrite(ctx context.Context, dbVisFunc, sourceVisFunc, destinationVisFunc func() error) error {
	var writeMode string
	if v.writeMode != nil {
		writeMode = v.writeMode()
	} else {
		key := VisibilityAdminDeletionKey("visibilityAdminDelete")
		if value := ctx.Value(key); value != nil && value.(bool) {
			writeMode = v.chooseVisibilityModeForAdmin()
		}
	}

	switch writeMode {
	// only perform as triple manager during migration by setting write mode to triple,
	// other time perform as a dual visibility manager and db
	case common.AdvancedVisibilityMigrationWritingModeOff:
		if v.dbVisibilityManager != nil {
			return dbVisFunc()
		}
		v.logger.Warn("basic visibility is not available to write, fall back to advanced visibility")
		return destinationVisFunc()
	case common.AdvancedVisibilityMigrationWritingModeSource:
		if v.sourceVisibilityManager != nil {
			return sourceVisFunc()
		}
		v.logger.Warn("advanced visibility is not available to write, fall back to basic visibility")
		return dbVisFunc()
	case common.AdvancedVisibilityMigrationWritingModeDestination:
		if v.destinationVisibilityManager != nil {
			return destinationVisFunc()
		}
		v.logger.Warn("advanced visibility is not available to write, fall back to basic visibility")
		return dbVisFunc()
	case common.AdvancedVisibilityWritingModeDual:
		// by default the AdvancedVisibilityMigrationWritingMode is set to Dual
		// this will enable double writes to both source and destination advanced visibility stores
		// it is necessary during the migration period to ensure that both stores are in sync
		if v.destinationVisibilityManager != nil && v.sourceVisibilityManager != nil {
			if err := sourceVisFunc(); err != nil {
				return err
			}
			return destinationVisFunc()
		} else if v.destinationVisibilityManager != nil {
			v.logger.Warn("Source advanced visibility is not available to write, fall back to destination advanced visibility")
			return destinationVisFunc()
		} else if v.sourceVisibilityManager != nil {
			v.logger.Warn("Destination advanced visibility is not available to write, fall back to source advanced visibility")
			return sourceVisFunc()
		} else {
			v.logger.Warn("advanced visibility is not available to write, fall back to basic visibility")
			return dbVisFunc()
		}
	case common.AdvancedVisibilityMigrationWritingModeTriple:
		if v.destinationVisibilityManager != nil && v.sourceVisibilityManager != nil && v.dbVisibilityManager != nil {
			if err := sourceVisFunc(); err != nil {
				return err
			}

			if err := destinationVisFunc(); err != nil {
				return err
			}
			return dbVisFunc()
		}
		v.logger.Warn("advanced visibility is not available to write")
		return dbVisFunc()
	default:
		return &types.InternalServiceError{
			Message: fmt.Sprintf("Unknown visibility migration writing mode: %s", writeMode),
		}
	}
}

// For Pinot Migration uses. It will be a temporary usage
type userParameters struct {
	operation    string
	domainName   string
	workflowType string
	workflowID   string
	closeStatus  int // if it is -1, then will have --open flag in comparator workflow
	customQuery  string
	earliestTime int64
	latestTime   int64
}

// For Visibility Migration uses. It will be a temporary usage
// logUserQueryParameters will log user queries' parameters so that a comparator workflow can consume
func (v *visibilityTripleManager) logUserQueryParameters(userParam userParameters, domain string, override bool) {
	// Don't log if it is not enabled
	// don't log if it is a call from Pinot Response Comparator workflow
	if !v.logCustomerQueryParameter(domain) || override {
		return
	}

	randNum := rand.Intn(10)
	if randNum != 5 { // Intentionally to have 1/10 chance to log custom query parameters
		return
	}

	v.logger.Info("Logging user query parameters for visibility migration response comparator...",
		tag.OperationName(userParam.operation),
		tag.WorkflowDomainName(userParam.domainName),
		tag.WorkflowType(userParam.workflowType),
		tag.WorkflowID(userParam.workflowID),
		tag.WorkflowCloseStatus(userParam.closeStatus),
		tag.VisibilityQuery(filterAttrPrefix(userParam.customQuery)),
		tag.EarliestTime(userParam.earliestTime),
		tag.LatestTime(userParam.latestTime))

}

// This is for only logUserQueryParameters (for Pinot Response comparator) usage.
// Be careful because there's a low possibility that there'll be false positive cases (shown in unit tests)
func filterAttrPrefix(str string) string {
	str = strings.Replace(str, "`Attr.", "", -1)
	return strings.Replace(str, "`", "", -1)
}

func (v *visibilityTripleManager) getShadowMgrForDoubleRead(domain string) VisibilityManager {
	// invalid cases:
	// case0: when it is not double read
	if !v.readModeIsDouble(domain) {
		return nil
	}
	// case1: when it is double read, and both advanced visibility are not available
	// case2: when it is double read, and only one of advanced visibility is available
	if v.destinationVisibilityManager == nil || v.sourceVisibilityManager == nil {
		return nil
	}

	// Valid cases:
	// case3: when it is double read, and both advanced visibility are available, and read mode is from destination
	if v.readModeIsFromDestination(domain) {
		return v.sourceVisibilityManager
	}
	// case4: when it is double read, and both advanced visibility are available, and read mode is from source
	if v.readModeIsFromSource(domain) {
		return v.destinationVisibilityManager
	}
	// exclude all other cases
	return nil
}

func (v *visibilityTripleManager) ListOpenWorkflowExecutions(
	ctx context.Context,
	request *ListWorkflowExecutionsRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		closeStatus:  -1, // is open. Will have --open flag in comparator workflow
		earliestTime: request.EarliestTime,
		latestTime:   request.LatestTime,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ListOpenWorkflowExecutions, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	// return result from primary
	return manager.ListOpenWorkflowExecutions(ctx, request)
}

func (v *visibilityTripleManager) ListClosedWorkflowExecutions(
	ctx context.Context,
	request *ListWorkflowExecutionsRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		closeStatus:  6, // 6 means not set closeStatus.
		earliestTime: request.EarliestTime,
		latestTime:   request.LatestTime,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ListClosedWorkflowExecutions, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.ListClosedWorkflowExecutions(ctx, request)
}

func (v *visibilityTripleManager) ListOpenWorkflowExecutionsByType(
	ctx context.Context,
	request *ListWorkflowExecutionsByTypeRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		workflowType: request.WorkflowTypeName,
		closeStatus:  -1, // is open. Will have --open flag in comparator workflow
		earliestTime: request.EarliestTime,
		latestTime:   request.LatestTime,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ListOpenWorkflowExecutionsByType, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.ListOpenWorkflowExecutionsByType(ctx, request)
}

func (v *visibilityTripleManager) ListClosedWorkflowExecutionsByType(
	ctx context.Context,
	request *ListWorkflowExecutionsByTypeRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		workflowType: request.WorkflowTypeName,
		closeStatus:  6, // 6 means not set closeStatus.
		earliestTime: request.EarliestTime,
		latestTime:   request.LatestTime,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ListClosedWorkflowExecutionsByType, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.ListClosedWorkflowExecutionsByType(ctx, request)
}

func (v *visibilityTripleManager) ListOpenWorkflowExecutionsByWorkflowID(
	ctx context.Context,
	request *ListWorkflowExecutionsByWorkflowIDRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		workflowID:   request.WorkflowID,
		closeStatus:  -1,
		earliestTime: request.EarliestTime,
		latestTime:   request.LatestTime,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ListOpenWorkflowExecutionsByWorkflowID, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.ListOpenWorkflowExecutionsByWorkflowID(ctx, request)
}

func (v *visibilityTripleManager) ListClosedWorkflowExecutionsByWorkflowID(
	ctx context.Context,
	request *ListWorkflowExecutionsByWorkflowIDRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		workflowID:   request.WorkflowID,
		closeStatus:  6, // 6 means not set closeStatus.
		earliestTime: request.EarliestTime,
		latestTime:   request.LatestTime,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ListClosedWorkflowExecutionsByWorkflowID, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.ListClosedWorkflowExecutionsByWorkflowID(ctx, request)
}

func (v *visibilityTripleManager) ListClosedWorkflowExecutionsByStatus(
	ctx context.Context,
	request *ListClosedWorkflowExecutionsByStatusRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		closeStatus:  int(request.Status),
		earliestTime: request.EarliestTime,
		latestTime:   request.LatestTime,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ListClosedWorkflowExecutionsByStatus, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.ListClosedWorkflowExecutionsByStatus(ctx, request)
}

func (v *visibilityTripleManager) GetClosedWorkflowExecution(
	ctx context.Context,
	request *GetClosedWorkflowExecutionRequest,
) (*GetClosedWorkflowExecutionResponse, error) {
	earlistTime := int64(0) // this is to get all closed workflow execution
	latestTime := time.Now().UnixNano()

	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		closeStatus:  6, // 6 means not set closeStatus.
		earliestTime: earlistTime,
		latestTime:   latestTime,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.GetClosedWorkflowExecution, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.GetClosedWorkflowExecution(ctx, request)
}

func (v *visibilityTripleManager) ListWorkflowExecutions(
	ctx context.Context,
	request *ListWorkflowExecutionsByQueryRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		closeStatus:  6, // 6 means not set closeStatus.
		customQuery:  request.Query,
		earliestTime: -1,
		latestTime:   -1,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ListWorkflowExecutions, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.ListWorkflowExecutions(ctx, request)
}

func (v *visibilityTripleManager) ScanWorkflowExecutions(
	ctx context.Context,
	request *ListWorkflowExecutionsByQueryRequest,
) (*ListWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.LIST),
		domainName:   request.Domain,
		closeStatus:  6, // 6 means not set closeStatus.
		customQuery:  request.Query,
		earliestTime: -1,
		latestTime:   -1,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.ScanWorkflowExecutions, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.ScanWorkflowExecutions(ctx, request)
}

func (v *visibilityTripleManager) CountWorkflowExecutions(
	ctx context.Context,
	request *CountWorkflowExecutionsRequest,
) (*CountWorkflowExecutionsResponse, error) {
	override := ctx.Value(ContextKey)
	v.logUserQueryParameters(userParameters{
		operation:    string(Operation.COUNT),
		domainName:   request.Domain,
		closeStatus:  6, // 6 means not set closeStatus.
		customQuery:  request.Query,
		earliestTime: -1,
		latestTime:   -1,
	}, request.Domain, override != nil)

	// get another manager for double read
	shadowMgr := v.getShadowMgrForDoubleRead(request.Domain)
	// call the API for latency comparison
	if shadowMgr != nil {
		go shadow(shadowMgr.CountWorkflowExecutions, request, v.logger)
	}

	manager := v.chooseVisibilityManagerForRead(ctx, request.Domain)
	return manager.CountWorkflowExecutions(ctx, request)
}

func (v *visibilityTripleManager) chooseVisibilityManagerForRead(ctx context.Context, domain string) VisibilityManager {
	if override := ctx.Value(ContextKey); override == VisibilityOverridePrimary {
		v.logger.Info("Visibility Migration log: Primary visibility manager was chosen for read.")
		return v.sourceVisibilityManager
	} else if override == VisibilityOverrideSecondary {
		v.logger.Info("Visibility Migration log: Secondary visibility manager was chosen for read.")
		return v.destinationVisibilityManager
	}

	var visibilityMgr VisibilityManager
	if v.readModeIsFromSource(domain) {
		if v.sourceVisibilityManager != nil {
			visibilityMgr = v.sourceVisibilityManager
		} else {
			visibilityMgr = v.dbVisibilityManager
			v.logger.Warn("domain is configured to read from advanced visibility but it's not available, fall back to basic visibility",
				tag.WorkflowDomainName(domain))
		}
	} else if v.readModeIsFromDestination(domain) {
		if v.destinationVisibilityManager != nil {
			visibilityMgr = v.destinationVisibilityManager
		} else {
			visibilityMgr = v.dbVisibilityManager
			v.logger.Warn("domain is configured to read from advanced visibility but it's not available, fall back to basic visibility",
				tag.WorkflowDomainName(domain))
		}
	} else {
		if v.dbVisibilityManager != nil {
			visibilityMgr = v.dbVisibilityManager
		} else {
			visibilityMgr = v.destinationVisibilityManager
			v.logger.Warn("domain is configured to read from basic visibility but it's not available, fall back to advanced visibility",
				tag.WorkflowDomainName(domain))
		}
	}
	return visibilityMgr
}

func shadow[ReqT any, ResT any](f func(ctx context.Context, request ReqT) (ResT, error), request ReqT, logger log.Logger) {
	ctxNew, cancel := context.WithTimeout(context.Background(), 2*time.Minute) // don't want f to run too long

	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			logger.Info(fmt.Sprintf("Recovered in Shadow function in double read: %v", r))
		}
	}()

	_, err := f(ctxNew, request)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in Shadow function in double read: %s", err.Error()))
	}
}
