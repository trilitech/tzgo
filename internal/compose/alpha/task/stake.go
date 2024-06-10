// Copyright (c) 2024 Blockwatch Data Inc.
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

var _ alpha.TaskBuilder = (*StakeTask)(nil)

func init() {
	alpha.RegisterTask("stake", NewStakeTask)
}

type StakeTask struct {
	BaseTask
	Amount int64
}

func NewStakeTask() alpha.TaskBuilder {
	return &StakeTask{}
}

func (t *StakeTask) Type() string {
	return "stake"
}

func (t *StakeTask) Build(ctx compose.Context, task alpha.Task) (*codec.Op, *rpc.CallOptions, error) {
	if err := t.parse(ctx, task); err != nil {
		return nil, nil, errors.Wrap(err, "parse")
	}
	opts := rpc.NewCallOptions()
	opts.Signer = signer.NewFromKey(t.Key)
	op := codec.NewOp().
		WithSource(t.Source).
		WithStake(t.Amount).
		WithLimits([]tezos.Limits{rpc.DefaultStakeLimits}, 0)
	return op, opts, nil
}

func (t *StakeTask) Validate(ctx compose.Context, task alpha.Task) error {
	return t.parse(ctx, task)
}

func (t *StakeTask) parse(ctx compose.Context, task alpha.Task) (err error) {
	if err = t.BaseTask.parse(ctx, task); err != nil {
		return err
	}
	t.Amount = int64(task.Amount)
	return
}
