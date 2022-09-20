TRUNCATE proto_config CASCADE;
ALTER TABLE proto_config RENAME COLUMN trigger_dca_spread TO token_a_drip_trigger_spread;
ALTER TABLE proto_config RENAME COLUMN base_withdrawal_spread TO token_b_withdrawal_spread;
ALTER TABLE proto_config ADD COLUMN token_b_referral_spread integer NOT NULL;
