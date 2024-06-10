// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package mavryk

// Limits represents all resource limits defined for an operation in Tezos.
type Limits struct {
	Fee          int64
	GasLimit     int64
	StorageLimit int64
}

// Add adds two limits z = x + y and returns the sum z without changing any of the inputs.
func (x Limits) Add(y Limits) Limits {
	x.Fee += y.Fee
	x.GasLimit += y.GasLimit
	x.StorageLimit += y.StorageLimit
	return x
}

// Costs represents all costs paid by an operation in Tezos. Its contents depends on
// operation type and activity. Consensus and voting operations have no cost,
// user operations have variable cost. For transactions with internal results costs
// are a summary.
type Costs struct {
	Fee            int64 // the total fee paid in mumav
	Burn           int64 // total amount of mumav burned (not included in fee)
	GasUsed        int64 // gas used
	StorageUsed    int64 // new storage bytes allocated
	StorageBurn    int64 // mumav burned for allocating new storage (not included in fee)
	AllocationBurn int64 // mumav burned for allocating a new account (not included in fee)
}

// Add adds two costs z = x + y and returns the sum z without changing any of the inputs.
func (x Costs) Add(y Costs) Costs {
	x.Fee += y.Fee
	x.Burn += y.Burn
	x.GasUsed += y.GasUsed
	x.StorageUsed += y.StorageUsed
	x.StorageBurn += y.StorageBurn
	x.AllocationBurn += y.AllocationBurn
	return x
}
