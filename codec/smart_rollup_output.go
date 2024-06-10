// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// SmartRollupExecuteOutboxMessage represents "smart_rollup_execute_outbox_message" operation
type SmartRollupExecuteOutboxMessage struct {
	Manager
	Rollup   mavryk.Address               `json:"rollup"`
	Cemented mavryk.SmartRollupCommitHash `json:"cemented_commitment"`
	Proof    mavryk.HexBytes              `json:"output_proof"`
}

func (o SmartRollupExecuteOutboxMessage) Kind() mavryk.OpType {
	return mavryk.OpTypeSmartRollupExecuteOutboxMessage
}

func (o SmartRollupExecuteOutboxMessage) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteByte(',')
	o.Manager.EncodeJSON(buf)
	buf.WriteString(`,"rollup":`)
	buf.WriteString(strconv.Quote(o.Rollup.String()))
	buf.WriteString(`,"cemented_commitment":`)
	buf.WriteString(strconv.Quote(o.Cemented.String()))
	buf.WriteString(`,"output_proof":`)
	buf.WriteString(strconv.Quote(o.Proof.String()))
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o SmartRollupExecuteOutboxMessage) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	buf.Write(o.Rollup.Hash()) // 20 byte only
	buf.Write(o.Cemented[:])
	writeBytesWithLen(buf, o.Proof)
	return nil
}

func (o *SmartRollupExecuteOutboxMessage) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.Manager.DecodeBuffer(buf, p); err != nil {
		return
	}
	o.Rollup = mavryk.NewAddress(mavryk.AddressTypeSmartRollup, buf.Next(20))
	o.Cemented = mavryk.NewSmartRollupCommitHash(buf.Next(32))
	o.Proof, err = readBytesWithLen(buf)
	return
}

func (o SmartRollupExecuteOutboxMessage) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *SmartRollupExecuteOutboxMessage) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
