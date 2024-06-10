// Copyright (c) 2020-2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package mavryk

var (
	ProtoAlpha     = MustParseProtocolHash("ProtoALphaALphaALphaALphaALphaALphaALphaALphaDdp3zK")
	ProtoGenesis   = MustParseProtocolHash("PrihK96nBAFSxVL1GLJTVhu9YnzkMFiBeuJRPA8NwuZVZCE1L6i")
	ProtoBootstrap = MustParseProtocolHash("Ps9mPmXaRzmzk35gbAYNCAw6UXdE2qoABTHbN2oEEc1qM7CwT9P")
	ProtoV001      = MustParseProtocolHash("PtAtLasomUEW99aVhVTrqjCHjJSpFUa8uHNEAEamx9v2SNeTaNp")

	// aliases
	PtAtLas = ProtoV001

	Mainnet  = MustParseChainIdHash("NetXdQprcVkpaWU")
	Basenet  = MustParseChainIdHash("NetXnHfVqm9iesp")
	Atlasnet = MustParseChainIdHash("NetXvyTAafh8goH")

	Versions = map[ProtocolHash]int{
		ProtoGenesis:   0,
		ProtoBootstrap: 0,
		ProtoV001:      18,
		ProtoAlpha:     19,
	}

	Deployments = map[ChainIdHash]ProtocolHistory{
		Mainnet: {
			{ProtoGenesis, 0, 0, 0, 0, 5, 4096, 256},   // 0
			{ProtoBootstrap, 0, 1, 1, 0, 5, 4096, 256}, // 0
			{PtAtLas, 2, 2, 28082, 0, 5, 4096, 256},    // v18
			// {PtAtLas, 0, 5070849, -1, 703, 5, 16384, 1024},       // v18
		},
		Basenet: {
			{ProtoGenesis, 0, 0, 0, 0, 3, 4096, 256},   // 0
			{ProtoBootstrap, 0, 1, 1, 0, 3, 4096, 256}, // 0
			// {Proxford, 0, 2957313, -1, 625, 3, 8192, 512},     // v18
		},
		Atlasnet: {
			{ProtoGenesis, 0, 0, 0, 0, 3, 4096, 256},   // 0
			{ProtoBootstrap, 0, 1, 1, 0, 3, 8192, 512}, // 0
			{PtAtLas, 0, 16385, -1, 2, 3, 8192, 512},   // v18
		},
	}
)

type Deployment struct {
	Protocol          ProtocolHash
	StartOffset       int64
	StartHeight       int64
	EndHeight         int64
	StartCycle        int64
	PreservedCycles   int64
	BlocksPerCycle    int64
	BlocksPerSnapshot int64
}

type ProtocolHistory []Deployment

func (h ProtocolHistory) Clone() ProtocolHistory {
	clone := make(ProtocolHistory, len(h))
	copy(clone, h)
	return clone
}

func (h ProtocolHistory) AtBlock(height int64) (d Deployment) {
	d = h.Last()
	for i := len(h) - 1; i >= 0; i-- {
		if h[i].StartHeight <= height {
			d = h[i]
			break
		}
	}
	return
}

func (h ProtocolHistory) AtCycle(cycle int64) (d Deployment) {
	d = h.Last()
	for i := len(h) - 1; i >= 0; i-- {
		if h[i].StartCycle <= cycle {
			d = h[i]
			break
		}
	}
	return
}

func (h ProtocolHistory) AtProtocol(proto ProtocolHash) (d Deployment) {
	d = h.Last()
	for _, v := range h {
		if v.Protocol == proto {
			d = v
			break
		}
	}
	return
}

func (h *ProtocolHistory) Add(d Deployment) {
	(*h) = append((*h), d)
}

func (h ProtocolHistory) Last() (d Deployment) {
	if l := len(h); l > 0 {
		d = h[l-1]
	}
	return
}
