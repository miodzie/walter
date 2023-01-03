package database

// I'm not sure how I feel about this, but ok.
import "github.com/miodzie/walter/mods/rss/internal/usecases/adapters/storage/sqlite"

var Migration string

func init() {
	Migration = sqlite.Migration
}
