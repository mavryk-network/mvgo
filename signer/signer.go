// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package signer

import (
	"context"

	"github.com/mavryk-network/mvgo/codec"
	"github.com/mavryk-network/mvgo/mavryk"
)

type Signer interface {
	// Return a list of addresses the signer manages.
	ListAddresses(context.Context) ([]mavryk.Address, error)

	// Returns the public key for a managed address. Required for reveal ops.
	GetKey(context.Context, mavryk.Address) (mavryk.Key, error)

	// Sign an arbitrary text message wrapped into a failing noop
	SignMessage(context.Context, mavryk.Address, string) (mavryk.Signature, error)

	// Sign an operation.
	SignOperation(context.Context, mavryk.Address, *codec.Op) (mavryk.Signature, error)

	// Sign a block header.
	SignBlock(context.Context, mavryk.Address, *codec.BlockHeader) (mavryk.Signature, error)
}
