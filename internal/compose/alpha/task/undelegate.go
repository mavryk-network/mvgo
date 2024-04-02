// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc, abdul@blockwatch.cc

package task

import (
	"github.com/mavryk-network/mvgo/codec"
	"github.com/mavryk-network/mvgo/internal/compose"
	"github.com/mavryk-network/mvgo/internal/compose/alpha"
	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/rpc"
	"github.com/mavryk-network/mvgo/signer"

	"github.com/pkg/errors"
)

var _ alpha.TaskBuilder = (*UndelegateTask)(nil)

func init() {
	alpha.RegisterTask("undelegate", NewUndelegateTask)
}

type UndelegateTask struct {
	BaseTask
}

func NewUndelegateTask() alpha.TaskBuilder {
	return &UndelegateTask{}
}

func (t *UndelegateTask) Type() string {
	return "undelegate"
}

func (t *UndelegateTask) Build(ctx compose.Context, task alpha.Task) (*codec.Op, *rpc.CallOptions, error) {
	if err := t.parse(ctx, task); err != nil {
		return nil, nil, errors.Wrap(err, "parse")
	}
	opts := rpc.NewCallOptions()
	opts.Signer = signer.NewFromKey(t.Key)
	opts.IgnoreLimits = true
	op := codec.NewOp().
		WithSource(t.Source).
		WithUndelegation().
		WithLimits([]mavryk.Limits{rpc.DefaultDelegationLimitsEOA}, 0)
	return op, opts, nil
}

func (t *UndelegateTask) Validate(ctx compose.Context, task alpha.Task) error {
	return t.parse(ctx, task)
}

func (t *UndelegateTask) parse(ctx compose.Context, task alpha.Task) error {
	return t.BaseTask.parse(ctx, task)
}
