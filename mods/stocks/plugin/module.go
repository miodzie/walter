// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"walter"
	"walter/mods/stocks"
)

type StockMod struct {
	stocks.Ticker
	walter.Actions
	running bool
}

func New(ticker stocks.Ticker) *StockMod {
	return &StockMod{Ticker: ticker}
}

func (mod *StockMod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	mod.Actions = actions
	for mod.running {
		msg := <-stream
		msg.Command("stock", mod.getPrice)
	}

	return nil
}
