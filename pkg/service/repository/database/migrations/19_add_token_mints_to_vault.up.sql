TRUNCATE vault CASCADE;

ALTER TABLE vault add column token_a_mint varchar(255) NOT NULL;
ALTER table vault add column token_b_mint varchar(255) NOT NULL;

ALTER TABLE vault ADD FOREIGN KEY ("token_a_mint") REFERENCES "token" ("pubkey");
ALTER TABLE vault ADD FOREIGN KEY ("token_b_mint") REFERENCES "token" ("pubkey");
