// Copyright (c) 2020-2021 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/micheline"
)

// Ensure Origination implements the TypedOperation interface.
var _ TypedOperation = (*Origination)(nil)

// Origination represents a contract creation operation
type Origination struct {
	Manager
	ManagerPubkey  mavryk.Address    `json:"manager_pubkey"` // proto v1 & >= v4
	ManagerPubkey2 mavryk.Address    `json:"managerPubkey"`  // proto v2, v3
	Balance        int64             `json:"balance,string"`
	Spendable      *bool             `json:"spendable"`   // true when missing before v5 Babylon
	Delegatable    *bool             `json:"delegatable"` // true when missing before v5 Babylon
	Delegate       *mavryk.Address   `json:"delegate"`
	Script         *micheline.Script `json:"script"`
}

func (o Origination) ManagerAddress() mavryk.Address {
	if o.ManagerPubkey2.IsValid() {
		return o.ManagerPubkey2
	}
	return o.ManagerPubkey
}

// Costs returns operation cost to implement TypedOperation interface.
func (o Origination) Costs() mavryk.Costs {
	res := o.Metadata.Result
	cost := mavryk.Costs{
		Fee:         o.Manager.Fee,
		GasUsed:     res.Gas(),
		StorageUsed: res.PaidStorageSizeDiff,
	}
	var i int
	for _, v := range res.BalanceUpdates {
		if v.Kind != CONTRACT {
			continue
		}
		if res.PaidStorageSizeDiff > 0 && i == 0 {
			burn := v.Amount()
			cost.StorageBurn += -burn
			cost.Burn += -burn
			i++
			continue
		}
		if len(res.OriginatedContracts) > 0 && i == 1 {
			burn := v.Amount()
			cost.AllocationBurn += -burn
			cost.Burn += -burn
			i++
		}
	}
	return cost
}
