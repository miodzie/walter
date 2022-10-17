CREATE TABLE IF NOT EXISTS feeds (
  name        TEXT UNIQUE NOT NULL,
  url         TEXT UNIQUE NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS feed_subscriptions (
  feed_id  INT NOT NULL,
  channel  TEXT NOT NULL,
  user     TEXT NOT NULL,
  keywords TEXT NOT NULL,
  seen TEXT NOT NULL DEFAULT '',
  UNIQUE(feed_id, channel, user),
  FOREIGN KEY(feed_id) REFERENCES feeds(id)
) STRICT;
