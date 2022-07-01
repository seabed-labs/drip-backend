TRUNCATE token_swap CASCADE;

ALTER TABLE token_swap RENAME COLUMN "pair" TO "token_pair_id";

ALTER TABLE token_swap ADD "id" uuid NOT NULL;
ALTER TABLE token_swap DROP CONSTRAINT token_swap_pkey;
ALTER TABLE token_swap ADD CONSTRAINT id_pkey PRIMARY KEY (id);


ALTER TABLE token_swap ADD CONSTRAINT pubkey_token_a_mint_token_b_mint_unique UNIQUE ("pubkey", "token_a_mint", "token_b_mint");

