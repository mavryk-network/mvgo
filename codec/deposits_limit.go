// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// SetDepositsLimit represents "set_deposits_limit" operation
type SetDepositsLimit struct {
	Manager
	Limit *mavryk.N `json:"limit,omitempty"`
}

func (o SetDepositsLimit) Kind() mavryk.OpType {
	return mavryk.OpTypeSetDepositsLimit
}

func (o SetDepositsLimit) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteByte(',')
	o.Manager.EncodeJSON(buf)
	if o.Limit != nil {
		buf.WriteString(`,"limit":`)
		buf.WriteString(strconv.Quote(o.Limit.String()))
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o SetDepositsLimit) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	if o.Limit == nil {
		buf.WriteByte(0x00)
	} else {
		buf.WriteByte(0xff)
		o.Limit.EncodeBuffer(buf)
	}
	return nil
}

func (o *SetDepositsLimit) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.Manager.DecodeBuffer(buf, p); err != nil {
		return err
	}
	var ok bool
	ok, err = readBool(buf.Next(1))
	if err != nil {
		return err
	}
	if ok {
		var limit mavryk.N
		if err = limit.DecodeBuffer(buf); err != nil {
			return err
		}
		o.Limit = &limit
	}
	return nil
}

func (o SetDepositsLimit) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *SetDepositsLimit) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
