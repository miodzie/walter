package irc

import "github.com/miodzie/seras"

// seras.Modable
func (con *Connection) Mods() []seras.Module {
	return con.mods
}
func (con *Connection) AddMods(mods []seras.Module) {
	con.mods = append(con.mods, mods...)
}

func (con *Connection) ModList() map[string]interface{} {
	return con.config.ModConfig
	mods := make(map[string]interface{})
	for _, m := range con.config.Mods {
		mods[m] = struct{}{}
	}

	return mods
}
