CREATE TABLE IF NOT EXISTS messages
(
    channel VARCHAR(255)   NOT NULL,
    user    VARCHAR(255)   NOT NULL,
    content VARCHAR(255)   NOT NULL,
    raw     VARCHAR(10000) NOT NULL,
    created VARCHAR(30)    NOT NULL
) STRICT;
