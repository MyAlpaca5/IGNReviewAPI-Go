-- create Generalized Inverted Index
CREATE INDEX IF NOT EXISTS users_username_idx ON users USING GIN (username);