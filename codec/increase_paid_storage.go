// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// IncreasePaidStorage represents "increase_paid_storage" operation
type IncreasePaidStorage struct {
	Manager
	Amount      mavryk.Z       `json:"amount"`
	Destination mavryk.Address `json:"destination"`
}

func (o IncreasePaidStorage) Kind() mavryk.OpType {
	return mavryk.OpTypeIncreasePaidStorage
}

func (o IncreasePaidStorage) MarshalJSON() ([]byte, error) {
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
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o IncreasePaidStorage) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	o.Amount.EncodeBuffer(buf)
	buf.Write(o.Destination.EncodePadded())
	return nil
}

func (o *IncreasePaidStorage) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.Manager.DecodeBuffer(buf, p); err != nil {
		return err
	}
	if err = o.Amount.DecodeBuffer(buf); err != nil {
		return err
	}
	return o.Destination.Decode(buf.Next(22))
}

func (o IncreasePaidStorage) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *IncreasePaidStorage) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
