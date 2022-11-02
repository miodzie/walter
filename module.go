// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package walter

type Module interface {
	Name() string
	Start(Stream, Actions) error
	Stop()
}

type Actions interface {
	Messenger
	MessageFormatter
	Admin
}

type ModuleManager struct {
	modules []Module
	actions Actions
	streams map[string]chan Message
}

func NewModManager(mods []Module, actions Actions) (*ModuleManager, error) {
	manager := &ModuleManager{
		modules: mods,
		actions: actions,
		streams: make(map[string]chan Message),
	}

	return manager, nil
}

func (manager *ModuleManager) Run(stream Stream) error {
	for _, mod := range manager.modules {
		modStream := make(chan Message)
		manager.streams[mod.Name()] = modStream
		go mod.Start(modStream, manager.actions)
	}

	for msg := range stream {
		for _, ch := range manager.streams {
			go func(ch chan Message) { ch <- msg }(ch)
		}
	}

	return nil
}

func (manager *ModuleManager) Stop() {
	for _, mod := range manager.modules {
		mod.Stop()
	}
	for _, stream := range manager.streams {
		close(stream)
	}
}
