TRUNCATE vault CASCADE;
TRUNCATE vault_period CASCADE;
TRUNCATE token_swap CASCADE;

ALTER TABLE token_swap ADD "token_a_mint" varchar(255) NOT NULL;
ALTER TABLE token_swap ADD "token_b_mint" varchar(255) NOT NULL;
