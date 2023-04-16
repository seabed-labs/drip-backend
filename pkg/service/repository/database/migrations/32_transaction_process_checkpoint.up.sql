
CREATE TABLE "transaction_processing_checkpoint" (
    "signature" VARCHAR(255) NOT NULL,
    "slot" NUMERIC NOT NULL,
    PRIMARY KEY("signature")
);


CREATE TABLE "transaction_update_queue_item" (
    "signature" VARCHAR(255) NOT NULL,
    "tx_json" json NOT NULL,
    "time" timestamp NOT NULL DEFAULT NOW(),
    "priority" int NOT NULL DEFAULT 3,
    "try" int NOT NULL DEFAULT 0,
    "max_try" int NOT NULL DEFAULT 3,
    "retry_time" TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY("signature")
);

CREATE INDEX index_transaction_update_queue_time ON transaction_update_queue_item(priority);
CREATE INDEX index_transaction_update_queue_retry_time ON transaction_update_queue_item(retry_time);

