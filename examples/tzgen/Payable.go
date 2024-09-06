// Code generated by tzgen - DO NOT EDIT.
// This file is a binding generated from Payable smart contract at address KT1FPA7vN4cBk24df7VxUu9DstRcc7am3qnf.
// Any manual changes will be lost.

package main

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/trilitech/tzgo/contract"
	"github.com/trilitech/tzgo/contract/bind"
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

// Payable is a generated binding to a Tezos smart contract.
type Payable struct {
	bind.Contract
	builder PayableBuilder
	rpc     bind.RPC
	script  *micheline.Script
}

// PayableSession is a generated binding to a Tezos smart contract, that will
// use Opts for every call.
type PayableSession struct {
	*Payable
	Opts *rpc.CallOptions
}

// PayableBuilder is a generated struct that builds micheline.Parameters from
// go types.
type PayableBuilder struct{}

// NewPayable creates a new Payable handle, bound to the provided address
// with the given rpc.
//
// Returns an error if the contract was not found at the given address.
func NewPayable(ctx context.Context, address tezos.Address, client *rpc.Client) (*Payable, error) {
	script, err := client.GetContractScript(ctx, address)
	if err != nil {
		return nil, err
	}

	return &Payable{
		Contract: contract.NewContract(address, client),
		rpc:      client,
		script:   script,
	}, nil
}

// Session returns a new PayableSession with the configured rpc.CallOptions.
func (_p *Payable) Session(opts *rpc.CallOptions) *PayableSession {
	return &PayableSession{Payable: _p, Opts: opts}
}

// Builder returns the builder struct for this contract.
func (_p *Payable) Builder() PayableBuilder {
	return _p.builder
}

// Storage queries the current storage of the contract.
func (_p *Payable) Storage(ctx context.Context) (*big.Int, error) {
	return _p.StorageAt(ctx, rpc.Head)
}

// StorageAt queries the contract's storage at the given block.
func (_p *Payable) StorageAt(ctx context.Context, block rpc.BlockID) (*big.Int, error) {
	var storage *big.Int
	prim, err := _p.rpc.GetContractStorage(ctx, _p.Contract.Address(), block)
	if err != nil {
		return storage, errors.Wrap(err, "failed to get storage")
	}

	err = bind.UnmarshalPrim(prim, &storage)
	if err != nil {
		return storage, errors.Wrap(err, "failed to unmarshal storage")
	}
	return storage, nil
}

// DeployPayable deploys a Payable contract by using client and opts, and PayableMicheline.
//
// Returns the receipt and a handle to the Payable deployed contract.
func DeployPayable(ctx context.Context, opts *rpc.CallOptions, client *rpc.Client, storage *big.Int) (*rpc.Receipt, *Payable, error) {
	var script *micheline.Script
	err := json.Unmarshal([]byte(PayableMicheline), &script)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshall contract's script")
	}

	prim, err := bind.MarshalPrim(storage, false)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to marshal storage")
	}
	script.Storage = prim

	c := contract.NewEmptyContract(client).WithScript(script)
	receipt, err := c.Deploy(ctx, opts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to deploy contract")
	}
	return receipt, &Payable{Contract: c, rpc: client}, nil
}

// region Entrypoints

// SendTz makes a call to the `send_tz` contract entry.
//
// send_tz()
func (_p *Payable) SendTz(ctx context.Context, opts *rpc.CallOptions) (*rpc.Receipt, error) {
	params, err := _p.builder.SendTz()
	if err != nil {
		return nil, err
	}
	return _p.Contract.Call(ctx, &contract.TxArgs{Params: params}, opts)
}

// SendTz makes a call to the `send_tz` contract entry.
//
// send_tz()
func (_p *PayableSession) SendTz(ctx context.Context) (*rpc.Receipt, error) {
	return _p.Payable.SendTz(ctx, _p.Opts)
}

// SendTz builds `send_tz` contract entry's parameters.
//
// send_tz()
func (PayableBuilder) SendTz() (micheline.Parameters, error) {
	prim, err := bind.MarshalParams(false)
	if err != nil {
		return micheline.Parameters{}, errors.Wrap(err, "failed to marshal params")
	}
	return micheline.Parameters{Entrypoint: PayableSendTzEntry, Value: prim}, nil
}

// endregion

// Payable entry names
const (
	PayableSendTzEntry = "send_tz"
)

const PayableMicheline = `{"code":[{"prim":"parameter","args":[{"prim":"unit","annots":["%send_tz"]}]},{"prim":"storage","args":[{"prim":"mutez"}]},{"prim":"code","args":[[{"prim":"UNPAIR"},{"prim":"DROP"},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"AMOUNT"},{"prim":"COMPARE"},{"prim":"GT"},{"prim":"NOT"},{"prim":"IF","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"r1"}]},{"prim":"PUSH","args":[{"prim":"string"},{"string":"INVALID_CONDITION"}]},{"prim":"PAIR"},{"prim":"FAILWITH"}],[]]},{"prim":"AMOUNT"},{"prim":"SWAP"},{"prim":"DROP"},{"prim":"NIL","args":[{"prim":"operation"}]},{"prim":"PAIR"}]]}],"storage":{}}`

var (
	_ = big.NewInt
	_ = micheline.NewPrim
	_ = bind.MarshalParams
	_ = time.Now
)
