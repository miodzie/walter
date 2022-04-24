package seras

import "database/sql"

type ModuleManager struct {
	modules []Module
	streams []chan Message
	actions Actions
	dbs     map[string]*sql.DB
}

func NewModManager(mods []Module, actions Actions) *ModuleManager {
	manager := &ModuleManager{
		modules: mods,
		actions: actions,
	}

	return manager
}

func (manager *ModuleManager) Run(stream Stream) {
	// Init mod streams, start them up.
	for _, mod := range manager.modules {
		modStream := make(chan Message)
		manager.streams = append(manager.streams, modStream)
		mod.Loop(modStream, manager.actions)
	}

	// Collect messages from stream, broadcast to mods.
	for msg := range stream {
		for _, ch := range manager.streams {
			ch <- msg
		}
	}
}

func (manager *ModuleManager) Stop() {
	for _, mod := range manager.modules {
		mod.Stop()
	}
	for _, stream := range manager.streams {
		close(stream)
	}
}

