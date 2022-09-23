TRUNCATE proto_config CASCADE;

ALTER TABLE proto_config ADD COLUMN admin varchar(255) NOT NULL;

ALTER TABLE vault ADD COLUMN max_slippage_bps integer NOT NULL;
