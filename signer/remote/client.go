// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package remote

import (
	"context"
	"net/http"

	"github.com/mavryk-network/mvgo/codec"
	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/rpc"
	"github.com/mavryk-network/mvgo/signer"
)

var _ signer.Signer = (*RemoteSigner)(nil)

type RemoteSigner struct {
	c     *rpc.Client
	addrs []mavryk.Address
	auth  mavryk.PrivateKey
}

// New creates a new remote signer client and initializes it with the remote url.
// Users may pass an optional http client with a custom configuration, otherwise
// the http.DefaultClient is used.
func New(url string, client *http.Client) (*RemoteSigner, error) {
	c, err := rpc.NewClient(url, client)
	if err != nil {
		return nil, err
	}
	return &RemoteSigner{c: c}, nil
}

func (s *RemoteSigner) WithAddress(addr mavryk.Address) *RemoteSigner {
	s.addrs = append(s.addrs, addr)
	return s
}

func (s *RemoteSigner) WithAuthKey(sk mavryk.PrivateKey) *RemoteSigner {
	s.auth = sk
	return s
}

// AuthorizedKeys returns a list of addresses the remote signer accepts for
// authenticating requests.
func (s RemoteSigner) AuthorizedKeys(ctx context.Context) ([]mavryk.Address, error) {
	type response struct {
		Addrs []mavryk.Address `json:"authorized_keys"`
	}
	var resp response
	err := s.c.Get(ctx, "/authorized_keys", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Addrs, nil
}

// ListAddresses returns a list of addresses the remote signer can produce signatures for.
func (s RemoteSigner) ListAddresses(ctx context.Context) ([]mavryk.Address, error) {
	return s.addrs, nil
}

// GetKey returns the public key associated with address.
func (s RemoteSigner) GetKey(ctx context.Context, address mavryk.Address) (mavryk.Key, error) {
	type response struct {
		Pk mavryk.Key `json:"public_key"`
	}
	var resp response
	err := s.c.Get(ctx, "/keys/"+address.String(), &resp)
	return resp.Pk, err
}

// SignMessage signs msg for address by wrapping it into a failing noop operation
// with zero branch hash. This prevents unintended signature of message bytes that
// represent a valid transaction.
//
// Note that most remote signers for Tezos do not support signing of operation kinds other
// than baking related operations.
func (s RemoteSigner) SignMessage(ctx context.Context, address mavryk.Address, msg string) (mavryk.Signature, error) {
	op := codec.NewOp().
		WithBranch(mavryk.ZeroBlockHash).
		WithContents(&codec.FailingNoop{
			Arbitrary: msg,
		})
	return s.SignOperation(ctx, address, op)
}

// SignOperation signs operation op for address using the configured remote signer's
// REST API. For endorsements this call requires branch_id to be present.
//
// Note that most remote signers for Tezos do not support signing of operation kinds other
// than baking related operations.
func (s RemoteSigner) SignOperation(ctx context.Context, address mavryk.Address, op *codec.Op) (mavryk.Signature, error) {
	type response struct {
		Sig mavryk.Signature `json:"signature"`
	}
	var resp response
	err := s.c.Post(ctx, "/keys/"+address.String(), mavryk.HexBytes(op.WatermarkedBytes()), &resp)
	return resp.Sig, err
}

// SignOperation signs a block header for address using the configured remote signer's
// REST API. This call requires branch_id to be present.
func (s RemoteSigner) SignBlock(ctx context.Context, address mavryk.Address, head *codec.BlockHeader) (mavryk.Signature, error) {
	type response struct {
		Sig mavryk.Signature `json:"signature"`
	}
	var resp response
	err := s.c.Post(ctx, "/keys/"+address.String(), mavryk.HexBytes(head.WatermarkedBytes()), &resp)
	return resp.Sig, err
}
