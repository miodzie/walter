// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package mods

import (
	"fmt"
	"github.com/miodzie/walter"
	art "github.com/miodzie/walter/mods/art/plugin"
	"github.com/miodzie/walter/mods/dong"
	"github.com/miodzie/walter/mods/logger"
	"github.com/miodzie/walter/mods/logger/drivers"
	rss_plugin "github.com/miodzie/walter/mods/rss/plugin"
	sed "github.com/miodzie/walter/mods/sed/plugin"
)

var factories map[string]Factory

type Factory interface {
	Create(config any) (walter.Module, error)
}

func init() {
	factories = make(map[string]Factory)
	Register("dong", &dong.ModFactory{})
	Register("sed", &sed.ModFactory{})
	Register("logger", &logger.ModFactory{DefaultLogger: drivers.NewMultiLogger(drivers.ConsoleLogger{})})
	Register("rss", &rss_plugin.ModFactory{})
	Register("visionary", &art.VisionaryFactory{})
	Register("artist", &art.ArtistFactory{})
}

func Register(name string, factory Factory) {
	factories[name] = factory
}

func MakeFromConfig(name string, config interface{}) (walter.Module, error) {
	f, ok := factories[name]
	if !ok {
		return nil, fmt.Errorf("unknown module: `%s`", name)
	}

	return f.Create(config)
}
