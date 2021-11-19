// Copyright (c) 2020-2021 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package micheline

import (
	"blockwatch.cc/tzgo/tezos"
)

type ConstantDict map[string]Prim

func (d *ConstantDict) Add(address tezos.ExprHash, value Prim) {
	if *d == nil {
		*d = make(ConstantDict)
	}
	(*d)[address.String()] = value
}

func (d ConstantDict) Has(address tezos.ExprHash) bool {
	if d == nil {
		return false
	}
	_, ok := d[address.String()]
	return ok
}

func (d ConstantDict) Get(address tezos.ExprHash) (Prim, bool) {
	if d == nil {
		return InvalidPrim, false
	}
	p, ok := d[address.String()]
	return p, ok
}

func (d ConstantDict) GetString(address string) (Prim, bool) {
	if d == nil {
		return InvalidPrim, false
	}
	p, ok := d[address]
	return p, ok
}
