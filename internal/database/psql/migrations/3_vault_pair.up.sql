TRUNCATE vault CASCADE;

ALTER TABLE vault DROP COLUMN token_a_mint;
ALTER TABLE vault DROP COLUMN token_b_mint;
ALTER TABLE vault ADD COLUMN token_pair_id uuid NOT NULL;

ALTER TABLE "vault" ADD FOREIGN KEY ("token_pair_id") REFERENCES "token_pair" ("id");

