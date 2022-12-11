CREATE TABLE "oracle_config" (
  "pubkey" varchar(64) PRIMARY KEY,
  "enabled" bool NOT NULL,
  "source" smallint NOT NULL,
  "update_authority" varchar(255) NOT NULL,
  "token_a_mint" varchar(255) NOT NULL,
  "token_a_price" varchar(255) NOT NULL,
  "token_b_mint" varchar(255) NOT NULL,
  "token_b_price" varchar(255) NOT NULL
);

ALTER TABLE vault ADD COLUMN oracle_config varchar(255);
ALTER TABLE vault ADD COLUMN max_price_deviation_bps integer NOT NULL DEFAULT 0;
ALTER TABLE vault ADD FOREIGN KEY ("oracle_config") REFERENCES "oracle_config" ("pubkey");
