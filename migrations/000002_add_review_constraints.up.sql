ALTER TABLE reviews ADD CONSTRAINT review_score_check CHECK (review_score > 0 AND review_score <= 10);
