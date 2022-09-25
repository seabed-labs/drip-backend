
CREATE TABLE "orca_whirlpool" (
    "pubkey" varchar(255) PRIMARY KEY,
    "whirlpools_config" varchar(255) NOT NULL,
    "token_mint_a" varchar(255) NOT NULL,
    "token_vault_a" varchar(255) NOT NULL,
    "token_mint_b" varchar(255) NOT NULL,
    "token_vault_b" varchar(255) NOT NULL,
    "oracle" varchar(255) NOT NULL,

    "tick_spacing" integer NOT NULL,
    "fee_rate" integer NOT NULL,
    "protocol_fee_rate" integer NOT NULL,
    "tick_current_index" integer NOT NULL,

    "protocol_fee_owed_a" numeric NOT NULL,
    "protocol_fee_owed_b" numeric NOT NULL,
    "reward_last_updated_timestamp" numeric NOT NULL,
    "liquidity" numeric NOT NULL,
    "sqrt_price" numeric NOT NULL,
    "fee_growth_global_a" numeric NOT NULL,
    "fee_growth_global_b" numeric NOT NULL,

    "token_pair_id" uuid NOT NULL
);

ALTER TABLE "orca_whirlpool" ADD FOREIGN KEY ("token_pair_id") REFERENCES "token_pair" ("id");