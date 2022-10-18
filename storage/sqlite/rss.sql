CREATE TABLE IF NOT EXISTS feeds
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    url  TEXT UNIQUE NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS feed_subscriptions
(
    id       INTEGER PRIMARY KEY,
    feed_id  INT  NOT NULL,
    channel  TEXT NOT NULL,
    user     TEXT NOT NULL,
    keywords TEXT NOT NULL,
    ignore   TEXT NOT NULL DEFAULT '',
    seen     TEXT NOT NULL DEFAULT '',
    UNIQUE (feed_id, channel, user),
    FOREIGN KEY (feed_id) REFERENCES feeds (id)
) STRICT;
