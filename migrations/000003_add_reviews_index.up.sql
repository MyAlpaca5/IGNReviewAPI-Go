-- create Generalized Inverted Index
CREATE INDEX IF NOT EXISTS reviews_name_idx ON reviews USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS reviews_genre_list_idx ON reviews USING GIN (genre_list);