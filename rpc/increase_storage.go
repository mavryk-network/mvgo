// Copyright (c) 2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import "github.com/mavryk-network/mvgo/mavryk"

// Ensure IncreasePaidStorage implements the TypedOperation interface.
var _ TypedOperation = (*IncreasePaidStorage)(nil)

// IncreasePaidStorage represents a transaction operation
type IncreasePaidStorage struct {
	Manager
	Destination mavryk.Address `json:"destination"`
	Amount      int64          `json:"amount,string"`
}

// Costs returns operation cost to implement TypedOperation interface.
func (t IncreasePaidStorage) Costs() mavryk.Costs {
	res := t.Metadata.Result
	cost := mavryk.Costs{
		Fee:     t.Manager.Fee,
		GasUsed: res.Gas(),
	}
	if !t.Result().IsSuccess() {
		return cost
	}
	for _, v := range res.BalanceUpdates {
		if v.Kind != CONTRACT {
			continue
		}
		burn := v.Amount()
		if burn >= 0 {
			continue
		}
		cost.StorageBurn += -burn
	}
	return cost
}
