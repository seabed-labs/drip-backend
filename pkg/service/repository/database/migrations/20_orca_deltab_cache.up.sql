CREATE TABLE "orca_whirlpool_delta_b_quote" (
    "vault_pubkey" varchar(255) NOT NULL,
    "whirlpool_pubkey" varchar(255) NOT NULL,
    "token_pair_id" uuid NOT NULL,
    "delta_b" numeric NOT NULL,
    "last_updated" timestamp NOT NULL DEFAULT NOW(),
    UNIQUE(vault_pubkey, whirlpool_pubkey),
    PRIMARY KEY(vault_pubkey, whirlpool_pubkey)
);

ALTER TABLE "orca_whirlpool_delta_b_quote" ADD FOREIGN KEY ("vault_pubkey") REFERENCES "vault" ("pubkey");
ALTER TABLE "orca_whirlpool_delta_b_quote" ADD FOREIGN KEY ("token_pair_id") REFERENCES "token_pair" ("id");