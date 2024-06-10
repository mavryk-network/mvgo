// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package signer

import (
	"context"
	"errors"

	"github.com/mavryk-network/mvgo/codec"
	"github.com/mavryk-network/mvgo/mavryk"
)

var ErrAddressMismatch = errors.New("signer: address mismatch")

type MemorySigner struct {
	key mavryk.PrivateKey
}

func NewFromKey(k mavryk.PrivateKey) *MemorySigner {
	return &MemorySigner{
		key: k,
	}
}

func (s MemorySigner) ListAddresses(_ context.Context) ([]mavryk.Address, error) {
	return []mavryk.Address{s.key.Address()}, nil
}

func (s MemorySigner) GetKey(_ context.Context, addr mavryk.Address) (mavryk.Key, error) {
	pk := s.key.Public()
	if !pk.Address().Equal(addr) {
		return mavryk.InvalidKey, ErrAddressMismatch
	}
	return pk, nil
}

func (s MemorySigner) SignMessage(_ context.Context, addr mavryk.Address, msg string) (mavryk.Signature, error) {
	if !s.key.Address().Equal(addr) {
		return mavryk.InvalidSignature, ErrAddressMismatch
	}
	op := codec.NewOp().
		WithBranch(mavryk.ZeroBlockHash).
		WithContents(&codec.FailingNoop{
			Arbitrary: msg,
		})
	digest := mavryk.Digest(op.Bytes())
	return s.key.Sign(digest[:])
}

func (s MemorySigner) SignOperation(_ context.Context, addr mavryk.Address, op *codec.Op) (mavryk.Signature, error) {
	if !s.key.Address().Equal(addr) {
		return mavryk.InvalidSignature, ErrAddressMismatch
	}
	err := op.Sign(s.key)
	return op.Signature, err
}

func (s MemorySigner) SignBlock(_ context.Context, addr mavryk.Address, head *codec.BlockHeader) (mavryk.Signature, error) {
	if !s.key.Address().Equal(addr) {
		return mavryk.InvalidSignature, ErrAddressMismatch
	}
	err := head.Sign(s.key)
	return head.Signature, err
}
