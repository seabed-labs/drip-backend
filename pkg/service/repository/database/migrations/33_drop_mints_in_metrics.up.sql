ALTER TABLE deposit_metric DROP COLUMN token_a_mint;
ALTER TABLE drip_metric DROP COLUMN token_a_mint;
ALTER TABLE drip_metric DROP COLUMN token_b_mint;
ALTER TABLE withdrawal_metric DROP COLUMN token_a_mint;
ALTER TABLE withdrawal_metric DROP COLUMN token_b_mint;

ALTER  TABLE account_update_queue_item ADD COLUMN reason TEXT NOT NULL DEFAULT 'default';
ALTER  TABLE transaction_update_queue_item ADD COLUMN reason TEXT NOT NULL DEFAULT 'default';
