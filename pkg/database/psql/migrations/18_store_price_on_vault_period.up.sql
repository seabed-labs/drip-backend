TRUNCATE vault_period CASCADE;
ALTER TABLE vault_period ADD COLUMN price_b_over_a numeric DEFAULT NULL;
