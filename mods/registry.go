package mods

import (
	"fmt"
	"github.com/miodzie/seras"
)

type Factory interface {
	Create(config interface{}) (seras.Module, error)
}

var factories map[string]Factory

func init() {
}

func Register(name string, factory Factory) {
	factories[name] = factory
}

func Make(name string, config interface{}) (seras.Module, error) {
	f, ok := factories[name]
	if !ok {
		return nil, fmt.Errorf("unknown module: `%s`", name)
	}

	return f.Create(config)
}
