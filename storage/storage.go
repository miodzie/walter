package storage

import (
	"database/sql"
	"fmt"
)

var storages map[string]*sql.DB

func InitFromConfig() {

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
