TRUNCATE orca_whirlpool CASCADE;

ALTER TABLE orca_whirlpool DROP CONSTRAINT pubkey_token_mint_a_token_mint_b_unique;
ALTER TABLE orca_whirlpool ADD CONSTRAINT pubkey_token_pair_id UNIQUE ("pubkey", "token_pair_id");

