// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// ActivateAccount represents "activate_account" operation
type ActivateAccount struct {
	Simple
	PublicKeyHash mavryk.Address  `json:"pkh"`
	Secret        mavryk.HexBytes `json:"secret"`
}

func (o ActivateAccount) Kind() mavryk.OpType {
	return mavryk.OpTypeActivateAccount
}

func (o ActivateAccount) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteString(`,"pkh":`)
	buf.WriteString(strconv.Quote(o.PublicKeyHash.String()))
	buf.WriteString(`,"secret":`)
	buf.WriteString(strconv.Quote(o.Secret.String()))
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o ActivateAccount) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	buf.Write(o.PublicKeyHash[1:]) // only place where a 20 byte address is used (!)
	buf.Write(o.Secret.Bytes())
	return nil
}

func (o *ActivateAccount) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	if err := ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return err
	}
	o.PublicKeyHash = mavryk.NewAddress(mavryk.AddressTypeEd25519, buf.Next(20))
	if !o.PublicKeyHash.IsValid() {
		return fmt.Errorf("invalid address %q", o.PublicKeyHash)
	}
	return o.Secret.ReadBytes(buf, 20)
}

func (o ActivateAccount) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *ActivateAccount) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}
