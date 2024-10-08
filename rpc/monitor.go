// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"time"

	"github.com/mavryk-network/mvgo/mavryk"
)

var ErrMonitorClosed = errors.New("monitor closed")

type Monitor interface {
	New() interface{}
	Send(ctx context.Context, val interface{})
	Err(error)
	Closed() <-chan struct{}
	Close()
}

// BootstrappedBlock represents bootstrapped block stream message
type BootstrappedBlock struct {
	Block     mavryk.BlockHash `json:"block"`
	Timestamp time.Time        `json:"timestamp"`
}

type BootstrapMonitor struct {
	result chan *BootstrappedBlock
	closed chan struct{}
	err    error
}

// make sure BootstrapMonitor implements Monitor interface
var _ Monitor = (*BootstrapMonitor)(nil)

func NewBootstrapMonitor() *BootstrapMonitor {
	return &BootstrapMonitor{
		result: make(chan *BootstrappedBlock),
		closed: make(chan struct{}),
	}
}

func (m *BootstrapMonitor) New() interface{} {
	return &BootstrappedBlock{}
}

func (m *BootstrapMonitor) Send(ctx context.Context, val interface{}) {
	select {
	case <-m.closed:
		return
	default:
	}
	select {
	case <-ctx.Done():
	case <-m.closed:
	case m.result <- val.(*BootstrappedBlock):
	}
}

func (m *BootstrapMonitor) Recv(ctx context.Context) (*BootstrappedBlock, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-m.closed:
		err := m.err
		if err == nil {
			err = ErrMonitorClosed
		}
		return nil, err
	case res, ok := <-m.result:
		if !ok {
			if m.err != nil {
				return nil, m.err
			}
			return nil, io.EOF
		}
		return res, nil
	}
}

func (m *BootstrapMonitor) Err(err error) {
	m.err = err
	m.Close()
}

func (m *BootstrapMonitor) Closed() <-chan struct{} {
	return m.closed
}

func (m *BootstrapMonitor) Close() {
	select {
	case <-m.closed:
		return
	default:
	}
	close(m.closed)
	close(m.result)
}

// BlockHeaderLogEntry is a log entry returned for a new block when monitoring
type BlockHeaderLogEntry struct {
	Hash           mavryk.BlockHash      `json:"hash"`
	Level          int64                 `json:"level"`
	Proto          int                   `json:"proto"`
	Predecessor    mavryk.BlockHash      `json:"predecessor"`
	Timestamp      time.Time             `json:"timestamp"`
	ValidationPass int                   `json:"validation_pass"`
	OperationsHash mavryk.OpListListHash `json:"operations_hash"`
	Fitness        []mavryk.HexBytes     `json:"fitness"`
	Context        mavryk.ContextHash    `json:"context"`
	ProtocolData   mavryk.HexBytes       `json:"protocol_data"`
}

func (h *BlockHeader) LogEntry() *BlockHeaderLogEntry {
	return &BlockHeaderLogEntry{
		Hash:           h.Hash,
		Level:          h.Level,
		Proto:          h.Proto,
		Predecessor:    h.Predecessor,
		Timestamp:      h.Timestamp,
		ValidationPass: h.ValidationPass,
		OperationsHash: h.OperationsHash,
		Fitness:        h.Fitness,
		Context:        h.Context,
		ProtocolData:   mavryk.HexBytes(h.ProtocolData()),
	}
}

func (l BlockHeaderLogEntry) Round() int {
	return int(binary.BigEndian.Uint32(l.ProtocolData[32:]))
}

func (l BlockHeaderLogEntry) PayloadHash() (h mavryk.PayloadHash) {
	copy(h[:], l.ProtocolData[:])
	return
}

func (l BlockHeaderLogEntry) Pow() (h mavryk.HexBytes) {
	h.UnmarshalBinary(l.ProtocolData[36:44])
	return
}

func (b *Block) LogEntry() *BlockHeaderLogEntry {
	e := b.Header.LogEntry()
	e.Hash = b.Hash
	return e
}

type BlockHeaderMonitor struct {
	result chan *BlockHeaderLogEntry
	closed chan struct{}
	err    error
}

// make sure BlockHeaderMonitor implements Monitor interface
var _ Monitor = (*BlockHeaderMonitor)(nil)

func NewBlockHeaderMonitor() *BlockHeaderMonitor {
	return &BlockHeaderMonitor{
		result: make(chan *BlockHeaderLogEntry),
		closed: make(chan struct{}),
	}
}

func (m *BlockHeaderMonitor) New() interface{} {
	return &BlockHeaderLogEntry{}
}

func (m *BlockHeaderMonitor) Send(ctx context.Context, val interface{}) {
	select {
	case <-m.closed:
		return
	default:
	}
	select {
	case <-ctx.Done():
	case <-m.closed:
	case m.result <- val.(*BlockHeaderLogEntry):
	}
}

func (m *BlockHeaderMonitor) Recv(ctx context.Context) (*BlockHeaderLogEntry, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-m.closed:
		err := m.err
		if err == nil {
			err = ErrMonitorClosed
		}
		return nil, err
	case res, ok := <-m.result:
		if !ok {
			if m.err != nil {
				return nil, m.err
			}
			return nil, io.EOF
		}
		return res, nil
	}
}

func (m *BlockHeaderMonitor) Err(err error) {
	m.err = err
	m.Close()
}

func (m *BlockHeaderMonitor) Close() {
	select {
	case <-m.closed:
		return
	default:
	}
	close(m.closed)
	close(m.result)
}

func (m *BlockHeaderMonitor) Closed() <-chan struct{} {
	return m.closed
}

// MempoolMonitor is a monitor for the Tezos mempool. Note that the connection
// resets every time a new head is attached to the chain. MempoolMonitor is
// closed with an error in this case and cannot be reused after close.
//
// The Tezos mempool re-evaluates all operations and potentially updates their state
// when the head block changes. This applies to operations in lists branch_delayed
// and branch_refused. After reorg, operations already included in a previous block
// may enter the mempool again.
type MempoolMonitor struct {
	result chan *[]*Operation
	closed chan struct{}
	err    error
}

// make sure MempoolMonitor implements Monitor interface
var _ Monitor = (*MempoolMonitor)(nil)

func NewMempoolMonitor() *MempoolMonitor {
	return &MempoolMonitor{
		result: make(chan *[]*Operation),
		closed: make(chan struct{}),
	}
}

func (m *MempoolMonitor) New() interface{} {
	slice := make([]*Operation, 0)
	return &slice
}

func (m *MempoolMonitor) Send(ctx context.Context, val interface{}) {
	select {
	case <-m.closed:
		return
	default:
	}
	select {
	case <-ctx.Done():
	case <-m.closed:
	case m.result <- val.(*[]*Operation):
	}
}

func (m *MempoolMonitor) Recv(ctx context.Context) ([]*Operation, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-m.closed:
		err := m.err
		if err == nil {
			err = ErrMonitorClosed
		}
		return nil, err
	case res, ok := <-m.result:
		if !ok {
			if m.err != nil {
				return nil, m.err
			}
			return nil, io.EOF
		}
		return *res, nil
	}
}

func (m *MempoolMonitor) Err(err error) {
	m.err = err
	m.Close()
}

func (m *MempoolMonitor) Close() {
	select {
	case <-m.closed:
		return
	default:
	}
	close(m.closed)
	close(m.result)
}

func (m *MempoolMonitor) Closed() <-chan struct{} {
	return m.closed
}

// NetworkPeerLogEntry represents peer log entry
type NetworkPeerLogEntry struct {
	NetworkAddress
	Kind      string    `json:"kind"`
	Timestamp time.Time `json:"timestamp"`
}

type NetworkPeerMonitor struct {
	result chan *NetworkPeerLogEntry
	closed chan struct{}
	err    error
}

// make sure NetworkPeerMonitor implements Monitor interface
var _ Monitor = (*NetworkPeerMonitor)(nil)

func NewNetworkPeerMonitor() *NetworkPeerMonitor {
	return &NetworkPeerMonitor{
		result: make(chan *NetworkPeerLogEntry),
		closed: make(chan struct{}),
	}
}

func (m *NetworkPeerMonitor) New() interface{} {
	return &NetworkPeerLogEntry{}
}

func (m *NetworkPeerMonitor) Send(ctx context.Context, val interface{}) {
	select {
	case <-m.closed:
		return
	default:
	}
	select {
	case <-ctx.Done():
	case <-m.closed:
	case m.result <- val.(*NetworkPeerLogEntry):
	}
}

func (m *NetworkPeerMonitor) Recv(ctx context.Context) (*NetworkPeerLogEntry, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-m.closed:
		err := m.err
		if err == nil {
			err = ErrMonitorClosed
		}
		return nil, err
	case res, ok := <-m.result:
		if !ok {
			if m.err != nil {
				return nil, m.err
			}
			return nil, io.EOF
		}
		return res, nil
	}
}

func (m *NetworkPeerMonitor) Err(err error) {
	m.err = err
	m.Close()
}

func (m *NetworkPeerMonitor) Close() {
	select {
	case <-m.closed:
		return
	default:
	}
	close(m.closed)
	close(m.result)
}

func (m *NetworkPeerMonitor) Closed() <-chan struct{} {
	return m.closed
}

// NetworkPointLogEntry represents point's log entry
type NetworkPointLogEntry struct {
	Kind      NetworkPointState `json:"kind"`
	Timestamp time.Time         `json:"timestamp"`
}

type NetworkPointMonitor struct {
	result chan *NetworkPointLogEntry
	closed chan struct{}
	err    error
}

// make sure NetworkPointMonitor implements Monitor interface
var _ Monitor = (*NetworkPointMonitor)(nil)

func NewNetworkPointMonitor() *NetworkPointMonitor {
	return &NetworkPointMonitor{
		result: make(chan *NetworkPointLogEntry),
		closed: make(chan struct{}),
	}
}

func (m *NetworkPointMonitor) New() interface{} {
	return &NetworkPointLogEntry{}
}

func (m *NetworkPointMonitor) Send(ctx context.Context, val interface{}) {
	select {
	case <-m.closed:
		return
	default:
	}
	select {
	case <-ctx.Done():
	case <-m.closed:
	case m.result <- val.(*NetworkPointLogEntry):
	}
}

func (m *NetworkPointMonitor) Recv(ctx context.Context) (*NetworkPointLogEntry, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-m.closed:
		err := m.err
		if err == nil {
			err = ErrMonitorClosed
		}
		return nil, err
	case res, ok := <-m.result:
		if !ok {
			if m.err != nil {
				return nil, m.err
			}
			return nil, io.EOF
		}
		return res, nil
	}
}

func (m *NetworkPointMonitor) Err(err error) {
	m.err = err
	m.Close()
}

func (m *NetworkPointMonitor) Close() {
	select {
	case <-m.closed:
		return
	default:
	}
	close(m.closed)
	close(m.result)
}

func (m *NetworkPointMonitor) Closed() <-chan struct{} {
	return m.closed
}

// MonitorBootstrapped reads from the bootstrapped blocks stream http://protocol.mavryk.org/mainnet/api/rpc.html#get-monitor-bootstrapped
func (c *Client) MonitorBootstrapped(ctx context.Context, monitor *BootstrapMonitor) error {
	return c.GetAsync(ctx, "monitor/bootstrapped", monitor)
}

// MonitorBlockHeader reads from the chain heads stream http://protocol.mavryk.org/mainnet/api/rpc.html#get-monitor-heads-chain-id
func (c *Client) MonitorBlockHeader(ctx context.Context, monitor *BlockHeaderMonitor) error {
	return c.GetAsync(ctx, "monitor/heads/main", monitor)
}

// MonitorMempool reads from the chain heads stream http://protocol.mavryk.org/mainnet/api/rpc.html#get-monitor-heads-chain-id
func (c *Client) MonitorMempool(ctx context.Context, monitor *MempoolMonitor) error {
	return c.GetAsync(ctx, "chains/main/mempool/monitor_operations", monitor)
}

// MonitorNetworkPointLog monitors network events related to an `IP:addr`.
// https://protocol.mavryk.org/mainnet/api/rpc.html#get-network-peers-peer-id-log
func (c *Client) MonitorNetworkPointLog(ctx context.Context, address string, monitor *NetworkPointMonitor) error {
	return c.GetAsync(ctx, "network/points/"+address+"/log?monitor", monitor)
}

// MonitorNetworkPeerLog monitors network events related to a given peer.
// https://protocol.mavryk.org/mainnet/api/rpc.html#get-network-peers-peer-id-log
func (c *Client) MonitorNetworkPeerLog(ctx context.Context, peerID string, monitor *NetworkPeerMonitor) error {
	return c.GetAsync(ctx, "network/peers/"+peerID+"/log?monitor", monitor)
}
