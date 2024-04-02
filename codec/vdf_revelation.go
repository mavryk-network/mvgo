// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// VdfRevelation represents "vdf_revelation" operation
type VdfRevelation struct {
	Simple
	Solution mavryk.HexBytes `json:"solution"`
}

func (o VdfRevelation) Kind() mavryk.OpType {
	return mavryk.OpTypeVdfRevelation
}

func (o VdfRevelation) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteString(`,"solution":`)
	buf.WriteString(strconv.Quote(o.Solution.String()))
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o VdfRevelation) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	buf.Write(o.Solution.Bytes())
	return nil
}

func (o *VdfRevelation) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	return o.Solution.ReadBytes(buf, 200)
}

func (o VdfRevelation) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *VdfRevelation) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
