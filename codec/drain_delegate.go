// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// DrainDelegate represents "drain_delegate" operation
type DrainDelegate struct {
	Simple
	ConsensusKey mavryk.Address `json:"consensus_key"`
	Delegate     mavryk.Address `json:"delegate"`
	Destination  mavryk.Address `json:"destination"`
}

func (o DrainDelegate) Kind() mavryk.OpType {
	return mavryk.OpTypeDrainDelegate
}

func (o DrainDelegate) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteString(`,"consensus_key":`)
	buf.WriteString(strconv.Quote(o.ConsensusKey.String()))
	buf.WriteString(`,"delegate":`)
	buf.WriteString(strconv.Quote(o.Delegate.String()))
	buf.WriteString(`,"destination":`)
	buf.WriteString(strconv.Quote(o.Destination.String()))
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o DrainDelegate) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	buf.Write(o.ConsensusKey.Encode())
	buf.Write(o.Delegate.Encode())
	buf.Write(o.Destination.Encode())
	return nil
}

func (o *DrainDelegate) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.ConsensusKey.Decode(buf.Next(21)); err != nil {
		return
	}
	if err = o.Delegate.Decode(buf.Next(21)); err != nil {
		return
	}
	if err = o.Destination.Decode(buf.Next(21)); err != nil {
		return
	}
	return nil
}

func (o DrainDelegate) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *DrainDelegate) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
