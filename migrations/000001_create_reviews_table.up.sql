CREATE TABLE IF NOT EXISTS reviews (
    id serial PRIMARY KEY,
    name text NOT NULL,
    description text,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    review_url text NOT NULL DEFAULT 'https://www.ign.com/articles',
    review_score real NOT NULL DEFAULT 0.0,
    media_type text,
    genre_list text[],
    creator_list text[]
);