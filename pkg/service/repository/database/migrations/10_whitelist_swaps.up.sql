CREATE TABLE "vault_whitelist" (
    "id" uuid PRIMARY KEY,
    "vault_pubkey" varchar(64) NOT NULL,
    "token_swap_pubkey" varchar(64) NOT NULL
);

ALTER TABLE vault_whitelist ADD CONSTRAINT vault_pubkey_token_swap_pubkey_unique UNIQUE ("vault_pubkey", "token_swap_pubkey");
