package mavryk_test

import (
	"testing"

	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/rpc"
)

type (
	Block            = rpc.Block
	BlockHeader      = rpc.BlockHeader
	BlockMetadata    = rpc.BlockMetadata
	VotingPeriodInfo = rpc.VotingPeriodInfo
	LevelInfo        = rpc.LevelInfo
)

var (
	ProtoGenesis   = mavryk.ProtoGenesis
	ProtoBootstrap = mavryk.ProtoBootstrap
	PtAtLas        = mavryk.PtAtLas
	PtBoreas       = mavryk.PtBoreas

	Mainnet     = mavryk.Mainnet
	NewParams   = mavryk.NewParams
	Deployments = mavryk.Deployments
)

func TestParams(t *testing.T) {
	var lastProto mavryk.ProtocolHash

	// walk test blocks
	for _, v := range paramBlocks {
		// update test state
		isProtoUpgrade := !lastProto.Equal(v.Protocol)
		if isProtoUpgrade {
			lastProto = v.Protocol
		}

		// prepare block
		block := Block{
			Protocol: v.Protocol,
			ChainId:  Mainnet,
			Header: BlockHeader{
				Level: v.LevelInfo.Level,
			},
			Metadata: v,
		}
		height, cycle := block.GetLevel(), block.GetCycle()

		// prepare params
		p := NewParams().
			WithChainId(Mainnet).
			WithProtocol(v.Protocol).
			WithDeployment(Deployments[Mainnet].AtProtocol(v.Protocol))

		// load expected result
		check := paramResults[height]

		checkParams(t, p, height, cycle, check)
	}
}

func TestParamsStatic(t *testing.T) {
	for height, check := range paramResults {
		p := NewParams().WithChainId(Mainnet).AtBlock(height)
		checkParams(t, p, height, check.Cycle, check)
	}
}

func TestDefaultParams(t *testing.T) {
	for n, p := range map[string]*mavryk.Params{
		"main": mavryk.DefaultParams,
		"base": mavryk.BasenetParams,
	} {
		if p.Network == "" {
			t.Errorf("%s params: Empty network name", n)
		}
		if !p.ChainId.IsValid() {
			t.Errorf("%s params: zero network id", n)
		}
		if !p.Protocol.IsValid() {
			t.Errorf("%s params: zero protocol", n)
		}
		if have, want := p.Version, mavryk.Versions[p.Protocol]; have != want {
			t.Errorf("%s params: version mismatch: have=%d want=%d", n, have, want)
		}
		if p.MinimalBlockDelay == 0 {
			t.Errorf("%s params: zero MinimalBlockDelay", n)
		}
		if p.CostPerByte == 0 {
			t.Errorf("%s params: zero CostPerByte", n)
		}
		if p.OriginationSize == 0 {
			t.Errorf("%s params: zero OriginationSize", n)
		}
		if p.BlocksPerCycle == 0 {
			t.Errorf("%s params: zero BlocksPerCycle", n)
		}
		if p.ConsensusRightsDelay == 0 {
			t.Errorf("%s params: zero PreservedCycles", n)
		}
		if p.BlocksPerSnapshot == 0 {
			t.Errorf("%s params: zero BlocksPerSnapshot", n)
		}
		if p.HardGasLimitPerOperation == 0 {
			t.Errorf("%s params: zero HardGasLimitPerOperation", n)
		}
		if p.HardGasLimitPerBlock == 0 {
			t.Errorf("%s params: zero HardGasLimitPerBlock", n)
		}
		if p.HardStorageLimitPerOperation == 0 {
			t.Errorf("%s params: zero HardStorageLimitPerOperation", n)
		}
		if p.MaxOperationDataLength == 0 {
			t.Errorf("%s params: zero MaxOperationDataLength", n)
		}
		if p.MaxOperationsTTL == 0 {
			t.Errorf("%s params: zero MaxOperationsTTL", n)
		}
		if p.MaxOperationDataLength == 0 {
			t.Errorf("%s params: zero MaxOperationDataLength", n)
		}
		if p.OperationTagsVersion < 0 || p.OperationTagsVersion > 3 {
			t.Errorf("%s params: unknown OperationTagsVersion %d", n, p.OperationTagsVersion)
		}
		if p.StartHeight == 0 {
			t.Errorf("%s params: zero StartHeight", n)
		}
		if p.EndHeight == 0 {
			t.Errorf("%s params: zero EndHeight", n)
		}
		if p.StartHeight > p.BlocksPerCycle && p.StartCycle == 0 {
			t.Errorf("%s params: zero StartCycle", n)
		}
	}
}

func checkParams(t *testing.T, p *mavryk.Params, height, cycle int64, check paramResult) {
	// test param functions
	if !p.ContainsHeight(height) {
		t.Errorf("v%03d ContainsHeight(%d) failed", p.Version, height)
	}
	if !p.ContainsCycle(cycle) {
		t.Errorf("v%03d %d ContainsCycle(%d) failed", p.Version, height, cycle)
	}
	if have, want := p.IsCycleStart(height), check.IsCycleStart(); have != want {
		t.Errorf("v%03d IsCycleStart(%d) mismatch: have=%t want=%t", p.Version, height, have, want)
	}
	if have, want := p.IsCycleEnd(height), check.IsCycleEnd(); have != want {
		t.Errorf("v%03d IsCycleEnd(%d) mismatch: have=%t want=%t", p.Version, height, have, want)
	}
	if have, want := p.CycleFromHeight(height), check.Cycle; have != want {
		t.Errorf("v%03d CycleFromHeight(%d) mismatch: have=%d want=%d", p.Version, height, have, want)
	}
	cstart := p.CycleStartHeight(cycle)
	cend := p.CycleEndHeight(cycle)
	cpos := p.CyclePosition(height)
	if cstart < 0 {
		t.Errorf("v%03d %d negative cycle start %d", p.Version, height, cstart)
	}
	if cend < 0 {
		t.Errorf("v%03d %d negative cycle end %d", p.Version, height, cend)
	}
	if cpos < 0 {
		t.Errorf("v%03d %d negative cycle pos %d", p.Version, height, cpos)
	}
	if cstart >= cend {
		t.Errorf("v%03d %d cycle start %d > end %d", p.Version, height, cstart, cend)
	}
	if cstart+cpos != height {
		t.Errorf("v%03d %d cycle pos %d + start %d != height", p.Version, height, cstart, cpos)
	}

	if have, want := p.IsSnapshotBlock(height), check.IsSnapshot(); have != want {
		t.Errorf("v%03d IsSnapshotBlock(%d) mismatch: have=%t want=%t", p.Version, height, have, want)
	}
	if have, want := p.SnapshotIndex(height), check.Snap; have != want {
		t.Errorf("v%03d SnapshotIndex(%d) mismatch: have=%d want=%d", p.Version, height, have, want)
	}
	if have, want := p.SnapshotBlock(cycle, 0), height; have > want {
		t.Errorf("v%03d SnapshotBlock(%d) mismatch: have=%d > want=%d", p.Version, height, have, want)
	}
}

type paramResult struct {
	Cycle int64
	Snap  int
	Flags byte // 16 Snapshot | 8 CycleStart | 4 CycleEnd | 2 VoteStart | 1 VoteEnd
}

func (p paramResult) IsSnapshot() bool {
	return (p.Flags>>4)&0x1 > 0
}

func (p paramResult) IsCycleStart() bool {
	return (p.Flags>>3)&0x1 > 0
}

func (p paramResult) IsCycleEnd() bool {
	return (p.Flags>>2)&0x1 > 0
}

func (p paramResult) IsVoteStart() bool {
	return (p.Flags>>1)&0x1 > 0
}

func (p paramResult) IsVoteEnd() bool {
	return p.Flags&0x1 > 0
}

var paramResults = map[int64]paramResult{
	0:       {0, -1, 0},            // genesis
	1:       {0, -1, 8},            // bootstrap
	5070849: {703, -1, 8 + 2},      // v018 start
	5726208: {742, 15, 16 + 4 + 1}, // --> end
	5726209: {743, 15, 8 + 2},      // v019 start
}

var paramBlocks = []BlockMetadata{
	{
		// genesis
		Protocol:         ProtoGenesis,
		NextProtocol:     ProtoBootstrap,
		LevelInfo:        &LevelInfo{},
		VotingPeriodInfo: &VotingPeriodInfo{},
	}, {
		// bootstrap
		Protocol:     ProtoBootstrap,
		NextProtocol: PtAtLas,
		LevelInfo: &LevelInfo{
			Level: 1,
		},
		VotingPeriodInfo: &VotingPeriodInfo{},
	}, {
		// v18 start
		Protocol:     PtAtLas,
		NextProtocol: PtAtLas,
		LevelInfo: &LevelInfo{
			Level:              5070849,
			Cycle:              703,
			CyclePosition:      0,
			ExpectedCommitment: false,
		},
		VotingPeriodInfo: &VotingPeriodInfo{
			Position:  0,
			Remaining: 81912,
		},
	}, {
		// v18 end
		Protocol:     PtAtLas,
		NextProtocol: PtBoreas,
		LevelInfo: &LevelInfo{
			Level:              5726208,
			Cycle:              742,
			CyclePosition:      16383,
			ExpectedCommitment: true,
		},
		VotingPeriodInfo: &VotingPeriodInfo{
			Position:  81912,
			Remaining: 0,
		},
	}, {
		// v19 start
		Protocol:     PtBoreas,
		NextProtocol: PtBoreas,
		LevelInfo: &LevelInfo{
			Level:              5726209,
			Cycle:              743,
			CyclePosition:      0,
			ExpectedCommitment: false,
		},
		VotingPeriodInfo: &VotingPeriodInfo{
			Position:  0,
			Remaining: 81912,
		},
	},
}
