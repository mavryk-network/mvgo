// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc, abdul@blockwatch.cc

package task

import (
	"github.com/mavryk-network/tzgo/codec"
	"github.com/mavryk-network/tzgo/internal/compose"
	"github.com/mavryk-network/tzgo/internal/compose/alpha"
	"github.com/mavryk-network/tzgo/rpc"
	"github.com/mavryk-network/tzgo/signer"
	"github.com/mavryk-network/tzgo/tezos"

	"github.com/pkg/errors"
)

var _ alpha.TaskBuilder = (*RegisterBakerTask)(nil)

func init() {
	alpha.RegisterTask("register_baker", NewRegisterBakerTask)
}

type RegisterBakerTask struct {
	BaseTask
}

func NewRegisterBakerTask() alpha.TaskBuilder {
	return &RegisterBakerTask{}
}

func (t *RegisterBakerTask) Type() string {
	return "register_baker"
}

func (t *RegisterBakerTask) Build(ctx compose.Context, task alpha.Task) (*codec.Op, *rpc.CallOptions, error) {
	if err := t.parse(ctx, task); err != nil {
		return nil, nil, errors.Wrap(err, "parse")
	}
	opts := rpc.NewCallOptions()
	opts.Signer = signer.NewFromKey(t.Key)
	opts.IgnoreLimits = true
	op := codec.NewOp().
		WithSource(t.Source).
		WithRegisterBaker().
		WithLimits([]tezos.Limits{rpc.DefaultBakerRegistrationLimits}, 0)
	return op, opts, nil
}

func (t *RegisterBakerTask) Validate(ctx compose.Context, task alpha.Task) error {
	return t.parse(ctx, task)
}

func (t *RegisterBakerTask) parse(ctx compose.Context, task alpha.Task) error {
	return t.BaseTask.parse(ctx, task)
}
