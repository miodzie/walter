// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package irc

import "github.com/miodzie/seras"

// seras.Modable

func (con *Connection) Modules() []seras.Module {
	return con.mods
}
func (con *Connection) AddMods(mods []seras.Module) {
	con.mods = append(con.mods, mods...)
}

func (con *Connection) ModList() map[string]interface{} {
	return con.Mods
}
