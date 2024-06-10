// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/micheline"
)

// Transaction represents "transaction" operation
type Transaction struct {
	Manager
	Amount      mavryk.N              `json:"amount"`
	Destination mavryk.Address        `json:"destination"`
	Parameters  *micheline.Parameters `json:"parameters,omitempty"`
}

func (o Transaction) Kind() mavryk.OpType {
	return mavryk.OpTypeTransaction
}

func (o Transaction) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteByte(',')
	o.Manager.EncodeJSON(buf)
	buf.WriteString(`,"amount":`)
	buf.WriteString(strconv.Quote(o.Amount.String()))
	buf.WriteString(`,"destination":`)
	buf.WriteString(strconv.Quote(o.Destination.String()))
	if o.Parameters != nil {
		buf.WriteString(`,"parameters":`)
		b, _ := o.Parameters.MarshalJSON()
		buf.Write(b)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o Transaction) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	o.Amount.EncodeBuffer(buf)
	buf.Write(o.Destination.EncodePadded())
	if o.Parameters != nil {
		buf.WriteByte(0xff)
		o.Parameters.EncodeBuffer(buf)
	} else {
		buf.WriteByte(0x0)
	}
	return nil
}

func (o *Transaction) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.Manager.DecodeBuffer(buf, p); err != nil {
		return err
	}
	if err = o.Amount.DecodeBuffer(buf); err != nil {
		return
	}
	if err = o.Destination.Decode(buf.Next(22)); err != nil {
		return
	}
	var ok bool
	ok, err = readBool(buf.Next(1))
	if err != nil {
		return
	}
	if ok {
		param := &micheline.Parameters{}
		if err = param.DecodeBuffer(buf); err != nil {
			return err
		}
		o.Parameters = param
	}
	return nil
}

func (o Transaction) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *Transaction) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
