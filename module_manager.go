package seras

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: implement management of dbs
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
		dbs:     make(map[string]*sql.DB),
	}
	// Init mod databases.
	for _, mod := range mods {
		fmt.Println(mod.Name())
		// Defaulting to sqlite3 is fine for now.
		path := fmt.Sprintf("storage/%s.sqlite", mod.Name())
		db, err := sql.Open("sqlite3", path)
		if err != nil {
			// TODO: Return err
			panic(err)
		}
		// TODO: Check for duplicate names.
		manager.dbs[mod.Name()] = db
		mod.setDB(db)
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
