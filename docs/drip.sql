CREATE TABLE "vault" (
  "pubkey" varchar(64) PRIMARY KEY,
  "proto_config" varchar(255) NOT NULL,
  "token_a_mint" varchar(255) NOT NULL,
  "token_b_mint" varchar(255) NOT NULL,
  "token_a_account" varchar(255) NOT NULL,
  "token_b_account" varchar(255) NOT NULL,
  "treasury_token_b_account" varchar(255) NOT NULL,
  "last_dca_period" numeric NOT NULL,
  "drip_amount" numeric NOT NULL,
  "dca_activation_timestamp" int64 NOT NULL
);

CREATE TABLE "vault_period" (
  "pubkey" varchar(255) PRIMARY KEY,
  "vault" varchar(255) NOT NULL,
  "period_id" numeric NOT NULL,
  "twap" numeric NOT NULL,
  "dar" numeric NOT NULL
);

CREATE TABLE "proto_config" (
  "pubkey" varchar(255) PRIMARY KEY,
  "granularity" numeric NOT NULL,
  "trigger_dca_spread" smallint NOT NULL,
  "base_withdrawal_spread" smallint NOT NULL
);

CREATE TABLE "position" (
  "pubkey" varchar(255) PRIMARY KEY,
  "vault" varchar(255) NOT NULL,
  "authority" varchar(255) NOT NULL,
  "deposited_token_a_amount" numeric NOT NULL,
  "withdrawn_token_b_amount" numeric NOT NULL,
  "deposit_timestamp" int64 NOT NULL,
  "dca_period_id_before_deposit" numeric NOT NULL,
  "number_of_swaps" numeric NOT NULL,
  "periodic_drip_amount" numeric NOT NULL,
  "is_closed" boolean NOT NULL
);

CREATE TABLE "user_position" (
  "pubkey" varchar(255) PRIMARY KEY,
  "mint" varchar(255) NOT NULL,
  "amount" boolean NOT NULL
);

CREATE TABLE "token" (
  "pubkey" varchar(255) PRIMARY KEY,
  "symbol" varchar(255),
  "decimals" smallint NOT NULL,
  "icon_url" varchar(255)
);

CREATE TABLE "token_price" (
  "id" SERIAL PRIMARY KEY,
  "base" varchar(255),
  "quote" varchar(255),
  "close" numeric NOT NULL,
  "date" datetime NOT NULL,
  "source" varchar(255)
);

CREATE TABLE "source_reference" (
  "value" varchar(255) PRIMARY KEY
);

CREATE TABLE "token_pair" (
  "id" uuid PRIMARY KEY,
  "token_a" varchar(255) NOT NULL,
  "token_b" varchar(255) NOT NULL
);

CREATE TABLE "token_swap" (
  "pubkey" varchar(255) PRIMARY KEY,
  "mint" varchar(255) NOT NULL,
  "authority" varchar(255) NOT NULL,
  "fee_account" varchar(255) NOT NULL,
  "token_a_account" varchar(255) NOT NULL,
  "token_b_account" varchar(255) NOT NULL,
  "pair" uuid NOT NULL
);

ALTER TABLE "vault" ADD FOREIGN KEY ("proto_config") REFERENCES "proto_config" ("pubkey");

ALTER TABLE "vault_period" ADD FOREIGN KEY ("vault") REFERENCES "vault" ("pubkey");

ALTER TABLE "position" ADD FOREIGN KEY ("vault") REFERENCES "vault" ("pubkey");

ALTER TABLE "user_position" ADD FOREIGN KEY ("mint") REFERENCES "position" ("authority");

ALTER TABLE "token_price" ADD FOREIGN KEY ("base") REFERENCES "token" ("pubkey");

ALTER TABLE "token_price" ADD FOREIGN KEY ("quote") REFERENCES "token" ("pubkey");

ALTER TABLE "token_price" ADD FOREIGN KEY ("source") REFERENCES "source_reference" ("value");

ALTER TABLE "token_pair" ADD FOREIGN KEY ("token_a") REFERENCES "token" ("pubkey");

ALTER TABLE "token_pair" ADD FOREIGN KEY ("token_b") REFERENCES "token" ("pubkey");

ALTER TABLE "token_swap" ADD FOREIGN KEY ("pair") REFERENCES "token_pair" ("id");
