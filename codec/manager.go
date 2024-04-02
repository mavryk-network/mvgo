// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// Simple is an empty helper struct that's used to fulfil the Operation interface
// for anonymous, consensus and voting operations which do not contain fees and
// counter.
type Simple struct{}

func (o *Simple) Limits() mavryk.Limits {
	return mavryk.Limits{}
}

func (o Simple) GetCounter() int64 {
	return -1
}

func (o *Simple) WithCounter(int64) {}

func (o *Simple) WithLimits(mavryk.Limits) {}

func (o *Simple) WithSource(mavryk.Address) {}

// Manager contains fields common for all manager operations
type Manager struct {
	Source       mavryk.Address `json:"source"`
	Fee          mavryk.N       `json:"fee"`
	Counter      mavryk.N       `json:"counter"`
	GasLimit     mavryk.N       `json:"gas_limit"`
	StorageLimit mavryk.N       `json:"storage_limit"`
}

func (o Manager) Limits() mavryk.Limits {
	return mavryk.Limits{
		Fee:          o.Fee.Int64(),
		GasLimit:     o.GasLimit.Int64(),
		StorageLimit: o.StorageLimit.Int64(),
	}
}

func (o Manager) EncodeJSON(buf *bytes.Buffer) error {
	buf.WriteString(`"source":`)
	buf.WriteString(strconv.Quote(o.Source.String()))
	buf.WriteString(`,"fee":`)
	buf.WriteString(strconv.Quote(o.Fee.String()))
	buf.WriteString(`,"counter":`)
	buf.WriteString(strconv.Quote(o.Counter.String()))
	buf.WriteString(`,"gas_limit":`)
	buf.WriteString(strconv.Quote(o.GasLimit.String()))
	buf.WriteString(`,"storage_limit":`)
	buf.WriteString(strconv.Quote(o.StorageLimit.String()))
	return nil
}

func (o Manager) EncodeBuffer(buf *bytes.Buffer, _ *mavryk.Params) error {
	buf.Write(o.Source.Encode())
	o.Fee.EncodeBuffer(buf)
	o.Counter.EncodeBuffer(buf)
	o.GasLimit.EncodeBuffer(buf)
	o.StorageLimit.EncodeBuffer(buf)
	return nil
}

func (o *Manager) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = o.Source.Decode(buf.Next(21)); err != nil {
		return
	}
	if err = o.Fee.DecodeBuffer(buf); err != nil {
		return err
	}
	if err = o.Counter.DecodeBuffer(buf); err != nil {
		return err
	}
	if err = o.GasLimit.DecodeBuffer(buf); err != nil {
		return err
	}
	if err = o.StorageLimit.DecodeBuffer(buf); err != nil {
		return err
	}
	return nil
}

func (o *Manager) WithSource(addr mavryk.Address) {
	o.Source = addr
}

func (o *Manager) WithCounter(c int64) {
	o.Counter.SetInt64(c)
}

func (o Manager) GetCounter() int64 {
	return o.Counter.Int64()
}

func (o *Manager) WithLimits(limits mavryk.Limits) {
	o.Fee.SetInt64(limits.Fee)
	o.GasLimit.SetInt64(limits.GasLimit)
	o.StorageLimit.SetInt64(limits.StorageLimit)
}
