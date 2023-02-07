// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package storage

import (
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/miodzie/walter/storage/sqlite"
)

var storages map[string]*sql.DB

type Config struct {
	Storage map[string]Storage
}

type Storage struct {
	Type string
	File string
}

func InitFromConfig(path string) error {
	storages = make(map[string]*sql.DB)
	var config Config
	_, err := toml.DecodeFile(path, &config)

	for name, strg := range config.Storage {
		// Factories are dumb, just switch it.
		if strg.Type == "sqlite" {
			db, err := sqlite.Setup(strg.File)
			if err != nil {
				return err
			}
			Register(name, db)
		}
	}

	return err
}

func Register(name string, db *sql.DB) {
	storages[name] = db
}

func Get(name string) (*sql.DB, error) {
	db, ok := storages[name]
	if !ok {
		return nil, fmt.Errorf("database `%s` not registered", name)
	}

	return db, nil
}
