
CREATE TABLE if not exists users
(
    id INTEGER PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS apps
(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    secret TEXT NOT NULL
);