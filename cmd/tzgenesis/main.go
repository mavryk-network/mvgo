package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mavryk-network/mvgo/mavryk"
)

var (
	proto = mavryk.MustParseProtocolHash("ProtoGenesisGenesisGenesisGenesisGenesisGenesk612im")
	key   = mavryk.MustParseKey("edpkuSLWfVU1Vq7Jg9FucPyKmma6otcMHac9zG4oU1KMHSTBpJuGQ2")
	block mavryk.BlockHash
	name  = "TEZOS"
	tm    string
)

func init() {
	rand.Read(block[:])
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
	}
}
func run() error {
	flag.Var(&proto, "proto", "Genesis protocol hash")
	flag.Var(&key, "key", "Genesis pubkey")
	flag.Var(&block, "block", "Genesis block")
	flag.StringVar(&name, "name", "TEZOS", "Chain name")
	flag.StringVar(&tm, "time", time.Now().UTC().Format(time.RFC3339), "Genesis timestamp")
	flag.Parse()

	ts, err := time.Parse(time.RFC3339, tm)
	if err != nil {
		return fmt.Errorf("Parsing timestamp %q: %v", tm, err)
	}

	genesis := Genesis{
		Genesis: GenesisInfo{
			Timestamp: ts.Format(time.RFC3339),
			Block:     block,
			Protocol:  proto,
		},
		Params: GenesisParams{
			Values: GenesisValues{
				Key: key,
			},
		},
		ChainName:        name,
		SandboxChainName: "SANDBOXED_" + name,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(genesis)
}

type Genesis struct {
	Genesis          GenesisInfo   `json:"genesis"`
	Params           GenesisParams `json:"genesis_parameters"`
	ChainName        string        `json:"chain_name"`
	SandboxChainName string        `json:"sandboxed_chain_name"`
}

type GenesisInfo struct {
	Timestamp string              `json:"timestamp"`
	Block     mavryk.BlockHash    `json:"block"`
	Protocol  mavryk.ProtocolHash `json:"protocol"`
}

type GenesisParams struct {
	Values GenesisValues `json:"values"`
}

type GenesisValues struct {
	Key mavryk.Key `json:"genesis_pubkey"`
}
