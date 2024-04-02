// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/micheline"
)

// Ensure SmartRollup types implement the TypedOperation interface.
var (
	_ TypedOperation = (*SmartRollupOriginate)(nil)
	_ TypedOperation = (*SmartRollupAddMessages)(nil)
	_ TypedOperation = (*SmartRollupCement)(nil)
	_ TypedOperation = (*SmartRollupPublish)(nil)
	_ TypedOperation = (*SmartRollupRefute)(nil)
	_ TypedOperation = (*SmartRollupTimeout)(nil)
	_ TypedOperation = (*SmartRollupExecuteOutboxMessage)(nil)
	_ TypedOperation = (*SmartRollupRecoverBond)(nil)
)

type SmartRollupResult struct {
	Address          *mavryk.Address               `json:"address,omitempty"`            // v016, smart_rollup_originate
	Size             *mavryk.Z                     `json:"size,omitempty"`               // v016, smart_rollup_originate
	InboxLevel       int64                         `json:"inbox_level,omitempty"`        // v016, smart_rollup_cement
	StakedHash       *mavryk.SmartRollupCommitHash `json:"staked_hash,omitempty"`        // v016, smart_rollup_publish
	PublishedAtLevel int64                         `json:"published_at_level,omitempty"` // v016, smart_rollup_publish
	GameStatus       *GameStatus                   `json:"game_status,omitempty"`        // v016, smart_rollup_refute, smart_rollup_timeout
	Commitment       *mavryk.SmartRollupCommitHash `json:"commitment_hash,omitempty"`    // v017, smart_rollup_cement
}

type SmartRollupOriginate struct {
	Manager
	PvmKind          mavryk.PvmKind  `json:"pvm_kind"`
	Kernel           mavryk.HexBytes `json:"kernel"`
	OriginationProof mavryk.HexBytes `json:"origination_proof"`
	ParametersTy     micheline.Prim  `json:"parameters_ty"`
}

type SmartRollupAddMessages struct {
	Manager
	Messages []mavryk.HexBytes `json:"message"`
}

type SmartRollupCement struct {
	Manager
	Rollup     mavryk.Address                `json:"rollup"`
	Commitment *mavryk.SmartRollupCommitHash `json:"commitment,omitempty"` // deprecated in v17
}

type SmartRollupCommitment struct {
	CompressedState mavryk.SmartRollupStateHash  `json:"compressed_state"`
	InboxLevel      int64                        `json:"inbox_level"`
	Predecessor     mavryk.SmartRollupCommitHash `json:"predecessor"`
	NumberOfTicks   mavryk.Z                     `json:"number_of_ticks"`
}

type SmartRollupPublish struct {
	Manager
	Rollup     mavryk.Address        `json:"rollup"`
	Commitment SmartRollupCommitment `json:"commitment"`
}

type SmartRollupRefute struct {
	Manager
	Rollup     mavryk.Address        `json:"rollup"`
	Opponent   mavryk.Address        `json:"opponent"`
	Refutation SmartRollupRefutation `json:"refutation"`
}

type SmartRollupRefutation struct {
	Kind         string                        `json:"refutation_kind"`
	PlayerHash   *mavryk.SmartRollupCommitHash `json:"player_commitment_hash,omitempty"`
	OpponentHash *mavryk.SmartRollupCommitHash `json:"opponent_commitment_hash,omitempty"`
	Choice       *mavryk.Z                     `json:"choice,omitempty"`
	Step         *SmartRollupRefuteStep        `json:"step,omitempty"`
}

// Step can either be
//
// - []SmartRollupTick
// - SmartRollupInputProof
// - smth else?
//
// There is no indication in the outer parts of the refutation struct that
// suggests how to decode this.
type SmartRollupRefuteStep struct {
	Ticks []SmartRollupTick
	Proof *SmartRollupProof
}

type SmartRollupProof struct {
	PvmStep    mavryk.HexBytes        `json:"pvm_step,omitempty"`
	InputProof *SmartRollupInputProof `json:"input_proof,omitempty"`
}

func (s *SmartRollupRefuteStep) UnmarshalJSON(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	switch buf[0] {
	case '[':
		s.Ticks = make([]SmartRollupTick, 0)
		return json.Unmarshal(buf, &s.Ticks)
	case '{':
		s.Proof = &SmartRollupProof{}
		return json.Unmarshal(buf, s.Proof)
	default:
		return fmt.Errorf("Invalid refute step data %q", string(buf))
	}
}

func (s SmartRollupRefuteStep) MarshalJSON() ([]byte, error) {
	if s.Ticks != nil {
		return json.Marshal(s.Ticks)
	}
	if s.Proof != nil {
		return json.Marshal(s.Proof)
	}
	return nil, nil
}

type SmartRollupTick struct {
	State mavryk.SmartRollupStateHash `json:"state"`
	Tick  mavryk.Z                    `json:"tick"`
}

type SmartRollupInputProof struct {
	Kind    string          `json:"input_proof_kind"`
	Level   int64           `json:"level"`
	Counter mavryk.Z        `json:"message_counter"`
	Proof   mavryk.HexBytes `json:"serialized_proof"`
}

type SmartRollupTimeout struct {
	Manager
	Rollup  mavryk.Address `json:"rollup"`
	Stakers struct {
		Alice mavryk.Address `json:"alice"`
		Bob   mavryk.Address `json:"bob"`
	} `json:"stakers"`
}

type SmartRollupExecuteOutboxMessage struct {
	Manager
	Rollup             mavryk.Address               `json:"rollup"`
	CementedCommitment mavryk.SmartRollupCommitHash `json:"cemented_commitment"`
	OutputProof        mavryk.HexBytes              `json:"output_proof"`
}

type SmartRollupRecoverBond struct {
	Manager
	Rollup mavryk.Address `json:"rollup"`
	Staker mavryk.Address `json:"staker"`
}

type GameStatus struct {
	Status string          `json:"status,omitempty"`
	Kind   string          `json:"kind,omitempty"`
	Reason string          `json:"reason,omitempty"`
	Player *mavryk.Address `json:"player,omitempty"`
}

func (s *GameStatus) UnmarshalJSON(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	switch buf[0] {
	case '"':
		s.Status = string(buf[1 : len(buf)-1])
	case '{':
		type alias *GameStatus
		type wrapper struct {
			S alias `json:"result"`
		}
		a := wrapper{alias(s)}
		_ = json.Unmarshal(buf, &a)
	default:
		return fmt.Errorf("Invalid game status data %q", string(buf))
	}
	return nil
}
