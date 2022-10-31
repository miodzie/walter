// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package logger

import (
	"errors"
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/log"
)

type Mod struct {
	running bool
	logger  Logger
}

func New(logger Logger) *Mod {
	return &Mod{logger: logger}
}

func (mod *Mod) Name() string {
	return "logger"
}

func (mod *Mod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		err := mod.logger.Log(msg)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}

type ModFactory struct {
	DefaultLogger Logger
}

func (m ModFactory) Create(logger interface{}) (seras.Module, error) {
	if m.DefaultLogger != nil {
		return New(m.DefaultLogger), nil
	}
	l, ok := logger.(Logger)
	if !ok {
		return nil, errors.New("passed logger is not of type Logger")
	}
	return New(l), nil
}
