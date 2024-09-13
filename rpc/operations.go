// Copyright (c) 2020-2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/micheline"
)

type MetadataMode string

const (
	MetadataModeUnset  MetadataMode = ""
	MetadataModeNever  MetadataMode = "never"
	MetadataModeAlways MetadataMode = "always"
)

// Operation represents a single operation or batch of operations included in a block
type Operation struct {
	Protocol  mavryk.ProtocolHash `json:"protocol"`
	ChainID   mavryk.ChainIdHash  `json:"chain_id"`
	Hash      mavryk.OpHash       `json:"hash"`
	Branch    mavryk.BlockHash    `json:"branch"`
	Contents  OperationList       `json:"contents"`
	Signature mavryk.Signature    `json:"signature"`
	Errors    []OperationError    `json:"error,omitempty"`    // mempool only
	Metadata  string              `json:"metadata,omitempty"` // contains `too large` when stripped, this is BAD!!
}

// TotalCosts returns the sum of costs across all batched and internal operations.
func (o Operation) TotalCosts() mavryk.Costs {
	var c mavryk.Costs
	for _, op := range o.Contents {
		c = c.Add(op.Costs())
	}
	return c
}

// Costs returns ta list of individual costs for all batched operations.
func (o Operation) Costs() []mavryk.Costs {
	list := make([]mavryk.Costs, len(o.Contents))
	for i, op := range o.Contents {
		list[i] = op.Costs()
	}
	return list
}

// TypedOperation must be implemented by all operations
type TypedOperation interface {
	Kind() mavryk.OpType
	Meta() OperationMetadata
	Result() OperationResult
	Costs() mavryk.Costs
	Limits() mavryk.Limits
}

// OperationError represents data describing an error conditon that lead to a
// failed operation execution.
type OperationError struct {
	GenericError
	Contract *mavryk.Address `json:"contract,omitempty"`
	Raw      json.RawMessage `json:"-"`
}

// OperationMetadata contains execution receipts for successful and failed
// operations.
type OperationMetadata struct {
	BalanceUpdates BalanceUpdates  `json:"balance_updates"` // fee-related
	Result         OperationResult `json:"operation_result"`

	// transaction only
	InternalResults []*InternalResult `json:"internal_operation_results,omitempty"`

	// endorsement only
	Delegate            mavryk.Address `json:"delegate"`
	Slots               []int          `json:"slots,omitempty"`
	EndorsementPower    int            `json:"endorsement_power,omitempty"`    // v12+
	PreendorsementPower int            `json:"preendorsement_power,omitempty"` // v12+

	// some rollup ops only, FIXME: is this correct here or is this field in result?
	Level int64 `json:"level"`

	// v18 slashing ops may block a baker
	ForbiddenDelegate mavryk.Address `json:"forbidden_delegate"` // v18+
}

// Address returns the delegate address for endorsements.
func (m OperationMetadata) Address() mavryk.Address {
	return m.Delegate
}

// OperationResult contains receipts for executed operations, both success and failed.
// This type is a generic container for all possible results. Which fields are actually
// used depends on operation type and performed actions.
type OperationResult struct {
	Status               mavryk.OpStatus  `json:"status"`
	BalanceUpdates       BalanceUpdates   `json:"balance_updates"`
	ConsumedGas          int64            `json:"consumed_gas,string"`      // deprecated in v015
	ConsumedMilliGas     int64            `json:"consumed_milligas,string"` // v007+
	Errors               []OperationError `json:"errors,omitempty"`
	Allocated            bool             `json:"allocated_destination_contract"` // tx only
	Storage              *micheline.Prim  `json:"storage,omitempty"`              // tx, orig
	OriginatedContracts  []mavryk.Address `json:"originated_contracts"`           // orig only
	StorageSize          int64            `json:"storage_size,string"`            // tx, orig, const
	PaidStorageSizeDiff  int64            `json:"paid_storage_size_diff,string"`  // tx, orig
	BigmapDiff           json.RawMessage  `json:"big_map_diff,omitempty"`         // tx, orig, <v013
	LazyStorageDiff      json.RawMessage  `json:"lazy_storage_diff,omitempty"`    // v008+ tx, orig
	GlobalAddress        mavryk.ExprHash  `json:"global_address"`                 // const
	TicketUpdatesCorrect []TicketUpdate   `json:"ticket_updates"`                 // v015
	TicketReceipts       []TicketUpdate   `json:"ticket_receipt"`                 // v015, name on internal

	// v013 tx rollup
	TxRollupResult

	// v016 smart rollup
	SmartRollupResult

	// v019 DAL
	DalResult
}

// Always use this helper to retrieve Ticket updates. This is because due to
// lack of quality control Tezos Lima protocol ended up with 2 distinct names
// for ticket updates in external call receipts versus internal call receipts.
func (r OperationResult) TicketUpdates() []TicketUpdate {
	if len(r.TicketUpdatesCorrect) > 0 {
		return r.TicketUpdatesCorrect
	}
	return r.TicketReceipts
}

func (r OperationResult) BigmapEvents() micheline.BigmapEvents {
	if r.LazyStorageDiff != nil {
		res := make(micheline.LazyEvents, 0)
		_ = json.Unmarshal(r.LazyStorageDiff, &res)
		return res.BigmapEvents()
	}
	if r.BigmapDiff != nil {
		res := make(micheline.BigmapEvents, 0)
		_ = json.Unmarshal(r.BigmapDiff, &res)
		return res
	}
	return nil
}

func (r OperationResult) IsSuccess() bool {
	return r.Status == mavryk.OpStatusApplied
}

func (r OperationResult) Gas() int64 {
	if r.ConsumedMilliGas > 0 {
		var corr int64
		if r.ConsumedMilliGas%1000 > 0 {
			corr++
		}
		return r.ConsumedMilliGas/1000 + corr
	}
	return r.ConsumedGas
}

func (r OperationResult) MilliGas() int64 {
	if r.ConsumedMilliGas > 0 {
		return r.ConsumedMilliGas
	}
	return r.ConsumedGas * 1000
}

func (o OperationError) MarshalJSON() ([]byte, error) {
	return o.Raw, nil
}

func (o *OperationError) UnmarshalJSON(data []byte) error {
	type alias OperationError
	if err := json.Unmarshal(data, (*alias)(o)); err != nil {
		return err
	}
	o.Raw = make([]byte, len(data))
	copy(o.Raw, data)
	return nil
}

// Generic is the most generic operation type.
type Generic struct {
	OpKind   mavryk.OpType     `json:"kind"`
	Metadata OperationMetadata `json:"metadata"`
}

// Kind returns the operation's type. Implements TypedOperation interface.
func (e Generic) Kind() mavryk.OpType {
	return e.OpKind
}

// Meta returns an empty operation metadata to implement TypedOperation interface.
func (e Generic) Meta() OperationMetadata {
	return e.Metadata
}

// Result returns an empty operation result to implement TypedOperation interface.
func (e Generic) Result() OperationResult {
	return e.Metadata.Result
}

// Costs returns empty operation costs to implement TypedOperation interface.
func (e Generic) Costs() mavryk.Costs {
	return mavryk.Costs{}
}

// Limits returns empty operation limits to implement TypedOperation interface.
func (e Generic) Limits() mavryk.Limits {
	return mavryk.Limits{}
}

// Manager represents data common for all manager operations.
type Manager struct {
	Generic
	Source       mavryk.Address `json:"source"`
	Fee          int64          `json:"fee,string"`
	Counter      int64          `json:"counter,string"`
	GasLimit     int64          `json:"gas_limit,string"`
	StorageLimit int64          `json:"storage_limit,string"`
}

// Limits returns manager operation limits to implement TypedOperation interface.
func (e Manager) Limits() mavryk.Limits {
	return mavryk.Limits{
		Fee:          e.Fee,
		GasLimit:     e.GasLimit,
		StorageLimit: e.StorageLimit,
	}
}

// OperationList is a slice of TypedOperation (interface type) with custom JSON unmarshaller
type OperationList []TypedOperation

// Contains returns true when the list contains an operation of kind typ.
func (o OperationList) Contains(typ mavryk.OpType) bool {
	for _, v := range o {
		if v.Kind() == typ {
			return true
		}
	}
	return false
}

func (o OperationList) Select(typ mavryk.OpType, n int) TypedOperation {
	var cnt int
	for _, v := range o {
		if v.Kind() != typ {
			continue
		}
		if cnt == n {
			return v
		}
		cnt++
	}
	return nil
}

func (o OperationList) Len() int {
	return len(o)
}

func (o OperationList) N(n int) TypedOperation {
	if n < 0 {
		n += len(o)
	}
	return o[n]
}

// UnmarshalJSON implements json.Unmarshaler
func (e *OperationList) UnmarshalJSON(data []byte) error {
	if len(data) <= 2 {
		return nil
	}

	if data[0] != '[' {
		return fmt.Errorf("rpc: expected operation array")
	}

	// fmt.Printf("Decoding ops: %s\n", string(data))
	dec := json.NewDecoder(bytes.NewReader(data))

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		return fmt.Errorf("rpc: %v", err)
	}

	for dec.More() {
		// peek into `{"kind":"...",` field
		start := int(dec.InputOffset()) + 9
		// after first JSON object, decoder pos is at `,`
		if data[start] == '"' {
			start += 1
		}
		end := start + bytes.IndexByte(data[start:], '"')
		kind := mavryk.ParseOpType(string(data[start:end]))
		var op TypedOperation
		switch kind {
		// anonymous operations
		case mavryk.OpTypeActivateAccount:
			op = &Activation{}
		case mavryk.OpTypeDoubleBakingEvidence:
			op = &DoubleBaking{}
		case mavryk.OpTypeDoubleEndorsementEvidence,
			mavryk.OpTypeDoublePreendorsementEvidence,
			mavryk.OpTypeDoubleAttestationEvidence,
			mavryk.OpTypeDoublePreattestationEvidence:
			op = &DoubleEndorsement{}
		case mavryk.OpTypeSeedNonceRevelation:
			op = &SeedNonce{}
		case mavryk.OpTypeDrainDelegate:
			op = &DrainDelegate{}

		// consensus operations
		case mavryk.OpTypeEndorsement,
			mavryk.OpTypeEndorsementWithSlot,
			mavryk.OpTypePreendorsement,
			mavryk.OpTypeAttestation,
			mavryk.OpTypeAttestationWithDal,
			mavryk.OpTypePreattestation:
			op = &Endorsement{}

		// amendment operations
		case mavryk.OpTypeProposals:
			op = &Proposals{}
		case mavryk.OpTypeBallot:
			op = &Ballot{}

		// manager operations
		case mavryk.OpTypeTransaction:
			op = &Transaction{}
		case mavryk.OpTypeOrigination:
			op = &Origination{}
		case mavryk.OpTypeDelegation:
			op = &Delegation{}
		case mavryk.OpTypeReveal:
			op = &Reveal{}
		case mavryk.OpTypeRegisterConstant:
			op = &ConstantRegistration{}
		case mavryk.OpTypeSetDepositsLimit:
			op = &SetDepositsLimit{}
		case mavryk.OpTypeIncreasePaidStorage:
			op = &IncreasePaidStorage{}
		case mavryk.OpTypeVdfRevelation:
			op = &VdfRevelation{}
		case mavryk.OpTypeTransferTicket:
			op = &TransferTicket{}
		case mavryk.OpTypeUpdateConsensusKey:
			op = &UpdateConsensusKey{}

			// DEPRECATED: tx rollup operations, kept for testnet backward compatibility
		case mavryk.OpTypeTxRollupOrigination,
			mavryk.OpTypeTxRollupSubmitBatch,
			mavryk.OpTypeTxRollupCommit,
			mavryk.OpTypeTxRollupReturnBond,
			mavryk.OpTypeTxRollupFinalizeCommitment,
			mavryk.OpTypeTxRollupRemoveCommitment,
			mavryk.OpTypeTxRollupRejection,
			mavryk.OpTypeTxRollupDispatchTickets:
			op = &TxRollup{}

		case mavryk.OpTypeSmartRollupOriginate:
			op = &SmartRollupOriginate{}
		case mavryk.OpTypeSmartRollupAddMessages:
			op = &SmartRollupAddMessages{}
		case mavryk.OpTypeSmartRollupCement:
			op = &SmartRollupCement{}
		case mavryk.OpTypeSmartRollupPublish:
			op = &SmartRollupPublish{}
		case mavryk.OpTypeSmartRollupRefute:
			op = &SmartRollupRefute{}
		case mavryk.OpTypeSmartRollupTimeout:
			op = &SmartRollupTimeout{}
		case mavryk.OpTypeSmartRollupExecuteOutboxMessage:
			op = &SmartRollupExecuteOutboxMessage{}
		case mavryk.OpTypeSmartRollupRecoverBond:
			op = &SmartRollupRecoverBond{}
		case mavryk.OpTypeDalPublishCommitment:
			op = &DalPublishCommitment{}

		default:
			return fmt.Errorf("rpc: unsupported op %q", string(data[start:end]))
		}

		if err := dec.Decode(op); err != nil {
			return fmt.Errorf("rpc: operation kind %s: %v", kind, err)
		}
		(*e) = append(*e, op)
	}

	return nil
}

// GetBlockOperationHash returns a single operation hashes included in block
// https://protocol.mavryk.org/active/rpc.html#get-block-id-operation-hashes-list-offset-operation-offset
func (c *Client) GetBlockOperationHash(ctx context.Context, id BlockID, l, n int) (mavryk.OpHash, error) {
	var hash mavryk.OpHash
	u := fmt.Sprintf("chains/main/blocks/%s/operation_hashes/%d/%d", id, l, n)
	err := c.Get(ctx, u, &hash)
	return hash, err
}

// GetBlockOperationHashes returns a list of list of operation hashes included in block
// https://protocol.mavryk.org/active/rpc.html#get-block-id-operation-hashes
func (c *Client) GetBlockOperationHashes(ctx context.Context, id BlockID) ([][]mavryk.OpHash, error) {
	hashes := make([][]mavryk.OpHash, 0)
	u := fmt.Sprintf("chains/main/blocks/%s/operation_hashes", id)
	if err := c.Get(ctx, u, &hashes); err != nil {
		return nil, err
	}
	return hashes, nil
}

// GetBlockOperationListHashes returns a list of operation hashes included in block
// at a specified list position (i.e. validation pass) [0..3]
// https://protocol.mavryk.org/active/rpc.html#get-block-id-operation-hashes-list-offset
func (c *Client) GetBlockOperationListHashes(ctx context.Context, id BlockID, l int) ([]mavryk.OpHash, error) {
	hashes := make([]mavryk.OpHash, 0)
	u := fmt.Sprintf("chains/main/blocks/%s/operation_hashes/%d", id, l)
	if err := c.Get(ctx, u, &hashes); err != nil {
		return nil, err
	}
	return hashes, nil
}

// GetBlockOperation returns information about a single validated Tezos operation group
// (i.e. a single operation or a batch of operations) at list l and position n
// https://protocol.mavryk.org/active/rpc.html#get-block-id-operations-list-offset-operation-offset
func (c *Client) GetBlockOperation(ctx context.Context, id BlockID, l, n int) (*Operation, error) {
	var op Operation
	u := fmt.Sprintf("chains/main/blocks/%s/operations/%d/%d", id, l, n)
	if c.MetadataMode != "" {
		u += "?metadata=" + string(c.MetadataMode)
	}
	if err := c.Get(ctx, u, &op); err != nil {
		return nil, err
	}
	return &op, nil
}

// GetBlockOperationList returns information about all validated Tezos operation group
// inside operation list l (i.e. validation pass) [0..3].
// https://protocol.mavryk.org/active/rpc.html#get-block-id-operations-list-offset
func (c *Client) GetBlockOperationList(ctx context.Context, id BlockID, l int) ([]Operation, error) {
	ops := make([]Operation, 0)
	u := fmt.Sprintf("chains/main/blocks/%s/operations/%d", id, l)
	if c.MetadataMode != "" {
		u += "?metadata=" + string(c.MetadataMode)
	}
	if err := c.Get(ctx, u, &ops); err != nil {
		return nil, err
	}
	return ops, nil
}

// GetBlockOperations returns information about all validated Tezos operation groups
// from all operation lists in block.
// https://protocol.mavryk.org/active/rpc.html#get-block-id-operations
func (c *Client) GetBlockOperations(ctx context.Context, id BlockID) ([][]Operation, error) {
	ops := make([][]Operation, 0)
	u := fmt.Sprintf("chains/main/blocks/%s/operations", id)
	if c.MetadataMode != "" {
		u += "?metadata=" + string(c.MetadataMode)
	}
	if err := c.Get(ctx, u, &ops); err != nil {
		return nil, err
	}
	return ops, nil
}
