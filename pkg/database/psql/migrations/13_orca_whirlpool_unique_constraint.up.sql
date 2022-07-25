TRUNCATE orca_whirlpool CASCADE;

--  TODO(Mocha): We can make pubkey the primary key
ALTER TABLE orca_whirlpool ADD "id" uuid NOT NULL;
ALTER TABLE orca_whirlpool DROP CONSTRAINT orca_whirlpool_pkey;
ALTER TABLE orca_whirlpool ADD CONSTRAINT orca_whirlpool_pkey PRIMARY KEY (id);

ALTER TABLE orca_whirlpool ADD CONSTRAINT pubkey_token_mint_a_token_mint_b_unique UNIQUE ("pubkey", "token_mint_a", "token_mint_b");
