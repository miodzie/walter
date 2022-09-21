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
			return &feed, fmt.Errorf("GetByName %s: no such feed", name)
		}
		return &feed, err
	}

	return &feed, nil
}

func (r *RssRepository) AddSub(sub *rss.Subscription) error {
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

func (r *RssRepository) SubByUserFeedNameChannel(user, feedName, channel string) (*rss.Subscription, error) {
	feed, err := r.FeedByName(feedName)
	if feed == nil || err != nil {
		return nil, fmt.Errorf("failed to locate feed with name: `%s`", feedName)
	}
	q := "SELECT rowid, * FROM feed_subscriptions WHERE feed_id = ? AND user = ? AND channel = ?"
	var sub rss.Subscription

	row := r.db.QueryRow(q, feed.Id, user, channel)
	if err := row.Scan(&sub.Id, &sub.Channel, &sub.User, &sub.Keywords, &sub.Seen); err != nil {
		if err == sql.ErrNoRows {
			return &sub, fmt.Errorf("subscription not found")
		}
		return &sub, err
	}

	return &sub, nil
}

func (r *RssRepository) SubsByFeedId(id uint64) ([]*rss.Subscription, error) {
	rows, err := r.db.Query("SELECT rowid, * FROM feed_subscriptions WHERE feed_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*rss.Subscription
	for rows.Next() {
		var sub rss.Subscription
		if err := rows.Scan(&sub.Id, &sub.FeedId, &sub.Channel, &sub.User, &sub.Keywords, &sub.Seen); err != nil {
			return nil, err
		}
		subs = append(subs, &sub)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}
