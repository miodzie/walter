// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package mods

import (
	"github.com/miodzie/seras"
)

// List [mod name]config
type List map[string]interface{}

func CreateFromList(list List) ([]seras.Module, error) {
	var mods []seras.Module

	for name, config := range list {
		mod, err := MakeFromConfig(name, config)
		if err != nil {
			return nil, err
		}
		mods = append(mods, mod)
	}

	return mods, nil
}
