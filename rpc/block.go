// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/mavryk-network/mvgo/mavryk"
)

// Block holds information about a Tezos block
type Block struct {
	Protocol   mavryk.ProtocolHash `json:"protocol"`
	ChainId    mavryk.ChainIdHash  `json:"chain_id"`
	Hash       mavryk.BlockHash    `json:"hash"`
	Header     BlockHeader         `json:"header"`
	Metadata   BlockMetadata       `json:"metadata"`
	Operations [][]*Operation      `json:"operations"`
}

func (b Block) GetLevel() int64 {
	return b.Header.Level
}

func (b Block) GetTimestamp() time.Time {
	return b.Header.Timestamp
}

func (b Block) GetVersion() int {
	return b.Header.Proto
}

func (b Block) GetCycle() int64 {
	if b.Metadata.LevelInfo != nil {
		return b.Metadata.LevelInfo.Cycle
	}
	if b.Metadata.Level != nil {
		return b.Metadata.Level.Cycle
	}
	return 0
}

func (b Block) GetLevelInfo() LevelInfo {
	if b.Metadata.LevelInfo != nil {
		return *b.Metadata.LevelInfo
	}
	if b.Metadata.Level != nil {
		return *b.Metadata.Level
	}
	return LevelInfo{}
}

// only works for mainnet when before Edo or for all nets after Edo
// due to fixed constants used
func (b Block) GetVotingInfo() VotingPeriodInfo {
	if b.Metadata.VotingPeriodInfo != nil {
		return *b.Metadata.VotingPeriodInfo
	}
	if b.Metadata.Level != nil {
		return VotingPeriodInfo{
			Position:  b.Metadata.Level.VotingPeriodPosition,
			Remaining: 32768 - b.Metadata.Level.VotingPeriodPosition,
			VotingPeriod: VotingPeriod{
				Index:         b.Metadata.Level.VotingPeriod,
				Kind:          *b.Metadata.VotingPeriodKind,
				StartPosition: b.Metadata.Level.VotingPeriod * 32768,
			},
		}
	}
	return VotingPeriodInfo{}
}

func (b Block) GetVotingPeriodKind() mavryk.VotingPeriodKind {
	if b.Metadata.VotingPeriodInfo != nil {
		return b.Metadata.VotingPeriodInfo.VotingPeriod.Kind
	}
	if b.Metadata.VotingPeriodKind != nil {
		return *b.Metadata.VotingPeriodKind
	}
	return mavryk.VotingPeriodInvalid
}

func (b Block) GetVotingPeriod() int64 {
	if b.Metadata.VotingPeriodInfo != nil {
		return b.Metadata.VotingPeriodInfo.VotingPeriod.Index
	}
	if b.Metadata.Level != nil {
		return b.Metadata.Level.VotingPeriod
	}
	return 0
}

func (b Block) IsProtocolUpgrade() bool {
	return !b.Metadata.Protocol.Equal(b.Metadata.NextProtocol)
}

// InvalidBlock represents invalid block hash along with the errors that led to it being declared invalid
type InvalidBlock struct {
	Block mavryk.BlockHash `json:"block"`
	Level int64            `json:"level"`
	Error Errors           `json:"error"`
}

// BlockHeader is a part of the Tezos block data
type BlockHeader struct {
	Level                     int64                 `json:"level"`
	Proto                     int                   `json:"proto"`
	Predecessor               mavryk.BlockHash      `json:"predecessor"`
	Timestamp                 time.Time             `json:"timestamp"`
	ValidationPass            int                   `json:"validation_pass"`
	OperationsHash            mavryk.OpListListHash `json:"operations_hash"`
	Fitness                   []mavryk.HexBytes     `json:"fitness"`
	Context                   mavryk.ContextHash    `json:"context"`
	PayloadHash               mavryk.PayloadHash    `json:"payload_hash"`
	PayloadRound              int                   `json:"payload_round"`
	Priority                  int                   `json:"priority"`
	ProofOfWorkNonce          mavryk.HexBytes       `json:"proof_of_work_nonce"`
	SeedNonceHash             *mavryk.NonceHash     `json:"seed_nonce_hash"`
	Signature                 mavryk.Signature      `json:"signature"`
	Content                   *BlockContent         `json:"content,omitempty"`
	LiquidityBakingEscapeVote bool                  `json:"liquidity_baking_escape_vote"`
	LiquidityBakingToggleVote mavryk.FeatureVote    `json:"liquidity_baking_toggle_vote"`
	AdaptiveIssuanceVote      mavryk.FeatureVote    `json:"adaptive_issuance_vote"`

	// only present when header is fetched explicitly
	Hash     mavryk.BlockHash    `json:"hash"`
	Protocol mavryk.ProtocolHash `json:"protocol"`
	ChainId  mavryk.ChainIdHash  `json:"chain_id"`
}

func (h BlockHeader) LbVote() mavryk.FeatureVote {
	if h.LiquidityBakingToggleVote.IsValid() {
		return h.LiquidityBakingToggleVote
	}
	// sic! bool flag has opposite meaning
	if h.LiquidityBakingEscapeVote {
		return mavryk.FeatureVoteOff
	}
	return mavryk.FeatureVoteOn
}

func (h BlockHeader) AiVote() mavryk.FeatureVote {
	return h.AdaptiveIssuanceVote
}

// ProtocolData exports protocol-specific extra header fields as binary encoded data.
// Used to produce compliant block monitor data streams.
//
// octez-codec describe 018-Proxford.block_header.protocol_data binary schema
// +---------------------------------------+----------+-------------------------------------+
// | Name                                  | Size     | Contents                            |
// +=======================================+==========+=====================================+
// | payload_hash                          | 32 bytes | bytes                               |
// +---------------------------------------+----------+-------------------------------------+
// | payload_round                         | 4 bytes  | signed 32-bit integer               |
// +---------------------------------------+----------+-------------------------------------+
// | proof_of_work_nonce                   | 8 bytes  | bytes                               |
// +---------------------------------------+----------+-------------------------------------+
// | ? presence of field "seed_nonce_hash" | 1 byte   | boolean (0 for false, 255 for true) |
// +---------------------------------------+----------+-------------------------------------+
// | seed_nonce_hash                       | 32 bytes | bytes                               |
// +---------------------------------------+----------+-------------------------------------+
// | per_block_votes                       | 1 byte   | signed 8-bit integer                |
// +---------------------------------------+----------+-------------------------------------+
// | signature                             | 64 bytes | bytes                               |
// +---------------------------------------+----------+-------------------------------------+

func (h BlockHeader) ProtocolData() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(h.PayloadHash.Bytes())
	binary.Write(buf, binary.BigEndian, uint32(h.PayloadRound))
	buf.Write(h.ProofOfWorkNonce)
	if h.SeedNonceHash != nil {
		buf.WriteByte(0xff)
		buf.Write(h.SeedNonceHash.Bytes())
	} else {
		buf.WriteByte(0x0)
	}
	// broken, how to merge multiple flags is undocumented
	buf.WriteByte(h.LbVote().Tag() | (h.AiVote().Tag() << 2))
	if h.Signature.IsValid() {
		buf.Write(h.Signature.Data) // raw, no tag!
	}
	return buf.Bytes()
}

// BlockContent is part of block 1 header that seeds the initial context
type BlockContent struct {
	Command    string              `json:"command"`
	Protocol   mavryk.ProtocolHash `json:"hash"`
	Fitness    []mavryk.HexBytes   `json:"fitness"`
	Parameters *GenesisData        `json:"protocol_parameters"`
}

// OperationListLength is a part of the BlockMetadata
type OperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op"`
}

// BlockLevel is a part of BlockMetadata
type LevelInfo struct {
	Level              int64 `json:"level"`
	LevelPosition      int64 `json:"level_position"`
	Cycle              int64 `json:"cycle"`
	CyclePosition      int64 `json:"cycle_position"`
	ExpectedCommitment bool  `json:"expected_commitment"`

	// <v008
	VotingPeriod         int64 `json:"voting_period"`
	VotingPeriodPosition int64 `json:"voting_period_position"`
}

type VotingPeriod struct {
	Index         int64                   `json:"index"`
	Kind          mavryk.VotingPeriodKind `json:"kind"`
	StartPosition int64                   `json:"start_position"`
}

type VotingPeriodInfo struct {
	Position     int64        `json:"position"`
	Remaining    int64        `json:"remaining"`
	VotingPeriod VotingPeriod `json:"voting_period"`
}

// BlockMetadata is a part of the Tezos block data
type BlockMetadata struct {
	Protocol               mavryk.ProtocolHash    `json:"protocol"`
	NextProtocol           mavryk.ProtocolHash    `json:"next_protocol"`
	MaxOperationsTTL       int                    `json:"max_operations_ttl"`
	MaxOperationDataLength int                    `json:"max_operation_data_length"`
	MaxBlockHeaderLength   int                    `json:"max_block_header_length"`
	MaxOperationListLength []*OperationListLength `json:"max_operation_list_length"`
	Baker                  mavryk.Address         `json:"baker"`
	Proposer               mavryk.Address         `json:"proposer"`
	NonceHash              mavryk.NonceHash       `json:"nonce_hash"`
	ConsumedGas            int64                  `json:"consumed_gas,string"`
	Deactivated            []mavryk.Address       `json:"deactivated"`
	BalanceUpdates         BalanceUpdates         `json:"balance_updates"`

	// <v008
	Level            *LevelInfo               `json:"level"`
	VotingPeriodKind *mavryk.VotingPeriodKind `json:"voting_period_kind"`

	// v008+
	LevelInfo        *LevelInfo        `json:"level_info"`
	VotingPeriodInfo *VotingPeriodInfo `json:"voting_period_info"`

	// v010+
	ImplicitOperationsResults []ImplicitResult `json:"implicit_operations_results"`
	LiquidityBakingEscapeEma  int64            `json:"liquidity_baking_escape_ema"`

	// v015+
	ProposerConsensusKey mavryk.Address `json:"proposer_consensus_key"`
	BakerConsensusKey    mavryk.Address `json:"baker_consensus_key"`

	// v019+
	DalAttestation mavryk.Z `json:"dal_attestation"`
}

func (m *BlockMetadata) GetLevel() int64 {
	if m.LevelInfo != nil {
		return m.LevelInfo.Level
	}
	if m.Level == nil {
		return 0
	}
	return m.Level.Level
}

// GetBlock returns information about a Tezos block
// https://protocol.mavryk.org/mainnet/api/rpc.html#get-block-id
func (c *Client) GetBlock(ctx context.Context, id BlockID) (*Block, error) {
	var block Block
	u := fmt.Sprintf("chains/main/blocks/%s", id)
	if c.MetadataMode != "" {
		u += "?metadata=" + string(c.MetadataMode)
	}
	if err := c.Get(ctx, u, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// GetBlockHeight returns information about a Tezos block
// https://protocol.mavryk.org/mainnet/api/rpc.html#get-block-id
func (c *Client) GetBlockHeight(ctx context.Context, height int64) (*Block, error) {
	return c.GetBlock(ctx, BlockLevel(height))
}

// GetTips returns hashes of the current chain tip blocks, first in the array is the
// current main chain.
// https://protocol.mavryk.org/mainnet/api/rpc.html#chains-chain-id-blocks
func (c *Client) GetTips(ctx context.Context, depth int, head mavryk.BlockHash) ([][]mavryk.BlockHash, error) {
	if depth == 0 {
		depth = 1
	}
	tips := make([][]mavryk.BlockHash, 0, 10)
	var u string
	if head.IsValid() {
		u = fmt.Sprintf("chains/main/blocks?length=%d&head=%s", depth, head)
	} else {
		u = fmt.Sprintf("chains/main/blocks?length=%d", depth)
	}
	if err := c.Get(ctx, u, &tips); err != nil {
		return nil, err
	}
	return tips, nil
}

// GetHeadBlock returns the chain's head block.
// https://protocol.mavryk.org/mainnet/api/rpc.html#chains-chain-id-blocks
func (c *Client) GetHeadBlock(ctx context.Context) (*Block, error) {
	return c.GetBlock(ctx, Head)
}

// GetGenesisBlock returns main chain genesis block.
// https://protocol.mavryk.org/mainnet/api/rpc.html#chains-chain-id-blocks
func (c *Client) GetGenesisBlock(ctx context.Context) (*Block, error) {
	return c.GetBlock(ctx, Genesis)
}

// GetTipHeader returns the head block's header.
// https://protocol.mavryk.org/mainnet/api/rpc.html#chains-chain-id-blocks
func (c *Client) GetTipHeader(ctx context.Context) (*BlockHeader, error) {
	var head BlockHeader
	u := "chains/main/blocks/head/header"
	if err := c.Get(ctx, u, &head); err != nil {
		return nil, err
	}
	return &head, nil
}

// GetBlockHeader returns a block header.
// https://protocol.mavryk.org/mainnet/api/rpc.html#chains-chain-id-blocks
func (c *Client) GetBlockHeader(ctx context.Context, id BlockID) (*BlockHeader, error) {
	var head BlockHeader
	u := fmt.Sprintf("chains/main/blocks/%s/header", id)
	if err := c.Get(ctx, u, &head); err != nil {
		return nil, err
	}
	return &head, nil
}

// GetBlockMetadata returns a block metadata.
// https://protocol.mavryk.org/mainnet/api/rpc.html#chains-chain-id-blocks
func (c *Client) GetBlockMetadata(ctx context.Context, id BlockID) (*BlockMetadata, error) {
	var meta BlockMetadata
	u := fmt.Sprintf("chains/main/blocks/%s/metadata", id)
	if c.MetadataMode != "" {
		u += "?metadata=" + string(c.MetadataMode)
	}
	if err := c.Get(ctx, u, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// GetBlockHash returns the main chain's block header.
// https://protocol.mavryk.org/mainnet/api/rpc.html#chains-chain-id-blocks
func (c *Client) GetBlockHash(ctx context.Context, id BlockID) (hash mavryk.BlockHash, err error) {
	u := fmt.Sprintf("chains/main/blocks/%s/hash", id)
	err = c.Get(ctx, u, &hash)
	return
}

// GetBlockPredHashes returns count parent blocks before block with given hash.
// https://protocol.mavryk.org/mainnet/api/rpc.html#get-chains-chain-id-blocks
func (c *Client) GetBlockPredHashes(ctx context.Context, hash mavryk.BlockHash, count int) ([]mavryk.BlockHash, error) {
	if count <= 0 {
		count = 1
	}
	blockIds := make([][]mavryk.BlockHash, 0, count)
	u := fmt.Sprintf("chains/main/blocks?length=%d&head=%s", count, hash)
	if err := c.Get(ctx, u, &blockIds); err != nil {
		return nil, err
	}
	return blockIds[0], nil
}

// GetInvalidBlocks lists blocks that have been declared invalid along with the errors that led to them being declared invalid.
// https://protocol.mavryk.org/mainnet/api/rpc.html#get-chains-chain-id-invalid-blocks
func (c *Client) GetInvalidBlocks(ctx context.Context) ([]*InvalidBlock, error) {
	var invalidBlocks []*InvalidBlock
	if err := c.Get(ctx, "chains/main/invalid_blocks", &invalidBlocks); err != nil {
		return nil, err
	}
	return invalidBlocks, nil
}

// GetInvalidBlock returns a single invalid block with the errors that led to it being declared invalid.
// https://protocol.mavryk.org/mainnet/api/rpc.html#get-chains-chain-id-invalid-blocks-block-hash
func (c *Client) GetInvalidBlock(ctx context.Context, blockID mavryk.BlockHash) (*InvalidBlock, error) {
	var invalidBlock InvalidBlock
	u := fmt.Sprintf("chains/main/invalid_blocks/%s", blockID)
	if err := c.Get(ctx, u, &invalidBlock); err != nil {
		return nil, err
	}
	return &invalidBlock, nil
}
