/*
 * Copyright 2022-present miodzie. All rights reserved.
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

CREATE TABLE IF NOT EXISTS feeds
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    url  TEXT UNIQUE NOT NULL
);

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
);
