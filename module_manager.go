package seras

import (
	"fmt"
	"strings"
)

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
	var list []string
	for _, mod := range mods {
		list = append(list, mod.Name())
	}
	fmt.Printf("Modules: %s\n", strings.Join(list, ", "))

	return manager, nil
}

func (manager *ModuleManager) Run(stream Stream) error {
	// Init mod streams, start them up.
	for _, mod := range manager.modules {
		modStream := make(chan Message)
		manager.streams[mod.Name()] = modStream
		go mod.Start(modStream, manager.actions)
	}

	// Collect messages from stream, broadcast to mods.
	for msg := range stream {
		for _, ch := range manager.streams {
			ch <- msg
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
