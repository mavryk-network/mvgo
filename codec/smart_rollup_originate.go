// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/micheline"
)

// SmartRollupOriginate represents "smart_rollup_originate" operation
type SmartRollupOriginate struct {
	Manager
	Pvm    mavryk.PvmKind  `json:"pvm_kind"`
	Kernel mavryk.HexBytes `json:"kernel"`
	Proof  mavryk.HexBytes `json:"origination_proof"`
	Type   micheline.Prim  `json:"parameters_ty"`
}

func (o SmartRollupOriginate) Kind() mavryk.OpType {
	return mavryk.OpTypeSmartRollupOriginate
}

func (o SmartRollupOriginate) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteByte(',')
	o.Manager.EncodeJSON(buf)
	buf.WriteString(`,"pvm_kind":`)
	buf.WriteString(strconv.Quote(o.Pvm.String()))
	buf.WriteString(`,"kernel":`)
	buf.WriteString(strconv.Quote(o.Kernel.String()))
	buf.WriteString(`,"origination_proof":`)
	buf.WriteString(strconv.Quote(o.Proof.String()))
	buf.WriteString(`,"parameters_ty":`)
	o.Type.EncodeJSON(buf)
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o SmartRollupOriginate) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	binary.Write(buf, enc, o.Pvm)
	writeBytesWithLen(buf, o.Kernel)
	writeBytesWithLen(buf, o.Proof)
	writePrimWithLen(buf, o.Type)
	return nil
}

func (o *SmartRollupOriginate) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.Manager.DecodeBuffer(buf, p); err != nil {
		return
	}
	var b byte
	if b, err = readByte(buf.Next(1)); err != nil {
		return
	} else if typ := mavryk.PvmKind(b); !typ.IsValid() {
		err = fmt.Errorf("Unsupported PVM type %d", b)
		return
	} else {
		o.Pvm = typ
	}
	if o.Kernel, err = readBytesWithLen(buf); err != nil {
		return
	}
	if o.Proof, err = readBytesWithLen(buf); err != nil {
		return
	}
	if o.Type, err = readPrimWithLen(buf); err != nil {
		return
	}
	return
}

func (o SmartRollupOriginate) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *SmartRollupOriginate) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
