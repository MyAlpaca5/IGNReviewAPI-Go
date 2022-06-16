CREATE TABLE IF NOT EXISTS tokens (
    token bytea PRIMARY KEY,
    userID int NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp NOT NULL DEFAULT (NOW() + interval '1 day'),
    role int NOT NULL
);