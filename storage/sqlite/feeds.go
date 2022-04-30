package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/miodzie/seras/mods/rss"
)

type FeedRepository struct {
}

func (repo *FeedRepository) All() ([]rss.Feed, error) {
	rows, err := db.Query("SELECT rowid, * FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []rss.Feed
	for rows.Next() {
		var feed rss.Feed
		if err := rows.Scan(&feed.Id, &feed.Name, &feed.Url); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (repo *FeedRepository) Save(feed *rss.Feed) error {
    result, err := db.Exec("INSERT INTO feeds (name, url) VALUES(?, ?)", feed.Name, feed.Url)
    if err != nil {
        return fmt.Errorf("Save: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    feed.Id = uint64(id)

	return nil
}

func (repo *FeedRepository) GetByName(name string) (rss.Feed, error) {
	var feed rss.Feed
	row := db.QueryRow("SELECT rowid, * FROM feeds WHERE name = ?", name)
	if err := row.Scan(&feed.Id, &feed.Name, &feed.Url); err != nil {
		if err == sql.ErrNoRows {
			return feed, fmt.Errorf("GetByName %s: no such feed", name)
		}
		return feed, err
	}

	return feed, nil
}
