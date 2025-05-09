// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc, abdul@blockwatch.cc

package task

import (
	"github.com/trilitech/tzgo/codec"
	"github.com/trilitech/tzgo/internal/compose"
	"github.com/trilitech/tzgo/internal/compose/alpha"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/signer"
	"github.com/trilitech/tzgo/tezos"

	"github.com/pkg/errors"
)

var _ alpha.TaskBuilder = (*SetBakerParamsTask)(nil)

func init() {
	alpha.RegisterTask("set_baker_params", NewSetBakerParamsTask)
}

type SetBakerParamsTask struct {
	BaseTask
	Limit int64 // the parameter `limit_of_staking_over_baking_millionth`
	Edge  int64 // the parameter `edge_of_baking_over_staking_billionth`
}

func NewSetBakerParamsTask() alpha.TaskBuilder {
	return &SetBakerParamsTask{}
}

func (t *SetBakerParamsTask) Type() string {
	return "set_baker_params"
}

func (t *SetBakerParamsTask) Build(ctx compose.Context, task alpha.Task) (*codec.Op, *rpc.CallOptions, error) {
	if err := t.parse(ctx, task); err != nil {
		return nil, nil, errors.Wrap(err, "parse")
	}
	opts := rpc.NewCallOptions()
	opts.Signer = signer.NewFromKey(t.Key)
	opts.IgnoreLimits = true
	op := codec.NewOp().
		WithSource(t.Source).
		WithSetBakerParams(t.Edge, t.Limit).
		WithLimits([]tezos.Limits{rpc.DefaultBakerParamUpdateLimits}, 0)
	return op, opts, nil
}

func (t *SetBakerParamsTask) Validate(ctx compose.Context, task alpha.Task) error {
	return t.parse(ctx, task)
}

func (t *SetBakerParamsTask) parse(ctx compose.Context, task alpha.Task) error {
	if err := t.BaseTask.parse(ctx, task); err != nil {
		return err
	}
	t.Edge = int64(task.Edge)
	t.Limit = int64(task.Limit)
	return nil
}
