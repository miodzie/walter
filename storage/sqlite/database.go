package sqlite

import (
	"database/sql"

	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed rss.sql
var migration string

func Setup(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Configure sqlite3 to be le dankest.
	_, err = db.Exec("PRAGMA journal_mode = WAL; PRAGMA busy_timeout = 5000; PRAGMA foreign_keys = ON;")
	if err != nil {
		return db, err
	}

	return db, migrate(db)
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(migration)

	return err
}
