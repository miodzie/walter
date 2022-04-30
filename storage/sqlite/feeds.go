package sqlite

import "github.com/miodzie/seras/mods/rss"

type FeedRepository struct {
}

func (repo *FeedRepository) All() ([]rss.Feed, error) {
	rows, err := db.Query("SELECT * FROM feeds")
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
