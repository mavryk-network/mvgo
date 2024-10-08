// Copyright (c) 2020-2024 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import "github.com/mavryk-network/mvgo/mavryk"

// Ensure Endorsement implements the TypedOperation interface.
var _ TypedOperation = (*Endorsement)(nil)

// Endorsement represents an endorsement operation
type Endorsement struct {
	Generic
	Level          int64               `json:"level"`                 // <= v008, v012+
	Endorsement    *InlinedEndorsement `json:"endorsement,omitempty"` // v009+
	Slot           int                 `json:"slot"`                  // v009+
	Round          int                 `json:"round"`                 // v012+
	PayloadHash    mavryk.PayloadHash  `json:"block_payload_hash"`    // v012+
	DalAttestation mavryk.Z            `json:"dal_attestation"`       // v019+
}

func (e Endorsement) GetLevel() int64 {
	if e.Endorsement != nil {
		return e.Endorsement.Operations.Level
	}
	return e.Level
}

// InlinedEndorsement represents and embedded endorsement
type InlinedEndorsement struct {
	Branch     mavryk.BlockHash `json:"branch"`     // the double block
	Operations Endorsement      `json:"operations"` // only level and kind are set
	Signature  mavryk.Signature `json:"signature"`
}
