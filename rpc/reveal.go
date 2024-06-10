// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import "github.com/mavryk-network/mvgo/mavryk"

// Ensure Reveal implements the TypedOperation interface.
var _ TypedOperation = (*Reveal)(nil)

// Reveal represents a reveal operation
type Reveal struct {
	Manager
	PublicKey mavryk.Key `json:"public_key"`
}

// Costs returns operation cost to implement TypedOperation interface.
func (r Reveal) Costs() mavryk.Costs {
	return mavryk.Costs{
		Fee:     r.Manager.Fee,
		GasUsed: r.Metadata.Result.Gas(),
	}
}
