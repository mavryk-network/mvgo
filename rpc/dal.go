// Copyright (c) 2024 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import "github.com/mavryk-network/mvgo/mavryk"

// Ensure DAL types implement the TypedOperation interface.
var (
	_ TypedOperation = (*DalPublishCommitment)(nil)
)

type DalPublishCommitment struct {
	Manager
	SlotHeader struct {
		Index      byte            `json:"slot_index"`
		Commitment string          `json:"commitment"`
		Proof      mavryk.HexBytes `json:"commitment_proof"`
	} `json:"slot_header"`
}

type DalResult struct {
	SlotHeader struct {
		Version    string `json:"version"`
		Level      int64  `json:"level"`
		Index      byte   `json:"index"`
		Commitment string `json:"commitment"`
	} `json:"slot_header"`
}
