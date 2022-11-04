ALTER TABLE vault ALTER COLUMN dca_activation_timestamp TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE position ALTER COLUMN deposit_timestamp TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE orca_whirlpool_delta_b_quote ALTER COLUMN last_updated TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE account_update_queue_item ALTER COLUMN time TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE account_update_queue_item ALTER COLUMN retry_time TYPE TIMESTAMP WITH TIME ZONE;