package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/miodzie/seras/mods/rss"
)

type RssRepository struct {
	db *sql.DB
}

func NewRssRepository(db *sql.DB) *RssRepository {
	return &RssRepository{db: db}
}

func (r *RssRepository) Feeds() ([]*rss.Feed, error) {
	rows, err := r.db.Query("SELECT rowid, * FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []*rss.Feed
	for rows.Next() {
		var feed rss.Feed
		if err := rows.Scan(&feed.Id, &feed.Name, &feed.Url); err != nil {
			return nil, err
		}
		feeds = append(feeds, &feed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (r *RssRepository) AddFeed(feed *rss.Feed) error {
	result, err := r.db.Exec("INSERT INTO feeds (name, url) VALUES(?, ?)", feed.Name, feed.Url)
	if err != nil {
		return fmt.Errorf("FeedRepository.add: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	feed.Id = uint64(id)

	return nil
}

func (r *RssRepository) FeedByName(name string) (*rss.Feed, error) {
	var feed rss.Feed
	row := r.db.QueryRow("SELECT rowid, * FROM feeds WHERE name = ?", name)
	if err := row.Scan(&feed.Id, &feed.Name, &feed.Url); err != nil {
		if err == sql.ErrNoRows {
			return &feed, fmt.Errorf("unable to locate feed with name: %s", name)
		}
		return &feed, err
	}

	return &feed, nil
}

func (r *RssRepository) AddSub(sub *rss.Subscription) error {
	// TODO: Check for duplicate for better error response.
	q := "INSERT INTO feed_subscriptions (feed_id, channel, user, keywords, seen) VALUES(?,?,?,?,?)"
	result, err := r.db.Exec(q, sub.FeedId, sub.Channel, sub.User, sub.Keywords, sub.Seen)
	if err != nil {
		return fmt.Errorf("SubscriptionRepository.add: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sub.Id = uint64(id)

	return nil
}

func (r *RssRepository) UpdateSub(sub *rss.Subscription) error {
	q := "UPDATE feed_subscriptions SET feed_id = ?, channel = ?, user = ?, keywords = ?, seen = ? WHERE rowid = ?"
	_, err := r.db.Exec(q, sub.FeedId, sub.Channel, sub.User, sub.Keywords, sub.Seen, sub.Id)
	if err != nil {
		return fmt.Errorf("SubscriptionRepository.Update: %v", err)
	}

	return nil
}

func (r *RssRepository) RemoveSub(subscription *rss.Subscription) error {
	_, err := r.db.Exec("DELETE FROM feed_subscriptions WHERE rowid = ?", subscription.Id)
	return err
}

func (r *RssRepository) Subs(search rss.SubSearchOpt) ([]*rss.Subscription, error) {
	var subs []*rss.Subscription
	var args []interface{}
	query := "SELECT rowid, * from feed_subscriptions WHERE 1 = 1"
	if search.User != "" {
		query += " AND user = ?"
		args = append(args, search.User)
	}
	if search.Channel != "" {
		query += " AND channel = ?"
		args = append(args, search.Channel)
	}
	if search.FeedId != 0 {
		query += " AND feed_id = ?"
		args = append(args, search.FeedId)
	}
	if search.FeedName != "" {
		feed, err := r.FeedByName(search.FeedName)
		if err != nil {
			return subs, err
		}
		query += " AND feed_id = ?"
		args = append(args, feed.Id)
	}
	fmt.Println(query, args)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		sub, err := r.scanSub(rows)
		if err != nil {
			return subs, err
		}
		subs = append(subs, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *RssRepository) scanSub(rows *sql.Rows) (*rss.Subscription, error) {
	var sub rss.Subscription
	err := rows.Scan(&sub.Id, &sub.FeedId, &sub.Channel, &sub.User, &sub.Keywords, &sub.Seen)
	if err != nil {
		return nil, err
	}

	feeds, _ := r.Feeds()
	for _, f := range feeds {
		if f.Id == sub.FeedId {
			sub.Feed = f
			break
		}
	}

	return &sub, nil
}
