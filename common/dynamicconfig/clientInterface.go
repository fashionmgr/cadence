// Copyright (c) 2017 Uber Technologies, Inc.
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

//go:generate mockgen -package $GOPACKAGE -source $GOFILE -destination clientInterface_mock.go -self_package github.com/uber/cadence/common/dynamicconfig

package dynamicconfig

import (
	"time"

	"github.com/uber/cadence/common/types"
)

const (
	ConfigStoreClient = "configstore"
	FileBasedClient   = "filebased"
	InMemoryClient    = "memory"
	NopClient         = "nop"
)

// Client allows fetching values from a dynamic configuration system NOTE: This does not have async
// options right now. In the interest of keeping it minimal, we can add when requirement arises.
type Client interface {
	GetValue(name Key) (interface{}, error)
	GetValueWithFilters(name Key, filters map[Filter]interface{}) (interface{}, error)

	GetIntValue(name IntKey, filters map[Filter]interface{}) (int, error)
	GetFloatValue(name FloatKey, filters map[Filter]interface{}) (float64, error)
	GetBoolValue(name BoolKey, filters map[Filter]interface{}) (bool, error)
	GetStringValue(name StringKey, filters map[Filter]interface{}) (string, error)
	GetMapValue(name MapKey, filters map[Filter]interface{}) (map[string]interface{}, error)
	GetDurationValue(name DurationKey, filters map[Filter]interface{}) (time.Duration, error)
	GetListValue(name ListKey, filters map[Filter]interface{}) ([]interface{}, error)
	// UpdateValue takes value as map and updates by overriding. It doesn't support update with filters.
	UpdateValue(name Key, value interface{}) error
	RestoreValue(name Key, filters map[Filter]interface{}) error
	ListValue(name Key) ([]*types.DynamicConfigEntry, error)
}

var NotFoundError = &types.EntityNotExistsError{
	Message: "unable to find key",
}
