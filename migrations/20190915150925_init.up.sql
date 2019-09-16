CREATE TABLE IF NOT EXISTS "user"
(
    id            SERIAL PRIMARY KEY,
    first_name    TEXT,
    middle_name   TEXT,
    last_name     TEXT,
    username      TEXT NOT NULL,
    password_hash TEXT NOT NULL
);