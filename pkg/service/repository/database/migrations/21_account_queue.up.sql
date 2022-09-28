CREATE TABLE "account_update_queue_item" (
    "pubkey" varchar(255) NOT NULL,
    "program_id" varchar(255) NOT NULL,
    "time" timestamp NOT NULL DEFAULT NOW(),
    UNIQUE(pubkey, program_id),
    PRIMARY KEY(pubkey)
);

CREATE INDEX index_account_update_queue_time ON account_update_queue_item(time);
