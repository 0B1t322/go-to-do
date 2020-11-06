CREATE TABLE IF NOT EXISTS users (
    _id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    token TEXT );

CREATE TABLE IF NOT EXISTS tasks (
    _id         INTEGER PRIMARY KEY,
    user_id     INTEGER NOT NULL,
    name        TEXT NOT NULL,
    decription  TEXT NOT NULL,
    DATE        TEXT NOT NULL,
    DONE        BOOL NOT NULL);