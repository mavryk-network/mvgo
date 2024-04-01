// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"mavrykdynamics/tzgo/tezos"
)

// Ensure DAL types implement the TypedOperation interface.
var (
	_ TypedOperation = (*DalPublishSlotHeader)(nil)
	_ TypedOperation = (*DalAttestation)(nil)
)

type DalPublishSlotHeader struct {
	Manager
	SlotHeader struct {
		Level      int64          `json:"level"`
		Index      byte           `json:"index"`
		Commitment string         `json:"commitment"`
		Proof      tezos.HexBytes `json:"commitment_proof"`
	} `json:"slot_header"`
}

type DalAttestation struct {
	Generic
	Attestor    tezos.Address `json:"attestor"`
	Attestation tezos.Z       `json:"attestation"`
	Level       int64         `json:"level"`
}
