// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"fmt"
	"walter"
)

func (mod *StockMod) getPrice(msg walter.Message) {
	if len(msg.Arguments) <= 1 {
		return
	}
	security := msg.Arguments[1]
	price, _ := mod.PriceFor(security)

	mod.Reply(msg, fmt.Sprintf("$%s: %.2f", security, price))
}
