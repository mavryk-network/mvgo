// Copyright (c) 2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import "github.com/mavryk-network/mvgo/mavryk"

// Ensure DrainDelegate implements the TypedOperation interface.
var _ TypedOperation = (*DrainDelegate)(nil)

// DrainDelegate represents a transaction operation
type DrainDelegate struct {
	Generic
	ConsensusKey mavryk.Address `json:"consensus_key"`
	Delegate     mavryk.Address `json:"delegate"`
	Destination  mavryk.Address `json:"destination"`
}
