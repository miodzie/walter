package sqlite

import (
	"database/sql"

	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

//go:embed feeds.sql
var migration string

func Setup(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	return migrate()
}

func migrate() error {
	_, err := db.Exec(migration)

	return err
}
