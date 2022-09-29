ALTER TABLE account_update_queue_item ADD COLUMN try int NOT NULL DEFAULT 0;
ALTER TABLE account_update_queue_item ADD COLUMN max_try int NOT NULL DEFAULT 3;
ALTER TABLE account_update_queue_item ADD COLUMN retry_time timestamp;

CREATE INDEX index_account_update_queue_retry_time ON account_update_queue_item(retry_time);
