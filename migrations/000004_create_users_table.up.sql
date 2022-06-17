CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    username text UNIQUE NOT NULL,
    password bytea NOT NULL,
    email citext,
    role int NOT NULL
);