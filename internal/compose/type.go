// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc, abdul@blockwatch.cc

package compose

import (
	"fmt"
	"time"

	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/micheline"
)

func ParseValue(typ micheline.OpCode, value string) (any, error) {
	switch typ {
	case micheline.T_STRING:
		return value, nil
	case micheline.T_ADDRESS:
		return mavryk.ParseAddress(value)
	case micheline.T_NAT, micheline.T_MUMAV, micheline.T_INT:
		return mavryk.ParseZ(value)
	case micheline.T_TIMESTAMP:
		return time.Parse(time.RFC3339, value)
	case micheline.T_BYTES:
		var h mavryk.HexBytes
		if err := h.UnmarshalText([]byte(value)); err != nil {
			return nil, err
		}
		return h.Bytes(), nil
	case micheline.T_KEY:
		return mavryk.DecodeKey([]byte(value))
	case micheline.T_SIGNATURE:
		return mavryk.ParseSignature(value)
	case micheline.T_CHAIN_ID:
		return mavryk.ParseChainIdHash(value)
	default:
		return nil, fmt.Errorf("cannot parsed typ %q is ", typ)
	}
}
