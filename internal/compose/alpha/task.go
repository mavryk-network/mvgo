// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc, abdul@blockwatch.cc

package alpha

import (
	"fmt"

	"github.com/mavryk-network/tzgo/codec"
	"github.com/mavryk-network/tzgo/internal/compose"
	"github.com/mavryk-network/tzgo/rpc"
)

type TaskBuilder interface {
	Type() string
	Validate(compose.Context, Task) error
	Build(compose.Context, Task) (*codec.Op, *rpc.CallOptions, error)
}

type TaskFactory func() TaskBuilder

var (
	taskRegistry map[string]TaskFactory = make(map[string]TaskFactory)
)

func RegisterTask(typ string, fn TaskFactory) {
	taskRegistry[typ] = fn
}

func NewTask(typ string) (TaskBuilder, error) {
	fn, ok := taskRegistry[typ]
	if !ok {
		return nil, fmt.Errorf("unsupported task type %s", typ)
	}
	return fn(), nil
}
