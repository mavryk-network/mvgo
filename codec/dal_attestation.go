// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// DalAttestation represents "dal_attestation" operation
type DalAttestation struct {
	Simple
	Attestor    mavryk.Address `json:"attestor"`
	Attestation mavryk.Z       `json:"attestation"`
	Level       int32          `json:"level"`
}

func (o DalAttestation) Kind() mavryk.OpType {
	return mavryk.OpTypeDalAttestation
}

func (o DalAttestation) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteString(`,"attestor":`)
	buf.WriteString(strconv.Quote(o.Attestor.String()))
	buf.WriteString(`,"attestation":`)
	buf.WriteString(strconv.Quote(o.Attestation.String()))
	buf.WriteString(`,"level":`)
	buf.WriteString(strconv.Itoa(int(o.Level)))
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o DalAttestation) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	buf.Write(o.Attestor.Encode())
	buf.Write(o.Attestation.Bytes())
	binary.Write(buf, enc, o.Level)
	return nil
}

func (o *DalAttestation) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.Attestor.Decode(buf.Next(21)); err != nil {
		return
	}
	if err = o.Attestation.DecodeBuffer(buf); err != nil {
		return
	}
	o.Level, err = readInt32(buf.Next(4))
	return
}

func (o DalAttestation) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *DalAttestation) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
