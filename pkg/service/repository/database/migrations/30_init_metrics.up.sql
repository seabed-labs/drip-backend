-- No foreign keys are used intentionally
-- the data in these tables is pushed to primarily from the metrics server
-- the event server/primary backend populates price

CREATE TABLE "deposit_metric" (
    "signature" VARCHAR(255) NOT NULL,
    "ix_index" INTEGER NOT NULL,
    "ix_name" VARCHAR(32) NOT NULL,
    "ix_version" INTEGER NOT NULL,
    "slot" INTEGER NOT NULL,
    "time" TIMESTAMP WITH TIME ZONE NOT NULL,

    "vault" VARCHAR(255) NOT NULL,
    "token_a_mint" VARCHAR(255) NOT NULL,
    "referrer" VARCHAR(255) NULL,

    "token_a_deposit_amount" NUMERIC NOT NULL,
    "token_a_usd_price_day" NUMERIC NULL,

    PRIMARY KEY("signature", "ix_index")
);

CREATE TABLE "drip_metric" (
    "signature" VARCHAR(255) NOT NULL,
    "ix_index" INTEGER NOT NULL,
    "ix_name" VARCHAR(32) NOT NULL,
    "ix_version" INTEGER NOT NULL,
    "slot" INTEGER NOT NULL,
    "time" TIMESTAMP WITH TIME ZONE NOT NULL,

    "vault" VARCHAR(255) NOT NULL,
    "token_a_mint" VARCHAR(255) NOT NULL,
    "token_b_mint" VARCHAR(255) NOT NULL,

    "vault_token_a_swapped_amount" NUMERIC NOT NULL,
    "vault_token_b_received_amount" NUMERIC NOT NULL,
    "keeper_token_a_received_amount" NUMERIC NOT NULL,
    "token_a_usd_price_day" NUMERIC NULL,
    "token_b_usd_price_day" NUMERIC NULL,

    PRIMARY KEY("signature", "ix_index")
);

CREATE TABLE "withdrawal_metric" (
    "signature" VARCHAR(255) NOT NULL,
    "ix_index" INTEGER NOT NULL,
    "ix_name" VARCHAR(32) NOT NULL,
    "ix_version" INTEGER NOT NULL,
    "slot" INTEGER NOT NULL,
    "time" TIMESTAMP WITH TIME ZONE NOT NULL,

    "vault" VARCHAR(255) NOT NULL,
    "token_a_mint" VARCHAR(255) NULL,
    "token_b_mint" VARCHAR(255) NOT NULL,

    "user_token_a_withdraw_amount" NUMERIC NOT NULL,
    "user_token_b_withdraw_amount" NUMERIC NOT NULL,
    "treasury_token_b_received_amount" NUMERIC NOT NULL,
    "referral_token_b_received_amount" NUMERIC NOT NULL,
    "token_a_usd_price_day" NUMERIC NULL,
    "token_b_usd_price_day" NUMERIC NULL,

    PRIMARY KEY("signature", "ix_index")
);

