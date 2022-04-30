CREATE TABLE IF NOT EXISTS feeds (
  id          INT AUTO_INCREMENT NOT NULL,
  name        VARCHAR(255) UNIQUE NOT NULL,
  url         VARCHAR(255) UNIQUE NOT NULL,
  PRIMARY KEY ('id')
);

CREATE TABLE IF NOT EXISTS feed_subscriptions (
  id       INT AUTO_INCREMENT NOT NULL,
  feed_id  INT NOT NULL,
  channel  VARCHAR(255) NOT NULL,
  user     VARCHAR(255) NOT NULL,
  keywords VARCHAR(255) NOT NULL,
  UNIQUE(feed_id, channel, user),
  FOREIGN KEY(feed_id) REFERENCES feeds(id)
);
