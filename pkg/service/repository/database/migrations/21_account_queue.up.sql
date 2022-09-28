CREATE TABLE "account_update_queue_item" (
    "pubkey" varchar(255) NOT NULL,
    "program_id" varchar(255) NOT NULL,
    "time" timestamp NOT NULL DEFAULT NOW(),
    "priority" int NOT NULL DEFAULT 3,
    UNIQUE(pubkey, program_id),
    PRIMARY KEY(pubkey)
);

CREATE INDEX index_account_update_queue_time ON account_update_queue_item(time);
CREATE INDEX index_account_update_queue_program_id ON account_update_queue_item(program_id);
