package seras

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ModuleManager struct {
	modules []Module
	streams []chan Message
	actions Actions
	dbs     map[string]*gorm.DB
}

func NewModManager(mods []Module, actions Actions) (*ModuleManager, error) {
	manager := &ModuleManager{
		modules: mods,
		actions: actions,
		dbs:     make(map[string]*gorm.DB),
	}
	// Init mod databases.
	for _, mod := range mods {
		// Defaulting to sqlite3 is fine for now.
		path := fmt.Sprintf("storage/%s.db", mod.Name())
		db, err := gorm.Open(sqlite.Open(path))
		if err != nil {
			return &ModuleManager{}, err
		}
		if _, ok := manager.dbs[mod.Name()]; ok {
			return &ModuleManager{}, errors.New("duplicate module name")
		}
		manager.dbs[mod.Name()] = db
		// mod.setDB(db)
		fmt.Printf("Module loaded: %s\n", mod.Name())
	}

	return manager, nil
}

func (manager *ModuleManager) Run(stream Stream) error {
	// Init mod streams, start them up.
	for _, mod := range manager.modules {
		modStream := make(chan Message)
		manager.streams = append(manager.streams, modStream)
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
