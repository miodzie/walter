package sqlite

import (
	"fmt"

	"github.com/miodzie/seras/mods/rss"
)

type SubscriptionRepository struct {
}

func (repo *SubscriptionRepository) Add(sub *rss.Subscription) error {
	q := "INSERT INTO feed_subscriptions (feed_id, channel, user, keywords) VALUES(?,?,?,?)"
	result, err := db.Exec(q, sub.FeedId, sub.Channel, sub.User, sub.Keywords)
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

func (repo *SubscriptionRepository) ByFeedId(id uint64) ([]*rss.Subscription, error) {
	rows, err := db.Query("SELECT rowid, * FROM feed_subscriptions WHERE feed_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*rss.Subscription
	for rows.Next() {
		var sub *rss.Subscription
		if err := rows.Scan(&sub.Id, &sub.FeedId, &sub.Channel, &sub.User, &sub.Keywords); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}
