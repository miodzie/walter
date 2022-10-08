CREATE TABLE IF NOT EXISTS feeds (
  name        VARCHAR(255) UNIQUE NOT NULL,
  url         VARCHAR(255) UNIQUE NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS feed_subscriptions (
  feed_id  INT NOT NULL,
  channel  VARCHAR(255) NOT NULL,
  user     VARCHAR(255) NOT NULL,
  keywords VARCHAR(255) NOT NULL,
  seen VARCHAR(10000) NOT NULL DEFAULT '',
  UNIQUE(feed_id, channel, user),
  FOREIGN KEY(feed_id) REFERENCES feeds(id)
) STRICT;
