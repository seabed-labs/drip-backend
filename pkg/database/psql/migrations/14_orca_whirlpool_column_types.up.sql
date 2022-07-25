ALTER TABLE "orca_whirlpool"
    ALTER COLUMN tick_spacing TYPE integer,
    ALTER COLUMN fee_rate TYPE integer,
    ALTER COLUMN protocol_fee_rate TYPE integer,
    ALTER COLUMN protocol_fee_owed_a TYPE numeric,
    ALTER COLUMN protocol_fee_owed_b TYPE numeric,
    ALTER COLUMN reward_last_updated_timestamp TYPE numeric;