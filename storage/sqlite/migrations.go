package sqlite

import (
	"database/sql"
	"github.com/miodzie/walter/mods/rss/ports/database"
)

func Migrate(db *sql.DB) error {
	_, err := db.Exec(database.Migration)

	return err
}
