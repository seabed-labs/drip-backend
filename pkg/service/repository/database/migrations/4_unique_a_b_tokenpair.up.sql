TRUNCATE token_pair CASCADE;

ALTER TABLE token_pair ADD UNIQUE (token_a, token_b);