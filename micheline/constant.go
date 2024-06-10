// Copyright (c) 2020-2021 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package micheline

import "github.com/mavryk-network/mvgo/mavryk"

type ConstantDict map[string]Prim

func (d *ConstantDict) Add(address mavryk.ExprHash, value Prim) {
	if *d == nil {
		*d = make(ConstantDict)
	}
	(*d)[address.String()] = value
}

func (d ConstantDict) Has(address mavryk.ExprHash) bool {
	if d == nil {
		return false
	}
	_, ok := d[address.String()]
	return ok
}

func (d ConstantDict) Get(address mavryk.ExprHash) (Prim, bool) {
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

func (p Prim) Constants() []mavryk.ExprHash {
	c := make([]mavryk.ExprHash, 0)
	p.Walk(func(p Prim) error {
		if p.IsConstant() {
			if h, err := mavryk.ParseExprHash(p.Args[0].String); err == nil {
				c = append(c, h)
			}
		}
		return nil
	})
	return c
}
