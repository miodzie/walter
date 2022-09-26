package mods

import (
	"fmt"
	"github.com/miodzie/seras"
	"net/http"
)

type ModConfig interface{}

// [name]createMod
var factories map[string]func(ModConfig) (seras.Module, error)

func init() {
	factories = make(map[string]func(ModConfig) (seras.Module, error))
}

func Register(name string, creator func(ModConfig) (seras.Module, error)) {
	factories[name] = creator
	http.Client{}
}

func Make(name string, config ModConfig) (seras.Module, error) {
	f, ok := factories[name]
	if !ok {
		return nil, fmt.Errorf("unknown module: `%s`", name)
	}

	return f(config)
}
