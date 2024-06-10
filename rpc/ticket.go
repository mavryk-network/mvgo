// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/micheline"
)

type Ticket struct {
	Ticketer mavryk.Address `json:"ticketer"`
	Content  micheline.Prim `json:"content"`
	Type     micheline.Prim `json:"content_type"`
}

type TicketBalanceUpdate struct {
	Account mavryk.Address `json:"account"`
	Amount  mavryk.Z       `json:"amount"`
}

type TicketUpdate struct {
	Ticket  Ticket                `json:"ticket_token"`
	Updates []TicketBalanceUpdate `json:"updates"`
}
