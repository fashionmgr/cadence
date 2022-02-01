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

package history

import (
	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/membership"
	"github.com/uber/cadence/common/service"
)

// PeerResolver is used to resolve history peers.
// Those are deployed instances of Cadence history services that participate in the cluster ring.
// The resulting peer is simply an address of form ip:port where RPC calls can be routed to.
type PeerResolver struct {
	numberOfShards int
	resolver       membership.Resolver
	addressMapper  AddressMapperFn
}

type AddressMapperFn func(string) (string, error)

// NewPeerResolver creates a new history peer resolver.
func NewPeerResolver(numberOfShards int, resolver membership.Resolver, addressMapper AddressMapperFn) PeerResolver {
	return PeerResolver{
		numberOfShards: numberOfShards,
		resolver:       resolver,
		addressMapper:  addressMapper,
	}
}

// FromWorkflowID resolves the history peer responsible for a given workflowID.
// WorkflowID is converted to logical shardID using a consistent hash function.
// FromShardID is used for further resolving.
func (pr PeerResolver) FromWorkflowID(workflowID string) (string, error) {
	shardID := common.WorkflowIDToHistoryShard(workflowID, pr.numberOfShards)
	return pr.FromShardID(shardID)
}

// FromDomainID resolves the history peer responsible for a given domainID.
// DomainID is converted to logical shardID using a consistent hash function.
// FromShardID is used for further resolving.
func (pr PeerResolver) FromDomainID(domainID string) (string, error) {
	shardID := common.DomainIDToHistoryShard(domainID, pr.numberOfShards)
	return pr.FromShardID(shardID)
}

// FromShardID resolves the history peer responsible for a given logical shardID.
// It uses our membership provider to lookup which instance currently owns the given shard.
// FromHostAddress is used for further resolving.
func (pr PeerResolver) FromShardID(shardID int) (string, error) {
	shardIDString := string(rune(shardID))
	host, err := pr.resolver.Lookup(service.History, shardIDString)
	if err != nil {
		return "", err
	}
	return pr.FromHostAddress(host.GetAddress())
}

// FromHostAddress resolves the final history peer responsible for the given host address.
// The address may be used as is, or processed with additional address mapper.
// In case of gRPC transport, the port within the address is replaced with gRPC port.
func (pr PeerResolver) FromHostAddress(hostAddress string) (string, error) {
	if pr.addressMapper == nil {
		return hostAddress, nil
	}
	return pr.addressMapper(hostAddress)
}
