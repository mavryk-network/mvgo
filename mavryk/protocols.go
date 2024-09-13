// Copyright (c) 2020-2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package mavryk

import "sync"

var (
	ProtoAlpha     = MustParseProtocolHash("ProtoALphaALphaALphaALphaALphaALphaALphaALphaDdp3zK")
	ProtoGenesis   = MustParseProtocolHash("PrihK96nBAFSxVL1GLJTVhu9YnzkMFiBeuJRPA8NwuZVZCE1L6i")
	ProtoBootstrap = MustParseProtocolHash("Ps9mPmXaRzmzk35gbAYNCAw6UXdE2qoABTHbN2oEEc1qM7CwT9P")
	ProtoV001      = MustParseProtocolHash("PtAtLasomUEW99aVhVTrqjCHjJSpFUa8uHNEAEamx9v2SNeTaNp")
	ProtoV002      = MustParseProtocolHash("Pt8h9rz3r9F3Yx3wzJqF42sxsyVoo6kL4FBoJSWRzKmDvjXjHwV")

	// aliases
	PtAtLas  = ProtoV001
	PtBoreas = ProtoV002

	Mainnet   = MustParseChainIdHash("NetXdQprcVkpaWU")
	Basenet   = MustParseChainIdHash("NetXnHfVqm9iesp")
	Atlasnet  = MustParseChainIdHash("NetXvyTAafh8goH")
	Boreasnet = MustParseChainIdHash("NetXvyTAafh8goH")

	versionsMtx = sync.RWMutex{}
	Versions    = map[ProtocolHash]int{
		ProtoGenesis:   0,
		ProtoBootstrap: 0,
		ProtoV001:      18,
		ProtoV002:      19,
		ProtoAlpha:     19,
	}

	Deployments = map[ChainIdHash]ProtocolHistory{
		Mainnet: {
			{ProtoGenesis, 0, 0, 0, 0, 5, 4096, 256},   // 0
			{ProtoBootstrap, 0, 1, 1, 0, 5, 4096, 256}, // 0
			// {PtAtLas, 2, 2, 28082, 0, 5, 4096, 256},    // v18
			{PtAtLas, 0, 5070849, 5726208, 703, 5, 16384, 1024}, // v18
			{PtBoreas, 0, 5726209, -1, 743, 2, 24576, 24576},    // v19
		},
		Basenet: {
			{ProtoGenesis, 0, 0, 0, 0, 3, 4096, 256},          // 0
			{ProtoBootstrap, 0, 1, 1, 0, 3, 4096, 256},        // 0
			{PtAtLas, 0, 5316609, 6422528, 913, 3, 8192, 512}, // v18
			{PtBoreas, 0, 6422529, -1, 1048, 2, 12288, 12288}, // v19
		},
		Atlasnet: {
			{ProtoGenesis, 0, 0, 0, 0, 3, 4096, 256},   // 0
			{ProtoBootstrap, 0, 1, 1, 0, 3, 8192, 512}, // 0
			{PtAtLas, 0, 16385, -1, 2, 3, 8192, 512},   // v18
		},
		Boreasnet: {
			{ProtoGenesis, 0, 0, 0, 0, 3, 8192, 512},    // 0
			{ProtoBootstrap, 0, 1, 1, 0, 3, 8192, 512},  // 0
			{PtAtLas, 2, 2, 8192, 0, 3, 8192, 512},      // v18
			{PtBoreas, 0, 8193, -1, 1, 2, 12288, 12288}, // v19
		},
	}
)

type Deployment struct {
	Protocol             ProtocolHash
	StartOffset          int64
	StartHeight          int64
	EndHeight            int64
	StartCycle           int64
	ConsensusRightsDelay int64
	BlocksPerCycle       int64
	BlocksPerSnapshot    int64
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
